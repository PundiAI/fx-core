package app_test

import (
	"os"
	"path/filepath"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v5/app"
	v5 "github.com/functionx/fx-core/v5/app/upgrades/v5"
	"github.com/functionx/fx-core/v5/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v5/types"
)

func Test_TestnetUpgrade(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test: ", t.Name())

	fxtypes.SetConfig(true)
	fxtypes.SetChainId(fxtypes.TestnetChainId) // only for testnet

	testCases := []struct {
		name                  string
		fromVersion           int
		toVersion             int
		LocalStoreBlockHeight uint64
		plan                  upgradetypes.Plan
	}{
		{
			name:        "upgrade v5.0.x",
			fromVersion: 4,
			toVersion:   5,
			plan: upgradetypes.Plan{
				Name: v5.Upgrade.UpgradeName,
				Info: "local test upgrade v5.0.x",
			},
		},
	}

	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(fxtypes.GetDefaultNodeHome(), "data"))
	require.NoError(t, err)

	makeEncodingConfig := app.MakeEncodingConfig()
	myApp := app.New(log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowAll()),
		db, nil, false, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		makeEncodingConfig, app.EmptyAppOptions{})
	// todo default DefaultStoreLoader  New module verification failed
	myApp.SetStoreLoader(upgradetypes.UpgradeStoreLoader(myApp.LastBlockHeight()+1, v5.Upgrade.StoreUpgrades()))
	err = myApp.LoadLatestVersion()
	require.NoError(t, err)

	ctx := newContext(t, myApp)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.plan.Height = ctx.BlockHeight()

			myApp.UpgradeKeeper.ApplyUpgrade(ctx, testCase.plan)

			checkVersionMap(t, ctx, myApp, getConsensusVersion(testCase.toVersion))
		})
	}

	myApp.EthKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.BscKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.TronKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.PolygonKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.AvalancheKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
}

func newContext(t *testing.T, myApp *app.App) sdk.Context {
	chainId := fxtypes.MainnetChainId
	if os.Getenv("CHAIN_ID") == fxtypes.TestnetChainId {
		chainId = fxtypes.TestnetChainId
	}
	ctx := myApp.NewUncachedContext(false, tmproto.Header{
		ChainID: chainId, Height: myApp.LastBlockHeight(),
	})
	// set the first validator to proposer
	validators := myApp.StakingKeeper.GetAllValidators(ctx)
	assert.True(t, len(validators) > 0)
	var pubKey cryptotypes.PubKey
	assert.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubKey))
	ctx = ctx.WithProposer(pubKey.Address().Bytes())
	return ctx
}

func checkVersionMap(t *testing.T, ctx sdk.Context, myApp *app.App, versionMap module.VersionMap) {
	vm := myApp.UpgradeKeeper.GetModuleVersionMap(ctx)
	for k, v := range vm {
		require.Equal(t, versionMap[k], v, k)
	}
}

func getConsensusVersion(appVersion int) (versionMap module.VersionMap) {
	// moduleName: v1,v2,v3
	historyVersions := map[string][]uint64{
		"auth":         {0, 1, 2, 3},
		"authz":        {0, 0, 1, 2},
		"avalanche":    {0, 0, 1, 2, 3},
		"bank":         {0, 1, 2, 3},
		"bsc":          {1, 2, 3, 4},
		"capability":   {1},
		"crisis":       {1},
		"crosschain":   {1},
		"distribution": {1, 2},
		"erc20":        {0, 1, 2, 3},
		"evidence":     {1},
		"evm":          {0, 0, 3, 5},
		"eth":          {0, 0, 1, 2, 3},
		"feegrant":     {0, 0, 1, 2},
		"feemarket":    {0, 0, 3, 4},
		"genutil":      {1},
		"gov":          {0, 1, 2, 3},
		"gravity":      {1, 1, 2},
		"ibc":          {1, 2},
		"migrate":      {0, 1},
		"mint":         {1},
		"other":        {1},
		"params":       {1},
		"polygon":      {1, 2, 3, 4},
		"slashing":     {1, 2},
		"staking":      {0, 1, 2, 3},
		"transfer":     {1, 1, 2}, // ibc-transfer
		"fxtransfer":   {0, 0, 1}, // fx-ibc-transfer
		"tron":         {1, 2, 3, 4},
		"upgrade":      {0, 0, 1, 2},
		"vesting":      {1},
		"arbitrum":     {0, 0, 0, 1},
		"optimism":     {0, 0, 0, 1},
	}
	versionMap = make(map[string]uint64)
	for key, versions := range historyVersions {
		if len(versions) <= appVersion-1 {
			// If not exist, select the last one
			versionMap[key] = versions[len(versions)-1]
		} else {
			versionMap[key] = versions[appVersion-1]
		}
		// If the value is zero, the current version does not exist
		if versionMap[key] == 0 {
			delete(versionMap, key)
		}
	}
	return versionMap
}
