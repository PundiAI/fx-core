package store_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v7/app"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func Benchmark_Subspace(b *testing.B) {
	storeKey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")

	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(tkey, storetypes.StoreTypeIAVL, nil)
	assert.NoError(b, ms.LoadLatestVersion())

	encodingConfig := app.MakeEncodingConfig()
	subspace := types.NewSubspace(encodingConfig.Codec, encodingConfig.Amino, storeKey, tkey, "sub")
	subspace.WithKeyTable(types.NewKeyTable().RegisterParamSet(&erc20types.Params{}))

	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())
	params := erc20types.DefaultParams()
	subspace.SetParamSet(ctx, &params)

	b.Run("A", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var enableErc20 bool
			subspace.Get(ctx, erc20types.ParamStoreKeyEnableErc20, &enableErc20)
			assert.Equal(b, true, enableErc20)
		}
	})

	b.Run("B", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newParams := new(erc20types.Params)
			subspace.GetParamSet(ctx, newParams)
			assert.Equal(b, true, params.EnableErc20)
		}
	})
}
