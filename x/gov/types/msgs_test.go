package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/gov/types"
)

var (
	coinsPos  = sdk.NewCoins(sdk.NewInt64Coin(fxtypes.DefaultDenom, 1000))
	coins     = sdk.NewInt64Coin(fxtypes.DefaultDenom, 1000)
	coinsZero = sdk.NewCoins()

	timeDuration        = time.Second * 60 * 60 * 24 * 7
	invalidTimeDuration = -time.Second * 60 * 60 * 24 * 7

	msgType = sdk.MsgTypeURL(&types.MsgUpdateParams{})
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
		msg := types.NewMsgUpdateParams(authtypes.NewModuleAddress(govtypes.ModuleName).String(), *types.NewParam(tc.MsgType, tc.MinDeposit, tc.MinInitialDeposit, tc.VotingPeriod, tc.Quorum, tc.MaxDepositPeriod, tc.Threshold, tc.VetoThreshold))
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %s", tc.Name)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %s", tc.Name)
		}
	}
}

func TestMsgUpdateParamsGetSignBytes(t *testing.T) {
	msg := types.NewMsgUpdateParams("gov", *types.DefaultParams())
	res := msg.GetSignBytes()
	expected := `{"type":"gov/MsgUpdateParams","value":{"authority":"gov","params":{"max_deposit_period":"172800000000000","min_deposit":[{"amount":"10000000","denom":"stake"}],"min_initial_deposit":{"amount":"1000000000000000000000","denom":"FX"},"msg_type":"/fx.evm.v1.MsgCallContract","quorum":"0.334000000000000000","threshold":"0.500000000000000000","veto_threshold":"0.334000000000000000","voting_period":"172800000000000"}}}`
	require.Equal(t, expected, string(res))
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

func TestMsgUpdateEGFParamsGetSignBytes(t *testing.T) {
	msg := types.NewMsgUpdateEGFParams("gov", *types.DefaultEGFParams())
	res := msg.GetSignBytes()
	expected := `{"type":"gov/MsgUpdateEGFParams","value":{"authority":"gov","params":{"claim_ratio":"0.100000000000000000","egf_deposit_threshold":{"amount":"10000000000000000000000","denom":"FX"}}}}`
	require.Equal(t, expected, string(res))
}
