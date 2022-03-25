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
func (k Keeper) PostTxProcessing(ctx sdk.Context, from common.Address, to *common.Address, receipt *ethtypes.Receipt) error {
	if ctx.BlockHeight() < fxtype.IntrarelayerSupportBlock() || !k.HasInit(ctx) {
		return nil
	}
	params := k.GetParams(ctx)
	if !params.EnableEVMHook {
		return sdkerrors.Wrap(types.ErrInternalTokenPair, "EVM Hook is currently disabled")
	}
	//process relay token
	if err := k.RelayTokenProcessing(ctx, from, to, receipt); err != nil {
		return err
	}
	//process relay transfer cross(cross-chain,ibc...)
	if err := k.RelayTransferCrossProcessing(ctx, from, to, receipt); err != nil {
		return err
	}
	return nil
}
