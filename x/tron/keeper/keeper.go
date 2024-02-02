package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/tron/types"
)

type Keeper struct {
	crosschainkeeper.Keeper
	erc20Keeper crosschaintypes.Erc20Keeper
}

func NewKeeper(keeper crosschainkeeper.Keeper, erc20Keeper crosschaintypes.Erc20Keeper) Keeper {
	return Keeper{
		Keeper:      keeper,
		erc20Keeper: erc20Keeper,
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
