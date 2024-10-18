package v8

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/erc20/types"
)

var (
	KeyPrefixTokenPair        = []byte{0x01}
	KeyPrefixTokenPairByERC20 = []byte{0x02}
	KeyPrefixTokenPairByDenom = []byte{0x03}
	KeyPrefixIBCTransfer      = []byte{0x04}
	KeyPrefixAliasDenom       = []byte{0x05}
	ParamsKey                 = []byte{0x06}
	KeyPrefixOutgoingTransfer = []byte{0x07}
)

func GetRemovedStoreKeys() [][]byte {
	return [][]byte{
		KeyPrefixTokenPair, KeyPrefixTokenPairByERC20, KeyPrefixTokenPairByDenom, KeyPrefixAliasDenom,
	}
}

func LegacyIsNativeERC20(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, denom string) bool {
	store := ctx.KVStore(storeKey)
	id := store.Get(append(KeyPrefixTokenPairByDenom, []byte(denom)...))
	bz := store.Get(append(KeyPrefixTokenPair, id...))
	if len(bz) == 0 {
		return false
	}
	var tokenPair types.ERC20Token
	cdc.MustUnmarshal(bz, &tokenPair)
	return tokenPair.IsNativeERC20()
}
