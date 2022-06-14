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
	RegisterLegacyAminoCodec(ModuleCdc)
}

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateOracleBridger{},
		&MsgAddOracleDelegate{},
		&MsgEditOracle{},
		&MsgWithdrawReward{},

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
		&UpdateChainOraclesProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*ExternalClaim)(nil), nil)

	cdc.RegisterConcrete(&MsgCreateOracleBridger{}, fmt.Sprintf("%s/%s", ModuleName, "MsgCreateOracleBridger"), nil)
	cdc.RegisterConcrete(&MsgAddOracleDelegate{}, fmt.Sprintf("%s/%s", ModuleName, "MsgAddOracleDelegate"), nil)
	cdc.RegisterConcrete(&MsgEditOracle{}, fmt.Sprintf("%s/%s", ModuleName, "MsgEditOracle"), nil)
	cdc.RegisterConcrete(&MsgWithdrawReward{}, fmt.Sprintf("%s/%s", ModuleName, "MsgWithdrawReward"), nil)

	cdc.RegisterConcrete(&MsgOracleSetConfirm{}, fmt.Sprintf("%s/%s", ModuleName, "MsgOracleSetConfirm"), nil)
	cdc.RegisterConcrete(&MsgOracleSetUpdatedClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgOracleSetUpdatedClaim"), nil)

	cdc.RegisterConcrete(&MsgBridgeTokenClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgBridgeTokenClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToFxClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToFxClaim"), nil)

	cdc.RegisterConcrete(&MsgSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgCancelSendToExternal{}, fmt.Sprintf("%s/%s", ModuleName, "MsgCancelSendToExternal"), nil)
	cdc.RegisterConcrete(&MsgSendToExternalClaim{}, fmt.Sprintf("%s/%s", ModuleName, "MsgSendToExternalClaim"), nil)

	cdc.RegisterConcrete(&MsgRequestBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgRequestBatch"), nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, fmt.Sprintf("%s/%s", ModuleName, "MsgConfirmBatch"), nil)

	// register Proposal
	cdc.RegisterConcrete(&UpdateChainOraclesProposal{}, fmt.Sprintf("%s/%s", ModuleName, "UpdateChainOraclesProposal"), nil)
}
