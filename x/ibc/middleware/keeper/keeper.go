package keeper

import (
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/functionx/fx-core/v8/x/ibc/middleware/types"
)

type Keeper struct {
	cdc              codec.Codec
	evmKeeper        types.EvmKeeper
	crossChainKeeper types.CrossChainKeeper
}

func NewKeeper(cdc codec.Codec, evmKeeper types.EvmKeeper, crossChainKeeper types.CrossChainKeeper) Keeper {
	return Keeper{
		cdc:              cdc,
		evmKeeper:        evmKeeper,
		crossChainKeeper: crossChainKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+ibcexported.ModuleName+"-"+"middleware")
}
