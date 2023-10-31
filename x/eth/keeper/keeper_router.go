package keeper

import crosschainkeeper "github.com/functionx/fx-core/v6/x/crosschain/keeper"

func NewModuleHandler(keeper crosschainkeeper.Keeper) *crosschainkeeper.ModuleHandler {
	return &crosschainkeeper.ModuleHandler{
		QueryServer:    keeper,
		MsgServer:      NewMsgServerImpl(keeper),
		ProposalServer: keeper,
	}
}
