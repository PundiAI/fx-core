package v8

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/app/upgrades/store"
)

var (
	// Deprecated: DenomToTokenKey prefixes the index of asset denom to external token
	DenomToTokenKey = []byte{0x26}

	// Deprecated: TokenToDenomKey prefixes the index of assets external token to denom
	TokenToDenomKey = []byte{0x27}

	// Deprecated: remove data in upgrade
	BridgeCallFromMsgKey = []byte{0x51}
)

// Deprecated: GetTokenToDenomKey returns the following key format
func GetTokenToDenomKey(denom string) []byte {
	return append(TokenToDenomKey, []byte(denom)...)
}

func GetRemovedStoreKeys() [][]byte {
	return [][]byte{
		DenomToTokenKey,
		TokenToDenomKey,
		BridgeCallFromMsgKey,
	}
}

func Migrate(ctx sdk.Context, storeKey storetypes.StoreKey) error {
	store.RemoveStoreKeys(ctx, storeKey, GetRemovedStoreKeys())
	return nil
}
