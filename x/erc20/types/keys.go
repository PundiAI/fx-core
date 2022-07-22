package types

import (
	"fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// constants
const (
	// module name
	ModuleName = "erc20"

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

//GetIBCTransferKey [sourcePort/sourceChannel/sequence]
func GetIBCTransferKey(sourcePort, sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%s/%d", sourcePort, sourceChannel, sequence)
	return append(KeyPrefixIBCTransfer, []byte(key)...)
}
