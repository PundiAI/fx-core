package types

import (
	fxtypes "github.com/pundiai/fx-core/v8/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "arbitrum"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

func init() {
	fxtypes.RegisterExternalAddress(ModuleName, ethtypes.EthereumAddress{})
}
