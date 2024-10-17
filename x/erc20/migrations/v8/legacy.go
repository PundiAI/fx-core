package v8

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func LegacyIsNativeERC20(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, denom string) bool {
	store := ctx.KVStore(storeKey)
	id := store.Get(append(types.KeyPrefixTokenPairByDenom, []byte(denom)...))
	bz := store.Get(append(types.KeyPrefixTokenPair, id...))
	if len(bz) == 0 {
		return false
	}
	var tokenPair types.ERC20Token
	cdc.MustUnmarshal(bz, &tokenPair)
	return tokenPair.IsNativeERC20()
}
