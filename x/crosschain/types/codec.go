package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.NewLegacyAmino()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),

		&MsgSetOrchestratorAddress{},
		&MsgAddOracleDeposit{},
		&MsgOracleSetConfirm{},
		&MsgOracleSetUpdatedClaim{},

		&MsgBridgeTokenClaim{},

		&MsgSendToFxClaim{},

		&MsgSendToExternal{},
		&MsgCancelSendToExternal{},
		&MsgSendToExternalClaim{},

		&MsgRequestBatch{},
		&MsgConfirmBatch{},
	)

	registry.RegisterInterface(
		"gravity.v1beta1.ExternalClaim",
		(*ExternalClaim)(nil),
		&MsgSendToExternalClaim{},
		&MsgSendToFxClaim{},
		&MsgBridgeTokenClaim{},
		&MsgOracleSetUpdatedClaim{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&InitCrossChainParamsProposal{},
		&UpdateChainOraclesProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*ExternalClaim)(nil), nil)

	cdc.RegisterConcrete(&MsgSetOrchestratorAddress{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSetOrchestratorAddress"), nil)
	cdc.RegisterConcrete(&MsgAddOracleDeposit{}, fmt.Sprintf("%s/%s", ModuleName, "MsgAddOracleDeposit"), nil)
	cdc.RegisterConcrete(&MsgOracleSetConfirm{}, fmt.Sprintf("%s/%s", ModuleName, "MsgOracleSetConfirm"), nil)
	cdc.RegisterConcrete(&MsgOracleSetUpdatedClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgOracleSetUpdatedClaim"), nil)

	cdc.RegisterConcrete(&MsgBridgeTokenClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeTokenClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToFxClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToFxClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgCancelSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgCancelSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgSendToExternalClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternalClaim"), nil)

	cdc.RegisterConcrete(&MsgRequestBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRequestBatch"), nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConfirmBatch"), nil)

	cdc.RegisterConcrete(&OracleSet{}, fmt.Sprintf("%s/%s", ModuleName, "OracleSet"), nil)
	cdc.RegisterConcrete(&OutgoingTxBatch{}, fmt.Sprintf("%s/%s", ModuleName, "OutgoingTxBatch"), nil)
	cdc.RegisterConcrete(&OutgoingTransferTx{}, fmt.Sprintf("%s/%s", ModuleName, "OutgoingTransferTx"), nil)
	cdc.RegisterConcrete(&ExternalToken{}, fmt.Sprintf("%s/%s", ModuleName, "ExternalToken"), nil)
	cdc.RegisterConcrete(&IDSet{}, fmt.Sprintf("%s/%s", ModuleName, "IDSet"), nil)
	cdc.RegisterConcrete(&Attestation{}, fmt.Sprintf("%s/%s", ModuleName, "Attestation"), nil)

	// register Proposal
	cdc.RegisterConcrete(&InitCrossChainParamsProposal{}, fmt.Sprintf("%s/%s", ModuleName, "InitCrossChainParamsProposal"), nil)
	cdc.RegisterConcrete(&UpdateChainOraclesProposal{}, fmt.Sprintf("%s/%s", ModuleName, "UpdateChainOraclesProposal"), nil)
}
