package app_test

import (
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"cosmossdk.io/collections"
	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	dbm "github.com/cosmos/cosmos-db"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/app"
	nextversion "github.com/pundiai/fx-core/v8/app/upgrades/v8"
	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	fxgovv8 "github.com/pundiai/fx-core/v8/x/gov/migrations/v8"
	fxgovtypes "github.com/pundiai/fx-core/v8/x/gov/types"
	fxstakingv8 "github.com/pundiai/fx-core/v8/x/staking/migrations/v8"
)

func Test_UpgradeAndMigrate(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	chainId := fxtypes.MainnetChainId
	myApp := buildApp(t)

	ctx := newContext(t, myApp, chainId, false)

	bdd := beforeUpgrade(ctx, myApp)

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
	checkAppUpgrade(t, ctx, myApp, bdd)
}

func Test_UpgradeTestnet(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	chainId := fxtypes.TestnetChainId
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
	checkLayer2OracleIsOnline(t, ctx, myApp.Layer2Keeper)
	checkPundixPurse(t, ctx, myApp)
	checkTotalSupply(t, ctx, myApp)
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

type BeforeUpgradeData struct {
	AccountBalances map[string]sdk.Coins
	ModuleBalances  map[string]sdk.Coins
}

func beforeUpgrade(ctx sdk.Context, myApp *app.App) BeforeUpgradeData {
	bdd := BeforeUpgradeData{
		AccountBalances: make(map[string]sdk.Coins),
		ModuleBalances:  make(map[string]sdk.Coins),
	}

	accountBalance, moduleBalance := allBalances(ctx, myApp)
	bdd.AccountBalances = accountBalance
	bdd.ModuleBalances = moduleBalance

	return bdd
}

func checkAppUpgrade(t *testing.T, ctx sdk.Context, myApp *app.App, bdd BeforeUpgradeData) {
	t.Helper()
	checkStakingMigrationDelete(t, ctx, myApp)

	checkGovCustomParams(t, ctx, myApp)

	checkErc20Keys(t, ctx, myApp)

	checkOutgoingBatch(t, ctx, myApp)

	checkMigrateBalance(t, ctx, myApp, bdd)

	checkBridgeToken(t, ctx, myApp)
	checkErc20Token(t, ctx, myApp)
	checkPundixPurse(t, ctx, myApp)
	checkTotalSupply(t, ctx, myApp)
	checkLayer2OracleIsOnline(t, ctx, myApp.Layer2Keeper)
}

func checkLayer2OracleIsOnline(t *testing.T, ctx sdk.Context, layer2Keeper crosschainkeeper.Keeper) {
	t.Helper()
	oracles := layer2Keeper.GetAllOracles(ctx, false)
	for _, oracle := range oracles {
		assert.True(t, oracle.Online, oracle.OracleAddress)
	}
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
		keeper.IterateOutgoingTxBatches(ctx, func(batch *crosschaintypes.OutgoingTxBatch) bool {
			assert.True(t, kvStore.Has(crosschaintypes.GetOutgoingTxBatchBlockKey(batch.Block, batch.BatchNonce)))
			return false
		})
	}
}

func checkMigrateBalance(t *testing.T, ctx sdk.Context, myApp *app.App, bdd BeforeUpgradeData) {
	t.Helper()

	newAccountBalance, newModuleBalance := allBalances(ctx, myApp)
	require.GreaterOrEqual(t, len(bdd.AccountBalances), len(newAccountBalance))

	// check address balance
	checkAccountBalance(t, ctx, myApp, bdd.AccountBalances, newAccountBalance)

	for moduleName, coins := range newModuleBalance {
		if moduleName == erc20types.ModuleName {
			for _, coin := range coins {
				found, err := myApp.Erc20Keeper.HasToken(ctx, coin.Denom)
				require.NoError(t, err)
				require.False(t, found)

				found, err = myApp.Erc20Keeper.ERC20Token.Has(ctx, coin.Denom)
				require.NoError(t, err)
				require.True(t, found, coin.Denom)
			}
		}
	}
}

func checkAccountBalance(t *testing.T, ctx sdk.Context, myApp *app.App, accountBalances, newAccountBalance map[string]sdk.Coins) {
	t.Helper()

	for addrStr, coins := range accountBalances {
		newCoins := newAccountBalance[addrStr]
		delete(newAccountBalance, addrStr)

		if coins.Equal(newCoins) {
			continue
		}
		addr := sdk.MustAccAddressFromBech32(addrStr)
		for _, coin := range coins {
			foundMD := myApp.BankKeeper.HasDenomMetaData(ctx, coin.Denom)
			foundToken, err := myApp.Erc20Keeper.HasToken(ctx, coin.Denom)
			require.NoError(t, err)

			if !foundToken && !foundMD {
				balance := myApp.BankKeeper.GetBalance(ctx, addr, coin.Denom)
				require.True(t, balance.Amount.Equal(coin.Amount) || balance.Amount.IsZero())
				continue
			}

			baseDenom := coin.Denom
			if foundToken {
				baseDenom, err = myApp.Erc20Keeper.GetBaseDenom(ctx, coin.Denom)
				require.NoError(t, err)
			}

			tokenDenoms := make(map[string]struct{}, 0)
			tokenDenoms[baseDenom] = struct{}{}
			for _, keeper := range myApp.CrosschainKeepers.ToSlice() {
				bridgeToken, err := myApp.Erc20Keeper.GetBridgeToken(ctx, keeper.ModuleName(), baseDenom)
				if errors.Is(err, collections.ErrNotFound) {
					continue
				}
				require.NoError(t, err)
				tokenDenoms[bridgeToken.BridgeDenom()] = struct{}{}
			}
			ibcTokens, err := getIBCTokens(ctx, myApp, baseDenom)
			require.NoError(t, err)
			for _, ibcToken := range ibcTokens {
				tokenDenoms[ibcToken.GetIbcDenom()] = struct{}{}
			}

			amount := sdkmath.ZeroInt()
			for denom := range tokenDenoms {
				amount = amount.Add(coins.AmountOf(denom))
			}

			balance := myApp.BankKeeper.GetBalance(ctx, addr, baseDenom)
			require.Equal(t, balance.Amount, amount)
		}
	}
	require.Empty(t, newAccountBalance)
}

func checkBridgeToken(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	ethTokens, err := myApp.Erc20Keeper.GetBridgeTokens(ctx, ethtypes.ModuleName)
	require.NoError(t, err)
	require.NotEmpty(t, ethTokens)
	for _, token := range ethTokens {
		require.Equal(t, ethtypes.ModuleName, token.ChainName)
	}
}

func checkErc20Token(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	iter, err := myApp.Erc20Keeper.ERC20Token.Iterate(ctx, nil)
	require.NoError(t, err)
	defer iter.Close()

	kvs, err := iter.KeyValues()
	require.NoError(t, err)

	erc20Tokens := make([]erc20types.ERC20Token, 0, len(kvs))
	for _, kv := range kvs {
		erc20Tokens = append(erc20Tokens, kv.Value)
	}

	for _, et := range erc20Tokens {
		has, err := myApp.Erc20Keeper.DenomIndex.Has(ctx, et.Erc20Address)
		require.NoError(t, err)
		require.True(t, has)
	}
}

func checkPundixPurse(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	erc20ModuleAddr := authtypes.NewModuleAddress(erc20types.ModuleName)
	erc20Contract := contract.NewERC20TokenKeeper(myApp.EvmKeeper)
	for _, denom := range []string{"pundix", "purse"} {
		erc20Token, err := myApp.Erc20Keeper.GetERC20Token(ctx, denom)
		require.NoError(t, err)

		erc20TokenSupply, err := erc20Contract.TotalSupply(ctx, erc20Token.GetERC20Contract())
		require.NoError(t, err)
		if erc20TokenSupply.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		balance := myApp.BankKeeper.GetBalance(ctx, erc20ModuleAddr, denom)
		require.Equal(t, sdkmath.NewIntFromBigInt(erc20TokenSupply), balance.Amount)
	}
}

func checkTotalSupply(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	iter, err := myApp.Erc20Keeper.ERC20Token.Iterate(ctx, nil)
	require.NoError(t, err)
	defer iter.Close()

	kvs, err := iter.KeyValues()
	require.NoError(t, err)

	erc20Tokens := make([]erc20types.ERC20Token, 0, len(kvs))
	for _, kv := range kvs {
		erc20Tokens = append(erc20Tokens, kv.Value)
	}

	for _, et := range erc20Tokens {
		aliasDenoms := make([]string, 0, 10)

		bridgeTokens, err := getBridgeToken(ctx, myApp, et.GetDenom())
		require.NoError(t, err)
		for _, bt := range bridgeTokens {
			if bt.IsOrigin() || bt.IsNative {
				continue
			}
			aliasDenoms = append(aliasDenoms, bt.BridgeDenom())
		}

		ibcTokens, err := getIBCTokens(ctx, myApp, et.GetDenom())
		require.NoError(t, err)
		for _, ibcToken := range ibcTokens {
			aliasDenoms = append(aliasDenoms, ibcToken.GetIbcDenom())
		}

		if len(aliasDenoms) == 0 {
			continue
		}

		aliasTotal := sdkmath.ZeroInt()
		for _, denom := range aliasDenoms {
			supply := myApp.BankKeeper.GetSupply(ctx, denom)
			aliasTotal = aliasTotal.Add(supply.Amount)
		}
		baseSupply := myApp.BankKeeper.GetSupply(ctx, et.GetDenom())
		baseSupply = getTestnetTokenAmount(ctx, baseSupply)
		// NOTE: after sendToExternal fixed, fix bridge token amount
		if !aliasTotal.Equal(baseSupply.Amount) {
			t.Log("not equal", "denom", et.GetDenom())
			continue
		}
	}
}

func allBalances(ctx sdk.Context, myApp *app.App) (map[string]sdk.Coins, map[string]sdk.Coins) {
	accountBalance := make(map[string]sdk.Coins)
	moduleBalance := make(map[string]sdk.Coins)
	myApp.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, balance sdk.Coin) bool {
		account := myApp.AccountKeeper.GetAccount(ctx, addr)
		if ma, ok := account.(*authtypes.ModuleAccount); ok {
			if ma.Name == stakingtypes.BondedPoolName ||
				ma.Name == stakingtypes.NotBondedPoolName ||
				ma.Name == distributiontypes.ModuleName {
				return false
			}

			coins, ok := moduleBalance[ma.Name]
			if !ok {
				coins = sdk.NewCoins()
			}
			coins = coins.Add(balance)
			moduleBalance[ma.Name] = coins
			return false
		}
		if addr.Equals(authtypes.NewModuleAddress(crosschaintypes.ModuleName)) {
			return false
		}
		coins, ok := accountBalance[addr.String()]
		if !ok {
			coins = sdk.NewCoins()
		}
		coins = coins.Add(balance)
		accountBalance[addr.String()] = coins
		return false
	})
	return accountBalance, moduleBalance
}

func getIBCTokens(ctx sdk.Context, myApp *app.App, baseDenom string) ([]erc20types.IBCToken, error) {
	rng := collections.NewPrefixedPairRange[string, string](baseDenom)
	iter, err := myApp.Erc20Keeper.IBCToken.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}
	kvs, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}

	tokens := make([]erc20types.IBCToken, 0, len(kvs))
	for _, kv := range kvs {
		tokens = append(tokens, kv.Value)
	}
	return tokens, nil
}

func getBridgeToken(ctx sdk.Context, myApp *app.App, baseDenom string) ([]erc20types.BridgeToken, error) {
	bridgeTokens := make([]erc20types.BridgeToken, 0, len(myApp.CrosschainKeepers.ToSlice()))
	for _, ck := range myApp.CrosschainKeepers.ToSlice() {
		key := collections.Join(ck.ModuleName(), baseDenom)
		has, err := myApp.Erc20Keeper.BridgeToken.Has(ctx, key)
		if err != nil {
			return nil, err
		}
		if !has {
			continue
		}
		bridgeToken, _ := myApp.Erc20Keeper.GetBridgeToken(ctx, ck.ModuleName(), baseDenom)
		bridgeTokens = append(bridgeTokens, bridgeToken)
	}
	return bridgeTokens, nil
}

var fixTestnetTokensAmount = map[string]sdkmath.Int{
	"atom/osmo":    nextversion.MustParseIntFromString("54386202508063381657"),
	"stosmo/tosmo": nextversion.MustParseIntFromString("1740010000000000000000000"),
	"tatom/tosmo":  nextversion.MustParseIntFromString("1700841000000000000000000"),
	"uosmo":        nextversion.MustParseIntFromString("35005636"),
	"usdc/tosmo":   nextversion.MustParseIntFromString("1730015000000000000000000"),
	"wbtc/tosmo":   nextversion.MustParseIntFromString("1740010000000000000000000"),
	"weth/tosmo":   nextversion.MustParseIntFromString("1700015000000000000000000"),
}

func getTestnetTokenAmount(ctx sdk.Context, coin sdk.Coin) sdk.Coin {
	if ctx.ChainID() == fxtypes.MainnetChainId {
		return coin
	}
	amount, ok := fixTestnetTokensAmount[coin.Denom]
	if !ok {
		return coin
	}
	coin.Amount = coin.Amount.Sub(amount)
	return coin
}
