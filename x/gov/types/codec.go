package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"

	"github.com/functionx/fx-core/v7/x/gov/legacy"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.NewAminoCodec(codec.NewLegacyAmino())

func init() {
	RegisterLegacyAminoCodec(ModuleCdc.LegacyAmino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
}

// RegisterInterfaces register implementations
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgUpdateFXParams{},
		&MsgUpdateEGFParams{},
		&MsgUpdateStore{},
		&MsgUpdateSwitchParams{},
		&legacy.MsgUpdateParams{}, // nolint: staticcheck
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateFXParams{}, "gov/MsgUpdateFXParams", nil)
	cdc.RegisterConcrete(&MsgUpdateEGFParams{}, "gov/MsgUpdateEGFParams", nil)
	cdc.RegisterConcrete(&MsgUpdateStore{}, "gov/MsgUpdateStore", nil)
	cdc.RegisterConcrete(&MsgUpdateSwitchParams{}, "gov/MsgUpdateSwitchParams", nil)
}
