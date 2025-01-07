package types

import (
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "tron"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

func init() {
	fxtypes.RegisterExternalAddress(ModuleName, tronAddress{})
}
