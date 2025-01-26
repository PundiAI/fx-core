package legacy

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/pundiai/fx-core/v8/types/legacy/gravity"
)

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetOrchestratorAddress{},
		&MsgAddOracleDeposit{},
		&MsgCancelSendToExternal{},
		&MsgIncreaseBridgeFee{},
		&MsgRequestBatch{},

		&MsgValsetConfirm{},
		&MsgSendToEth{},
		&MsgRequestBatch{},
		&MsgConfirmBatch{},
		&MsgDepositClaim{},
		&MsgWithdrawClaim{},
		&MsgSetOrchestratorAddress{},
		&MsgCancelSendToEth{},
		&MsgValsetUpdatedClaim{},
		&MsgFxOriginatedTokenClaim{},

		&MsgGrantPrivilege{},
		&MsgEditConsensusPubKey{},

		&MsgUpdateFXParams{},
		&MsgUpdateEGFParams{},
		&MsgUpdateParams{},

		&MsgConvertERC20{},
		&MsgConvertDenom{},
		&MsgRegisterCoin{},
		&MsgRegisterERC20{},
		&MsgUpdateDenomAlias{},

		&MsgTransfer{},

		&gravity.MsgRequestBatch{},
		&gravity.MsgSetOrchestratorAddress{},
	)

	registry.RegisterImplementations(
		(*govv1betal.Content)(nil),
		&InitEvmParamsProposal{},

		&RegisterCoinProposal{},
		&RegisterERC20Proposal{},
		&ToggleTokenConversionProposal{},
		&UpdateDenomAliasProposal{},

		&InitCrossChainParamsProposal{},
		&UpdateChainOraclesProposal{},
	)
}

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&gravity.MsgSetOrchestratorAddress{}, "gravity/MsgSetOrchestratorAddress", nil)
	cdc.RegisterConcrete(&MsgValsetConfirm{}, "gravity/MsgValsetConfirm", nil)
	cdc.RegisterConcrete(&MsgSendToEth{}, "gravity/MsgSendToEth", nil)
	cdc.RegisterConcrete(&gravity.MsgRequestBatch{}, "gravity/MsgRequestBatch", nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, "gravity/MsgConfirmBatch", nil)
	cdc.RegisterConcrete(&MsgDepositClaim{}, "gravity/MsgDepositClaim", nil)
	cdc.RegisterConcrete(&MsgWithdrawClaim{}, "gravity/MsgWithdrawClaim", nil)
	cdc.RegisterConcrete(&MsgCancelSendToEth{}, "gravity/MsgCancelSendToEth", nil)
	cdc.RegisterConcrete(&MsgValsetUpdatedClaim{}, "gravity/MsgValsetUpdatedClaim", nil)
	cdc.RegisterConcrete(&MsgFxOriginatedTokenClaim{}, "gravity/MsgFxOriginatedTokenClaim", nil)

	cdc.RegisterConcrete(&MsgGrantPrivilege{}, "staking/MsgGrantPrivilege", nil)
	cdc.RegisterConcrete(&MsgEditConsensusPubKey{}, "staking/MsgEditConsensusPubKey", nil)

	cdc.RegisterConcrete(&MsgUpdateFXParams{}, "gov/MsgUpdateFXParams", nil)
	cdc.RegisterConcrete(&MsgUpdateEGFParams{}, "gov/MsgUpdateEGFParams", nil)

	cdc.RegisterConcrete(&MsgConvertERC20{}, "erc20/MsgConvertERC20", nil)
	cdc.RegisterConcrete(&MsgConvertDenom{}, "erc20/MsgConvertDenom", nil)
	cdc.RegisterConcrete(&MsgRegisterCoin{}, "erc20/MsgRegisterCoin", nil)
	cdc.RegisterConcrete(&MsgRegisterERC20{}, "erc20/MsgRegisterERC20", nil)
	cdc.RegisterConcrete(&MsgUpdateDenomAlias{}, "erc20/MsgUpdateDenomAlias", nil)

	cdc.RegisterConcrete(&RegisterCoinProposal{}, "erc20/RegisterCoinProposal", nil)
	cdc.RegisterConcrete(&RegisterERC20Proposal{}, "erc20/RegisterERC20Proposal", nil)
	cdc.RegisterConcrete(&ToggleTokenConversionProposal{}, "erc20/ToggleTokenConversionProposal", nil)
	cdc.RegisterConcrete(&UpdateDenomAliasProposal{}, "erc20/UpdateDenomAliasProposal", nil)

	cdc.RegisterConcrete(&UpdateChainOraclesProposal{}, "crosschain/UpdateChainOraclesProposal", nil)
	cdc.RegisterConcrete(&InitCrossChainParamsProposal{}, "crosschain/InitCrossChainParamsProposal", nil)

	cdc.RegisterConcrete(&MsgCancelSendToExternal{}, "crosschain/MsgCancelSendToExternal", nil)
	cdc.RegisterConcrete(&MsgIncreaseBridgeFee{}, "crosschain/MsgIncreaseBridgeFee", nil)
	cdc.RegisterConcrete(&MsgRequestBatch{}, "crosschain/MsgRequestBatch", nil)
	cdc.RegisterConcrete(&MsgSetOrchestratorAddress{}, "crosschain/MsgSetOrchestratorAddress", nil)
	cdc.RegisterConcrete(&MsgAddOracleDeposit{}, "crosschain/MsgAddOracleDeposit", nil)

	cdc.RegisterConcrete(&MsgTransfer{}, "fxtransfer/MsgTransfer", nil)
}
