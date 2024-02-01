package keeper_test

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v6/testutil/helpers"
	"github.com/functionx/fx-core/v6/x/crosschain/types"
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
		suite.Keeper().StoreOracleSet(suite.ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           uint64(i),
				ExternalAddress: fmt.Sprintf("0x%d", i),
			}},
			Height: uint64(i),
		})
	}

	wrapSDKContext := sdk.WrapSDKContext(suite.ctx)
	for _, testCase := range testCases {
		oracle := types.Oracle{
			OracleAddress:  testCase.OracleAddress.String(),
			BridgerAddress: testCase.BridgerAddress.String(),
			StartHeight:    testCase.StartHeight,
		}
		// save oracle
		suite.Keeper().SetOracle(suite.ctx, oracle)
		suite.Keeper().SetOracleByBridger(suite.ctx, testCase.BridgerAddress, oracle.GetOracle())

		response, err := suite.Keeper().LastPendingOracleSetRequestByAddr(wrapSDKContext,
			&types.QueryLastPendingOracleSetRequestByAddrRequest{
				BridgerAddress: testCase.BridgerAddress.String(),
			})
		require.NoError(suite.T(), err)
		require.EqualValues(suite.T(), testCase.ExpectOracleSetSize, len(response.OracleSets))
	}
}

func (suite *KeeperTestSuite) TestGetUnSlashedOracleSets() {
	height := 100
	index := 10
	for i := 1; i <= index; i++ {
		suite.Keeper().StoreOracleSet(suite.ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           tmrand.Uint64(),
				ExternalAddress: helpers.GenerateAddress().Hex(),
			}},
			Height: uint64(height + i),
		})
	}
	suite.Equal(uint64(0), suite.Keeper().GetLastSlashedOracleSetNonce(suite.ctx))

	sets := suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(height+index))
	require.EqualValues(suite.T(), index-1, sets.Len())

	suite.Keeper().SetLastSlashedOracleSetNonce(suite.ctx, 1)
	sets = suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(height+index))
	require.EqualValues(suite.T(), index-2, sets.Len())

	sets = suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(height+index+1))
	require.EqualValues(suite.T(), index-1, sets.Len())
}

func (suite *KeeperTestSuite) TestKeeper_IterateOracleSetConfirmByNonce() {
	index := tmrand.Intn(20) + 1
	for i := uint64(1); i <= uint64(index); i++ {
		for _, oracle := range suite.oracleAddrs {
			suite.Keeper().SetOracleSetConfirm(suite.ctx, oracle,
				&types.MsgOracleSetConfirm{
					Nonce:           i,
					BridgerAddress:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
					ExternalAddress: helpers.GenerateAddress().Hex(),
					Signature:       "",
					ChainName:       suite.chainName,
				},
			)
		}
	}

	index = tmrand.Intn(index) + 1
	var confirms []*types.MsgOracleSetConfirm
	suite.Keeper().IterateOracleSetConfirmByNonce(suite.ctx, uint64(index), func(confirm *types.MsgOracleSetConfirm) bool {
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
		Height:  uint64(suite.ctx.BlockHeight()),
	}
	suite.Keeper().StoreOracleSet(suite.ctx, oracleSet)

	for i, external := range suite.externalPris {
		externalAddress, signature := suite.SignOracleSetConfirm(external, oracleSet)
		suite.Keeper().SetOracleSetConfirm(suite.ctx, suite.oracleAddrs[i],
			&types.MsgOracleSetConfirm{
				Nonce:           oracleSet.Nonce,
				BridgerAddress:  suite.bridgerAddrs[i].String(),
				ExternalAddress: externalAddress,
				Signature:       hex.EncodeToString(signature),
				ChainName:       suite.chainName,
			},
		)
	}
	suite.Keeper().SetLastObservedOracleSet(suite.ctx,
		&types.OracleSet{
			Nonce: oracleSet.Nonce + 1,
		},
	)

	params := suite.Keeper().GetParams(suite.ctx)
	params.SignedWindow = 10
	err := suite.Keeper().SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	suite.Commit()
	for _, oracle := range suite.oracleAddrs {
		suite.NotNil(suite.Keeper().GetOracleSetConfirm(suite.ctx, oracleSet.Nonce, oracle))
	}

	suite.Commit(int64(params.SignedWindow + 1))
	for _, oracle := range suite.oracleAddrs {
		suite.Nil(suite.Keeper().GetOracleSetConfirm(suite.ctx, oracleSet.Nonce, oracle))
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
		suite.Keeper().StoreOracleSet(suite.ctx, &types.OracleSet{
			Nonce:   uint64(i),
			Members: member,
			Height:  uint64(i + 100),
		})
	}
	i := uint64(0)
	oracleSets := types.OracleSets{}
	suite.Keeper().IterateOracleSetByNonce(suite.ctx, 0, func(oracleSet *types.OracleSet) bool {
		i = i + 1
		suite.Equal(i, oracleSet.Nonce)
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 10)

	oracleSets = types.OracleSets{}
	suite.Keeper().IterateOracleSetByNonce(suite.ctx, 1, func(oracleSet *types.OracleSet) bool {
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 10)

	oracleSets = types.OracleSets{}
	suite.Keeper().IterateOracleSetByNonce(suite.ctx, 2, func(oracleSet *types.OracleSet) bool {
		oracleSets = append(oracleSets, oracleSet)
		return false
	})
	suite.Equal(len(oracleSets), 9)

	suite.Keeper().IterateOracleSets(suite.ctx, true, func(oracleSet *types.OracleSet) bool {
		suite.Equal(i, oracleSet.Nonce, oracleSet.Nonce)
		i = i - 1
		return false
	})

	suite.Keeper().IterateOracleSets(suite.ctx, false, func(oracleSet *types.OracleSet) bool {
		i = i + 1
		suite.Equal(i, oracleSet.Nonce)
		return false
	})
}
