package keeper_test

import (
	fxtypes "github.com/functionx/fx-core/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/x/gravity/types"
)

func TestSetOrchestratorValidator(t *testing.T) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := app.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances)))
	fxcore := app.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})
	accAddressList := app.AddTestAddrsIncremental(fxcore, ctx, 2, sdk.ZeroInt())

	val1OrchestratorAddr, val2OrchestratorAddr := accAddressList[0], accAddressList[1]
	val1Addr, found := fxcore.GravityKeeper.GetOrchestratorValidator(ctx, val1OrchestratorAddr)
	require.False(t, found)
	require.Empty(t, val1Addr)

	fxcore.GravityKeeper.SetOrchestratorValidator(ctx, sdk.ValAddress(val1OrchestratorAddr), val1OrchestratorAddr)

	val1Addr, found = fxcore.GravityKeeper.GetOrchestratorValidator(ctx, val1OrchestratorAddr)
	require.True(t, found)
	t.Log(val1Addr)
	require.EqualValues(t, sdk.ValAddress(val1OrchestratorAddr), val1Addr)

	val2Addr, found := fxcore.GravityKeeper.GetOrchestratorValidator(ctx, val2OrchestratorAddr)
	require.False(t, found)
	require.Empty(t, val2Addr)
}

func TestStoreValset(t *testing.T) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := app.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances)))
	fxcore := app.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})

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
	fxcore.GravityKeeper.StoreValset(ctx, valset)
	// checkout nonce has exists
	require.True(t, fxcore.GravityKeeper.HasValsetRequest(ctx, valset.Nonce))
	// check not exists nonce
	require.False(t, fxcore.GravityKeeper.HasValsetRequest(ctx, valset.Nonce-1))
	// check latest valset nonce
	require.EqualValues(t, valset.Nonce, fxcore.GravityKeeper.GetLatestValsetNonce(ctx))
	// check store valset
	storeValset := fxcore.GravityKeeper.GetValset(ctx, valset.Nonce)
	require.NotNil(t, storeValset)
	require.EqualValues(t, valset, storeValset)

	storeLatestValset := fxcore.GravityKeeper.GetLatestValset(ctx)
	require.NotNil(t, storeLatestValset)
	require.EqualValues(t, valset, storeLatestValset)
}

func TestSetLatestValsetNonce(t *testing.T) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := app.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances)))
	fxcore := app.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})

	require.NotPanics(t, func() {
		fxcore.GravityKeeper.SetLatestValsetNonce(ctx, 1)
	})
	require.EqualValues(t, 1, fxcore.GravityKeeper.GetLatestValsetNonce(ctx))

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
			fxcore.GravityKeeper.SetLatestValsetNonce(ctx, testData.newValsetNonce)
			require.EqualValues(t, testData.expect, fxcore.GravityKeeper.GetLatestValsetNonce(ctx))
		})
	}
}
func TestLastSlashedValsetNonce(t *testing.T) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := app.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances)))
	fxcore := app.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})

	keeper := fxcore.GravityKeeper
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
		t.Log(testCase.name)
		keeper.SetLastSlashedValsetNonce(ctx, testCase.setValue)
		require.EqualValues(t, testCase.expectValue, keeper.GetLastSlashedValsetNonce(ctx))
	}
}

func TestLastUnBondingBlockHeight(t *testing.T) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := app.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances)))
	fxcore := app.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})

	keeper := fxcore.GravityKeeper
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
		t.Log(testCase.name)
		keeper.SetLastUnBondingBlockHeight(ctx, testCase.setValue)
		require.EqualValues(t, testCase.expectValue, keeper.GetLastUnBondingBlockHeight(ctx))
	}
}
