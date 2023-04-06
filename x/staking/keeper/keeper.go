package keeper

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type Keeper struct {
	stakingkeeper.Keeper
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(k stakingkeeper.Keeper, storeKey storetypes.StoreKey) Keeper {
	return Keeper{
		Keeper:   k,
		storeKey: storeKey,
	}
}
