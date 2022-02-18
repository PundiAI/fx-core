package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"

	fxtype "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
)

var _ evmtypes.EvmHooks = (*Keeper)(nil)

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (k Keeper) PostTxProcessing(ctx sdk.Context, txHash common.Hash, logs []*ethtypes.Log) error {
	if ctx.BlockHeight() < fxtype.IntrarelayerSupportBlock() || !k.HasInit(ctx) {
		return nil
	}
	params := k.GetParams(ctx)
	if !params.EnableEVMHook {
		return sdkerrors.Wrap(types.ErrInternalTokenPair, "EVM Hook is currently disabled")
	}
	//process relay token
	if err := k.RelayTokenProcessing(ctx, txHash, logs); err != nil {
		return err
	}
	//process relay chain transfer
	if err := k.RelayTransferChainProcessing(ctx, txHash, logs); err != nil {
		return err
	}
	//process relay ibc transfer
	if err := k.RelayTransferIBCProcessing(ctx, txHash, logs); err != nil {
		return err
	}
	return nil
}
