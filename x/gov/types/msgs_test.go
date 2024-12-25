package types_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/x/gov/types"
)

func TestNewMsgUpdateStore(t *testing.T) {
	testCases := []struct {
		Name       string
		Stores     []types.UpdateStore
		ExpectPass bool
	}{
		{
			Name: "success",
			Stores: []types.UpdateStore{{
				Space:    "eth",
				Key:      "01",
				OldValue: "01",
				Value:    "01",
			}},
			ExpectPass: true,
		},
		{
			Name: "empty store space",
			Stores: []types.UpdateStore{{
				Space:    "",
				Key:      "01",
				OldValue: "01",
				Value:    "01",
			}},
			ExpectPass: false,
		},
		{
			Name: "empty key",
			Stores: []types.UpdateStore{{
				Space:    "eth",
				Key:      "",
				OldValue: "01",
				Value:    "01",
			}},
			ExpectPass: false,
		},
		{
			Name: "invalid key",
			Stores: []types.UpdateStore{{
				Space:    "eth",
				Key:      "-",
				OldValue: "01",
				Value:    "01",
			}},
			ExpectPass: false,
		},
		{
			Name: "empty old value",
			Stores: []types.UpdateStore{{
				Space:    "eth",
				Key:      "01",
				OldValue: "",
				Value:    "01",
			}},
			ExpectPass: true,
		},
		{
			Name: "invalid old value",
			Stores: []types.UpdateStore{{
				Space:    "eth",
				Key:      "01",
				OldValue: "-",
				Value:    "01",
			}},
			ExpectPass: false,
		},
		{
			Name: "empty value",
			Stores: []types.UpdateStore{{
				Space:    "eth",
				Key:      "01",
				OldValue: "01",
				Value:    "",
			}},
			ExpectPass: true,
		},
		{
			Name: "invalid value",
			Stores: []types.UpdateStore{{
				Space:    "eth",
				Key:      "01",
				OldValue: "01",
				Value:    "-",
			}},
			ExpectPass: false,
		},
	}
	for _, tc := range testCases {
		msg := types.NewMsgUpdateStore(authtypes.NewModuleAddress(govtypes.ModuleName).String(), tc.Stores)
		if tc.ExpectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %s", tc.Name)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %s", tc.Name)
		}
	}
}
