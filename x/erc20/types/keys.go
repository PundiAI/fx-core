package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// constants
const (
	ModuleName = "erc20"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// KVStore key prefixes
var (
	KeyPrefixTokenPair        = []byte{0x01}
	KeyPrefixTokenPairByERC20 = []byte{0x02}
	KeyPrefixTokenPairByDenom = []byte{0x03}
	KeyPrefixIBCTransfer      = []byte{0x04}
	KeyPrefixAliasDenom       = []byte{0x05}
	ParamsKey                 = []byte{0x06}
	KeyPrefixOutgoingTransfer = []byte{0x07}
)

// GetIBCTransferKey [sourceChannel/sequence]
func GetIBCTransferKey(sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%d", sourceChannel, sequence)
	return append(KeyPrefixIBCTransfer, []byte(key)...)
}

// GetOutgoingTransferKey [txID]
func GetOutgoingTransferKey(txID uint64) []byte {
	return append(KeyPrefixOutgoingTransfer, sdk.Uint64ToBigEndian(txID)...)
}
