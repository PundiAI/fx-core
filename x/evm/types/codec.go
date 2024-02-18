package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/functionx/fx-core/v7/x/evm/legacy"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.NewAminoCodec(codec.NewLegacyAmino())

func init() {
	RegisterLegacyAminoCodec(ModuleCdc.LegacyAmino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
}

// RegisterInterfaces registers the client interfaces to protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCallContract{},
	)

	registry.RegisterImplementations(
		(*govv1betal.Content)(nil),
		&legacy.InitEvmParamsProposal{}, // nolint: staticcheck
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCallContract{}, "evm/MsgCallContract", nil)
}
