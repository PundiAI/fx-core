package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	"github.com/pundiai/fx-core/v8/types/legacy"
)

// RegisterInterfaces register implementations
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),

		&legacy.MsgUpdateFXParams{},
		&legacy.MsgUpdateEGFParams{},
		&legacy.MsgUpdateParams{},

		&MsgUpdateStore{},
		&MsgUpdateSwitchParams{},
		&MsgUpdateCustomParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&legacy.MsgUpdateFXParams{}, "gov/MsgUpdateFXParams", nil)
	cdc.RegisterConcrete(&legacy.MsgUpdateEGFParams{}, "gov/MsgUpdateEGFParams", nil)

	cdc.RegisterConcrete(&MsgUpdateStore{}, "gov/MsgUpdateStore", nil)
	cdc.RegisterConcrete(&MsgUpdateSwitchParams{}, "gov/MsgUpdateSwitchParams", nil)
	cdc.RegisterConcrete(&MsgUpdateCustomParams{}, "gov/MsgUpdateCustomParams", nil)
}
