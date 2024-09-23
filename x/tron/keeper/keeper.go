package keeper

import (
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"

	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
	"github.com/functionx/fx-core/v8/x/tron/types"
)

type Keeper struct {
	crosschainkeeper.Keeper
}

func NewKeeper(keeper crosschainkeeper.Keeper) Keeper {
	return Keeper{
		Keeper: keeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func NewModuleHandler(keeper Keeper) *crosschainkeeper.ModuleHandler {
	return &crosschainkeeper.ModuleHandler{
		QueryServer: crosschainkeeper.NewQueryServerImpl(keeper.Keeper),
		MsgServer:   NewMsgServerImpl(keeper),
	}
}
