package types

import (
	"fmt"
)

// constants
const (
	ModuleName = "erc20"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// prefix bytes for the EVM persistent store
const (
	prefixTokenPair = iota + 1
	prefixTokenPairByERC20
	prefixTokenPairByDenom
	prefixIBCTransfer
	prefixAliasDenom
)

// KVStore key prefixes
var (
	KeyPrefixTokenPair        = []byte{prefixTokenPair}
	KeyPrefixTokenPairByERC20 = []byte{prefixTokenPairByERC20}
	KeyPrefixTokenPairByDenom = []byte{prefixTokenPairByDenom}
	KeyPrefixIBCTransfer      = []byte{prefixIBCTransfer}
	KeyPrefixAliasDenom       = []byte{prefixAliasDenom}
)

// GetIBCTransferKey [sourceChannel/sequence]
func GetIBCTransferKey(sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%d", sourceChannel, sequence)
	return append(KeyPrefixIBCTransfer, []byte(key)...)
}
