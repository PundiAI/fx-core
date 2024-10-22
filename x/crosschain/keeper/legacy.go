package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/crosschain/migrations/v8"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

// Deprecated: do not use, remove in v8
func (k Keeper) LegacyGetDenomBridgeToken(ctx sdk.Context, denom string) (*types.BridgeToken, bool) {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(v8.GetTokenToDenomKey(denom))
	if len(data) == 0 {
		return nil, false
	}
	return &types.BridgeToken{
		Denom: denom,
		Token: string(data),
	}, true
}
