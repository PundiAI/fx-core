package types

import (
	"cosmossdk.io/collections"
)

// constants
const (
	ModuleName = "erc20"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
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

	DenomIndexKey  = collections.NewPrefix(1)
	ParamsKey2     = collections.NewPrefix(2)
	ERC20TokenKey  = collections.NewPrefix(3)
	BridgeTokenKey = collections.NewPrefix(4)
	IBCTokenKey    = collections.NewPrefix(5)
	CacheKey       = collections.NewPrefix(6)
)
