package app_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/pundiai/fx-core/v8/app"
	fxcfg "github.com/pundiai/fx-core/v8/server/config"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/testutil/network"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func Test_ExportGenesisAndRunNode(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	fxtypes.SetConfig(true)
	myApp := helpers.NewApp()

	home := filepath.Join(os.Getenv("HOME"), "tmp/export")
	require.NoError(t, os.RemoveAll(home))
	require.NoError(t, os.MkdirAll(filepath.Join(home, "config"), 0o700))
	require.NoError(t, os.MkdirAll(filepath.Join(home, "data"), 0o700))

	genesisFile := filepath.Join(home, "config", "genesis.json")
	chainId := fxtypes.TestnetChainId
	exportHome := filepath.Join(os.Getenv("HOME"), "tmp")
	genesisDoc := exportGenesisDoc(t, exportHome)
	genesisDoc.ChainID = chainId
	updateGenesisState(t, home, myApp.AppCodec(), genesisDoc)
	require.NoError(t, genesisDoc.SaveAs(genesisFile))

	appCfg := fxcfg.DefaultConfig()
	appCfg.MinGasPrices = fmt.Sprintf("0%s", fxtypes.DefaultDenom)
	appCfg.GRPCWeb.Enable = false
	srvconfig.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
	srvconfig.WriteConfigFile(filepath.Join(home, "config/app.toml"), appCfg)

	clientCtx := client.Context{}.
		WithHomeDir(home).
		WithChainID(chainId).
		WithInterfaceRegistry(myApp.InterfaceRegistry()).
		WithCodec(myApp.AppCodec()).
		WithTxConfig(myApp.GetTxConfig()).
		WithAccountRetriever(authtypes.AccountRetriever{})

	srvCtx := server.NewDefaultContext()
	srvCtx.Logger = log.NewTestLogger(t)
	srvCtx.Config.Moniker = "moniker"
	srvCtx.Config.DBBackend = string(dbm.MemDBBackend)
	srvCtx.Config.Consensus = tmcfg.TestConsensusConfig()
	srvCtx.Config.RPC.PprofListenAddress = ""
	srvCtx.Config.Instrumentation.Prometheus = false
	srvCtx.Config.SetRoot(home)
	tmcfg.WriteConfigFile(filepath.Join(home, "config/config.toml"), srvCtx.Config)

	srvCtx.Viper.SetConfigFile(filepath.Join(home, "config/config.toml"))
	srvCtx.Viper.SetConfigFile(filepath.Join(home, "config/app.toml"))
	require.NoError(t, srvCtx.Viper.ReadInConfig())
	srvCtx.Viper.Set(crisis.FlagSkipGenesisInvariants, true)

	val := network.Validator{
		AppConfig: appCfg,
		ClientCtx: clientCtx,
		Ctx:       srvCtx,
	}
	ctx := context.Background()
	group, ctx := errgroup.WithContext(ctx)
	myAppConstructor := func(appConfig *fxcfg.Config, ctx *server.Context) servertypes.Application {
		return helpers.NewApp(func(opts *helpers.AppOpts) {
			opts.Home = ctx.Config.RootDir
		})
	}
	require.NoError(t, network.StartInProcess(ctx, group, myAppConstructor, &val))
	select {}
}

func exportGenesisDoc(t *testing.T, home string) *types.GenesisDoc {
	t.Helper()
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(home, "data"))
	require.NoError(t, err)

	myApp := helpers.NewApp(func(opts *helpers.AppOpts) {
		opts.Logger = log.NewTestLogger(t)
		opts.Home = home
		opts.DB = db
	})
	exportedApp, err := myApp.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err)
	genesisDoc := &types.GenesisDoc{
		GenesisTime:     time.Now(),
		ConsensusParams: app.CustomGenesisConsensusParams(),
		AppState:        exportedApp.AppState,
	}
	return genesisDoc
}

func updateGenesisState(t *testing.T, home string, cdc codec.Codec, genesisDoc *types.GenesisDoc) {
	t.Helper()
	appState := app.GenesisState{}
	err := json.Unmarshal(genesisDoc.AppState, &appState)
	require.NoError(t, err)

	newPubKey := newPrivValidatorKey(t, home)
	validator := updateStakingGenesisState(cdc, appState, newPubKey)
	updateBankGenesisState(cdc, appState, validator)
	updateSlashingGenesisState(cdc, appState, newPubKey)

	genesisDoc.AppState, err = json.Marshal(appState)
	require.NoError(t, err)
}

func updateStakingGenesisState(cdc codec.Codec, appState app.GenesisState, newPubKey *codectypes.Any) stakingtypes.Validator {
	stakingGenesisState := new(stakingtypes.GenesisState)
	cdc.MustUnmarshalJSON(appState[stakingtypes.ModuleName], stakingGenesisState)
	sort.Slice(stakingGenesisState.Validators, func(i, j int) bool {
		return stakingGenesisState.Validators[i].ConsensusPower(sdk.DefaultPowerReduction) > stakingGenesisState.Validators[j].ConsensusPower(sdk.DefaultPowerReduction)
	})

	validator := stakingGenesisState.Validators[0]
	validator.ConsensusPubkey = newPubKey

	for i := 1; i < len(stakingGenesisState.Validators); i++ {
		if stakingGenesisState.Validators[i].Status == stakingtypes.Bonded {
			stakingGenesisState.Validators[i].Status = stakingtypes.Unbonded
			stakingGenesisState.Validators[i].Jailed = true
		}
	}
	for i := 0; i < len(stakingGenesisState.LastValidatorPowers); i++ {
		if stakingGenesisState.LastValidatorPowers[i].Address == validator.OperatorAddress {
			stakingGenesisState.LastTotalPower = sdkmath.NewInt(stakingGenesisState.LastValidatorPowers[i].Power)
			stakingGenesisState.LastValidatorPowers = []stakingtypes.LastValidatorPower{
				stakingGenesisState.LastValidatorPowers[i],
			}
		}
	}
	stakingGenesisState.Validators[0] = validator
	appState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingGenesisState)
	return validator
}

func updateBankGenesisState(codec codec.Codec, appState app.GenesisState, validator stakingtypes.Validator) {
	bankGenesisState := new(banktypes.GenesisState)
	codec.MustUnmarshalJSON(appState[banktypes.ModuleName], bankGenesisState)
	bondedPoolAddr := authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String()
	notBondedPoolAddr := authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String()
	var notBoundedAmount sdk.Coins
	for i := 0; i < len(bankGenesisState.Balances); i++ {
		if bankGenesisState.Balances[i].Address == bondedPoolAddr {
			notBoundedAmount = bankGenesisState.Balances[i].Coins.Sub(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, validator.BondedTokens()))...)
			bankGenesisState.Balances[i].Coins = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, validator.BondedTokens()))
		}
		if bankGenesisState.Balances[i].Address == notBondedPoolAddr {
			bankGenesisState.Balances[i].Coins = bankGenesisState.Balances[i].Coins.Add(notBoundedAmount...)
		}
	}
	appState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesisState)
}

func updateSlashingGenesisState(cdc codec.Codec, appState app.GenesisState, newPubKey *codectypes.Any) {
	pubKey := newPubKey.GetCachedValue().(cryptotypes.PubKey)
	slashingGenesisState := new(slashingtypes.GenesisState)
	cdc.MustUnmarshalJSON(appState[slashingtypes.ModuleName], slashingGenesisState)
	slashingGenesisState.SigningInfos = append(slashingGenesisState.SigningInfos, slashingtypes.SigningInfo{
		Address: sdk.ConsAddress(pubKey.Address()).String(),
		ValidatorSigningInfo: slashingtypes.ValidatorSigningInfo{
			Address:             sdk.ConsAddress(pubKey.Address()).String(),
			StartHeight:         0,
			IndexOffset:         0,
			JailedUntil:         time.Now(),
			Tombstoned:          false,
			MissedBlocksCounter: 0,
		},
	})
	appState[slashingtypes.ModuleName] = cdc.MustMarshalJSON(slashingGenesisState)
}

func newPrivValidatorKey(t *testing.T, home string) *codectypes.Any {
	t.Helper()
	privKeyFile := filepath.Join(home, "config/priv_validator_key.json")
	privStateFile := filepath.Join(home, "data/priv_validator_state.json")
	secret := tmrand.Bytes(32)
	filePV := privval.NewFilePV(ed25519.GenPrivKeyFromSecret(secret), privKeyFile, privStateFile)
	filePV.Save()

	pubkey, err := cryptocodec.FromCmtPubKeyInterface(filePV.Key.PubKey)
	require.NoError(t, err)
	pubAny, err := codectypes.NewAnyWithValue(pubkey)
	require.NoError(t, err)
	return pubAny
}
