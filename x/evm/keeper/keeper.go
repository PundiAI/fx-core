package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

type Keeper struct {
	*evmkeeper.Keeper

	// access to account state
	accountKeeper types.AccountKeeper

	// has evm hooks
	hasHooks bool
}

func NewKeeper(ek *evmkeeper.Keeper, ak types.AccountKeeper) *Keeper {
	return &Keeper{
		Keeper:        ek,
		accountKeeper: ak,
	}
}

// WithChainID sets the chain id to the local variable in the keeper
func (k *Keeper) WithChainID(ctx sdk.Context) {
	cacheCtx, _ := ctx.CacheContext()
	cacheCtx = cacheCtx.WithChainID(fxtypes.ChainIdWithEIP155())

	k.Keeper.WithChainID(cacheCtx)
}

// SetHooks sets the hooks for the EVM module
// It should be called only once during initialization, it panic if called more than once.
func (k *Keeper) SetHooks(eh types.EvmHooks) *Keeper {
	k.Keeper.SetHooks(eh)
	k.hasHooks = true
	return k
}
