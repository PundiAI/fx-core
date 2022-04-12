package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type ExtensionOptionsWeb3TxI interface{}

// RegisterInterfaces registers the tendermint concrete client-related
// implementations and interfaces.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"fx.ethereum.v1.ExtensionOptionsWeb3Tx",
		(*ExtensionOptionsWeb3TxI)(nil),
		&ExtensionOptionsWeb3Tx{},
	)
}
