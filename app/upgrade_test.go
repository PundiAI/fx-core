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
	"time"

	"cosmossdk.io/collections"
	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	dbm "github.com/cosmos/cosmos-db"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankv2 "github.com/cosmos/cosmos-sdk/x/bank/migrations/v2"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
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
	arbitrumtypes "github.com/pundiai/fx-core/v8/x/arbitrum/types"
	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20v8 "github.com/pundiai/fx-core/v8/x/erc20/migrations/v8"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	fxgovv8 "github.com/pundiai/fx-core/v8/x/gov/migrations/v8"
	fxgovtypes "github.com/pundiai/fx-core/v8/x/gov/types"
	optimismtypes "github.com/pundiai/fx-core/v8/x/optimism/types"
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
	AccountBalances  map[string]sdk.Coins
	ModuleBalances   map[string]sdk.Coins
	Delegates        map[string]*stakingtypes.DelegationResponse
	DelegateRewards  map[string]sdk.DecCoins
	ERC20Token       map[string]erc20types.ERC20Token
	Metadata         map[string]banktypes.Metadata
	EvmBalances      map[string]map[string]*big.Int
	GovDepositAmount map[string]sdkmath.Int
}

func beforeUpgrade(t *testing.T, ctx sdk.Context, myApp *app.App) BeforeUpgradeData {
	t.Helper()
	bdd := BeforeUpgradeData{
		AccountBalances:  make(map[string]sdk.Coins),
		ModuleBalances:   make(map[string]sdk.Coins),
		Delegates:        make(map[string]*stakingtypes.DelegationResponse),
		DelegateRewards:  make(map[string]sdk.DecCoins),
		GovDepositAmount: make(map[string]sdkmath.Int),
	}

	accountBalance, moduleBalance := allBalances(ctx, myApp)
	bdd.AccountBalances = accountBalance
	bdd.ModuleBalances = moduleBalance

	delegatesMap, delegateRewardsMap := delegatesAndRewards(t, ctx, myApp)
	bdd.Delegates = delegatesMap
	bdd.DelegateRewards = delegateRewardsMap

	erc20Token, metadata := allOldDenom(ctx, myApp)
	bdd.ERC20Token = erc20Token
	bdd.Metadata = metadata

	evmBalance, err := allEvmBalance(ctx, myApp)
	require.NoError(t, err)
	bdd.EvmBalances = evmBalance

	depositAmount, err := allGovDepositAmount(ctx, myApp)
	require.NoError(t, err)
	bdd.GovDepositAmount = depositAmount

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

	checkWrapToken(t, ctx, myApp)
	checkDelegationData(t, bdd, ctx, myApp)
	checkStakingMigrationDelete(t, ctx, myApp)
	checkGovCustomParams(t, ctx, myApp)
	checkErc20Keys(t, ctx, myApp)
	checkOutgoingBatch(t, ctx, myApp)

	checkLayer2OracleIsOnline(t, ctx, myApp.Layer2Keeper)

	checkMigrateBalance(t, ctx, myApp, bdd)
	checkTotalSupply(t, ctx, myApp)
	checkEvmSupply(t, ctx, myApp)
	checkEvmBalance(t, ctx, myApp, bdd)

	checkErc20Token(t, ctx, myApp, bdd)
	checkNewErc20Token(t, ctx, myApp)
	checkDefaultDenom(t, ctx, myApp)
	checkMetadata(t, ctx, myApp, bdd)

	checkModulesData(t, ctx, myApp)
	checkBridgeAddress(t, ctx, myApp)
	checkGovProposal(t, ctx, myApp)
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
		var diff sdk.DecCoins
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
			// t.Logf("diff grate 1key:%s, diff:%s, swapRewards:%s, afterRewards:%s", delegateKey, diff, swapRewards.String(), afterRewards.String())
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
	checkAccountBalance(t, ctx, myApp, bdd.AccountBalances, newAccountBalance, bdd.GovDepositAmount)

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

func checkAccountBalance(t *testing.T, ctx sdk.Context, myApp *app.App, accountBalances, newAccountBalance map[string]sdk.Coins, govDepositAmount map[string]sdkmath.Int) {
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
			ibcTokens, err := myApp.Erc20Keeper.GetBaseIBCTokens(ctx, baseDenom)
			require.NoError(t, err)
			for _, ibcToken := range ibcTokens {
				tokenDenoms[ibcToken.GetIbcDenom()] = struct{}{}
			}

			amount := sdkmath.ZeroInt()
			for denom := range tokenDenoms {
				amount = amount.Add(coins.AmountOf(denom))
			}

			if depositAmount, ok := govDepositAmount[addrStr]; ok && baseDenom == fxtypes.DefaultDenom {
				amount = amount.Add(fxtypes.SwapAmount(depositAmount))
			}

			balance := myApp.BankKeeper.GetBalance(ctx, addr, baseDenom)
			require.Equal(t, balance.Amount, amount)
		}
	}
	require.Empty(t, newAccountBalance)
}

func checkEvmSupply(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	erc20Tokens, err := allErc20Token(ctx, myApp)
	require.NoError(t, err)

	erc20ModuleAddr := authtypes.NewModuleAddress(erc20types.ModuleName)
	erc20TokenKeeper := contract.NewERC20TokenKeeper(myApp.EvmKeeper)
	for _, et := range erc20Tokens {
		if et.GetDenom() == fxtypes.DefaultDenom {
			balance := myApp.BankKeeper.GetBalance(ctx, et.GetERC20Contract().Bytes(), et.GetDenom())
			totalSupply, err := erc20TokenKeeper.TotalSupply(ctx, et.GetERC20Contract())
			require.NoError(t, err)
			require.GreaterOrEqual(t, balance.Amount.String(), totalSupply.String(), et.GetDenom())
			continue
		}

		if et.ContractOwner == erc20types.OWNER_MODULE {
			balance := myApp.BankKeeper.GetBalance(ctx, erc20ModuleAddr, et.GetDenom())
			totalSupply, err := erc20TokenKeeper.TotalSupply(ctx, et.GetERC20Contract())
			require.NoError(t, err)
			require.Equal(t, sdkmath.NewIntFromBigInt(totalSupply).String(), balance.Amount.String(), et.GetDenom())
			continue
		}

		if et.ContractOwner == erc20types.OWNER_EXTERNAL {
			totalSupply := myApp.BankKeeper.GetSupply(ctx, et.GetDenom())
			balanceOf, err := erc20TokenKeeper.BalanceOf(ctx, et.GetERC20Contract(), common.BytesToAddress(erc20ModuleAddr.Bytes()))
			require.NoError(t, err)
			require.Equal(t, sdkmath.NewIntFromBigInt(balanceOf).String(), totalSupply.Amount.String(), et.GetDenom())
			continue
		}
		assert.Failf(t, "not check supply %s", et.GetDenom())
	}
}

func checkTotalSupply(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	erc20Tokens, err := allErc20Token(ctx, myApp)
	require.NoError(t, err)

	for _, et := range erc20Tokens {
		aliasDenoms := make([]string, 0, 10)
		nativeDenoms := make([]string, 0, 10)

		bridgeTokens, err := myApp.Erc20Keeper.GetBaseBridgeTokens(ctx, et.GetDenom())
		require.NoError(t, err)
		for _, bt := range bridgeTokens {
			if bt.IsOrigin() {
				continue
			}
			if bt.IsNative {
				nativeDenoms = append(nativeDenoms, bt.BridgeDenom())
			} else {
				aliasDenoms = append(aliasDenoms, bt.BridgeDenom())
			}
		}

		ibcTokens, err := myApp.Erc20Keeper.GetBaseIBCTokens(ctx, et.GetDenom())
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
		require.Equal(t, aliasTotal, baseSupply.Amount)

		nativeTotal := sdkmath.ZeroInt()
		for _, denom := range nativeDenoms {
			supply := myApp.BankKeeper.GetSupply(ctx, denom)
			nativeTotal = nativeTotal.Add(supply.Amount)
		}
		crossChainBalance := myApp.BankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(crosschaintypes.ModuleName), et.GetDenom())
		require.True(t, crossChainBalance.Amount.GTE(nativeTotal))
	}
}

func checkEvmBalance(t *testing.T, ctx sdk.Context, myApp *app.App, bdd BeforeUpgradeData) {
	t.Helper()

	erc20TokenKeeper := contract.NewERC20TokenKeeper(myApp.EvmKeeper)
	for erc20Addr, bals := range bdd.EvmBalances {
		baseDenom, err := myApp.Erc20Keeper.GetBaseDenom(ctx, erc20Addr)
		require.NoError(t, err)
		myApp.AccountKeeper.IterateAccounts(ctx, func(account sdk.AccountI) (stop bool) {
			addr := common.BytesToAddress(account.GetAddress())
			bal, ok := bals[addr.String()]
			if !ok {
				bal = big.NewInt(0)
			}
			if baseDenom == fxtypes.DefaultDenom {
				bal = fxtypes.SwapAmount(sdkmath.NewIntFromBigInt(bal)).BigInt()
			}
			newBal, err := erc20TokenKeeper.BalanceOf(ctx, common.HexToAddress(erc20Addr), addr)
			require.NoError(t, err)
			require.Equalf(t, bal.String(), newBal.String(), "address: %s", addr.String())
			return false
		})
	}
}

func checkErc20Token(t *testing.T, ctx sdk.Context, myApp *app.App, bdd BeforeUpgradeData) {
	t.Helper()

	for addr, et := range bdd.ERC20Token {
		baseDenom, err := myApp.Erc20Keeper.GetBaseDenom(ctx, addr)
		require.NoError(t, err)

		if baseDenom == fxtypes.DefaultDenom {
			require.Equal(t, fxtypes.LegacyFXDenom, et.Denom)
			continue
		}

		erc20Token, err := myApp.Erc20Keeper.GetERC20Token(ctx, baseDenom)
		require.NoError(t, err, baseDenom)
		require.Equal(t, et.GetErc20Address(), erc20Token.GetErc20Address())

		if baseDenom == "pundix" {
			bridgeToken, err := myApp.Erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, baseDenom)
			require.NoError(t, err)
			require.Equal(t, bridgeToken.BridgeDenom(), et.Denom)
			continue
		}
		if baseDenom == "purse" {
			ibcToken, err := myApp.Erc20Keeper.GetIBCToken(ctx, baseDenom, fxtypes.PundixChannel)
			require.NoError(t, err)
			require.Equal(t, ibcToken.GetIbcDenom(), et.Denom)
			continue
		}
		require.Equal(t, baseDenom, et.Denom)
	}
}

func checkMetadata(t *testing.T, ctx sdk.Context, myApp *app.App, bdd BeforeUpgradeData) {
	t.Helper()

	myApp.BankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		if len(metadata.DenomUnits) <= 1 {
			return false
		}
		require.NoError(t, metadata.Validate())
		require.NotEqual(t, metadata.Display, metadata.Base)
		if metadata.Base == fxtypes.DefaultDenom {
			require.Equal(t, metadata, fxtypes.NewDefaultMetadata())
		}
		return false
	})

	for denom, md := range bdd.Metadata {
		// one to one
		if len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0 {
			require.False(t, myApp.BankKeeper.HasDenomMetaData(ctx, denom))
			if md.Base != fxtypes.LegacyFXDenom {
				baseDenom, err := myApp.Erc20Keeper.GetBaseDenom(ctx, md.Base)
				require.NoError(t, err)
				if strings.HasPrefix(md.Base, "ibc/") {
					ibcToken, err := myApp.Erc20Keeper.GetIBCToken(ctx, baseDenom, fxtypes.PundixChannel)
					require.NoError(t, err)
					require.Equal(t, ibcToken.GetIbcDenom(), md.Base)
				} else {
					bridgeToken, err := myApp.Erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, baseDenom)
					require.NoError(t, err)
					require.Equal(t, bridgeToken.BridgeDenom(), md.Base)
				}
			}
			continue
		}

		// many to one
		for _, alias := range md.DenomUnits[0].Aliases {
			if strings.HasPrefix(alias, arbitrumtypes.ModuleName) ||
				strings.HasPrefix(alias, optimismtypes.ModuleName) {
				continue
			}
			baseDenom, err := myApp.Erc20Keeper.GetBaseDenom(ctx, alias)
			require.NoError(t, err)

			// ibc token
			if strings.HasPrefix(alias, ibctransfertypes.DenomPrefix) {
				channelID, ok := erc20v8.GetIBCDenomTrace(ctx, alias)
				require.True(t, ok)

				ibcToken, err := myApp.Erc20Keeper.GetIBCToken(ctx, baseDenom, channelID)
				require.NoError(t, err, baseDenom)
				require.Equal(t, ibcToken.GetIbcDenom(), alias)
				continue
			}

			// bridge token
			matched := false
			for _, k := range myApp.CrosschainKeepers.ToSlice() {
				if !strings.HasPrefix(alias, k.ModuleName()) {
					continue
				}
				bridgeToken, err := myApp.Erc20Keeper.GetBridgeToken(ctx, k.ModuleName(), baseDenom)
				require.NoError(t, err, baseDenom)
				require.Equal(t, bridgeToken.BridgeDenom(), alias)
				matched = true
			}
			require.Truef(t, matched, "skip %s alias %s", baseDenom, alias)
		}
	}
}

func checkNewErc20Token(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	erc20Token, err := allErc20Token(ctx, myApp)
	require.NoError(t, err)
	for _, et := range erc20Token {
		baseDenom, err := myApp.Erc20Keeper.GetBaseDenom(ctx, et.Erc20Address)
		require.NoError(t, err)

		bridgeTokens, err := myApp.Erc20Keeper.GetBaseBridgeTokens(ctx, baseDenom)
		require.NoError(t, err)
		for _, token := range bridgeTokens {
			bridgeBaseDenom, err := myApp.Erc20Keeper.GetBaseDenom(ctx, erc20types.NewBridgeDenom(token.ChainName, token.Contract))
			require.NoError(t, err)
			require.Equal(t, baseDenom, bridgeBaseDenom)
		}

		ibcTokens, err := myApp.Erc20Keeper.GetBaseIBCTokens(ctx, baseDenom)
		require.NoError(t, err)
		for _, token := range ibcTokens {
			ibcBaseDenom, err := myApp.Erc20Keeper.GetBaseDenom(ctx, token.GetIbcDenom())
			require.NoError(t, err)
			require.Equal(t, baseDenom, ibcBaseDenom)
		}
	}

	// check index to base
	iter, err := myApp.Erc20Keeper.DenomIndex.Iterate(ctx, nil)
	require.NoError(t, err)
	defer iter.Close()
	kvs, err := iter.KeyValues()
	require.NoError(t, err)
	for _, kv := range kvs {
		if common.IsHexAddress(kv.Key) {
			et, err := myApp.Erc20Keeper.GetERC20Token(ctx, kv.Value)
			require.NoError(t, err)
			require.Equal(t, kv.Key, et.GetErc20Address())
			continue
		}
		if strings.HasPrefix(kv.Key, "ibc/") {
			ibcTokens, err := myApp.Erc20Keeper.GetBaseIBCTokens(ctx, kv.Value)
			require.NoError(t, err)
			has := false
			for _, it := range ibcTokens {
				if kv.Key == it.GetIbcDenom() {
					has = true
					break
				}
			}
			require.True(t, has)
			continue
		}
		bridgeTokens, err := myApp.Erc20Keeper.GetBaseBridgeTokens(ctx, kv.Value)
		require.NoError(t, err)
		has := false
		for _, bt := range bridgeTokens {
			if kv.Key == erc20types.NewBridgeDenom(bt.ChainName, bt.Contract) {
				has = true
				break
			}
		}
		require.Truef(t, has, "key %s value %s", kv.Key, kv.Value)
	}
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

func checkDefaultDenom(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	has, err := myApp.Erc20Keeper.BridgeToken.Has(ctx, collections.Join(ethtypes.ModuleName, fxtypes.LegacyFXDenom))
	require.NoError(t, err)
	require.False(t, has)
	has, err = myApp.Erc20Keeper.ERC20Token.Has(ctx, fxtypes.LegacyFXDenom)
	require.NoError(t, err)
	require.False(t, has)

	bridgeToken, err := myApp.Erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, fxtypes.FXDenom)
	require.NoError(t, err)
	denom, err := myApp.Erc20Keeper.DenomIndex.Get(ctx, erc20types.NewBridgeDenom(ethtypes.ModuleName, bridgeToken.Contract))
	require.NoError(t, err)
	require.Equal(t, fxtypes.FXDenom, denom)
	has, err = myApp.Erc20Keeper.ERC20Token.Has(ctx, fxtypes.FXDenom)
	require.NoError(t, err)
	require.False(t, has)

	pundiaiERC20Token, err := myApp.Erc20Keeper.ERC20Token.Get(ctx, fxtypes.DefaultDenom)
	require.NoError(t, err)
	denom, err = myApp.Erc20Keeper.DenomIndex.Get(ctx, pundiaiERC20Token.Erc20Address)
	require.NoError(t, err)
	require.Equal(t, fxtypes.DefaultDenom, denom)
	bridgeToken, err = myApp.Erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, fxtypes.DefaultDenom)
	require.NoError(t, err)
	require.Equal(t, nextversion.GetPundiaiTokenAddr(ctx).String(), bridgeToken.Contract)
}

func checkBridgeAddress(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	moduleAccount := myApp.AccountKeeper.GetModuleAccount(ctx, crosschaintypes.ModuleName)
	assert.NotEmpty(t, moduleAccount)

	bridgeCallSender := authtypes.NewModuleAddress(crosschaintypes.BridgeCallSender)
	assert.True(t, myApp.AccountKeeper.HasAccount(ctx, bridgeCallSender))

	bridgeFeeCollector := authtypes.NewModuleAddress(crosschaintypes.BridgeFeeCollectorName)
	assert.True(t, myApp.AccountKeeper.HasAccount(ctx, bridgeFeeCollector))
}

func checkGovProposal(t *testing.T, ctx sdk.Context, myApp *app.App) {
	t.Helper()

	proposalCount := 0
	handle := func(_ collections.Pair[time.Time, uint64], _ uint64) (stop bool, err error) {
		proposalCount++
		return false, nil
	}
	err := myApp.GovKeeper.InactiveProposalsQueue.Walk(ctx, nil, handle)
	require.NoError(t, err)
	err = myApp.GovKeeper.ActiveProposalsQueue.Walk(ctx, nil, handle)
	require.NoError(t, err)
	require.Zero(t, proposalCount)

	depositCount := 0
	err = myApp.GovKeeper.Deposits.Walk(ctx, nil, func(_ collections.Pair[uint64, sdk.AccAddress], _ govv1.Deposit) (stop bool, err error) {
		depositCount++
		return false, nil
	})
	require.NoError(t, err)
	require.Zero(t, depositCount)

	voteCount := 0
	err = myApp.GovKeeper.Votes.Walk(ctx, nil, func(_ collections.Pair[uint64, sdk.AccAddress], _ govv1.Vote) (stop bool, err error) {
		voteCount++
		return false, nil
	})
	require.NoError(t, err)
	require.Zero(t, voteCount)
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

func allOldDenom(ctx sdk.Context, myApp *app.App) (map[string]erc20types.ERC20Token, map[string]banktypes.Metadata) {
	erc20Token := make(map[string]erc20types.ERC20Token)
	erc20Store := ctx.KVStore(myApp.GetKey(erc20types.StoreKey))
	iterator := storetypes.KVStorePrefixIterator(erc20Store, erc20v8.KeyPrefixTokenPair)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var tokenPair erc20types.ERC20Token
		myApp.AppCodec().MustUnmarshal(iterator.Value(), &tokenPair)
		erc20Token[tokenPair.GetErc20Address()] = tokenPair
	}

	metadata := make(map[string]banktypes.Metadata)
	storeService := runtime.NewKVStoreService(myApp.GetKey(banktypes.StoreKey))
	bankStore := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	oldDenomMetaDataStore := prefix.NewStore(bankStore, bankv2.DenomMetadataPrefix)
	oldDenomMetaDataIter := oldDenomMetaDataStore.Iterator(nil, nil)
	for ; oldDenomMetaDataIter.Valid(); oldDenomMetaDataIter.Next() {
		var md banktypes.Metadata
		myApp.AppCodec().MustUnmarshal(oldDenomMetaDataIter.Value(), &md)
		metadata[md.Base] = md
	}
	return erc20Token, metadata
}

func allErc20Token(ctx sdk.Context, myApp *app.App) ([]erc20types.ERC20Token, error) {
	iter, err := myApp.Erc20Keeper.ERC20Token.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	kvs, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}

	erc20Tokens := make([]erc20types.ERC20Token, 0, len(kvs))
	for _, kv := range kvs {
		erc20Tokens = append(erc20Tokens, kv.Value)
	}
	return erc20Tokens, nil
}

func allEvmBalance(ctx sdk.Context, myApp *app.App) (map[string]map[string]*big.Int, error) {
	erc20Tokens, _ := allOldDenom(ctx, myApp)
	accs := make([]sdk.AccAddress, 0, 100)
	myApp.AccountKeeper.IterateAccounts(ctx, func(account sdk.AccountI) (stop bool) {
		accs = append(accs, account.GetAddress())
		return false
	})
	evmBalance := make(map[string]map[string]*big.Int)
	for _, et := range erc20Tokens {
		if et.GetDenom() != fxtypes.LegacyFXDenom {
			continue
		}
		erc20Balances := make(map[string]*big.Int)
		erc20TokenKeeper := contract.NewERC20TokenKeeper(myApp.EvmKeeper)
		for _, acc := range accs {
			addr := common.BytesToAddress(acc.Bytes())
			balanceOf, err := erc20TokenKeeper.BalanceOf(ctx, et.GetERC20Contract(), addr)
			if err != nil {
				return nil, err
			}
			if balanceOf.Cmp(big.NewInt(0)) != 0 {
				erc20Balances[addr.String()] = balanceOf
			}
		}
		evmBalance[et.GetErc20Address()] = erc20Balances
	}
	return evmBalance, nil
}

func allGovDepositAmount(ctx sdk.Context, myApp *app.App) (map[string]sdkmath.Int, error) {
	depositAmount := make(map[string]sdkmath.Int)
	handle := func(key collections.Pair[time.Time, uint64], _ uint64) (stop bool, err error) {
		return false, myApp.GovKeeper.IterateDeposits(ctx, key.K2(), func(key collections.Pair[uint64, sdk.AccAddress], value govv1.Deposit) (bool, error) {
			acc := key.K2().String()
			amount, ok := depositAmount[acc]
			if !ok {
				amount = sdkmath.ZeroInt()
			}
			depositAmount[acc] = amount.Add(sdk.NewCoins(value.Amount...).AmountOf(fxtypes.LegacyFXDenom))
			return false, nil
		})
	}
	if err := myApp.GovKeeper.InactiveProposalsQueue.Walk(ctx, nil, handle); err != nil {
		return nil, err
	}
	return depositAmount, myApp.GovKeeper.ActiveProposalsQueue.Walk(ctx, nil, handle)
}
