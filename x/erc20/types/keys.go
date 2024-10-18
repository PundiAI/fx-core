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
	DenomIndexKey  = collections.NewPrefix(11)
	ParamsKey      = collections.NewPrefix(12)
	ERC20TokenKey  = collections.NewPrefix(13)
	BridgeTokenKey = collections.NewPrefix(14)
	IBCTokenKey    = collections.NewPrefix(15)
	CacheKey       = collections.NewPrefix(16)
)
