package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces register implementations
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgConvertCoin{},
		&MsgUpdateParams{},
		&MsgToggleTokenConversion{},
		&MsgRegisterNativeCoin{},
		&MsgRegisterNativeERC20{},
		&MsgRegisterBridgeToken{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgConvertCoin{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConvertCoin"), nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateParams"), nil)
	cdc.RegisterConcrete(&MsgToggleTokenConversion{}, fmt.Sprintf("%s/%s", ModuleName, "MsgToggleTokenConversion"), nil)
	cdc.RegisterConcrete(&MsgRegisterNativeCoin{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRegisterNativeCoin"), nil)
	cdc.RegisterConcrete(&MsgRegisterNativeERC20{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRegisterNativeERC20"), nil)
	cdc.RegisterConcrete(&MsgRegisterBridgeToken{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRegisterBridgeToken"), nil)
}
