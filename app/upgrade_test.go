package app_test

import (
	"os"
	"path/filepath"
	"testing"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	dbm "github.com/cosmos/cosmos-db"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/app"
	nextversion "github.com/functionx/fx-core/v8/app/upgrades/v8"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	fxgovv8 "github.com/functionx/fx-core/v8/x/gov/migrations/v8"
	fxgovtypes "github.com/functionx/fx-core/v8/x/gov/types"
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
	t.Helper()
	fxtypes.SetConfig(true)

	home := filepath.Join(os.Getenv("HOME"), "tmp")
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(home, "data"))
	require.NoError(t, err)

	myApp := helpers.NewApp(func(opts *helpers.AppOpts) {
		opts.Logger = log.NewLogger(os.Stdout)
		opts.DB = db
		opts.Home = home
	})
	return myApp
}

func newContext(t *testing.T, myApp *app.App, chainId string, deliveState bool) sdk.Context {
	t.Helper()
	header := tmproto.Header{
		ChainID: chainId,
		Height:  myApp.LastBlockHeight(),
		Time:    tmtime.Now(),
	}
	var ctx sdk.Context
	if deliveState {
		ctx = myApp.NewContextLegacy(false, header)
	} else {
		ctx = myApp.GetContextForCheckTx(nil).WithBlockHeader(header)
	}
	ctx = ctx.WithChainID(chainId)
	ctx = ctx.WithHeaderInfo(coreheader.Info{
		Height:  header.Height,
		Time:    header.Time,
		ChainID: header.ChainID,
	})
	// set the first validator to proposer
	validators, err := myApp.StakingKeeper.GetAllValidators(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, validators)
	var pubKey cryptotypes.PubKey
	require.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubKey))
	ctx = ctx.WithProposer(pubKey.Address().Bytes())
	return ctx
}

func checkAppUpgrade(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	checkStakingMigrationDelete(t, ctx, myApp)

	checkGovCustomParams(t, ctx, myApp)

	checkErc20Keys(t, ctx, myApp)
}

func checkErc20Keys(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	params, err := myApp.Erc20Keeper.Params.Get(ctx)
	require.NoError(t, err)

	require.True(t, params.EnableErc20)
}

func checkGovCustomParams(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	egfCustomParams, found := myApp.GovKeeper.GetCustomParams(ctx, sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{}))
	require.True(t, found)
	expectEGFParams := fxgovtypes.NewCustomParams(fxgovtypes.EGFCustomParamDepositRatio.String(), fxgovtypes.DefaultEGFCustomParamVotingPeriod, fxgovtypes.DefaultCustomParamQuorum40.String())
	assert.Equal(t, expectEGFParams, egfCustomParams)

	checkKeysIsDelete(t, ctx.KVStore(myApp.GetKey(govtypes.StoreKey)), fxgovv8.GetRemovedStoreKeys())
}

func checkStakingMigrationDelete(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	checkKeysIsDelete(t, ctx.KVStore(myApp.GetKey(stakingtypes.StoreKey)), fxstakingv8.GetRemovedStoreKeys())
}

func checkKeysIsDelete(t *testing.T, kvStore storetypes.KVStore, keys [][]byte) {
	t.Helper()
	require.NotEmpty(t, keys)
	checkFn := func(key []byte) {
		iterator := storetypes.KVStorePrefixIterator(kvStore, key)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			require.Failf(t, "key is not deleted", "prefix:%x, key:%x", key, iterator.Key())
		}
	}
	for _, removeKey := range keys {
		checkFn(removeKey)
	}
}
