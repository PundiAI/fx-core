package app_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/app"
	nextversion "github.com/functionx/fx-core/v7/app/upgrades/v7"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func Test_UpgradeAndMigrate(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	fxtypes.SetConfig(true)

	home := filepath.Join(os.Getenv("HOME"), "tmp")
	chainId := fxtypes.TestnetChainId // The upgrade test is not related to chainId, do not modify it
	fxtypes.SetChainId(chainId)

	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(home, "data"))
	require.NoError(t, err)

	makeEncodingConfig := app.MakeEncodingConfig()
	myApp := app.New(log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowAll()),
		db, nil, false, map[int64]bool{}, home, 0,
		makeEncodingConfig, app.EmptyAppOptions{}, baseapp.SetChainID(chainId))
	myApp.SetStoreLoader(upgradetypes.UpgradeStoreLoader(myApp.LastBlockHeight()+1, nextversion.Upgrade.StoreUpgrades()))
	require.NoError(t, myApp.LoadLatestVersion())

	ctx := newContext(t, myApp, chainId)

	require.NoError(t, myApp.UpgradeKeeper.ScheduleUpgrade(ctx, upgradetypes.Plan{
		Name:   nextversion.Upgrade.UpgradeName,
		Height: ctx.BlockHeight() + 1,
	}))

	header := ctx.BlockHeader()
	header.Height = header.Height + 1
	header.Time = time.Now().UTC()
	require.NotPanics(t, func() {
		myApp.BeginBlock(abci.RequestBeginBlock{
			Header: header,
		})
	})
	require.NotPanics(t, func() {
		myApp.EndBlock(abci.RequestEndBlock{
			Height: header.Height,
		})
	})
}

func newContext(t *testing.T, myApp *app.App, chainId string) sdk.Context {
	ctx := myApp.NewUncachedContext(false, tmproto.Header{
		ChainID: chainId,
		Height:  myApp.LastBlockHeight(),
	})
	// set the first validator to proposer
	validators := myApp.StakingKeeper.GetAllValidators(ctx)
	assert.True(t, len(validators) > 0)
	var pubKey cryptotypes.PubKey
	assert.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubKey))
	ctx = ctx.WithProposer(pubKey.Address().Bytes())
	return ctx
}
