package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app/fxcore"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func defaultModuleParams(oracles []string) types.Params {
	return types.Params{
		GravityId:                         "bsc",
		SignedWindow:                      20000,
		ExternalBatchTimeout:              43200000,
		AverageBlockTime:                  5000,
		AverageExternalBlockTime:          3000,
		SlashFraction:                     sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		IbcTransferTimeoutHeight:          10000,
		OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
		Oracles:                           oracles,
		DepositThreshold:                  sdk.NewCoin("FX", sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil))),
	}
}

func TestSetOracle(t *testing.T) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	oracleAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.ZeroInt())
	orchestratorAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.ZeroInt())

	var oracles []string
	for _, account := range oracleAddressList {
		oracles = append(oracles, account.String())
	}
	// set chain init params
	app.BscKeeper.SetParams(ctx, defaultModuleParams(oracles))

	oracleAddr1, oracleAddr2 := oracleAddressList[0], oracleAddressList[1]
	orchestratorAddr1, orchestratorAddr2 := orchestratorAddressList[0], orchestratorAddressList[1]
	var dbOracleAddr1 sdk.AccAddress
	var found bool
	dbOracleAddr1, found = app.BscKeeper.GetOracleAddressByOrchestratorKey(ctx, orchestratorAddr1)
	require.False(t, found)
	require.Empty(t, dbOracleAddr1)

	// 1. set oracle -> orchestrator  and  orchestrator -> oracle
	app.BscKeeper.SetOracleByOrchestrator(ctx, oracleAddr1, orchestratorAddr1)

	// 2. find oracle by orchestrator
	dbOracleAddr1, found = app.BscKeeper.GetOracleAddressByOrchestratorKey(ctx, orchestratorAddr1)
	require.True(t, found)
	require.EqualValues(t, oracleAddr1, dbOracleAddr1)

	// 2.1 find orchestrator by oracle
	//dbOrchestratorAddr1, found := app.BscKeeper.GetOrchestratorAddressByOracle(ctx, oracleAddr1)
	//require.True(t, found)
	//require.EqualValues(t, orchestratorAddr1, dbOrchestratorAddr1)

	// 3. find not exist orchestrator by oracle2
	//dbOrchestratorAddr2, found := app.BscKeeper.GetOrchestratorAddressByOracle(ctx, oracleAddr2)
	//require.False(t, found)
	//require.Nil(t, dbOrchestratorAddr2)

	// 3.1 set oracle2 -> orchestrator2
	app.BscKeeper.SetOracleByOrchestrator(ctx, oracleAddr2, orchestratorAddr2)

	// 3.2 find oracle2 by orchestrator2
	dbOrchestratorAddr2, found := app.BscKeeper.GetOracleAddressByOrchestratorKey(ctx, oracleAddr2)
	require.True(t, found)
	require.EqualValues(t, orchestratorAddr2, dbOrchestratorAddr2)
}

func TestLastPendingOracleSetRequestByAddr(t *testing.T) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	oracleAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.ZeroInt())
	orchestratorAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.ZeroInt())

	keeper := app.BscKeeper

	testCases := []struct {
		OracleAddress       sdk.AccAddress
		OrchestratorAddress sdk.AccAddress
		StartHeight         int64

		ExpectOracleSetSize int
	}{
		{
			OracleAddress:       oracleAddressList[0],
			OrchestratorAddress: orchestratorAddressList[0],
			StartHeight:         1,
			ExpectOracleSetSize: 3,
		},
		{
			OracleAddress:       oracleAddressList[1],
			OrchestratorAddress: orchestratorAddressList[1],
			StartHeight:         2,
			ExpectOracleSetSize: 2,
		},
		{
			OracleAddress:       oracleAddressList[2],
			OrchestratorAddress: orchestratorAddressList[2],
			StartHeight:         3,
			ExpectOracleSetSize: 1,
		},
	}

	for i := 1; i <= 3; i++ {
		keeper.StoreOracleSet(ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           uint64(i),
				ExternalAddress: fmt.Sprintf("0x%d", i),
			}},
			Height: uint64(i),
		})
	}

	wrapSDKContext := sdk.WrapSDKContext(ctx)
	for _, testCase := range testCases {
		oracle := types.Oracle{
			OracleAddress:       testCase.OracleAddress.String(),
			OrchestratorAddress: testCase.OrchestratorAddress.String(),
			StartHeight:         testCase.StartHeight,
		}
		// save oracle
		keeper.SetOracle(ctx, oracle)
		keeper.SetOracleByOrchestrator(ctx, oracle.GetOracle(), testCase.OrchestratorAddress)

		pendingOracleSetRequestByAddr, err := keeper.LastPendingOracleSetRequestByAddr(wrapSDKContext, &types.QueryLastPendingOracleSetRequestByAddrRequest{
			OrchestratorAddress: testCase.OrchestratorAddress.String(),
		})
		require.NoError(t, err)
		require.EqualValues(t, testCase.ExpectOracleSetSize, len(pendingOracleSetRequestByAddr.OracleSets))
	}
}

func TestLastPendingBatchRequestByAddr(t *testing.T) {

	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	oracleAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.ZeroInt())
	orchestratorAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.ZeroInt())

	keeper := app.BscKeeper

	testCases := []struct {
		Name                string
		OracleAddress       sdk.AccAddress
		OrchestratorAddress sdk.AccAddress
		StartHeight         int64
		ExpectStartHeight   uint64
	}{
		{
			Name:                "oracle start height with 1, expect oracle set block 3",
			OracleAddress:       oracleAddressList[0],
			OrchestratorAddress: orchestratorAddressList[0],
			StartHeight:         1,
			ExpectStartHeight:   3,
		},
		{
			Name:                "oracle start height with 2, expect oracle set block 2",
			OracleAddress:       oracleAddressList[1],
			OrchestratorAddress: orchestratorAddressList[1],
			StartHeight:         2,
			ExpectStartHeight:   3,
		},
		{
			Name:                "oracle start height with 3, expect oracle set block 1",
			OracleAddress:       oracleAddressList[2],
			OrchestratorAddress: orchestratorAddressList[2],
			StartHeight:         3,
			ExpectStartHeight:   3,
		},
	}
	for i := uint64(1); i <= 3; i++ {
		ctx = ctx.WithBlockHeight(int64(i))
		err := keeper.StoreBatch(ctx, &types.OutgoingTxBatch{
			Block:      i,
			BatchNonce: i,
			Transactions: types.OutgoingTransferTxs{{
				Id:          i,
				Sender:      fmt.Sprintf("0x%d", i),
				DestAddress: fmt.Sprintf("0x%d", i),
			}},
		})
		require.NoError(t, err)
	}

	wrapSDKContext := sdk.WrapSDKContext(ctx)
	for _, testCase := range testCases {
		oracle := types.Oracle{
			OracleAddress:       testCase.OracleAddress.String(),
			OrchestratorAddress: testCase.OrchestratorAddress.String(),
			StartHeight:         testCase.StartHeight,
		}
		// save oracle
		keeper.SetOracle(ctx, oracle)
		keeper.SetOracleByOrchestrator(ctx, oracle.GetOracle(), testCase.OrchestratorAddress)

		pendingLastPendingBatchRequestByAddr, err := keeper.LastPendingBatchRequestByAddr(wrapSDKContext, &types.QueryLastPendingBatchRequestByAddrRequest{
			OrchestratorAddress: testCase.OrchestratorAddress.String(),
		})
		require.NoError(t, err, testCase.Name)
		require.NotNil(t, pendingLastPendingBatchRequestByAddr, testCase.Name)
		require.NotNil(t, pendingLastPendingBatchRequestByAddr.Batch, testCase.Name)
		require.EqualValues(t, testCase.ExpectStartHeight, pendingLastPendingBatchRequestByAddr.Batch.Block, testCase.Name)
	}
}

func TestGetUnSlashedOracleSets(t *testing.T) {

	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	//oracleAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.ZeroInt())
	//orchestratorAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.ZeroInt())

	keeper := app.BscKeeper

	for i := 1; i <= 3; i++ {
		keeper.StoreOracleSet(ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           uint64(i),
				ExternalAddress: fmt.Sprintf("0x%d", i),
			}},
			Height: uint64(1000 + i),
		})
	}
	slashOracleSetHeight := 1003
	sets := keeper.GetUnSlashedOracleSets(ctx, uint64(slashOracleSetHeight))
	require.NotNil(t, sets)
	require.EqualValues(t, 2, sets.Len())

	keeper.SetLastSlashedOracleSetNonce(ctx, 1)
	slashOracleSetHeight = 1003
	sets = keeper.GetUnSlashedOracleSets(ctx, uint64(slashOracleSetHeight))
	require.NotNil(t, sets)
	require.EqualValues(t, 1, sets.Len())

	slashOracleSetHeight = 1004
	sets = keeper.GetUnSlashedOracleSets(ctx, uint64(slashOracleSetHeight))
	require.NotNil(t, sets)
	require.EqualValues(t, 2, sets.Len())

}
