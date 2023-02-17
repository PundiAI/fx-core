package keeper

import (
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v3/x/staking/types"
)

func (h Hooks) HookTransferEvent(ctx sdk.Context, relayTransfers []types.RelayTransfer) error {
	logger := h.k.Logger(ctx)
	for _, relay := range relayTransfers {
		logger.Info("relay lp token", "from", relay.From.String(), "to", relay.To, "amount", relay.Amount.String(),
			"lp-token", relay.TokenContract.String(), "validator", relay.Validator.String())

		// todo get rewards, unbond from, bond to

		telemetry.IncrCounterWithLabels(
			[]string{stakingtypes.ModuleName, "relay_transfer"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("lp_token", relay.TokenContract.String()),
				telemetry.NewLabel("amount", relay.Amount.String()),
			},
		)
	}
	return nil
}
