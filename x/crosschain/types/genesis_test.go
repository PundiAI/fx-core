package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidate(t *testing.T) {
	specs := map[string]struct {
		src    *GenesisState
		expErr bool
	}{
		"default params": {src: &GenesisState{Params: DefaultParams()}, expErr: false},
		"empty params":   {src: &GenesisState{}, expErr: true},
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
