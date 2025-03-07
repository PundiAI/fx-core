package app_test

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	dbm "github.com/cosmos/cosmos-db"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/app"
	nextversion "github.com/pundiai/fx-core/v8/app/upgrades/v8"
	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
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
	responsePreBlock, err := myApp.PreBlocker(ctx, nil)
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
		opts.Logger = log.NewLogger(os.Stdout, log.LevelOption(zerolog.InfoLevel))
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

	checkWrapToken(t, ctx, myApp)
}

func checkWrapToken(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	erc20TokenKeeper := contract.NewERC20TokenKeeper(myApp.EvmKeeper)
	wrapTokenAddress := nextversion.GetWFXAddress(ctx.ChainID())
	name, err := erc20TokenKeeper.Name(ctx, wrapTokenAddress)
	require.NoError(t, err)
	assert.EqualValues(t, nextversion.WrapName, name)
	symbol, err := erc20TokenKeeper.Symbol(ctx, wrapTokenAddress)
	require.NoError(t, err)
	assert.EqualValues(t, "WPUNDIAI", symbol)
	decimals, err := erc20TokenKeeper.Decimals(ctx, wrapTokenAddress)
	require.NoError(t, err)
	assert.EqualValues(t, 18, decimals)
	owner, err := erc20TokenKeeper.Owner(ctx, wrapTokenAddress)
	require.NoError(t, err)
	assert.EqualValues(t, common.BytesToAddress(myApp.AccountKeeper.GetModuleAddress(erc20types.ModuleName).Bytes()), owner)
	supply, err := erc20TokenKeeper.TotalSupply(ctx, wrapTokenAddress)
	require.NoError(t, err)
	assert.EqualValues(t, 1, supply.Cmp(big.NewInt(0)))
}
