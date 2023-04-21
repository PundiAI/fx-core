package app_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v4/app"
	v4 "github.com/functionx/fx-core/v4/app/upgrades/v4"
	"github.com/functionx/fx-core/v4/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v4/types"
	arbitrumtypes "github.com/functionx/fx-core/v4/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v4/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v4/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v4/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v4/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v4/x/eth/types"
	evmtypes "github.com/functionx/fx-core/v4/x/evm/types"
	fxgovtypes "github.com/functionx/fx-core/v4/x/gov/types"
	optimismtypes "github.com/functionx/fx-core/v4/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v4/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v4/x/tron/types"
)

func Test_Upgrade(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test: ", t.Name())

	fxtypes.SetConfig(true)

	testCases := []struct {
		name                  string
		fromVersion           int
		toVersion             int
		LocalStoreBlockHeight uint64
		plan                  upgradetypes.Plan
	}{
		{
			name:        "upgrade v4",
			fromVersion: 3,
			toVersion:   4,
			plan: upgradetypes.Plan{
				Name: v4.Upgrade.UpgradeName,
				Info: "local test upgrade v4",
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
	myApp.SetStoreLoader(upgradetypes.UpgradeStoreLoader(myApp.LastBlockHeight()+1, v4.Upgrade.StoreUpgrades()))
	err = myApp.LoadLatestVersion()
	require.NoError(t, err)

	ctx := newContext(t, myApp)

	// UpgradeBefore

	// check arbitrum and optimism register usdt ,weth
	checkDenomMetaData(t, ctx, myApp, true)
	// check params migrated from x/param to module erc20
	checkERC20MigrateParamStore(t, ctx, myApp, true)
	// check params migrated from x/param to module crosschain
	checkCrossChainMigrateParamStore(t, ctx, myApp, true)
	// check fxgovparams
	checkFXGovParams(t, ctx, myApp, true)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			checkVersionMap(t, ctx, myApp, getConsensusVersion(testCase.fromVersion))
			testCase.plan.Height = ctx.BlockHeight()

			myApp.UpgradeKeeper.ApplyUpgrade(ctx, testCase.plan)

			checkVersionMap(t, ctx, myApp, getConsensusVersion(testCase.toVersion))
		})
	}

	// UpgradeAfter
	checkDenomMetaData(t, ctx, myApp, false)
	checkERC20MigrateParamStore(t, ctx, myApp, false)
	checkCrossChainMigrateParamStore(t, ctx, myApp, false)
	checkFXGovParams(t, ctx, myApp, false)

	checkFIP20LogicUpgrade(t, ctx, myApp)
	checkWFXLogicUpgrade(t, ctx, myApp)
	checkCrossChainOracleDelegateInfo(t, myApp, ctx)

	myApp.EthKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.BscKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.TronKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.PolygonKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.AvalancheKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
}

func checkCrossChainOracleDelegateInfo(t *testing.T, myApp *app.App, ctx sdk.Context) {
	bscRemoveOracles := v4.GetBscRemoveOracles(ctx.ChainID())
	// list to map
	bscRemoveOraclesMap := make(map[string]bool)
	for _, oracle := range bscRemoveOracles {
		bscRemoveOraclesMap[oracle] = true
	}
	crosschainModules := []string{ethtypes.ModuleName, bsctypes.ModuleName, trontypes.ModuleName, polygontypes.ModuleName, avalanchetypes.ModuleName, arbitrumtypes.ModuleName, optimismtypes.ModuleName}
	for _, crosschainModule := range crosschainModules {
		oracles, err := myApp.CrosschainKeeper.Oracles(ctx, &crosschaintypes.QueryOraclesRequest{ChainName: crosschainModule})
		if err != nil {
			ctx.Logger().Error("query oracles error", "module", crosschainModule, "err", err)
			continue
		}
		for _, oracle := range oracles.Oracles {
			delegateAddress := oracle.GetDelegateAddress(crosschainModule)
			startingInfo := myApp.DistrKeeper.GetDelegatorStartingInfo(ctx, oracle.GetValidator(), delegateAddress)
			if crosschainModule == bsctypes.ModuleName && bscRemoveOraclesMap[oracle.GetOracle().String()] {
				require.EqualValues(t, uint64(0), startingInfo.Height)
				require.EqualValues(t, uint64(0), startingInfo.PreviousPeriod)
				require.True(t, startingInfo.Stake.IsNil())
				continue
			}
			require.True(t, startingInfo.Height > 0)
			require.True(t, startingInfo.PreviousPeriod > 0)
			require.EqualValues(t, sdk.NewDecFromInt(sdkmath.NewInt(10_000).MulRaw(1e18)).String(), startingInfo.Stake.String())

			// test can get rewards
			_, err = myApp.DistrKeeper.DelegationRewards(ctx, &distributiontypes.QueryDelegationRewardsRequest{
				DelegatorAddress: delegateAddress.String(),
				ValidatorAddress: oracle.GetValidator().String(),
			})
			require.NoError(t, err)
		}
	}
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

func checkDenomMetaData(t *testing.T, ctx sdk.Context, myApp *app.App, isUpgradeBefore bool) {
	denomAlias := v4.GetUpdateDenomAlias(ctx.ChainID())
	for _, da := range denomAlias {
		denomKey := da.Denom
		if isUpgradeBefore {
			denomKey = da.Denom
			_, found := myApp.BankKeeper.GetDenomMetaData(ctx, denomKey)
			assert.False(t, found)
			continue
		}
		// todo testnet not deployed weth
		if os.Getenv("CHAIN_ID") == fxtypes.TestnetChainId && da.Denom == "weth" {
			_, found := myApp.BankKeeper.GetDenomMetaData(ctx, denomKey)
			assert.False(t, found)
		} else {
			md, found := myApp.BankKeeper.GetDenomMetaData(ctx, denomKey)
			assert.True(t, found)
			assert.True(t, len(md.DenomUnits) > 0)
			assert.True(t, len(md.DenomUnits[0].Aliases) > 0)
			if isUpgradeBefore {
				assert.False(t, contain(md.DenomUnits[0].Aliases, da.Alias))
			} else {
				assert.True(t, contain(md.DenomUnits[0].Aliases, da.Alias))
			}
		}
	}
}

func checkFIP20LogicUpgrade(t *testing.T, ctx sdk.Context, myApp *app.App) {
	// check fip20 logic upgrade
	fipLogicAcc := myApp.EvmKeeper.GetAccount(ctx, fxtypes.GetFIP20().Address)
	require.True(t, fipLogicAcc.IsContract())

	fipLogic := fxtypes.GetFIP20()
	codeHash := crypto.Keccak256Hash(fipLogic.Code)
	require.Equal(t, codeHash.Bytes(), fipLogicAcc.CodeHash)

	code := myApp.EvmKeeper.GetCode(ctx, codeHash)
	require.Equal(t, fipLogic.Code, code)
}

func checkWFXLogicUpgrade(t *testing.T, ctx sdk.Context, myApp *app.App) {
	// check wfx logic upgrade
	wfxLogicAcc := myApp.EvmKeeper.GetAccount(ctx, fxtypes.GetWFX().Address)
	require.True(t, wfxLogicAcc.IsContract())

	wfxLogic := fxtypes.GetWFX()
	codeHash := crypto.Keccak256Hash(wfxLogic.Code)
	require.Equal(t, codeHash.Bytes(), wfxLogicAcc.CodeHash)

	code := myApp.EvmKeeper.GetCode(ctx, codeHash)
	require.Equal(t, wfxLogic.Code, code)
}

func checkFXGovParams(t *testing.T, ctx sdk.Context, myApp *app.App, isUpgradeBefore bool) {
	defaultParams := myApp.GovKeeper.GetParams(ctx, "")
	checkGovERC20ParamsUpgradeBefore(t, ctx, myApp, defaultParams, isUpgradeBefore)
	checkGovEVMParamsUpgradeBefore(t, ctx, myApp, defaultParams, isUpgradeBefore)
	checkGovEGFParamsUpgradeBefore(t, ctx, myApp, defaultParams, isUpgradeBefore)
}

func checkGovEGFParamsUpgradeBefore(t *testing.T, ctx sdk.Context, myApp *app.App, defaultParams fxgovtypes.Params, isUpgradeBefore bool) {
	if isUpgradeBefore {
		egfParams := myApp.GovKeeper.GetParams(ctx, "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal")
		assert.NoError(t, egfParams.ValidateBasic())
		assert.NoError(t, egfParams.ValidateBasic())
		assert.EqualValues(t, egfParams.MinDeposit, defaultParams.MinDeposit)
		assert.EqualValues(t, egfParams.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, fxgovtypes.DefaultMinInitialDeposit).String())
		assert.EqualValues(t, egfParams.MaxDepositPeriod, defaultParams.MaxDepositPeriod)
		assert.EqualValues(t, egfParams.VotingPeriod.String(), defaultParams.VotingPeriod.String())
		assert.EqualValues(t, egfParams.VetoThreshold, defaultParams.VetoThreshold)
		assert.EqualValues(t, egfParams.Threshold, defaultParams.Threshold)
		assert.EqualValues(t, egfParams.Quorum, defaultParams.Quorum)

		egf := myApp.GovKeeper.GetEGFParams(ctx)
		assert.Error(t, egf.ValidateBasic())
		assert.False(t, egf.EgfDepositThreshold.IsValid())
		assert.EqualValues(t, egf.ClaimRatio, "")
		return
	}
	egfParams := myApp.GovKeeper.GetParams(ctx, "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal")
	assert.NoError(t, egfParams.ValidateBasic())
	assert.EqualValues(t, egfParams.MinDeposit, defaultParams.MinDeposit)
	assert.EqualValues(t, egfParams.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, fxgovtypes.DefaultMinInitialDeposit).String())
	assert.EqualValues(t, egfParams.MaxDepositPeriod, defaultParams.MaxDepositPeriod)
	assert.EqualValues(t, egfParams.VotingPeriod.String(), fxgovtypes.DefaultEgfVotingPeriod.String())
	assert.EqualValues(t, egfParams.VetoThreshold, defaultParams.VetoThreshold)
	assert.EqualValues(t, egfParams.Threshold, defaultParams.Threshold)
	assert.EqualValues(t, egfParams.Quorum, defaultParams.Quorum)

	egf := myApp.GovKeeper.GetEGFParams(ctx)
	assert.NoError(t, egf.ValidateBasic())
	assert.EqualValues(t, egf.EgfDepositThreshold, sdk.NewCoin(fxtypes.DefaultDenom, fxgovtypes.DefaultEgfDepositThreshold))
	assert.EqualValues(t, egf.ClaimRatio, fxgovtypes.DefaultClaimRatio.String())
}

func checkGovEVMParamsUpgradeBefore(t *testing.T, ctx sdk.Context, myApp *app.App, defaultParams fxgovtypes.Params, isUpgradeBefore bool) {
	if isUpgradeBefore {
		evmParams := myApp.GovKeeper.GetParams(ctx, sdk.MsgTypeURL(&evmtypes.MsgCallContract{}))
		assert.NoError(t, evmParams.ValidateBasic())
		assert.EqualValues(t, evmParams.MinDeposit, defaultParams.MinDeposit)
		assert.EqualValues(t, evmParams.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, fxgovtypes.DefaultMinInitialDeposit).String())
		assert.EqualValues(t, evmParams.MaxDepositPeriod, defaultParams.MaxDepositPeriod)
		assert.EqualValues(t, evmParams.VotingPeriod.String(), defaultParams.VotingPeriod.String())
		assert.EqualValues(t, evmParams.VetoThreshold, defaultParams.VetoThreshold)
		assert.EqualValues(t, evmParams.Threshold, defaultParams.Threshold)
		assert.EqualValues(t, evmParams.Quorum, defaultParams.Quorum)
		return
	}
	evmParams := myApp.GovKeeper.GetParams(ctx, sdk.MsgTypeURL(&evmtypes.MsgCallContract{}))
	assert.NoError(t, evmParams.ValidateBasic())
	assert.EqualValues(t, evmParams.MinDeposit, defaultParams.MinDeposit)
	assert.EqualValues(t, evmParams.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, fxgovtypes.DefaultMinInitialDeposit).String())
	assert.EqualValues(t, evmParams.MaxDepositPeriod, defaultParams.MaxDepositPeriod)
	assert.EqualValues(t, evmParams.VotingPeriod.String(), fxgovtypes.DefaultEvmVotingPeriod.String())
	assert.EqualValues(t, evmParams.VetoThreshold, defaultParams.VetoThreshold)
	assert.EqualValues(t, evmParams.Threshold, defaultParams.Threshold)
	assert.EqualValues(t, evmParams.Quorum, fxgovtypes.DefaultEvmQuorum.String())
}

func checkGovERC20ParamsUpgradeBefore(t *testing.T, ctx sdk.Context, myApp *app.App, defaultParams fxgovtypes.Params, isUpgradeBefore bool) {
	if isUpgradeBefore {
		erc20MsgType := []string{
			"/fx.erc20.v1.RegisterCoinProposal",
			"/fx.erc20.v1.RegisterERC20Proposal",
			"/fx.erc20.v1.ToggleTokenConversionProposal",
			"/fx.erc20.v1.UpdateDenomAliasProposal",
			sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
			sdk.MsgTypeURL(&erc20types.MsgRegisterERC20{}),
			sdk.MsgTypeURL(&erc20types.MsgToggleTokenConversion{}),
			sdk.MsgTypeURL(&erc20types.MsgUpdateDenomAlias{}),
		}
		for _, erc20MsgType := range erc20MsgType {
			// registered Msg
			erc20params := myApp.GovKeeper.GetParams(ctx, erc20MsgType)
			assert.NoError(t, erc20params.ValidateBasic())
			assert.EqualValues(t, erc20params.MinDeposit, defaultParams.MinDeposit)
			assert.EqualValues(t, erc20params.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, fxgovtypes.DefaultMinInitialDeposit).String())
			assert.EqualValues(t, erc20params.MaxDepositPeriod, defaultParams.MaxDepositPeriod)
			assert.EqualValues(t, erc20params.VotingPeriod, defaultParams.VotingPeriod)
			assert.EqualValues(t, erc20params.VetoThreshold, defaultParams.VetoThreshold)
			assert.EqualValues(t, erc20params.Threshold, defaultParams.Threshold)
			assert.EqualValues(t, erc20params.Quorum, defaultParams.Quorum)
		}
		return
	}
	erc20MsgType := []string{
		"/fx.erc20.v1.RegisterCoinProposal",
		"/fx.erc20.v1.RegisterERC20Proposal",
		"/fx.erc20.v1.ToggleTokenConversionProposal",
		"/fx.erc20.v1.UpdateDenomAliasProposal",
		sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
		sdk.MsgTypeURL(&erc20types.MsgRegisterERC20{}),
		sdk.MsgTypeURL(&erc20types.MsgToggleTokenConversion{}),
		sdk.MsgTypeURL(&erc20types.MsgUpdateDenomAlias{}),
	}
	for _, erc20MsgType := range erc20MsgType {
		// registered Msg
		erc20params := myApp.GovKeeper.GetParams(ctx, erc20MsgType)
		assert.NoError(t, erc20params.ValidateBasic())
		assert.EqualValues(t, erc20params.MinDeposit, defaultParams.MinDeposit)
		assert.EqualValues(t, erc20params.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, fxgovtypes.DefaultMinInitialDeposit).String())
		assert.EqualValues(t, erc20params.MaxDepositPeriod, defaultParams.MaxDepositPeriod)
		assert.EqualValues(t, erc20params.VotingPeriod, defaultParams.VotingPeriod)
		assert.EqualValues(t, erc20params.VetoThreshold, defaultParams.VetoThreshold)
		assert.EqualValues(t, erc20params.Threshold, defaultParams.Threshold)
		assert.EqualValues(t, erc20params.Quorum, fxgovtypes.DefaultErc20Quorum.String())
	}
}

func checkERC20MigrateParamStore(t *testing.T, ctx sdk.Context, myApp *app.App, isUpgradeBefore bool) {
	subspace := myApp.GetSubspace(erc20types.ModuleName)
	var subspaceParams erc20types.Params
	if isUpgradeBefore {
		subspace = subspace.WithKeyTable(erc20types.ParamKeyTable())
		subspace.GetParamSet(ctx, &subspaceParams)
		params := myApp.Erc20Keeper.GetParams(ctx)
		assert.NotEqualValues(t, params.EnableErc20, subspaceParams.EnableErc20)
		assert.NotEqualValues(t, params.EnableEVMHook, subspaceParams.EnableEVMHook)
		assert.NotEqualValues(t, params.IbcTimeout, subspaceParams.IbcTimeout)
		return
	}
	subspace.GetParamSet(ctx, &subspaceParams)
	params := myApp.Erc20Keeper.GetParams(ctx)
	assert.EqualValues(t, params.EnableErc20, subspaceParams.EnableErc20)
	assert.EqualValues(t, params.EnableEVMHook, subspaceParams.EnableEVMHook)
	assert.EqualValues(t, params.IbcTimeout, subspaceParams.IbcTimeout)
}

func checkCrossChainMigrateParamStore(t *testing.T, ctx sdk.Context, myApp *app.App, isUpgradeBefore bool) {
	crosschainsModule := []string{avalanchetypes.ModuleName, bsctypes.ModuleName, ethtypes.ModuleName, polygontypes.ModuleName, trontypes.ModuleName}
	for _, moduleName := range crosschainsModule {
		subspace := myApp.GetSubspace(moduleName)
		var subspaceParams crosschaintypes.Params
		if isUpgradeBefore {
			subspace = subspace.WithKeyTable(crosschaintypes.ParamKeyTable())
			subspace.GetParamSet(ctx, &subspaceParams)
			response, err := myApp.CrosschainKeeper.Params(ctx, &crosschaintypes.QueryParamsRequest{ChainName: moduleName})
			assert.NoError(t, err)
			params := response.Params
			assert.NotEqualValues(t, params.GravityId, subspaceParams.GravityId)
			assert.NotEqualValues(t, params.AverageBlockTime, subspaceParams.AverageBlockTime)
			assert.NotEqualValues(t, params.AverageExternalBlockTime, subspaceParams.AverageExternalBlockTime)
			assert.NotEqualValues(t, params.ExternalBatchTimeout, subspaceParams.ExternalBatchTimeout)
			assert.NotEqualValues(t, params.SignedWindow, subspaceParams.SignedWindow)
			assert.NotEqualValues(t, params.SlashFraction, subspaceParams.SlashFraction)
			assert.NotEqualValues(t, params.OracleSetUpdatePowerChangePercent, subspaceParams.OracleSetUpdatePowerChangePercent)
			assert.NotEqualValues(t, params.IbcTransferTimeoutHeight, subspaceParams.IbcTransferTimeoutHeight)
			assert.NotEqualValues(t, params.DelegateThreshold, subspaceParams.DelegateThreshold)
			assert.NotEqualValues(t, params.DelegateMultiple, subspaceParams.DelegateMultiple)
			return
		}
		subspace.GetParamSet(ctx, &subspaceParams)
		response, err := myApp.CrosschainKeeper.Params(ctx, &crosschaintypes.QueryParamsRequest{ChainName: moduleName})
		assert.NoError(t, err)
		params := response.Params
		assert.EqualValues(t, params.GravityId, subspaceParams.GravityId)
		assert.EqualValues(t, params.AverageBlockTime, subspaceParams.AverageBlockTime)
		assert.EqualValues(t, params.AverageExternalBlockTime, subspaceParams.AverageExternalBlockTime)
		assert.EqualValues(t, params.ExternalBatchTimeout, subspaceParams.ExternalBatchTimeout)
		assert.EqualValues(t, params.SignedWindow, subspaceParams.SignedWindow)
		assert.EqualValues(t, params.SlashFraction, subspaceParams.SlashFraction)
		assert.EqualValues(t, params.OracleSetUpdatePowerChangePercent, subspaceParams.OracleSetUpdatePowerChangePercent)
		assert.EqualValues(t, params.IbcTransferTimeoutHeight, subspaceParams.IbcTransferTimeoutHeight)
		assert.EqualValues(t, params.DelegateThreshold, subspaceParams.DelegateThreshold)
		assert.EqualValues(t, params.DelegateMultiple, subspaceParams.DelegateMultiple)
	}
	defaultParams := crosschaintypes.DefaultParams()
	for _, newModule := range []string{arbitrumtypes.ModuleName, optimismtypes.ModuleName} {
		response, err := myApp.CrosschainKeeper.Params(ctx, &crosschaintypes.QueryParamsRequest{ChainName: newModule})
		assert.NoError(t, err)
		params := response.Params
		assert.EqualValues(t, params.GravityId, fmt.Sprintf("fx-%s-bridge", newModule))
		assert.EqualValues(t, params.AverageBlockTime, defaultParams.AverageBlockTime)
		assert.EqualValues(t, params.AverageExternalBlockTime, 2000)
		assert.EqualValues(t, params.ExternalBatchTimeout, defaultParams.ExternalBatchTimeout)
		assert.EqualValues(t, params.SignedWindow, defaultParams.SignedWindow)
		assert.EqualValues(t, params.SlashFraction, defaultParams.SlashFraction)
		assert.EqualValues(t, params.OracleSetUpdatePowerChangePercent, defaultParams.OracleSetUpdatePowerChangePercent)
		assert.EqualValues(t, params.IbcTransferTimeoutHeight, defaultParams.IbcTransferTimeoutHeight)
		assert.EqualValues(t, params.DelegateThreshold, defaultParams.DelegateThreshold)
		assert.EqualValues(t, params.DelegateMultiple, defaultParams.DelegateMultiple)
	}
}

func contain[T int | int64 | string](a []T, b T) bool {
	for i := range a {
		if a[i] == b {
			return true
		}
	}
	return false
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
