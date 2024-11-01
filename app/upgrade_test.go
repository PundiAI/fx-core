package app_test

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	dbm "github.com/cosmos/cosmos-db"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/app"
	nextversion "github.com/functionx/fx-core/v8/app/upgrades/v8"
	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	bsctypes "github.com/functionx/fx-core/v8/x/bsc/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
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

	checkOutgoingBatch(t, ctx, myApp)

	checkPundixChainMigrate(t, ctx, myApp)
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

func checkOutgoingBatch(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	for _, keeper := range myApp.CrosschainKeepers.ToSlice() {
		kvStore := ctx.KVStore(myApp.GetKey(keeper.ModuleName()))
		keeper.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
			assert.True(t, kvStore.Has(types.GetOutgoingTxBatchBlockKey(batch.Block, batch.BatchNonce)))
			return false
		})
	}
}

func checkPundixChainMigrate(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	pundixGenesisPath := path.Join(fxtypes.GetDefaultNodeHome(), "config/pundix_genesis.json")
	appState, err := nextversion.ReadGenesisState(pundixGenesisPath)
	require.NoError(t, err)

	checkPundixBank(t, ctx, myApp, appState[banktypes.ModuleName])
}

func checkPundixBank(t *testing.T, ctx sdk.Context, myApp *app.App, raw json.RawMessage) {
	t.Helper()

	var bankGenesis banktypes.GenesisState
	require.NoError(t, tmjson.Unmarshal(raw, &bankGenesis))
	erc20TokenKeeper := contract.NewERC20TokenKeeper(myApp.EvmKeeper)

	// pundix token
	pundixToken, err := myApp.Erc20Keeper.GetERC20Token(ctx, "pundix")
	require.NoError(t, err)
	pundixBridgeToken, err := myApp.Erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, pundixToken.GetDenom())
	require.NoError(t, err)
	pundixDenomHash := sha256.Sum256([]byte(fmt.Sprintf("%s/channel-0/%s", ibctransfertypes.ModuleName, pundixBridgeToken.BridgeDenom())))
	pundixIBCDenom := fmt.Sprintf("%s/%X", ibctransfertypes.DenomPrefix, pundixDenomHash[:])

	// purse token
	purseToken, err := myApp.Erc20Keeper.GetERC20Token(ctx, "purse")
	require.NoError(t, err)
	purseBscBridgeToken, err := myApp.Erc20Keeper.GetBridgeToken(ctx, bsctypes.ModuleName, purseToken.GetDenom())
	require.NoError(t, err)

	totalSupply, err := erc20TokenKeeper.TotalSupply(ctx, purseToken.GetERC20Contract())
	require.NoError(t, err)
	require.Equal(t, totalSupply.String(), bankGenesis.Supply.AmountOf(purseBscBridgeToken.BridgeDenom()).String())

	erc20Addr := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName))
	balOf, err := erc20TokenKeeper.BalanceOf(ctx, purseToken.GetERC20Contract(), erc20Addr)
	require.NoError(t, err)
	supply := myApp.BankKeeper.GetSupply(ctx, purseToken.GetDenom())
	require.Equal(t, balOf.String(), supply.Amount.String())

	pxEscrowAddr, err := nextversion.GetPxChannelEscrowAddr()
	require.NoError(t, err)

	for _, bal := range bankGenesis.Balances {
		if bal.Address == pxEscrowAddr {
			continue
		}
		bech32Addr, err := sdk.GetFromBech32(bal.Address, "px")
		require.NoError(t, err)
		account := myApp.AccountKeeper.GetAccount(ctx, bech32Addr)
		if _, ok := account.(sdk.ModuleAccountI); ok {
			continue
		}

		allBal := myApp.BankKeeper.GetAllBalances(ctx, bech32Addr)
		require.True(t, allBal.AmountOf(pundixToken.GetDenom()).GTE(bal.Coins.AmountOf(pundixIBCDenom)))
		require.True(t, allBal.AmountOf(purseToken.GetDenom()).GTE(bal.Coins.AmountOf(purseBscBridgeToken.BridgeDenom())))
	}
}
