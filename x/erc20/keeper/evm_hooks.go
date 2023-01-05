package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// Hooks wrapper struct for erc20 keeper
type Hooks struct {
	k Keeper
}

func (k Keeper) EVMHooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	if !h.k.GetEnableErc20(ctx) || !h.k.GetEnableEVMHook(ctx) {
		return nil
	}
	logger := h.k.Logger(ctx)

	msgTo := receipt.ContractAddress.String()
	if msg.To() != nil {
		msgTo = msg.To().String()
	}

	logger.Info("erc20 processing", "hash", receipt.TxHash.String(), "from", msg.From().String(), "to", msgTo)

	relayTransfers, relayTransferCrossChains, err := h.ParseEventLog(ctx, receipt.Logs, h.k.moduleAddress)
	if err != nil {
		logger.Error("erc20 processing", "hook-action", "parse event log", "error", err.Error())
		return err
	}

	// NOTE: PostTxProcessing doesn't trigger PostTxProcessing
	// NOTE: ConvertERC20NativeToken doesn't trigger PostTxProcessing

	// hook transfer event
	if err := h.HookTransfer(ctx, relayTransfers); err != nil {
		logger.Error("erc20 processing", "hook-action", "relay transfer event", "error", err.Error())
		return err
	}

	// hook transferCrossChain(cross-chain,ibc...) event
	if err := h.HookTransferCrossChain(ctx, relayTransferCrossChains, receipt.TxHash); err != nil {
		logger.Error("erc20 processing", "hook-action", "relay transfer cross chain event", "error", err.Error())
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeERC20Processing,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.From().String()),
		sdk.NewAttribute(types.AttributeKeyTo, msgTo),
		sdk.NewAttribute(types.AttributeKeyEvmTxHash, receipt.TxHash.String()),
	))

	return nil
}

func (h Hooks) ParseEventLog(
	ctx sdk.Context,
	logs []*ethtypes.Log,
	moduleAddress common.Address,
) ([]types.RelayTransfer, []types.RelayTransferCrossChain, error) {
	relayTransfers := make([]types.RelayTransfer, 0, len(logs))
	relayTransferCrossChains := make([]types.RelayTransferCrossChain, 0, len(logs))

	for _, log := range logs {
		tr, err := types.ParseTransferEvent(log)
		if err != nil {
			return nil, nil, sdkerrors.Wrapf(types.ErrUnexpectedEvent, "failed to parse transfer event: %s", err.Error())
		}
		tc, err := types.ParseTransferCrossChainEvent(log)
		if err != nil {
			return nil, nil, sdkerrors.Wrapf(types.ErrUnexpectedEvent, "failed to parse transfer cross chain event: %s", err.Error())
		}

		if (tr == nil || tr.To != moduleAddress) && tc == nil {
			continue
		}

		pair, found := h.k.GetTokenPairByAddress(ctx, log.Address)
		if !found {
			continue
		}
		if !pair.Enabled {
			return nil, nil, sdkerrors.Wrapf(types.ErrERC20TokenPairDisabled, "contract %s, denom %s", pair.Erc20Address, pair.Denom)
		}

		if tr != nil && tr.To == moduleAddress {
			relayTransfers = append(relayTransfers, types.RelayTransfer{
				From:          tr.From,
				Amount:        tr.Value,
				TokenContract: log.Address,
				Denom:         pair.Denom,
				ContractOwner: pair.ContractOwner,
			})
		}
		if tc != nil {
			relayTransferCrossChains = append(relayTransferCrossChains, types.RelayTransferCrossChain{
				TransferCrossChainEvent: tc,
				TokenContract:           log.Address,
				Denom:                   pair.Denom,
			})
		}
	}
	return relayTransfers, relayTransferCrossChains, nil
}
