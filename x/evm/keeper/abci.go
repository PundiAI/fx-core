package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlock sets the sdk Context and EIP155 chain id to the Keeper.
func (k *Keeper) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	k.WithChainID(ctx)
}
