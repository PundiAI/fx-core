package forks

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	fxtypes "github.com/functionx/fx-core/types"
)

func UpdateFXMetadata(ctx sdk.Context, bankKeeper bankKeeper.Keeper, key *types.KVStoreKey) {
	//delete fx
	deleteMetadata(ctx, key, strings.ToLower(fxtypes.DefaultDenom))
	//set FX
	md := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
	if err := md.Validate(); err != nil {
		panic("invalid FX metadata")
	}
	bankKeeper.SetDenomMetaData(ctx, md)
}

func deleteMetadata(ctx sdk.Context, key *types.KVStoreKey, base ...string) {
	store := ctx.KVStore(key)
	for _, b := range base {
		store.Delete(banktypes.DenomMetadataKey(b))
	}
}
