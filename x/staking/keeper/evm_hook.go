package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v3/x/staking/types"
)

// Hooks wrapper struct for staking keeper
type Hooks struct {
	k Keeper
}

func (k Keeper) EVMHooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	relayTransfers, err := h.ParseEventLog(ctx, receipt.Logs)
	if err != nil {
		h.k.Logger(ctx).Error("staking processing", "hook-action", "parse event log", "error", err.Error())
		return err
	}

	if len(relayTransfers) <= 0 {
		return nil
	}

	msgTo := receipt.ContractAddress.String()
	if msg.To() != nil {
		msgTo = msg.To().String()
	}
	h.k.Logger(ctx).Info("staking processing", "hash", receipt.TxHash.String(), "from", msg.From().String(), "to", msgTo)

	// hook transfer event
	if err = h.HookTransferEvent(ctx, relayTransfers); err != nil {
		h.k.Logger(ctx).Error("staking processing", "hook-action", "relay transfer event", "error", err.Error())
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeStakingProcessing,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.From().String()),
		sdk.NewAttribute(types.AttributeKeyTo, msgTo),
		sdk.NewAttribute(types.AttributeKeyEvmTxHash, receipt.TxHash.String()),
	))

	return nil
}

func (h Hooks) ParseEventLog(ctx sdk.Context, logs []*ethtypes.Log) ([]types.RelayTransfer, error) {
	relayTransfers := make([]types.RelayTransfer, 0, len(logs))
	for _, log := range logs {
		tr, err := types.ParseTransferEvent(log)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrUnexpectedEvent, "parse transfer event: %s", err.Error())
		}

		if tr == nil {
			continue
		}
		valAddr, found := h.k.GetLPTokenValidator(ctx, log.Address)
		if !found {
			continue
		}

		relayTransfers = append(relayTransfers, types.RelayTransfer{
			From:          tr.From,
			To:            tr.To,
			Amount:        tr.Value,
			TokenContract: log.Address,
			Validator:     valAddr,
		})
	}

	return relayTransfers, nil
}
