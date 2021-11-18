package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/functionx/fx-core/x/evm/types"
)

func (k Keeper) HandleInitEvmParamsProposal(ctx sdk.Context, p *types.InitEvmParamsProposal) error {
	// check duplicate init params.
	var evmDenom string
	k.paramSpace.GetIfExists(ctx, types.ParamStoreKeyEVMDenom, &evmDenom)
	if len(evmDenom) != 0 {
		return sdkerrors.Wrapf(types.ErrInvalid, "duplicate init evm params")
	}

	k.Logger(ctx).Info("handle init evm params...", "proposal", p.String())
	// init evm module params
	k.SetParams(ctx, *p.EvmParams)
	// init feeMarket module params
	k.feeMarketKeeper.SetParams(ctx, *p.FeemarketParams)

	// set feeMarket baseFee
	k.feeMarketKeeper.SetBaseFee(ctx, sdk.ZeroInt().BigInt())
	// set feeMarket blockGasUsed
	k.feeMarketKeeper.SetBlockGasUsed(ctx, 0)

	return nil
}
