package app_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/testutil/network"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func TestUseExportGenesisDataRunNode(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test: ", t.Name())

	// set sdk.Config and get network config
	networkConfig := testutil.DefaultNetworkConfig(app.MakeEncodingConfig())

	genesisDoc := GetGenesisDocFromAppData(t)
	// genesisFile := filepath.Join(fxtypes.GetDefaultNodeHome(), "config", "genesis.json")
	// assert.NoError(t, genesisDoc.SaveAs(genesisFile))
	appState := app.GenesisState{}
	assert.NoError(t, tmjson.Unmarshal(genesisDoc.AppState, &appState))

	networkConfig.TimeoutCommit = time.Millisecond
	networkConfig.NumValidators = 1
	networkConfig.EnableTMLogging = true
	networkConfig.GenesisState = appState

	myNetwork, err := network.New(t, t.TempDir(), networkConfig)
	assert.NoError(t, err)
	assert.Equal(t, myNetwork.Validators, 1)

	_, err = myNetwork.WaitForHeight(10)
	assert.NoError(t, err)
}

func GetGenesisDocFromAppData(t *testing.T) *types.GenesisDoc {
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(fxtypes.GetDefaultNodeHome(), "data"))
	require.NoError(t, err)

	appEncodingCfg := app.MakeEncodingConfig()
	// logger := log.NewNopLogger()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	myApp := app.New(logger, db,
		nil, true, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		appEncodingCfg, app.EmptyAppOptions{},
	)
	exportedApp, err := myApp.ExportAppStateAndValidators(false, []string{})
	assert.NoError(t, err)
	genesisDoc := &types.GenesisDoc{
		GenesisTime:     time.Now(),
		ChainID:         fxtypes.ChainId(),
		ConsensusParams: app.CustomGenesisConsensusParams(),
		Validators:      exportedApp.Validators,
		AppState:        exportedApp.AppState,
	}

	appState := app.GenesisState{}
	assert.NoError(t, tmjson.Unmarshal(genesisDoc.AppState, &appState))

	keyJSONBytes, err := os.ReadFile(filepath.Join(fxtypes.GetDefaultNodeHome(), "config", "priv_validator_key.json"))
	assert.NoError(t, err)

	pvKey := privval.FilePVKey{}
	err = tmjson.Unmarshal(keyJSONBytes, &pvKey)
	assert.NoError(t, err)

	encodingConfig := app.MakeEncodingConfig()
	cdc := encodingConfig.Codec

	// stakingGenesisState
	stakingGenesisState := new(stakingtypes.GenesisState)
	cdc.MustUnmarshalJSON(appState[stakingtypes.ModuleName], stakingGenesisState)
	sort.Slice(stakingGenesisState.Validators, func(i, j int) bool {
		return stakingGenesisState.Validators[i].ConsensusPower(sdk.DefaultPowerReduction) > stakingGenesisState.Validators[j].ConsensusPower(sdk.DefaultPowerReduction)
	})

	pubkey, err := cryptocodec.FromTmPubKeyInterface(pvKey.PubKey)
	assert.NoError(t, err)

	pubAny, err := codectypes.NewAnyWithValue(pubkey)
	assert.NoError(t, err)

	for i := 0; i < len(stakingGenesisState.Validators); i++ {
		if i == 0 {
			stakingGenesisState.Validators[0].ConsensusPubkey = pubAny
			stakingGenesisState.Validators[0].Tokens = stakingGenesisState.Validators[0].Tokens.Add(sdkmath.NewInt(190000).MulRaw(1e18))
			continue
		}
		if stakingGenesisState.Validators[i].Status == stakingtypes.Bonded {
			stakingGenesisState.Validators[i].Status = stakingtypes.Unbonded
			stakingGenesisState.Validators[i].Jailed = true
			stakingGenesisState.Validators[0].Tokens = stakingGenesisState.Validators[0].Tokens.Add(stakingGenesisState.Validators[i].Tokens)
			_, delegatorShares := stakingGenesisState.Validators[0].AddTokensFromDel(stakingGenesisState.Validators[i].Tokens)
			stakingGenesisState.Validators[0].DelegatorShares = delegatorShares
			stakingGenesisState.Validators[i].Tokens = sdkmath.ZeroInt()
			stakingGenesisState.Validators[i].DelegatorShares = sdk.ZeroDec()
		}
	}

	for i := 0; i < len(stakingGenesisState.LastValidatorPowers); i++ {
		if stakingGenesisState.LastValidatorPowers[i].Address == stakingGenesisState.Validators[0].OperatorAddress {
			stakingGenesisState.LastValidatorPowers[i].Power = stakingGenesisState.Validators[0].GetConsensusPower(sdk.DefaultPowerReduction)
			stakingGenesisState.LastValidatorPowers = []stakingtypes.LastValidatorPower{
				stakingGenesisState.LastValidatorPowers[i],
			}
		}
	}
	appState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingGenesisState)

	// genesisDoc.Validators
	validatorConsAddress := types.Address{}
	for i := 0; i < len(genesisDoc.Validators); i++ {
		if genesisDoc.Validators[i].Name == stakingGenesisState.Validators[0].Description.Moniker {
			validatorConsAddress = genesisDoc.Validators[i].Address
			genesisDoc.Validators[i].PubKey = pvKey.PubKey
			genesisDoc.Validators[i].Address = pvKey.Address
			genesisDoc.Validators[i].Power = stakingGenesisState.Validators[0].GetConsensusPower(sdk.DefaultPowerReduction)
			genesisDoc.Validators = []types.GenesisValidator{genesisDoc.Validators[i]}
			break
		}
	}

	// slashingGenesisState
	slashingGenesisState := new(slashingtypes.GenesisState)
	cdc.MustUnmarshalJSON(appState[slashingtypes.ModuleName], slashingGenesisState)

	for i := 0; i < len(slashingGenesisState.SigningInfos); i++ {
		if slashingGenesisState.SigningInfos[i].Address == sdk.ConsAddress(validatorConsAddress).String() {
			slashingGenesisState.SigningInfos[i].Address = sdk.ConsAddress(pvKey.Address.Bytes()).String()
			slashingGenesisState.SigningInfos[i].ValidatorSigningInfo.Address = sdk.ConsAddress(pvKey.Address.Bytes()).String()
			slashingGenesisState.SigningInfos = []slashingtypes.SigningInfo{slashingGenesisState.SigningInfos[i]}
			break
		}
	}
	appState[slashingtypes.ModuleName] = cdc.MustMarshalJSON(slashingGenesisState)

	genesisDoc.AppState, err = tmjson.Marshal(appState)
	assert.NoError(t, err)
	return genesisDoc
}
