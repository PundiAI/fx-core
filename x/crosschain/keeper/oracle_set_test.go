package keeper_test

import (
	"encoding/hex"
	"fmt"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestLastPendingOracleSetRequestByAddr() {
	testCases := []struct {
		OracleAddress  sdk.AccAddress
		BridgerAddress sdk.AccAddress
		StartHeight    int64

		ExpectOracleSetSize int
	}{
		{
			OracleAddress:       suite.oracleAddrs[0],
			BridgerAddress:      suite.bridgerAddrs[0],
			StartHeight:         1,
			ExpectOracleSetSize: 3,
		},
		{
			OracleAddress:       suite.oracleAddrs[1],
			BridgerAddress:      suite.bridgerAddrs[1],
			StartHeight:         2,
			ExpectOracleSetSize: 2,
		},
		{
			OracleAddress:       suite.oracleAddrs[2],
			BridgerAddress:      suite.bridgerAddrs[2],
			StartHeight:         3,
			ExpectOracleSetSize: 1,
		},
	}

	for i := 1; i <= 3; i++ {
		suite.Keeper().StoreOracleSet(suite.Ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           uint64(i),
				ExternalAddress: fmt.Sprintf("0x%d", i),
			}},
			Height: uint64(i),
		})
	}

	wrapSDKContext := suite.Ctx
	for _, testCase := range testCases {
		oracle := types.Oracle{
			OracleAddress:  testCase.OracleAddress.String(),
			BridgerAddress: testCase.BridgerAddress.String(),
			StartHeight:    testCase.StartHeight,
		}
		// save oracle
		suite.Keeper().SetOracle(suite.Ctx, oracle)
		suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, testCase.BridgerAddress, oracle.GetOracle())

		response, err := suite.QueryClient().LastPendingOracleSetRequestByAddr(wrapSDKContext,
			&types.QueryLastPendingOracleSetRequestByAddrRequest{
				BridgerAddress: testCase.BridgerAddress.String(),
			})
		suite.Require().NoError(err)
		suite.Require().EqualValues(testCase.ExpectOracleSetSize, len(response.OracleSets))
	}
}

func (suite *KeeperTestSuite) TestGetUnSlashedOracleSets() {
	height := 100
	index := 10
	for i := 1; i <= index; i++ {
		suite.Keeper().StoreOracleSet(suite.Ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           tmrand.Uint64(),
				ExternalAddress: helpers.GenHexAddress().Hex(),
			}},
			Height: uint64(height + i),
		})
	}
	suite.Equal(uint64(0), suite.Keeper().GetLastSlashedOracleSetNonce(suite.Ctx))

	sets := suite.Keeper().GetUnSlashedOracleSets(suite.Ctx, uint64(height+index))
	suite.Require().EqualValues(index-1, sets.Len())

	suite.Keeper().SetLastSlashedOracleSetNonce(suite.Ctx, 1)
	sets = suite.Keeper().GetUnSlashedOracleSets(suite.Ctx, uint64(height+index))
	suite.Require().EqualValues(index-2, sets.Len())

	sets = suite.Keeper().GetUnSlashedOracleSets(suite.Ctx, uint64(height+index+1))
	suite.Require().EqualValues(index-1, sets.Len())
}

func (suite *KeeperTestSuite) TestKeeper_IterateOracleSetConfirmByNonce() {
	index := tmrand.Intn(20) + 1
	for i := uint64(1); i <= uint64(index); i++ {
		for _, oracle := range suite.oracleAddrs {
			suite.Keeper().SetOracleSetConfirm(suite.Ctx, oracle,
				&types.MsgOracleSetConfirm{
					Nonce:           i,
					BridgerAddress:  helpers.GenAccAddress().String(),
					ExternalAddress: helpers.GenHexAddress().Hex(),
					Signature:       "",
					ChainName:       suite.chainName,
				},
			)
		}
	}

	index = tmrand.Intn(index) + 1
	var confirms []*types.MsgOracleSetConfirm
	suite.Keeper().IterateOracleSetConfirmByNonce(suite.Ctx, uint64(index), func(confirm *types.MsgOracleSetConfirm) bool {
		confirms = append(confirms, confirm)
		return false
	})
	suite.Equal(len(confirms), len(suite.oracleAddrs), index)
}

func (suite *KeeperTestSuite) TestKeeper_DeleteOracleSetConfirm() {
	member := make([]types.BridgeValidator, 0, len(suite.externalPris))
	for i, external := range suite.externalPris {
		externalAddr := suite.PubKeyToExternalAddr(external.PublicKey)
		member = append(member, types.BridgeValidator{
			Power:           uint64(i),
			ExternalAddress: externalAddr,
		})
	}
	oracleSet := &types.OracleSet{
		Nonce:   1,
		Members: member,
		Height:  uint64(suite.Ctx.BlockHeight()),
	}
	suite.Keeper().StoreOracleSet(suite.Ctx, oracleSet)

	for i, external := range suite.externalPris {
		externalAddress, signature := suite.SignOracleSetConfirm(external, oracleSet)
		suite.Keeper().SetOracleSetConfirm(suite.Ctx, suite.oracleAddrs[i],
			&types.MsgOracleSetConfirm{
				Nonce:           oracleSet.Nonce,
				BridgerAddress:  suite.bridgerAddrs[i].String(),
				ExternalAddress: externalAddress,
				Signature:       hex.EncodeToString(signature),
				ChainName:       suite.chainName,
			},
		)
	}
	suite.Keeper().SetLastObservedOracleSet(suite.Ctx,
		&types.OracleSet{
			Nonce: oracleSet.Nonce + 1,
		},
	)

	params := suite.Keeper().GetParams(suite.Ctx)
	params.SignedWindow = 10
	err := suite.Keeper().SetParams(suite.Ctx, &params)
	suite.Require().NoError(err)
	suite.Commit()
	for _, oracle := range suite.oracleAddrs {
		suite.NotNil(suite.Keeper().GetOracleSetConfirm(suite.Ctx, oracleSet.Nonce, oracle))
	}

	suite.Commit(int64(params.SignedWindow + 1))
	for _, oracle := range suite.oracleAddrs {
		suite.Nil(suite.Keeper().GetOracleSetConfirm(suite.Ctx, oracleSet.Nonce, oracle))
	}
}

func (suite *KeeperTestSuite) TestKeeper_IterateOracleSet() {
	member := make([]types.BridgeValidator, 0, len(suite.externalPris))
	for i, external := range suite.externalPris {
		member = append(member, types.BridgeValidator{
			Power:           uint64(i),
			ExternalAddress: crypto.PubkeyToAddress(external.PublicKey).String(),
		})
	}
	for i := 1; i <= 10; i++ {
		suite.Keeper().StoreOracleSet(suite.Ctx, &types.OracleSet{
			Nonce:   uint64(i),
			Members: member,
			Height:  uint64(i + 100),
		})
	}
	i := uint64(0)
	oracleSets := types.OracleSets{}
	suite.Keeper().IterateOracleSetByNonce(suite.Ctx, 0, func(oracleSet *types.OracleSet) bool {
		i = i + 1
		suite.Equal(i, oracleSet.Nonce)
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 10)

	oracleSets = types.OracleSets{}
	suite.Keeper().IterateOracleSetByNonce(suite.Ctx, 1, func(oracleSet *types.OracleSet) bool {
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 10)

	oracleSets = types.OracleSets{}
	suite.Keeper().IterateOracleSetByNonce(suite.Ctx, 2, func(oracleSet *types.OracleSet) bool {
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 9)

	suite.Keeper().IterateOracleSets(suite.Ctx, true, func(oracleSet *types.OracleSet) bool {
		suite.Equal(i, oracleSet.Nonce, oracleSet.Nonce)
		i = i - 1
		return false
	})

	suite.Keeper().IterateOracleSets(suite.Ctx, false, func(oracleSet *types.OracleSet) bool {
		i = i + 1
		suite.Equal(i, oracleSet.Nonce)
		return false
	})
}

func (suite *KeeperTestSuite) TestKeeper_UpdateOracleSetExecuted() {
	claim := types.MsgOracleSetUpdatedClaim{
		OracleSetNonce: 0,
		Members:        types.BridgeValidators{{Power: 1, ExternalAddress: helpers.GenExternalAddr(suite.chainName)}},
		ChainName:      suite.chainName,
	}
	err := suite.Keeper().UpdateOracleSetExecuted(suite.Ctx, &claim)
	suite.Require().NoError(err)
	oracleSet := suite.Keeper().GetLastObservedOracleSet(suite.Ctx)
	suite.Require().NotEmpty(oracleSet)
	suite.Equal(claim.OracleSetNonce, oracleSet.Nonce)
	suite.Equal(claim.Members, oracleSet.Members)

	claim.OracleSetNonce = 1
	err = suite.Keeper().UpdateOracleSetExecuted(suite.Ctx, &claim)
	suite.Require().Error(err)
	suite.ErrorIs(types.ErrInvalid.Wrapf("attested oracleSet (%v) does not exist in store", claim.OracleSetNonce), err)

	oracleSet = &types.OracleSet{Nonce: claim.OracleSetNonce, Members: claim.Members, Height: 10}
	suite.Keeper().StoreOracleSet(suite.Ctx, oracleSet)
	err = suite.Keeper().UpdateOracleSetExecuted(suite.Ctx, &claim)
	suite.Require().NoError(err)
	lastOracleSet := suite.Keeper().GetLastObservedOracleSet(suite.Ctx)
	suite.Require().NotEmpty(oracleSet)
	suite.Equal(oracleSet.Nonce, lastOracleSet.Nonce)
	suite.Equal(oracleSet.Members, lastOracleSet.Members)
	suite.Equal(oracleSet.Height, lastOracleSet.Height)
}
