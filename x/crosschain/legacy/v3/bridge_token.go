package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func MigrateBridgeToken(cdc codec.BinaryCodec, store sdk.KVStore) {
	iter := sdk.KVStorePrefixIterator(store, types.TokenToDenomKey)
	defer iter.Close()

	var bridgeTokens []types.BridgeToken
	for ; iter.Valid(); iter.Next() {
		var bridgeToken types.BridgeToken
		cdc.MustUnmarshal(iter.Value(), &bridgeToken)
		bridgeToken.Denom = string(iter.Key()[len(types.TokenToDenomKey):])
		bridgeTokens = append(bridgeTokens, bridgeToken)
	}

	for _, bridgeToken := range bridgeTokens {
		store.Set(types.GetTokenToDenomKey(bridgeToken.Denom), []byte(bridgeToken.Token))
		store.Set(types.GetDenomToTokenKey(bridgeToken.Token), []byte(bridgeToken.Denom))
	}
}
