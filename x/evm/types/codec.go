package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.NewAminoCodec(codec.NewLegacyAmino())

func init() {
	RegisterLegacyAminoCodec(ModuleCdc.LegacyAmino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
}

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCallContract{},
	)
	evmtypes.RegisterInterfaces(registry)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCallContract{}, "evm/MsgCallContract", nil)
}
