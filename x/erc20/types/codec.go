package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterInterfaces register implementations
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgConvertCoin{},
		&MsgConvertERC20{},
		&MsgConvertDenom{},

		&MsgUpdateParams{},
		&MsgRegisterCoin{},
		&MsgRegisterERC20{},
		&MsgToggleTokenConversion{},
		&MsgUpdateDenomAlias{},
		&MsgRegisterNativeCoin{},
		&MsgRegisterNativeERC20{},
		&MsgUpdateBridgeToken{},
	)
	registry.RegisterImplementations(
		(*govv1betal.Content)(nil),
		&RegisterCoinProposal{},
		&RegisterERC20Proposal{},
		&ToggleTokenConversionProposal{},
		&UpdateDenomAliasProposal{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgConvertCoin{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConvertCoin"), nil)
	cdc.RegisterConcrete(&MsgConvertERC20{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConvertERC20"), nil)
	cdc.RegisterConcrete(&MsgConvertDenom{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConvertDenom"), nil)

	cdc.RegisterConcrete(&MsgUpdateParams{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateParams"), nil)
	cdc.RegisterConcrete(&MsgRegisterCoin{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRegisterCoin"), nil)
	cdc.RegisterConcrete(&MsgRegisterERC20{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRegisterERC20"), nil)
	cdc.RegisterConcrete(&MsgToggleTokenConversion{}, fmt.Sprintf("%s/%s", ModuleName, "MsgToggleTokenConversion"), nil)
	cdc.RegisterConcrete(&MsgUpdateDenomAlias{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateDenomAlias"), nil)
	cdc.RegisterConcrete(&MsgRegisterNativeCoin{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRegisterNativeCoin"), nil)
	cdc.RegisterConcrete(&MsgRegisterNativeERC20{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRegisterNativeERC20"), nil)
	cdc.RegisterConcrete(&MsgUpdateBridgeToken{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateBridgeToken"), nil)

	cdc.RegisterConcrete(&RegisterCoinProposal{}, fmt.Sprintf("%s/%s", ModuleName, "RegisterCoinProposal"), nil)
	cdc.RegisterConcrete(&RegisterERC20Proposal{}, fmt.Sprintf("%s/%s", ModuleName, "RegisterERC20Proposal"), nil)
	cdc.RegisterConcrete(&ToggleTokenConversionProposal{}, fmt.Sprintf("%s/%s", ModuleName, "ToggleTokenConversionProposal"), nil)
	cdc.RegisterConcrete(&UpdateDenomAliasProposal{}, fmt.Sprintf("%s/%s", ModuleName, "UpdateDenomAliasProposal"), nil)
}
