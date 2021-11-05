package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app/fxcore"
	"github.com/functionx/fx-core/x/gravity/types"
)

func TestSetOrchestratorValidator(t *testing.T) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	accAddressList := fxcore.AddTestAddrsIncremental(app, ctx, 2, sdk.ZeroInt())

	val1OrchestratorAddr, val2OrchestratorAddr := accAddressList[0], accAddressList[1]
	val1Addr, found := app.GravityKeeper.GetOrchestratorValidator(ctx, val1OrchestratorAddr)
	require.False(t, found)
	require.Empty(t, val1Addr)

	app.GravityKeeper.SetOrchestratorValidator(ctx, sdk.ValAddress(val1OrchestratorAddr), val1OrchestratorAddr)

	val1Addr, found = app.GravityKeeper.GetOrchestratorValidator(ctx, val1OrchestratorAddr)
	require.True(t, found)
	t.Log(val1Addr)
	require.EqualValues(t, sdk.ValAddress(val1OrchestratorAddr), val1Addr)

	val2Addr, found := app.GravityKeeper.GetOrchestratorValidator(ctx, val2OrchestratorAddr)
	require.False(t, found)
	require.Empty(t, val2Addr)
}

func TestStoreValset(t *testing.T) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valset := &types.Valset{
		Nonce: 234,
		Members: []*types.BridgeValidator{
			{Power: 1, EthAddress: "0x1"},
			{Power: 1, EthAddress: "0x2"},
			{Power: 1, EthAddress: "0x3"},
		},
		Height: 234,
	}
	// store new valset
	app.GravityKeeper.StoreValset(ctx, valset)
	// checkout nonce has exists
	require.True(t, app.GravityKeeper.HasValsetRequest(ctx, valset.Nonce))
	// check not exists nonce
	require.False(t, app.GravityKeeper.HasValsetRequest(ctx, valset.Nonce-1))
	// check latest valset nonce
	require.EqualValues(t, valset.Nonce, app.GravityKeeper.GetLatestValsetNonce(ctx))
	// check store valset
	storeValset := app.GravityKeeper.GetValset(ctx, valset.Nonce)
	require.NotNil(t, storeValset)
	require.EqualValues(t, valset, storeValset)

	storeLatestValset := app.GravityKeeper.GetLatestValset(ctx)
	require.NotNil(t, storeLatestValset)
	require.EqualValues(t, valset, storeLatestValset)
}

func TestSetLatestValsetNonce(t *testing.T) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	require.NotPanics(t, func() {
		app.GravityKeeper.SetLatestValsetNonce(ctx, 1)
	})
	require.EqualValues(t, 1, app.GravityKeeper.GetLatestValsetNonce(ctx))

	type testCaseData struct {
		name           string
		newValsetNonce uint64
		expect         uint64
	}
	testDatas := []testCaseData{
		{
			name:           "nonce 1",
			newValsetNonce: 1,
			expect:         1,
		},
		{
			name:           "nonce 4567",
			newValsetNonce: 4567,
			expect:         4567,
		},
		{
			name:           "nonce 555",
			newValsetNonce: 555,
			expect:         555,
		},
	}

	for _, testData := range testDatas {
		t.Run(testData.name, func(t *testing.T) {
			app.GravityKeeper.SetLatestValsetNonce(ctx, testData.newValsetNonce)
			require.EqualValues(t, testData.expect, app.GravityKeeper.GetLatestValsetNonce(ctx))
		})
	}
}
func TestLastSlashedValsetNonce(t *testing.T) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	keeper := app.GravityKeeper
	// check no data
	require.EqualValues(t, 0, keeper.GetLastSlashedValsetNonce(ctx))
	type testCaseData struct {
		name        string
		setValue    uint64
		expectValue uint64
	}

	testCaseDatas := []testCaseData{
		{"set latest shasled valset nonce to 1000", 1000, 1000},
		{"set latest shasled valset nonce to 2000", 2000, 2000},
		{"set latest shasled valset nonce to 1500", 1500, 1500},
	}

	for _, testCase := range testCaseDatas {
		keeper.SetLastSlashedValsetNonce(ctx, testCase.setValue)
		require.EqualValues(t, testCase.expectValue, keeper.GetLastSlashedValsetNonce(ctx))
	}
}

func TestLastUnBondingBlockHeight(t *testing.T) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	keeper := app.GravityKeeper
	// check no data
	require.EqualValues(t, 0, keeper.GetLastUnBondingBlockHeight(ctx))
	type testCaseData struct {
		name        string
		setValue    uint64
		expectValue uint64
	}

	testCaseDatas := []testCaseData{
		{"set latest unBonding block height to 1000", 1000, 1000},
		{"set latest unBonding block height to 2000", 2000, 2000},
		{"set latest unBonding block height to 1500", 1500, 1500},
	}

	for _, testCase := range testCaseDatas {
		keeper.SetLastUnBondingBlockHeight(ctx, testCase.setValue)
		require.EqualValues(t, testCase.expectValue, keeper.GetLastUnBondingBlockHeight(ctx))
	}
}
