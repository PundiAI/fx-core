package v3

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func MigrateBridgeToken(cdc codec.BinaryCodec, store sdk.KVStore, moduleName string) {
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
		if strings.HasPrefix(bridgeToken.Denom, ibctransfertypes.DenomPrefix) {
			bridgeToken.Denom = fmt.Sprintf("%s%s", moduleName, bridgeToken.Token)
		}
		store.Set(types.GetTokenToDenomKey(bridgeToken.Denom), []byte(bridgeToken.Token))
		store.Set(types.GetDenomToTokenKey(bridgeToken.Token), []byte(bridgeToken.Denom))
	}
}
