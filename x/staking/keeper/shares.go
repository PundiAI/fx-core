package keeper

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/staking/types"
)

// SetAllowance sets the allowance of a spender for a delegator.
// shares must be non-negative.
func (k Keeper) SetAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress, shares *big.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAllowanceKey(valAddr, owner, spender), shares.Bytes())
}

// GetAllowance returns the allowance of a spender for a delegator.
func (k Keeper) GetAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress) *big.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAllowanceKey(valAddr, owner, spender))
	if len(bz) == 0 {
		return big.NewInt(0)
	}
	return big.NewInt(0).SetBytes(bz)
}

// IterateAllAllowance iterates over the allowances and performs a callback function.
func (k Keeper) IterateAllAllowance(ctx sdk.Context, handler func(valAddr sdk.ValAddress, owner, spender sdk.AccAddress, allowance *big.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.AllowanceKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := bytes.TrimPrefix(iter.Key(), types.AllowanceKey)
		valAddrLen := key[0]
		valAddr := sdk.ValAddress(key[1 : 1+valAddrLen])
		ownerAddrLen := key[1+valAddrLen]
		owner := sdk.AccAddress(key[2+valAddrLen : 2+valAddrLen+ownerAddrLen])
		spenderAddrLen := key[2+valAddrLen+ownerAddrLen]
		spender := sdk.AccAddress(key[3+valAddrLen+ownerAddrLen : 3+valAddrLen+ownerAddrLen+spenderAddrLen])
		shares := big.NewInt(0).SetBytes(iter.Value())
		if handler(valAddr, owner, spender, shares) {
			break
		}
	}
}
