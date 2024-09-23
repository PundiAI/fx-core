package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlock sets the sdk Context and EIP155 chain id to the Keeper.
func (k *Keeper) BeginBlock(ctx sdk.Context) error {
	// cache parameters that's common for the whole block.
	if _, err := k.EVMBlockConfig(ctx, k.ChainID()); err != nil {
		return err
	}

	return nil
}
