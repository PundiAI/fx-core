package keeper

import crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"

type Keeper struct {
	keeper crosschainkeeper.Keeper
}

func NewKeeper(keeper crosschainkeeper.Keeper) Keeper {
	return Keeper{keeper: keeper}
}

func (k Keeper) GetCrosschainKeeper() crosschainkeeper.Keeper {
	return k.keeper
}
