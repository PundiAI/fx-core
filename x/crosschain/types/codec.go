package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBondedOracle{},
		&MsgAddDelegate{},
		&MsgReDelegate{},
		&MsgEditBridger{},
		&MsgWithdrawReward{},

		&MsgSendToExternal{},

		&MsgUpdateParams{},
		&MsgUpdateChainOracles{},

		&MsgClaim{},

		&MsgConfirm{},
	)

	registry.RegisterInterface(
		"fx.gravity.crosschain.v1.ExternalClaim",
		(*ExternalClaim)(nil),
		&MsgSendToExternalClaim{},
		&MsgSendToFxClaim{},
		&MsgBridgeCallClaim{},
		&MsgBridgeTokenClaim{},
		&MsgOracleSetUpdatedClaim{},
		&MsgBridgeCallResultClaim{},
	)

	registry.RegisterInterface(
		"fx.gravity.crosschain.v1.Confirm",
		(*Confirm)(nil),
		&MsgConfirmBatch{},
		&MsgOracleSetConfirm{},
		&MsgBridgeCallConfirm{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*ExternalClaim)(nil), nil)
	cdc.RegisterInterface((*Confirm)(nil), nil)

	cdc.RegisterConcrete(&MsgBondedOracle{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBondedOracle"), nil)
	cdc.RegisterConcrete(&MsgAddDelegate{}, fmt.Sprintf("%s/%s", ModuleName, "MsgAddDelegate"), nil)
	cdc.RegisterConcrete(&MsgReDelegate{}, fmt.Sprintf("%s/%s", ModuleName, "MsgReDelegate"), nil)
	cdc.RegisterConcrete(&MsgEditBridger{}, fmt.Sprintf("%s/%s", ModuleName, "MsgEditBridger"), nil)
	cdc.RegisterConcrete(&MsgWithdrawReward{}, fmt.Sprintf("%s/%s", ModuleName, "MsgWithdrawReward"), nil)

	cdc.RegisterConcrete(&MsgOracleSetConfirm{}, fmt.Sprintf("%s/%s", ModuleName, "MsgOracleSetConfirm"), nil)
	cdc.RegisterConcrete(&MsgOracleSetUpdatedClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgOracleSetUpdatedClaim"), nil)

	cdc.RegisterConcrete(&MsgBridgeTokenClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeTokenClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToFxClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToFxClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgSendToExternalClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternalClaim"), nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConfirmBatch"), nil)

	cdc.RegisterConcrete(&MsgBridgeCallClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeCallClaim"), nil)
	cdc.RegisterConcrete(&MsgBridgeCallConfirm{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeCallConfirm"), nil)
	cdc.RegisterConcrete(&MsgBridgeCallResultClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeCallResultClaim"), nil)

	cdc.RegisterConcrete(&MsgUpdateParams{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateParams"), nil)
	cdc.RegisterConcrete(&MsgUpdateChainOracles{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateChainOracles"), nil)
}
