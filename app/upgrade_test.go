package app_test

import (
	"os"
	"path/filepath"
	"testing"

	coreheader "cosmossdk.io/core/header"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	dbm "github.com/cosmos/cosmos-db"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/app"
	nextversion "github.com/functionx/fx-core/v8/app/upgrades/v8"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	fxstakingv8 "github.com/functionx/fx-core/v8/x/staking/migrations/v8"
)

func Test_UpgradeAndMigrate(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	chainId := fxtypes.MainnetChainId
	myApp := buildApp(t)

	ctx := newContext(t, myApp, chainId, false)

	// 1. set upgrade plan
	require.NoError(t, myApp.UpgradeKeeper.ScheduleUpgrade(ctx, upgradetypes.Plan{
		Name:   nextversion.Upgrade.UpgradeName,
		Height: ctx.BlockHeight(),
	}))

	// 2. execute upgrade
	responsePreBlock, err := upgrade.PreBlocker(ctx, myApp.UpgradeKeeper)
	require.NoError(t, err)
	require.True(t, responsePreBlock.IsConsensusParamsChanged())

	// 3. check the status after the upgrade
	checkAppUpgrade(t, ctx, myApp)
}

func buildApp(t *testing.T) *app.App {
	fxtypes.SetConfig(true)

	home := filepath.Join(os.Getenv("HOME"), "tmp")
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(home, "data"))
	require.NoError(t, err)

	myApp := helpers.NewApp(func(opts *helpers.AppOpts) {
		opts.DB = db
		opts.Home = home
	})
	return myApp
}

func newContext(t *testing.T, myApp *app.App, chainId string, deliveState bool) sdk.Context {
	header := tmproto.Header{
		ChainID: chainId,
		Height:  myApp.LastBlockHeight(),
		Time:    tmtime.Now(),
	}
	var ctx sdk.Context
	if deliveState {
		ctx = myApp.NewContextLegacy(false, header)
	} else {
		ctx = myApp.NewUncachedContext(false, header)
	}
	ctx = ctx.WithHeaderInfo(coreheader.Info{
		Height:  header.Height,
		Time:    header.Time,
		ChainID: header.ChainID,
	})
	// set the first validator to proposer
	validators, err := myApp.StakingKeeper.GetAllValidators(ctx)
	require.NoError(t, err)
	assert.True(t, len(validators) > 0)
	var pubKey cryptotypes.PubKey
	assert.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubKey))
	ctx = ctx.WithProposer(pubKey.Address().Bytes())
	return ctx
}

func checkAppUpgrade(t *testing.T, ctx sdk.Context, myApp *app.App) {
	checkStakingMigrationDelete(t, ctx, myApp)
}

func checkStakingMigrationDelete(t *testing.T, ctx sdk.Context, myApp *app.App) {
	storeKey := myApp.GetKey(stakingtypes.StoreKey)
	kvStore := ctx.KVStore(storeKey)
	removeKeys := fxstakingv8.GetRemovedValidatorStoreKeys()
	require.Greater(t, len(removeKeys), 0)
	for _, removeKey := range removeKeys {
		iterator := storetypes.KVStorePrefixIterator(kvStore, removeKey)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			require.Failf(t, "key is not deleted", "prefix:%x, key:%x", removeKey, iterator.Key())
		}
	}
}
