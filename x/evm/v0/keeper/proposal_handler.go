package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/evm/v0/types"
)

func (k Keeper) HandleInitEvmParamsProposal(ctx sdk.Context, p *types.InitEvmParamsProposal) error {
	// check duplicate init params.
	if k.HasInit(ctx) {
		return sdkerrors.Wrapf(types.ErrInvalid, "duplicate init evm params")
	}

	k.Logger(ctx).Info("handle init evm params...", "proposal", p.String())
	// init evm module params
	k.SetParams(ctx, *p.EvmParams)
	// init feeMarket module params
	k.feeMarketKeeper.SetParams(ctx, *p.FeemarketParams)

	baseFee := sdk.ZeroInt()
	if !p.FeemarketParams.NoBaseFee && p.FeemarketParams.InitialBaseFee > 0 {
		baseFee = sdk.NewInt(p.FeemarketParams.InitialBaseFee)
	}
	// set feeMarket baseFee
	k.feeMarketKeeper.SetBaseFee(ctx, baseFee.BigInt())
	// set feeMarket blockGasUsed
	k.feeMarketKeeper.SetBlockGasUsed(ctx, 0)

	return nil
}
