package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/x/erc20/types"
)

func TestValidateGenesis(t *testing.T) {
	newGen := types.NewGenesisState(types.DefaultParams(), []types.TokenPair{})

	testCases := []struct {
		name     string
		genState *types.GenesisState
		expPass  bool
	}{
		{
			name:     "valid genesis constructor",
			genState: &newGen,
			expPass:  true,
		},
		{
			name:     "default",
			genState: types.DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "valid genesis",
			genState: &types.GenesisState{
				Params:     types.DefaultParams(),
				TokenPairs: []types.TokenPair{},
			},
			expPass: true,
		},
		{
			name: "valid genesis - with tokens pairs",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenPairs: []types.TokenPair{
					{
						Erc20Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
						Denom:        "usdt",
						Enabled:      true,
					},
				},
			},
			expPass: true,
		},
		{
			name: "invalid genesis - duplicated token pair",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenPairs: []types.TokenPair{
					{
						Erc20Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
						Denom:        "usdt",
						Enabled:      true,
					},
					{
						Erc20Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
						Denom:        "usdt",
						Enabled:      true,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - duplicated token pair",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenPairs: []types.TokenPair{
					{
						Erc20Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
						Denom:        "usdt",
						Enabled:      true,
					},
					{
						Erc20Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
						Denom:        "usdt2",
						Enabled:      true,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - duplicated token pair",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenPairs: []types.TokenPair{
					{
						Erc20Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
						Denom:        "usdt",
						Enabled:      true,
					},
					{
						Erc20Address: "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
						Denom:        "usdt",
						Enabled:      true,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - invalid token pair",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenPairs: []types.TokenPair{
					{
						Erc20Address: "0xinvalidaddress",
						Denom:        "bad",
						Enabled:      true,
					},
				},
			},
			expPass: false,
		},
		{
			// Voting period cant be zero
			name:     "empty genesis",
			genState: &types.GenesisState{Params: types.Params{IbcTimeout: 1}},
			expPass:  true,
		},
	}

	for _, tc := range testCases {
		err := tc.genState.Validate()
		if tc.expPass {
			assert.NoError(t, err, tc.name)
		} else {
			assert.Error(t, err, tc.name)
		}
	}
}
