package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func TestGenesisStateValidate(t *testing.T) {
	specs := map[string]struct {
		src    *types.GenesisState
		expErr bool
	}{
		"default params": {
			src: &types.GenesisState{
				Params: types.DefaultParams(),
			},
			expErr: false,
		},
		"empty params": {
			src:    &types.GenesisState{},
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
