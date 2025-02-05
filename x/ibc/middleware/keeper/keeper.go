package keeper

import (
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
)

type Keeper struct {
	cdc                       codec.Codec
	evmKeeper                 types.EvmKeeper
	crosschainKeeper          types.CrosschainKeeper
	crosschaniRouterMsgServer types.CrosschainRouterMsgServer
}

func NewKeeper(cdc codec.Codec, evmKeeper types.EvmKeeper, crosschainKeeper types.CrosschainKeeper, crosschaniRouterMsgServer types.CrosschainRouterMsgServer) Keeper {
	return Keeper{
		cdc:                       cdc,
		evmKeeper:                 evmKeeper,
		crosschainKeeper:          crosschainKeeper,
		crosschaniRouterMsgServer: crosschaniRouterMsgServer,
	}
}

//nolint:revive // SetCrosschainKeeper only for testing
func (k Keeper) SetCrosschainKeeper(crosschainKeeper types.CrosschainKeeper) Keeper {
	k.crosschainKeeper = crosschainKeeper
	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+ibcexported.ModuleName+"-"+"middleware")
}
