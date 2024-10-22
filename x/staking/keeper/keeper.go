package keeper

import (
	storetypes "cosmossdk.io/store/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type Keeper struct {
	*stakingkeeper.Keeper
	storeKey storetypes.StoreKey
}

func NewKeeper(k *stakingkeeper.Keeper, storeKey storetypes.StoreKey) *Keeper {
	return &Keeper{
		Keeper:   k,
		storeKey: storeKey,
	}
}
