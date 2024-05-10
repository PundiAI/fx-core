package types

import crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"

const (
	// ModuleName is the name of the module
	ModuleName = "arbitrum"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

func init() {
	crosschaintypes.RegisterExternalAddress(ModuleName, crosschaintypes.EthereumAddress{})
}
