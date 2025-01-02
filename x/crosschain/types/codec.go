package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetOrchestratorAddress{},
		&MsgAddOracleDeposit{},

		&MsgBondedOracle{},
		&MsgAddDelegate{},
		&MsgReDelegate{},
		&MsgEditBridger{},
		&MsgWithdrawReward{},

		&MsgOracleSetConfirm{},

		&MsgSendToExternal{},
		&MsgCancelSendToExternal{},
		&MsgIncreaseBridgeFee{},

		&MsgRequestBatch{},
		&MsgConfirmBatch{},

		&MsgBridgeCallConfirm{},

		&MsgUpdateParams{},
		&MsgUpdateChainOracles{},

		&MsgClaim{},
		&MsgSendToExternalClaim{},
		&MsgSendToFxClaim{},
		&MsgBridgeCallClaim{},
		&MsgBridgeTokenClaim{},
		&MsgOracleSetUpdatedClaim{},
		&MsgBridgeCallResultClaim{},

		&MsgConfirm{},
		&MsgConfirmBatch{},
		&MsgOracleSetConfirm{},
		&MsgBridgeCallConfirm{},
	)

	registry.RegisterImplementations(
		(*govv1betal.Content)(nil),
		&InitCrossChainParamsProposal{},
		&UpdateChainOraclesProposal{},
	)

	registry.RegisterImplementations(
		(*ExternalClaim)(nil),
		&MsgSendToExternalClaim{},
		&MsgSendToFxClaim{},
		&MsgBridgeCallClaim{},
		&MsgBridgeTokenClaim{},
		&MsgOracleSetUpdatedClaim{},
		&MsgBridgeCallResultClaim{},
	)

	registry.RegisterImplementations(
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
	cdc.RegisterConcrete(&MsgCancelSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgCancelSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgIncreaseBridgeFee{}, fmt.Sprintf("%s/%s", ModuleName, "MsgIncreaseBridgeFee"), nil)
	cdc.RegisterConcrete(&MsgSendToExternalClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternalClaim"), nil)

	cdc.RegisterConcrete(&MsgRequestBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRequestBatch"), nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConfirmBatch"), nil)

	cdc.RegisterConcrete(&MsgBridgeCallClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeCallClaim"), nil)
	cdc.RegisterConcrete(&MsgBridgeCallConfirm{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeCallConfirm"), nil)
	cdc.RegisterConcrete(&MsgBridgeCallResultClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeCallResultClaim"), nil)

	// register Proposal
	cdc.RegisterConcrete(&UpdateChainOraclesProposal{}, fmt.Sprintf("%s/%s", ModuleName, "UpdateChainOraclesProposal"), nil)
	cdc.RegisterConcrete(&InitCrossChainParamsProposal{}, fmt.Sprintf("%s/%s", ModuleName, "InitCrossChainParamsProposal"), nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateParams"), nil)
	cdc.RegisterConcrete(&MsgUpdateChainOracles{}, fmt.Sprintf("%s/%s", ModuleName, "MsgUpdateChainOracles"), nil)
}
