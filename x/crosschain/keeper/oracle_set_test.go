package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
	types2 "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func (suite *KeeperTestSuite) TestLastPendingOracleSetRequestByAddr() {
	testCases := []struct {
		OracleAddress  sdk.AccAddress
		BridgerAddress sdk.AccAddress
		StartHeight    int64

		ExpectOracleSetSize int
	}{
		{
			OracleAddress:       suite.oracles[0],
			BridgerAddress:      suite.bridgers[0],
			StartHeight:         1,
			ExpectOracleSetSize: 3,
		},
		{
			OracleAddress:       suite.oracles[1],
			BridgerAddress:      suite.bridgers[1],
			StartHeight:         2,
			ExpectOracleSetSize: 2,
		},
		{
			OracleAddress:       suite.oracles[2],
			BridgerAddress:      suite.bridgers[2],
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
				Power:           rand.Uint64(),
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
	index := rand.Intn(20) + 1
	for i := uint64(1); i <= uint64(index); i++ {
		for _, oracle := range suite.oracles {
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

	index = rand.Intn(index) + 1
	var confirms []*types.MsgOracleSetConfirm
	suite.Keeper().IterateOracleSetConfirmByNonce(suite.ctx, uint64(index), func(confirm *types.MsgOracleSetConfirm) bool {
		confirms = append(confirms, confirm)
		return false
	})
	suite.Equal(len(confirms), len(suite.oracles), index)
}

func (suite *KeeperTestSuite) TestKeeper_DeleteOracleSetConfirm() {
	var member []types.BridgeValidator
	for i, external := range suite.externals {
		externalAddr := crypto.PubkeyToAddress(external.PublicKey).String()
		if suite.chainName == trontypes.ModuleName {
			externalAddr = tronAddress.PubkeyToAddress(external.PublicKey).String()
		}

		member = append(member, types.BridgeValidator{
			Power:           uint64(i),
			ExternalAddress: externalAddr,
		})
	}
	oracleSet := &types.OracleSet{
		Nonce:   1,
		Members: member,
		Height:  100,
	}
	suite.Keeper().StoreOracleSet(suite.ctx, oracleSet)

	for i, external := range suite.externals {
		externalAddress := crypto.PubkeyToAddress(external.PublicKey).String()
		gravityId := suite.Keeper().GetGravityID(suite.ctx)
		checkpoint, err := oracleSet.GetCheckpoint(gravityId)
		suite.NoError(err)
		signature, err := types.NewEthereumSignature(checkpoint, external)
		suite.NoError(err)
		if trontypes.ModuleName == suite.chainName {
			externalAddress = tronAddress.PubkeyToAddress(suite.externals[i].PublicKey).String()

			checkpoint, err = trontypes.GetCheckpointOracleSet(oracleSet, gravityId)
			require.NoError(suite.T(), err)

			signature, err = trontypes.NewTronSignature(checkpoint, suite.externals[i])
			require.NoError(suite.T(), err)
		}

		suite.Keeper().SetOracleSetConfirm(suite.ctx, suite.oracles[i], &types.MsgOracleSetConfirm{
			Nonce:           oracleSet.Nonce,
			BridgerAddress:  suite.bridgers[i].String(),
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		})
	}

	params := suite.Keeper().GetParams(suite.ctx)
	params.SignedWindow = 10
	suite.Keeper().SetParams(suite.ctx, &params)
	height := suite.Keeper().GetSignedWindow(suite.ctx) + oracleSet.Height + 1
	for i := uint64(2); i <= height; i++ {
		suite.app.BeginBlock(types2.RequestBeginBlock{
			Header: tmproto.Header{Height: int64(i)},
		})
		suite.app.EndBlock(types2.RequestEndBlock{Height: int64(i)})
		suite.app.Commit()
	}

	for _, oracle := range suite.oracles {
		suite.Nil(suite.Keeper().GetOracleSetConfirm(suite.ctx, oracleSet.Nonce, oracle))
	}
}

func (suite *KeeperTestSuite) TestKeeper_IterateOracleSet() {
	var member []types.BridgeValidator
	for i, external := range suite.externals {
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
