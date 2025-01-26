package app_test

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"cosmossdk.io/collections"
	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	dbm "github.com/cosmos/cosmos-db"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/rs/zerolog"
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

	bdd := beforeUpgrade(t, ctx, myApp)

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
	responsePreBlock, err := myApp.PreBlocker(ctx, nil)
	require.NoError(t, err)
	require.True(t, responsePreBlock.IsConsensusParamsChanged())

	// 3. check the status after the upgrade
	checkFXBridgeDenom(t, ctx, myApp)
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

type BeforeUpgradeData struct {
	AccountBalances map[string]sdk.Coins
	ModuleBalances  map[string]sdk.Coins
	Delegates       map[string]*stakingtypes.DelegationResponse
	DelegateRewards map[string]sdk.DecCoins
}

func beforeUpgrade(t *testing.T, ctx sdk.Context, myApp *app.App) BeforeUpgradeData {
	t.Helper()
	bdd := BeforeUpgradeData{
		AccountBalances: make(map[string]sdk.Coins),
		ModuleBalances:  make(map[string]sdk.Coins),
		Delegates:       make(map[string]*stakingtypes.DelegationResponse),
		DelegateRewards: make(map[string]sdk.DecCoins),
	}

	accountBalance, moduleBalance := allBalances(ctx, myApp)
	bdd.AccountBalances = accountBalance
	bdd.ModuleBalances = moduleBalance

	delegatesMap, delegateRewardsMap := delegatesAndRewards(t, ctx, myApp)
	bdd.Delegates = delegatesMap
	bdd.DelegateRewards = delegateRewardsMap
	return bdd
}

func delegatesAndRewards(t *testing.T, ctx sdk.Context, myApp *app.App) (map[string]*stakingtypes.DelegationResponse, map[string]sdk.DecCoins) {
	t.Helper()

	delegatesMap := make(map[string]*stakingtypes.DelegationResponse)
	delegateRewardsMap := make(map[string]sdk.DecCoins)
	distrQuerier := distributionkeeper.NewQuerier(myApp.DistrKeeper)
	stakingQuerier := stakingkeeper.NewQuerier(myApp.StakingKeeper.Keeper)

	onlyKey := ""
	err := myApp.StakingKeeper.IterateAllDelegations(ctx, func(del stakingtypes.Delegation) (stop bool) {
		delegateKey := fmt.Sprintf("%s:%s", del.DelegatorAddress, del.ValidatorAddress)
		if onlyKey != "" && onlyKey != delegateKey {
			return false
		}
		delegation, err := stakingQuerier.Delegation(ctx, &stakingtypes.QueryDelegationRequest{
			DelegatorAddr: del.DelegatorAddress,
			ValidatorAddr: del.ValidatorAddress,
		})
		require.NoError(t, err)
		delegatesMap[delegateKey] = delegation.DelegationResponse
		cacheCtx, _ := ctx.CacheContext()
		rewards, err := distrQuerier.DelegationRewards(cacheCtx, &distributiontypes.QueryDelegationRewardsRequest{
			DelegatorAddress: del.DelegatorAddress,
			ValidatorAddress: del.ValidatorAddress,
		})
		require.NoError(t, err)
		delegateRewardsMap[delegateKey] = rewards.Rewards
		return onlyKey != "" && onlyKey == delegateKey
	})
	require.NoError(t, err)
	return delegatesMap, delegateRewardsMap
}

func checkAppUpgrade(t *testing.T, ctx sdk.Context, myApp *app.App, bdd BeforeUpgradeData) {
	t.Helper()

	checkDelegationData(t, bdd, ctx, myApp)

	checkWrapToken(t, ctx, myApp)

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
	checkMetadataValidate(t, ctx, myApp)
	checkPundiAIFXERC20Token(t, ctx, myApp)

	checkModulesData(t, ctx, myApp)
	checkFXBridgeDenom(t, ctx, myApp)
}

func checkDelegationData(t *testing.T, bdd BeforeUpgradeData, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	beforeDelegates, beforeDelegateRewards := bdd.Delegates, bdd.DelegateRewards
	afterDelegates, afterDelegateRewards := delegatesAndRewards(t, ctx, myApp)

	matchCount, notMatchCount := 0, 0
	notMatchValidators := make(map[string]uint64)
	matchValidators := make(map[string]uint64)
	diff1Count := 0
	for delegateKey, beforeDelegate := range beforeDelegates {
		afterDelegate, ok := afterDelegates[delegateKey]
		require.True(t, ok)
		assert.EqualValues(t, beforeDelegate.Delegation.String(), afterDelegate.Delegation.String(), delegateKey)

		valAddrStr := strings.Split(delegateKey, ":")[1]
		swapBalance := fxtypes.SwapCoin(beforeDelegate.Balance)
		isMatch := swapBalance.Sub(afterDelegate.Balance).Amount.LTE(sdkmath.NewInt(1))
		if swapBalance.Amount.GT(afterDelegate.Balance.Amount) {
			diff1Count++
		}
		assert.Truef(t, isMatch, "key:%s, swapBalance:%s, afterBalance:%s", delegateKey, swapBalance.String(), afterDelegate.Balance.String())
		if isMatch {
			matchCount++
			matchValidators[valAddrStr]++
		} else {
			notMatchCount++
			notMatchValidators[valAddrStr]++
			t.Errorf("not match key:%s, before:[%s],after:[%s]", delegateKey, beforeDelegate.Balance.String(), afterDelegate.Balance.String())
		}
	}
	assert.EqualValuesf(t, 0, notMatchCount, "match count:%d, not match count:%d, diff1Count:%d", matchCount, notMatchCount, diff1Count)
	type kv struct {
		Key   string
		Value uint64
	}

	notMatchValidatorsList := make([]kv, 0, len(notMatchValidators))
	for key, value := range notMatchValidators {
		notMatchValidatorsList = append(notMatchValidatorsList, kv{Key: key, Value: value})
	}
	sort.Slice(notMatchValidatorsList, func(i, j int) bool {
		return notMatchValidatorsList[i].Value > notMatchValidatorsList[j].Value
	})

	for _, valData := range notMatchValidatorsList {
		valAddr, err := sdk.ValAddressFromBech32(valData.Key)
		require.NoError(t, err)
		validator, err := myApp.StakingKeeper.GetValidator(ctx, valAddr)
		require.NoError(t, err)
		t.Logf("val moniker:%s, count:%d, isJailed:%t, token:%s,share:%s", validator.GetMoniker(), valData.Value, validator.IsJailed(), validator.Tokens.String(), validator.DelegatorShares.String())
	}

	rewardMatchCount, rewardNotMatchCount, rewardDiff1Count, rewardDiffGt1Count, diffGtBeforeCount, maxDiff, totalDiff := 0, 0, 0, 0, 0, sdkmath.LegacyNewDec(0), sdkmath.LegacyNewDec(0)
	for delegateKey, beforeRewards := range beforeDelegateRewards {
		afterRewards, ok := afterDelegateRewards[delegateKey]
		require.True(t, ok)
		swapRewards := fxtypes.SwapDecCoins(beforeRewards)
		require.LessOrEqual(t, swapRewards.Len(), 1)
		require.LessOrEqual(t, swapRewards.Len(), afterRewards.Len())
		if swapRewards.IsAllPositive() && swapRewards[0].Denom != fxtypes.DefaultDenom {
			t.Logf("reward not apundiai:%s", swapRewards.String())
			continue
		}
		diff := sdk.NewDecCoins()
		if !swapRewards.IsZero() && swapRewards[0].IsLT(afterRewards[0]) {
			diff = afterRewards.Sub(swapRewards)
			diffGtBeforeCount++
		} else {
			diff = swapRewards.Sub(afterRewards)
		}
		if diff.IsZero() {
			rewardMatchCount++
			continue
		} else if diff[0].Amount.LT(sdkmath.LegacyNewDec(1)) {
			rewardNotMatchCount++
		} else {
			rewardNotMatchCount++
			rewardDiffGt1Count++
			//t.Logf("diff grate 1key:%s, diff:%s, swapRewards:%s, afterRewards:%s", delegateKey, diff, swapRewards.String(), afterRewards.String())
		}
		rewardDiff1Count++
		totalDiff = totalDiff.Add(diff[0].Amount)
		if diff.IsAllPositive() && diff[0].Amount.GT(maxDiff) {
			maxDiff = diff[0].Amount
		}
	}
	t.Logf("reward match count:%d, not match count:%d, diff1Count:%d,rewardDiffGt1Count:%d,diffGtBeforeCount:%d, maxDiff:%s,totalDiff:%s", rewardMatchCount, rewardNotMatchCount, rewardDiff1Count, rewardDiffGt1Count, diffGtBeforeCount, maxDiff.String(), totalDiff.String())
}

func checkIBCTransferModule(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	escrowDenomsMap := nextversion.GetMigrateEscrowDenoms(ctx.ChainID())
	require.NotEmpty(t, escrowDenomsMap)

	for oldDenom, newDenom := range escrowDenomsMap {
		coin := myApp.IBCTransferKeeper.GetTotalEscrowForDenom(ctx, oldDenom)
		require.True(t, coin.IsZero())

		coin = myApp.IBCTransferKeeper.GetTotalEscrowForDenom(ctx, newDenom)
		require.False(t, coin.IsZero())
	}
}

func checkCrisisModule(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	constantFee, err := myApp.CrisisKeeper.ConstantFee.Get(ctx)
	require.NoError(t, err)
	require.EqualValues(t, sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(133).MulRaw(1e18)).String(), constantFee.String())
}

func checkWrapToken(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()
	erc20TokenKeeper := contract.NewERC20TokenKeeper(myApp.EvmKeeper)
	wrapTokenAddress := nextversion.GetWFXAddress(ctx.ChainID())
	name, err := erc20TokenKeeper.Name(ctx, wrapTokenAddress)
	require.NoError(t, err)
	assert.EqualValues(t, nextversion.WrapTokenName, name)
	symbol, err := erc20TokenKeeper.Symbol(ctx, wrapTokenAddress)
	require.NoError(t, err)
	assert.EqualValues(t, nextversion.WrapTokenSymbol, symbol)
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
		if moduleName == ethtypes.ModuleName {
			balance := myApp.BankKeeper.GetBalance(ctx, myApp.AccountKeeper.GetModuleAddress(moduleName), fxtypes.DefaultDenom)
			oldCoins := bdd.ModuleBalances[ethtypes.ModuleName]
			require.Equal(t, balance.Amount.String(), fxtypes.SwapAmount(oldCoins.AmountOf(fxtypes.LegacyFXDenom)).String())
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
		coins = fxtypes.SwapCoins(coins)
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

func checkMetadataValidate(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	myApp.BankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		if len(metadata.DenomUnits) <= 1 {
			return false
		}
		require.NoError(t, metadata.Validate())
		require.NotEqual(t, metadata.Display, metadata.Base)
		return false
	})
}

func checkPundiAIFXERC20Token(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	has, err := myApp.Erc20Keeper.ERC20Token.Has(ctx, fxtypes.DefaultDenom)
	require.NoError(t, err)
	require.True(t, has)
	has, err = myApp.Erc20Keeper.ERC20Token.Has(ctx, fxtypes.LegacyFXDenom)
	require.NoError(t, err)
	require.False(t, has)
}

func checkModulesData(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	checkCrisisModule(t, ctx, myApp)
	checkBankModule(t, ctx, myApp)
	checkEvmParams(t, ctx, myApp)
	checkIBCTransferModule(t, ctx, myApp)
	nextversion.CheckStakingModule(t, ctx, myApp.StakingKeeper.Keeper)
	checkMintModule(t, ctx, myApp)
	nextversion.CheckDistributionModule(t, ctx, myApp.DistrKeeper)
}

func checkEvmParams(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	params := myApp.EvmKeeper.GetParams(ctx)
	require.Equal(t, fxtypes.DefaultDenom, params.EvmDenom)
	require.Equal(t, evmtypes.DefaultHeaderHashNum, params.HeaderHashNum)
}

func checkBankModule(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	totalSupply := sdkmath.ZeroInt()
	myApp.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, balance sdk.Coin) bool {
		require.NotEqual(t, fxtypes.LegacyFXDenom, balance.Denom)
		if balance.Denom == fxtypes.DefaultDenom {
			totalSupply = totalSupply.Add(balance.Amount)
		}
		return false
	})

	supply := myApp.BankKeeper.GetSupply(ctx, fxtypes.DefaultDenom)
	require.Equal(t, totalSupply, supply.Amount)

	myApp.BankKeeper.IterateSendEnabledEntries(ctx, func(denom string, sendEnabled bool) (stop bool) {
		require.NotEqual(t, fxtypes.LegacyFXDenom, denom)
		return false
	})
}

func checkMintModule(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	params, err := myApp.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, fxtypes.DefaultDenom, params.MintDenom)
}

func checkFXBridgeDenom(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	bridgeToken, err := myApp.Erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, fxtypes.FXDenom)
	require.NoError(t, err)
	has, err := myApp.Erc20Keeper.BridgeToken.Has(ctx, collections.Join(ethtypes.ModuleName, fxtypes.LegacyFXDenom))
	require.NoError(t, err)
	require.False(t, has)

	denom, err := myApp.Erc20Keeper.DenomIndex.Get(ctx, erc20types.NewBridgeDenom(ethtypes.ModuleName, bridgeToken.Contract))
	require.NoError(t, err)
	require.Equal(t, fxtypes.FXDenom, denom)

	has, err = myApp.Erc20Keeper.ERC20Token.Has(ctx, fxtypes.LegacyFXDenom)
	require.NoError(t, err)
	require.False(t, has)

	pundiaiERC20Token, err := myApp.Erc20Keeper.ERC20Token.Get(ctx, fxtypes.DefaultDenom)
	require.NoError(t, err)
	denom, err = myApp.Erc20Keeper.DenomIndex.Get(ctx, pundiaiERC20Token.Erc20Address)
	require.NoError(t, err)
	require.Equal(t, fxtypes.DefaultDenom, denom)
}
