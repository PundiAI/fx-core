package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
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

	relayTransfers, relayTransferCrossChains, complete := h.ParseEventLog(ctx, receipt.Logs, h.k.moduleAddress)
	if !complete {
		return errors.New("parse event log failed")
	}

	// NOTE: PostTxProcessing doesn't trigger PostTxProcessing
	// NOTE: ConvertERC20NativeToken doesn't trigger PostTxProcessing

	// hook transfer event
	if err := h.HookTransfer(ctx, relayTransfers, receipt.TxHash); err != nil {
		return err
	}

	// hook transferCrossChain(cross-chain,ibc...) event
	if err := h.HookTransferCrossChain(ctx, relayTransferCrossChains, msg.From(), msg.To(), receipt.TxHash); err != nil {
		return err
	}
	return nil
}

func (h Hooks) ParseEventLog(ctx sdk.Context, logs []*ethtypes.Log, moduleAddress common.Address) ([]types.RelayTransfer, []types.RelayTransferCrossChain, bool) {
	fip20ABI := fxtypes.GetERC20().ABI

	relayTransfers := make([]types.RelayTransfer, 0, len(logs))
	relayTransferCrossChains := make([]types.RelayTransferCrossChain, 0, len(logs))

	for _, log := range logs {
		tr, toAddr, err := types.ParseTransferEvent(fip20ABI, log)
		if err != nil {
			return nil, nil, false
		}
		tc, err := types.ParseTransferCrossChainEvent(fip20ABI, log)
		if err != nil {
			return nil, nil, false
		}
		pair, found := h.k.GetTokenPairByAddress(ctx, log.Address)
		if !found || !pair.Enabled {
			continue
		}
		if toAddr == moduleAddress {
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
	return relayTransfers, relayTransferCrossChains, true
}
