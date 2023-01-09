package v2

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func iterateIBCTransferRelationLegacy(store sdk.KVStore, cb func(port, channel string, sequence uint64) bool) {
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixIBCTransfer)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		port, channel, sequence, found := parseIBCTransferKeyLegacy(iter.Key())
		if !found {
			panic(fmt.Sprintf("failed to parse ibc transfer key: %s", string(iter.Key())))
		}
		if cb(port, channel, sequence) {
			return
		}
	}
}

func deleteIBCTransferRelationLegacy(store sdk.KVStore, port, channel string, sequence uint64) {
	store.Delete(getIBCTransferKeyLegacy(port, channel, sequence))
}

// getIBCTransferKeyLegacy key -> [sourcePort/sourceChannel/sequence]
func getIBCTransferKeyLegacy(sourcePort, sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%s/%d", sourcePort, sourceChannel, sequence)
	return append(types.KeyPrefixIBCTransfer, []byte(key)...)
}

func parseIBCTransferKeyLegacy(key []byte) (string, string, uint64, bool) {
	if bytes.HasPrefix(key, types.KeyPrefixIBCTransfer) {
		return "", "", 0, false
	}
	keyStr := string(bytes.TrimPrefix(key, types.KeyPrefixIBCTransfer))

	split := strings.Split(keyStr, "/")
	if len(split) != 3 {
		return "", "", 0, false
	}

	port := split[0]
	channel := split[1]
	sequence, err := strconv.ParseUint(split[2], 10, 64)
	if err != nil {
		return "", "", 0, false
	}
	return port, channel, sequence, true
}
