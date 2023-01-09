package v2_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v3/app"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	v2 "github.com/functionx/fx-core/v3/x/gravity/legacy/v2"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

func TestMigrateParams(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
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

	err := v2.MigrateParams(amino, paramsStore, ethtypes.StoreKey)
	assert.NoError(t, err)
}
