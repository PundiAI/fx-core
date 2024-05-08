package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	"github.com/functionx/fx-core/v7/x/tron/types"
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
		QueryServer:    keeper,
		MsgServer:      NewMsgServerImpl(keeper),
		ProposalServer: keeper,
	}
}
