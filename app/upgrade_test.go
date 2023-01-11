package app_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	v2 "github.com/functionx/fx-core/v3/app/upgrades/v2"
	v3 "github.com/functionx/fx-core/v3/app/upgrades/v3"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func Test_Upgrade(t *testing.T) {
	if !helpers.IsLocalTest() {
		t.Skip("skipping local test", t.Name())
	}
	fxtypes.SetConfig(true)

	testCases := []struct {
		name                  string
		fromVersion           int
		toVersion             int
		LocalStoreBlockHeight uint64
		plan                  upgradetypes.Plan
	}{
		{
			name:        "upgrade v3",
			fromVersion: 2,
			toVersion:   3,
			plan: upgradetypes.Plan{
				Name: v3.Upgrade.UpgradeName,
				Info: "local test upgrade v3",
			},
			LocalStoreBlockHeight: 7654832,
		},
	}

	db, err := sdk.NewLevelDB("application", filepath.Join(fxtypes.GetDefaultNodeHome(), "data"))
	require.NoError(t, err)

	appEncodingCfg := app.MakeEncodingConfig()
	// logger := log.NewNopLogger()
	logger := log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowInfo())
	myApp := app.New(logger, db,
		nil, true, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		appEncodingCfg, app.EmptyAppOptions{},
	)
	ctx := myApp.NewUncachedContext(false, tmproto.Header{
		ChainID:         fxtypes.ChainId(),
		Height:          myApp.LastBlockHeight(),
		ProposerAddress: tmrand.Bytes(20),
	})
	validators := myApp.StakingKeeper.GetAllValidators(ctx)
	assert.True(t, len(validators) > 0)
	var pubkey cryptotypes.PubKey
	assert.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubkey))
	ctx = ctx.WithProposer(pubkey.Address().Bytes())

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.LocalStoreBlockHeight > 0 {
				require.Equal(t, ctx.BlockHeight(), int64(testCase.LocalStoreBlockHeight))
				for moduleName, keys := range v2.GetModuleKey() {
					kvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(moduleName))
					checkStoreKey(t, moduleName, keys, kvStore)
				}
			}
			checkVersionMap(t, ctx, myApp, getConsensusVersion(testCase.fromVersion))

			testCase.plan.Height = ctx.BlockHeight()
			myApp.UpgradeKeeper.ApplyUpgrade(ctx, testCase.plan)

			// reload store
			// err = upgradetypes.UpgradeStoreLoader(plan.Height, v3.Upgrade.StoreUpgrades())(myApp.CommitMultiStore())
			// require.NoError(t, err)

			checkVersionMap(t, ctx, myApp, getConsensusVersion(testCase.toVersion))

			if testCase.LocalStoreBlockHeight > 0 {
				for moduleName, keys := range v3.GetModuleKey() {
					kvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(moduleName))
					checkStoreKey(t, moduleName, keys, kvStore)
				}
			}
		})
	}
	// myApp.CommitMultiStore().Commit()
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
		"auth":         {1, 2},
		"authz":        {0, 1},
		"avalanche":    {0, 0, 1},
		"bank":         {1, 2},
		"bsc":          {1, 2, 3},
		"capability":   {1},
		"crisis":       {1},
		"crosschain":   {1},
		"distribution": {1, 2},
		"erc20":        {0, 1, 2},
		"evidence":     {1},
		"evm":          {0, 2, 3},
		"eth":          {0, 0, 1},
		"feegrant":     {0, 1},
		"feemarket":    {0, 3},
		"genutil":      {1},
		"gov":          {1, 2},
		"gravity":      {1, 1, 2},
		"ibc":          {1, 2},
		"migrate":      {0, 1},
		"mint":         {1},
		"other":        {1},
		"params":       {1},
		"polygon":      {1, 2, 3},
		"slashing":     {1, 2},
		"staking":      {1, 2},
		"transfer":     {1, 1, 2}, // ibc-transfer
		"fxtransfer":   {0, 0, 1}, // fx-ibc-transfer
		"tron":         {1, 2, 3},
		"upgrade":      {1},
		"vesting":      {1},
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

func checkStoreKey(t *testing.T, name string, keys map[byte][2]int, kvStores storetypes.KVStore) {
	iterator := kvStores.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		x, ok := keys[iterator.Key()[0]]
		assert.True(t, ok, fmt.Sprintf("%x", iterator.Key()[0]), iterator.Value())
		if ok {
			if x[0] == -1 && x[1] == -1 {
				// ignore
				continue
			}
			// set result
			keys[iterator.Key()[0]] = [2]int{x[0], x[1] + 1}
		}
	}
	iterator.Close()
	for k, x := range keys {
		assert.Equal(t, x[0], x[1], fmt.Sprintf("%s: %x", name, k))
	}
}
