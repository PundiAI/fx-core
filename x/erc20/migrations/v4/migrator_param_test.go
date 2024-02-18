package v4_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	v4 "github.com/functionx/fx-core/v7/x/erc20/migrations/v4"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

type mockSubspace struct {
	ps types.Params
}

func newMockSubspace(ps types.Params) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(_ sdk.Context, ps paramtypes.ParamSet) {
	*ps.(*types.Params) = ms.ps
}

func (ms mockSubspace) HasKeyTable() bool {
	return false
}

func (ms mockSubspace) WithKeyTable(_ paramtypes.KeyTable) paramtypes.Subspace {
	return paramtypes.Subspace{}
}

func TestStoreMigration(t *testing.T) {
	encCfg := simapp.MakeTestEncodingConfig()
	erc20Key := sdk.NewKVStoreKey(types.ModuleName)
	tBscKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(erc20Key, tBscKey)
	store := ctx.KVStore(erc20Key)

	legacySubspace := newMockSubspace(types.DefaultParams())

	testCases := []struct {
		name        string
		doMigration bool
	}{
		{
			name:        "without state migration",
			doMigration: false,
		},
		{
			name:        "with state migration",
			doMigration: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.doMigration {
				require.NoError(t, v4.MigratorParam(ctx, legacySubspace, erc20Key, encCfg.Codec))
			}
			if tc.doMigration {
				var res types.Params
				bz := store.Get(types.ParamsKey)
				require.NoError(t, encCfg.Codec.Unmarshal(bz, &res))
				require.Equal(t, legacySubspace.ps, res)
			} else {
				require.Equal(t, true, true)
			}
		})
	}
}
