package keeper

import (
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (h Hooks) HookTransferEvent(ctx sdk.Context, relayTransfers []types.RelayTransfer) error {
	fip20ABI := fxtypes.GetERC20().ABI
	for _, relay := range relayTransfers {
		h.k.Logger(ctx).Info("relay token", "from", relay.From.Hex(), "amount", relay.Amount.String(), "denom", relay.Denom, "token", relay.TokenContract)

		// create the corresponding sdk.Coin that is paired with FIP20
		coins := sdk.Coins{{Denom: relay.Denom, Amount: sdk.NewIntFromBigInt(relay.Amount)}}

		switch relay.ContractOwner {
		case types.OWNER_MODULE: // native coin
			if _, err := h.k.CallEVM(ctx, fip20ABI, h.k.moduleAddress, relay.TokenContract, true, "burn", h.k.moduleAddress, relay.Amount); err != nil {
				return err
			}

			if relay.Denom == fxtypes.DefaultDenom {
				if err := h.k.bankKeeper.SendCoinsFromAccountToModule(ctx, relay.TokenContract.Bytes(), types.ModuleName, coins); err != nil {
					return err
				}
			}
		case types.OWNER_EXTERNAL: // native coin
			if err := h.k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				return err
			}
		default:
			return types.ErrUndefinedOwner
		}

		// sender receive relay amount
		recipient := sdk.AccAddress(relay.From.Bytes())
		if err := h.k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins); err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRelayTransfer,
			sdk.NewAttribute(sdk.AttributeKeySender, relay.From.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, sdk.AccAddress(recipient.Bytes()).String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, relay.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, relay.Denom),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, relay.TokenContract.String()),
		))

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "relay_transfer"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", relay.TokenContract.String()),
				telemetry.NewLabel("denom", relay.Denom),
				telemetry.NewLabel("amount", relay.Amount.String()),
			},
		)
	}
	return nil
}
