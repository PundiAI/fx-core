package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"

	"github.com/functionx/fx-core/x/gravity/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	legacyMsgServer legacyMsgServer

	crosschainkeeper.EthereumMsgServer
	stakingKeeper types.StakingKeeper
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(msgServer legacyMsgServer) Keeper {
	return Keeper{legacyMsgServer: msgServer}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
