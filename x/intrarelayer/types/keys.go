package types

import (
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// constants
const (
	// module name
	ModuleName = "intrarelayer"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

// prefix bytes for the EVM persistent store
const (
	prefixTokenPair = iota + 1
	prefixTokenPairByFIP20
	prefixTokenPairByDenom
	prefixIBCTransfer
)

// KVStore key prefixes
var (
	KeyPrefixTokenPair        = []byte{prefixTokenPair}
	KeyPrefixTokenPairByFIP20 = []byte{prefixTokenPairByFIP20}
	KeyPrefixTokenPairByDenom = []byte{prefixTokenPairByDenom}
	KeyPrefixIBCTransfer      = []byte{prefixIBCTransfer}
)

//GetIBCTransferKey [sourcePort/sourceChannel/sequence]
func GetIBCTransferKey(sourcePort, sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%s/%d", sourcePort, sourceChannel, sequence)
	return append(KeyPrefixIBCTransfer, []byte(key)...)
}

// UInt64FromBytes create uint from binary big endian representation
func UInt64FromBytes(s []byte) uint64 {
	return binary.BigEndian.Uint64(s)
}

// UInt64Bytes uses the SDK byte marshaling to encode a uint64
func UInt64Bytes(n uint64) []byte {
	return sdk.Uint64ToBigEndian(n)
}
