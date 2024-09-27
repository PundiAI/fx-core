package keeper

import (
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/functionx/fx-core/v8/x/ibc/middleware/types"
)

type Keeper struct {
	cdc codec.Codec

	erc20Keeper types.Erc20Keeper
	evmKeeper   types.EvmKeeper
	refundHook  types.RefundHook
}

func NewKeeper(cdc codec.Codec, refundHook types.RefundHook, erc20Keeper types.Erc20Keeper, evmKeeper types.EvmKeeper) Keeper {
	return Keeper{
		cdc:         cdc,
		refundHook:  refundHook,
		erc20Keeper: erc20Keeper,
		evmKeeper:   evmKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+ibcexported.ModuleName+"-"+"middleware")
}
