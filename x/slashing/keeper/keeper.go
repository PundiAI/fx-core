package keeper

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

type Keeper struct {
	keeper.Keeper

	storeKey storetypes.StoreKey
}

func NewKeeper(k keeper.Keeper, key storetypes.StoreKey) Keeper {
	return Keeper{
		Keeper:   k,
		storeKey: key,
	}
}

func (k Keeper) DeleteConsensusPubKey(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(slashtypes.AddrPubkeyRelationKey(consAddr))
}

func (k Keeper) DeleteValidatorSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(slashtypes.ValidatorSigningInfoKey(consAddr))
}
