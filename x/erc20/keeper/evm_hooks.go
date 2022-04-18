package keeper

import (
	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Hooks wrapper struct for erc20 keeper
type Hooks struct {
	k Keeper
}

// Hooks Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	if ctx.BlockHeight() < fxtypes.EvmSupportBlock() {
		return nil
	}
	params := h.k.GetParams(ctx)
	if !params.EnableErc20 || !params.EnableEVMHook {
		return nil
	}
	//process relay token
	if err := h.k.RelayTokenProcessing(ctx, msg.From(), msg.To(), receipt); err != nil {
		return err
	}
	//process relay transfer cross chain(cross-chain,ibc...)
	if err := h.k.RelayTransferCrossChainProcessing(ctx, msg.From(), msg.To(), receipt); err != nil {
		return err
	}
	return nil
}

// TODO: Make sure that if ConvertERC20 is called, that the Hook doesn't trigger
// if it does, delete minting from ConvertErc20
