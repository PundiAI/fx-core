package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/gov/types"
)

var (
	coinsPos  = sdk.NewCoins(sdk.NewInt64Coin(fxtypes.DefaultDenom, 1000))
	coins     = sdk.NewInt64Coin(fxtypes.DefaultDenom, 1000)
	coinsZero = sdk.NewCoins()

	timeDuration        = time.Second * 60 * 60 * 24 * 7
	invalidTimeDuration = -time.Second * 60 * 60 * 24 * 7

	msgType = sdk.MsgTypeURL(&types.MsgUpdateFXParams{})
)

func TestNewMsgUpdateParams(t *testing.T) {
	testCases := []struct {
		Name              string
		MsgType           string
		MinDeposit        []sdk.Coin
		MinInitialDeposit sdk.Coin
		VotingPeriod      *time.Duration
		Quorum            string
		MaxDepositPeriod  *time.Duration
		Threshold         string
		VetoThreshold     string
		expectPass        bool
	}{
		{"success ", msgType, coinsPos, coins, &timeDuration, "0.2", &timeDuration, "0.5", "0.334", true},
		{"proto type not registered", "", coinsPos, coins, &timeDuration, "0.2", &timeDuration, "0.5", "0.334", true},
		{"invalid  minDeposit", msgType, coinsZero, coins, &timeDuration, "0.2", &timeDuration, "0.5", "0.334", false},
		{"invalid minInitialDeposit ", msgType, coinsPos, sdk.Coin{}, &timeDuration, "0.2", &timeDuration, "0.5", "0.334", false},
		{"invalid votingPeriod", msgType, coinsPos, coins, nil, "0.2", &timeDuration, "0.5", "0.334", false},
		{"invalid quorum", msgType, coinsPos, coins, &timeDuration, "-0.2", &timeDuration, "0.5", "0.334", false},
		{"invalid maxDepositPeriod", msgType, coinsPos, coins, &timeDuration, "0.2", &invalidTimeDuration, "0.5", "0.334", false},
		{"invalid threshold", msgType, coinsPos, coins, &timeDuration, "0.2", &timeDuration, "-0.5", "0.334", false},
		{"invalid vetoThreshold", msgType, coinsPos, coins, &timeDuration, "0.2", &timeDuration, "0.5", "", false},
	}
	for _, tc := range testCases {
		msg := types.NewMsgUpdateFXParams(authtypes.NewModuleAddress(govtypes.ModuleName).String(), *types.NewParam(tc.MsgType, tc.MinDeposit, tc.MinInitialDeposit, tc.VotingPeriod, tc.Quorum, tc.MaxDepositPeriod, tc.Threshold, tc.VetoThreshold, "", false, false, false))
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %s", tc.Name)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %s", tc.Name)
		}
	}
}

func TestNewMsgUpdateEGFParams(t *testing.T) {
	testCases := []struct {
		Name                string
		EgfDepositThreshold sdk.Coin
		ClaimRatio          string
		expectPass          bool
	}{
		{"success", coins, "0.1", true},
		{"invalid egfDepositThreshold ", sdk.Coin{}, "0.1", false},
		{"invalid claimRatio", coins, "-0.1", false},
	}
	for _, tc := range testCases {
		msg := types.NewMsgUpdateEGFParams(authtypes.NewModuleAddress(govtypes.ModuleName).String(), *types.NewEGFParam(tc.EgfDepositThreshold, tc.ClaimRatio))
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %s", tc.Name)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %s", tc.Name)
		}
	}
}

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
