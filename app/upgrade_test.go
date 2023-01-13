package app_test

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	v2 "github.com/functionx/fx-core/v3/app/upgrades/v2"
	v3 "github.com/functionx/fx-core/v3/app/upgrades/v3"
	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	gravityv2 "github.com/functionx/fx-core/v3/x/gravity/legacy/v2"
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
	myApp := app.New(
		log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowInfo()),
		db, nil, true, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		appEncodingCfg, app.EmptyAppOptions{},
	)
	ctx := newContext(t, myApp)

	checkStakingPool(t, ctx, myApp, true)

	checkTotalSupply(t, ctx, myApp)

	initEthOracleBalances(t, ctx, myApp)

	var totalSupplies sdk.Coins
	myApp.BankKeeper.IterateTotalSupply(ctx, func(coin sdk.Coin) bool {
		totalSupplies = append(totalSupplies, coin)
		return false
	})

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

	myApp.BankKeeper.IterateTotalSupply(ctx, func(coin sdk.Coin) bool {
		assert.Equal(t, totalSupplies.AmountOf(coin.Denom).String(), coin.Amount.String())
		return false
	})

	checkStakingPool(t, ctx, myApp, false)

	checkTotalSupply(t, ctx, myApp)

	myApp.EthKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))

	exportAppState(t, myApp)

	// optional: save to the database
	// myApp.CommitMultiStore().Commit()
}

func initEthOracleBalances(t *testing.T, ctx sdk.Context, myApp *app.App) {
	oracles := gravityv2.EthInitOracles(ctx.ChainID())
	for _, oracle := range oracles {
		addr := sdk.MustAccAddressFromBech32(oracle)
		err := myApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, distributiontypes.ModuleName, addr,
			sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18))))
		assert.NoError(t, err)
	}
}

func newContext(t *testing.T, myApp *app.App) sdk.Context {
	chainID := fxtypes.MainnetChainId
	if os.Getenv("CHAIN_ID") == fxtypes.TestnetChainId {
		chainID = fxtypes.TestnetChainId
	}
	ctx := myApp.NewUncachedContext(false, tmproto.Header{
		ChainID: chainID,
		Height:  myApp.LastBlockHeight(),
	})
	// set the first validator to proposer
	validators := myApp.StakingKeeper.GetAllValidators(ctx)
	assert.True(t, len(validators) > 0)
	var pubkey cryptotypes.PubKey
	assert.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubkey))
	ctx = ctx.WithProposer(pubkey.Address().Bytes())
	return ctx
}

func checkStakingPool(t *testing.T, ctx sdk.Context, myApp *app.App, isUpgradeBefore bool) {
	validators := myApp.StakingKeeper.GetAllValidators(ctx)
	totalBonded := sdk.ZeroInt()
	totalNotBounded := sdk.ZeroInt()
	for _, validator := range validators {
		if validator.IsBonded() {
			totalBonded = totalBonded.Add(validator.Tokens)
		} else {
			totalNotBounded = totalNotBounded.Add(validator.Tokens)
		}
	}

	undelegateAmount := sdk.ZeroInt()
	myApp.StakingKeeper.IterateUnbondingDelegations(ctx, func(index int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for _, entry := range ubd.Entries {
			undelegateAmount = undelegateAmount.Add(entry.Balance)
		}
		return false
	})

	bondDenom := myApp.StakingKeeper.BondDenom(ctx)
	bondedPool := myApp.StakingKeeper.GetBondedPool(ctx)
	notBondedPool := myApp.StakingKeeper.GetNotBondedPool(ctx)
	bondedPoolAmount := myApp.BankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount
	notBondedPoolAmount := myApp.BankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount

	if isUpgradeBefore {
		if ctx.ChainID() == fxtypes.TestnetChainId {
			totalBonded = totalBonded.Add(sdk.NewInt(90_000).MulRaw(1e18))
		} else {
			totalBonded = totalBonded.Add(sdk.NewInt(190_000).MulRaw(1e18))
		}
		assert.Equal(t, bondedPoolAmount.String(), totalBonded.String())
		assert.Equal(t, notBondedPoolAmount.String(), totalNotBounded.Add(undelegateAmount).String())
	} else {
		assert.Equal(t, bondedPoolAmount.String(), totalBonded.String())
		assert.Equal(t, notBondedPoolAmount.String(), totalNotBounded.Add(undelegateAmount).String())
	}
}

func checkTotalSupply(t *testing.T, ctx sdk.Context, myApp *app.App) {
	// chain total supply
	totalSupply := sdk.NewCoins()
	myApp.BankKeeper.IterateTotalSupply(ctx, func(coin sdk.Coin) bool {
		totalSupply = totalSupply.Add(coin)
		return false
	})

	// all lock token
	erc20Balances := myApp.BankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(erc20types.ModuleName))

	// all contract
	allTokenPairs := myApp.Erc20Keeper.GetAllTokenPairs(ctx)

	// NOTE: testnet upgrade evm, erc20 twice, history balance not remove
	erc20V1TestnetBalanceStr := []string{
		"1000000000000000100000000eth0x2870405E4ABF9FcCDc93d9cC83c09788296d8354", "1100773800000eth0xD69133f9A0206b3340d9622F2eBc4571022b3b5f",
		"1000293000000000000000000eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B", "1000100000000eth0xeC822cd1238d946Cf0f73be57359c5cAa5512a9D",
		"1002123400000000000000000ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4", "166500000000000000000polygon0x326C977E6efc84E512bB9C30f76E30c160eD06FB",
		"1001000000000tronTK1pM7NtkLohgRgKA6LeocW2znwJ8JtLrQ", "38100000000000000000tronTLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL", "145000000tronTXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj",
	}
	erc20V1TestnetBalances, err := sdk.ParseCoinsNormalized(strings.Join(erc20V1TestnetBalanceStr, ","))
	require.NoError(t, err)

	// contract totalSupply
	for _, pair := range allTokenPairs {
		res, err := myApp.Erc20Keeper.CallEVM(ctx, fxtypes.GetERC20().ABI, myApp.Erc20Keeper.ModuleAddress(),
			pair.GetERC20Contract(), false, "totalSupply")
		assert.NoError(t, err)

		var totalSupplyRes struct{ Value *big.Int }
		err = fxtypes.GetERC20().ABI.UnpackIntoInterface(&totalSupplyRes, "totalSupply", res.Ret)
		assert.NoError(t, err)

		denomBalance := erc20Balances.AmountOf(pair.GetDenom())
		if pair.GetDenom() == fxtypes.DefaultDenom {
			assert.True(t, 0 == denomBalance.Uint64())
			denomBalance = myApp.BankKeeper.GetBalance(ctx, pair.GetERC20Contract().Bytes(), pair.GetDenom()).Amount
		}
		if ctx.ChainID() == fxtypes.TestnetChainId {
			v1Balance := erc20V1TestnetBalances.AmountOf(pair.GetDenom())
			denomBalance = denomBalance.Sub(v1Balance)
		}
		assert.Equal(t, totalSupplyRes.Value.String(), denomBalance.BigInt().String(), pair.GetDenom())
	}

	// usdt totalSupply
	usdtTotalSupply := totalSupply.AmountOf("usdt")
	usdtMD, found := myApp.BankKeeper.GetDenomMetaData(ctx, "usdt")
	assert.True(t, found)
	chainUSDTTotalSupply := sdk.ZeroInt()
	for _, alias := range usdtMD.DenomUnits[0].Aliases {
		balance := myApp.BankKeeper.GetBalance(ctx, myApp.Erc20Keeper.ModuleAddress().Bytes(), alias)
		chainUSDTTotalSupply = chainUSDTTotalSupply.Add(balance.Amount)
	}
	assert.Equal(t, usdtTotalSupply.String(), chainUSDTTotalSupply.String(), "usdt")
}

func exportAppState(t *testing.T, myApp *app.App) {
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
