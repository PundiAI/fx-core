package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	fxtype "github.com/functionx/fx-core/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
)

var _ evmtypes.EvmHooks = (*Keeper)(nil)

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	if ctx.BlockHeight() < fxtype.IntrarelayerSupportBlock() || !k.HasInit(ctx) {
		return nil
	}
	params := k.GetParams(ctx)
	if !params.EnableIntrarelayer || !params.EnableEVMHook {
		return nil
	}
	//process relay token
	if err := k.RelayTokenProcessing(ctx, msg.From(), msg.To(), receipt); err != nil {
		return err
	}
	//process relay transfer cross(cross-chain,ibc...)
	if err := k.RelayTransferCrossProcessing(ctx, msg.From(), msg.To(), receipt); err != nil {
		return err
	}
	return nil
}
