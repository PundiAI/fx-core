package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// RegisterLegacyAminoCodec registers the necessary x/ibc transfer interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*MemoPacket)(nil), nil)

	cdc.RegisterConcrete(&IbcCallEvmPacket{}, fmt.Sprintf("%s/%s", CompatibleModuleName, "IbcCallEvmPacket"), nil)
}

// RegisterInterfaces register the ibc transfer module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"fx.ibc.applications.transfer.v1.MemoPacket",
		(*MemoPacket)(nil),
		&IbcCallEvmPacket{},
	)
}
