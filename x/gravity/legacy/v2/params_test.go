package v2_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v3/app"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	v2 "github.com/functionx/fx-core/v3/x/gravity/legacy/v2"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

func TestMigrateParams(t *testing.T) {
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.ModuleName)

	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(paramsStoreKey, sdk.StoreTypeIAVL, nil)
	assert.NoError(t, ms.LoadLatestVersion())

	amino := app.MakeEncodingConfig().Amino
	paramsStore := ms.GetKVStore(paramsStoreKey)
	oldStore := prefix.NewStore(paramsStore, append([]byte(types.ModuleName), '/'))
	gravityParams := v2.TestParams()
	for _, pair := range gravityParams.ParamSetPairs() {
		bz, err := amino.MarshalJSON(pair.Value)
		assert.NoError(t, err)
		oldStore.Set(pair.Key, bz)
	}

	err := v2.MigrateParams(amino, paramsStore, ethtypes.ModuleName)
	assert.NoError(t, err)

	newStore := prefix.NewStore(paramsStore, append([]byte(ethtypes.ModuleName), '/'))
	ethParams := &crosschaintypes.Params{}
	for _, pair := range ethParams.ParamSetPairs() {
		bz := newStore.Get(pair.Key)
		if len(bz) <= 0 {
			continue
		}
		if err := amino.UnmarshalJSON(bz, pair.Value); err != nil {
			panic(err)
		}
	}
	assert.EqualValues(t, &crosschaintypes.Params{
		GravityId:                         gravityParams.GravityId,
		AverageBlockTime:                  gravityParams.AverageBlockTime,
		ExternalBatchTimeout:              gravityParams.TargetBatchTimeout,
		AverageExternalBlockTime:          gravityParams.AverageEthBlockTime,
		SignedWindow:                      30_000,
		SlashFraction:                     sdk.MustNewDecFromStr("0.8"),
		OracleSetUpdatePowerChangePercent: gravityParams.ValsetUpdatePowerChangePercent,
		IbcTransferTimeoutHeight:          gravityParams.IbcTransferTimeoutHeight,
		Oracles:                           nil,
		DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)),
		DelegateMultiple:                  10,
	}, ethParams)
}
