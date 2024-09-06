package types

import crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"

const (
	// ModuleName is the name of the module
	ModuleName = "layer2"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

func init() {
	crosschaintypes.RegisterExternalAddress(ModuleName, crosschaintypes.EthereumAddress{})
}
