package types

import crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"

const (
	// ModuleName is the name of the module
	ModuleName = "tron"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

func init() {
	crosschaintypes.RegisterExternalAddress(ModuleName, tronAddress{})
}
