package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Hooks wrapper struct for erc20 keeper
type Hooks struct {
	k *Keeper
}

// NewHooks Return the wrapper struct
func NewHooks(k *Keeper) Hooks {
	return Hooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	if !h.k.GetEnableErc20(ctx) || !h.k.GetEnableEVMHook(ctx) {
		return nil
	}

	eventLog, complete := ParseEventLog(receipt, h.k.moduleAddress)
	if !complete {
		return errors.New("parse event log failed")
	}

	eventLog, err := h.k.TokenPairEnable(ctx, eventLog)
	if err != nil {
		return err
	}

	// NOTE: PostTxProcessing doesn't trigger PostTxProcessing
	// NOTE: ConvertERC20NativeToken doesn't trigger PostTxProcessing

	// hook relay token
	if err := h.k.HookRelayToken(ctx, eventLog.RelayToken, receipt); err != nil {
		return err
	}

	// hook transfer cross chain(cross-chain,ibc...)
	if err := h.k.HookTransferCrossChain(ctx, eventLog.TransferCrossChain, msg.From(), msg.To(), receipt); err != nil {
		return err
	}
	return nil
}
