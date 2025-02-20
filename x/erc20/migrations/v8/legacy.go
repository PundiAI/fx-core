package v8

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		KeyPrefixTokenPair, KeyPrefixTokenPairByERC20, KeyPrefixTokenPairByDenom, KeyPrefixIBCTransfer, KeyPrefixAliasDenom, KeyPrefixOutgoingTransfer,
	}
}

func GetBaseDenom(ctx sdk.Context, storeKey storetypes.StoreKey, alias string) (string, bool) {
	store := ctx.KVStore(storeKey)
	value := store.Get(append(KeyPrefixAliasDenom, []byte(alias)...))
	if value == nil {
		return "", false
	}
	return string(value), true
}
