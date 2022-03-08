package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/functionx/fx-core/x/evm/types"
)

func (k Keeper) HandleInitEvmParamsProposal(ctx sdk.Context, p *types.InitEvmParamsProposal) error {
	// check duplicate init params.
	if k.HasInit(ctx) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate init evm params")
	}

	k.Logger(ctx).Info("handle init evm params...", "proposal", p.String())

	if p.FeemarketParams.BaseFee.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "base fee cannot be negative")
	}
	// set feeMarket baseFee
	k.feeMarketKeeper.SetBaseFee(ctx, p.FeemarketParams.BaseFee.BigInt())
	// set feeMarket blockGasUsed
	k.feeMarketKeeper.SetBlockGasUsed(ctx, 0)
	// init feeMarket module params
	k.feeMarketKeeper.SetParams(ctx, *p.FeemarketParams)

	// init evm module params
	k.SetParams(ctx, *p.EvmParams)
	return nil
}
