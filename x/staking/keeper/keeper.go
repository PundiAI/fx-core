package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type Keeper struct {
	*stakingkeeper.Keeper
	storeKey storetypes.StoreKey
	cdc      codec.Codec
}

func NewKeeper(k *stakingkeeper.Keeper, storeKey storetypes.StoreKey, cdc codec.Codec) *Keeper {
	return &Keeper{
		Keeper:   k,
		storeKey: storeKey,
		cdc:      cdc,
	}
}
