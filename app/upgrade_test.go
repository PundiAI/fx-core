package app_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	v2 "github.com/functionx/fx-core/v3/app/upgrades/v2"
	v3 "github.com/functionx/fx-core/v3/app/upgrades/v3"
	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

func Test_Upgrade(t *testing.T) {
	if !helpers.IsLocalTest() {
		t.Skip("skipping local test", t.Name())
	}
	fxtypes.SetConfig(true)

	testCases := []struct {
		name                  string
		fromVersion           int
		toVersion             int
		LocalStoreBlockHeight uint64
		plan                  upgradetypes.Plan
	}{
		{
			name:        "upgrade v3",
			fromVersion: 2,
			toVersion:   3,
			plan: upgradetypes.Plan{
				Name: v3.Upgrade.UpgradeName,
				Info: "local test upgrade v3",
			},
			LocalStoreBlockHeight: 7654832,
		},
	}

	db, err := sdk.NewLevelDB("application", filepath.Join(fxtypes.GetDefaultNodeHome(), "data"))
	require.NoError(t, err)

	appEncodingCfg := app.MakeEncodingConfig()
	// logger := log.NewNopLogger()
	logger := log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowInfo())
	myApp := app.New(logger, db,
		nil, true, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		appEncodingCfg, app.EmptyAppOptions{},
	)
	ctx := myApp.NewUncachedContext(false, tmproto.Header{
		ChainID:         fxtypes.ChainId(),
		Height:          myApp.LastBlockHeight(),
		ProposerAddress: tmrand.Bytes(20),
	})
	validators := myApp.StakingKeeper.GetAllValidators(ctx)
	assert.True(t, len(validators) > 0)
	var pubkey cryptotypes.PubKey
	assert.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubkey))
	ctx = ctx.WithProposer(pubkey.Address().Bytes())

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.LocalStoreBlockHeight > 0 {
				require.Equal(t, ctx.BlockHeight(), int64(testCase.LocalStoreBlockHeight))
				for moduleName, keys := range v2.GetModuleKey() {
					kvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(moduleName))
					checkStoreKey(t, moduleName, keys, kvStore)
				}
			}
			checkVersionMap(t, ctx, myApp, getConsensusVersion(testCase.fromVersion))

			testCase.plan.Height = ctx.BlockHeight()
			myApp.UpgradeKeeper.ApplyUpgrade(ctx, testCase.plan)

			// reload store
			// err = upgradetypes.UpgradeStoreLoader(plan.Height, v3.Upgrade.StoreUpgrades())(myApp.CommitMultiStore())
			// require.NoError(t, err)

			checkVersionMap(t, ctx, myApp, getConsensusVersion(testCase.toVersion))

			checkDataAfterMigrateV3(t, ctx, myApp)

			if testCase.LocalStoreBlockHeight > 0 {
				for moduleName, keys := range v3.GetModuleKey() {
					kvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(moduleName))
					checkStoreKey(t, moduleName, keys, kvStore)
				}
			}
		})
	}
	exportedApp, err := myApp.ExportAppStateAndValidators(false, []string{})
	assert.NoError(t, err)
	genesisState := app.GenesisState{}
	assert.NoError(t, tmjson.Unmarshal(exportedApp.AppState, &genesisState))
	for moduleName, appState := range genesisState {
		for key := range v3.GetModuleKey() {
			if key == moduleName {
				t.Log("-------------------", moduleName)
				t.Log(string(appState))
			}
		}
	}
	conParamsBytes, err := tmjson.MarshalIndent(exportedApp.ConsensusParams, "", "  ")
	assert.NoError(t, err)
	t.Log(string(conParamsBytes))

	// myApp.CommitMultiStore().Commit()
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
		"auth":         {1, 2},
		"authz":        {0, 1},
		"avalanche":    {0, 0, 1},
		"bank":         {1, 2},
		"bsc":          {1, 2, 3},
		"capability":   {1},
		"crisis":       {1},
		"crosschain":   {1},
		"distribution": {1, 2},
		"erc20":        {0, 1, 2},
		"evidence":     {1},
		"evm":          {0, 2, 3},
		"eth":          {0, 0, 1},
		"feegrant":     {0, 1},
		"feemarket":    {0, 3},
		"genutil":      {1},
		"gov":          {1, 2},
		"gravity":      {1, 1, 2},
		"ibc":          {1, 2},
		"migrate":      {0, 1},
		"mint":         {1},
		"other":        {1},
		"params":       {1},
		"polygon":      {1, 2, 3},
		"slashing":     {1, 2},
		"staking":      {1, 2},
		"transfer":     {1, 1, 2}, // ibc-transfer
		"fxtransfer":   {0, 0, 1}, // fx-ibc-transfer
		"tron":         {1, 2, 3},
		"upgrade":      {1},
		"vesting":      {1},
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

func checkStoreKey(t *testing.T, name string, keys map[byte][2]int, kvStores sdk.KVStore) {
	iterator := kvStores.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		x, ok := keys[iterator.Key()[0]]
		assert.True(t, ok, fmt.Sprintf("%x", iterator.Key()[0]), iterator.Value())
		if ok {
			if x[0] == -1 && x[1] == -1 {
				// ignore
				continue
			}
			// set result
			keys[iterator.Key()[0]] = [2]int{x[0], x[1] + 1}
		}
	}
	iterator.Close()
	for k, x := range keys {
		assert.Equal(t, x[0], x[1], fmt.Sprintf("%s: %x", name, k))
	}
}

// check v3 upgrade data

func checkDataAfterMigrateV3(t *testing.T, ctx sdk.Context, myApp *app.App) {
	checkWFXLogicUpgrade(t, ctx, myApp)
	checkMetadataAliasNull(t, ctx, myApp)
	checkV3RegisterCoin(t, ctx, myApp)
	checkV3IBCTransferRelation(t, ctx, myApp)
	checkV3NewEvmParams(t, ctx, myApp)
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

func checkMetadataAliasNull(t *testing.T, ctx sdk.Context, myApp *app.App) {
	// check metadata alias null string
	myApp.BankKeeper.IterateAllDenomMetaData(ctx, func(md banktypes.Metadata) bool {
		if len(md.DenomUnits) != 2 {
			return false
		}
		if len(md.DenomUnits[1].Aliases) == 0 || len(md.DenomUnits[1].Aliases) > 1 {
			return false
		}
		for _, alias := range md.DenomUnits[1].Aliases {
			require.NotEqual(t, "null", alias)
		}
		return false
	})
}

func checkV3RegisterCoin(t *testing.T, ctx sdk.Context, myApp *app.App) {
	// check register coin
	mds := v3.GetMetadata(ctx.ChainID())
	codeHash := common.HexToHash("0xf8572bdecc4c287eec1a748169288993aaf0feed1f988dbefe28deb9321ee970")
	for _, md := range mds {
		tokenPair, found := myApp.Erc20Keeper.GetTokenPair(ctx, md.Base)
		require.True(t, found)
		require.Equal(t, tokenPair.GetDenom(), md.Base)
		require.True(t, tokenPair.Enabled)
		require.Equal(t, erc20types.OWNER_MODULE, tokenPair.ContractOwner)

		tokenPairAddress := tokenPair.GetERC20Contract()
		tokenPairAcc := myApp.EvmKeeper.GetAccount(ctx, tokenPairAddress)
		tokenPairCode := myApp.EvmKeeper.GetCode(ctx, common.BytesToHash(tokenPairAcc.CodeHash))
		require.Equal(t, codeHash, crypto.Keccak256Hash(tokenPairCode))
	}
}

func checkV3IBCTransferRelation(t *testing.T, ctx sdk.Context, myApp *app.App) {
	kvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(erc20types.StoreKey))

	iter := sdk.KVStorePrefixIterator(kvStore, erc20types.KeyPrefixIBCTransfer)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyStr := string(bytes.TrimPrefix(iter.Key(), erc20types.KeyPrefixIBCTransfer))

		channel, sequence, ok := parseIBCTransferKey(keyStr)
		require.True(t, ok)

		found := myApp.IBCKeeper.ChannelKeeper.HasPacketCommitment(ctx, ibctransfertypes.ModuleName, channel, sequence)
		require.True(t, found)
	}
}

func checkV3NewEvmParams(t *testing.T, ctx sdk.Context, myApp *app.App) {
	params := myApp.EvmKeeper.GetParams(ctx)
	defaultParams := evmtypes.DefaultParams()
	defaultParams.EvmDenom = fxtypes.DefaultDenom
	require.Equal(t, defaultParams.String(), params.String())
}

func parseIBCTransferKey(keyStr string) (string, uint64, bool) {
	split := strings.Split(keyStr, "/")
	if len(split) != 2 {
		return "", 0, false
	}

	channel := split[0]
	sequence, err := strconv.ParseUint(split[1], 10, 64)
	if err != nil {
		return "", 0, false
	}
	return channel, sequence, true
}
