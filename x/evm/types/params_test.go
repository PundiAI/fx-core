package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
)

func TestParamKeyTable(t *testing.T) {
	require.IsType(t, paramtypes.KeyTable{}, ParamKeyTable())
}

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{"default", DefaultParams(), false},
		{
			"valid",
			NewParams(true, true, 2929, 1884, 1344),
			false,
		},
		{
			"empty",
			Params{},
			false,
		},
		{
			"invalid eip",
			Params{
				ExtraEIPs: []int64{1},
			},
			true,
		},
		{
			"invalid chain config",
			NewParams(true, true, 2929, 1884, 1344),
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestParamsEIPs(t *testing.T) {
	actual := NewParams(true, true, 2929, 1884, 1344).EIPs()
	require.Equal(t, []int([]int{2929, 1884, 1344}), actual)
}

func TestParamsValidatePriv(t *testing.T) {
	require.Error(t, validateBool(""))
	require.NoError(t, validateBool(true))
	require.Error(t, validateEIPs(""))
	require.NoError(t, validateEIPs([]int64{1884}))
}

func TestIsLondon(t *testing.T) {
	testCases := []struct {
		name   string
		height int64
		result bool
	}{
		{
			"Before london block",
			5,
			false,
		},
		{
			"After london block",
			12_965_001,
			true,
		},
		{
			"london block",
			12_965_000,
			true,
		},
	}

	for _, tc := range testCases {
		ethCfg := params.MainnetChainConfig
		require.Equal(t, ethCfg.IsLondon(big.NewInt(tc.height)), tc.result)
	}
}
