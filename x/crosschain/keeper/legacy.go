package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/migrations/v8"
)

// Deprecated: do not use, remove in v8
func (k Keeper) LegacyGetDenomBridgeToken(ctx sdk.Context, denom string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(v8.GetTokenToDenomKey(denom))
	if len(data) == 0 {
		return "", false
	}
	return string(data), true
}
