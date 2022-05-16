package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	evmtypes "github.com/functionx/fx-core/x/evm/types"

	ethv0 "github.com/functionx/fx-core/app/ante/eth/v0"
	ethv1 "github.com/functionx/fx-core/app/ante/eth/v1"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper, EVM Keeper and Fee Market Keeper.
type HandlerOptions struct {
	AccountKeeper        evmtypes.AccountKeeper
	BankKeeper           evmtypes.BankKeeper
	EvmKeeper            ethv1.EVMKeeper
	EvmKeeperV0          ethv0.EVMKeeper
	FeeMarketKeeperV0    ethv0.FeeMarketKeeper
	SignModeHandler      authsigning.SignModeHandler
	SigGasConsumer       ante.SignatureVerificationGasConsumer
	MaxTxGasWanted       uint64
	BypassMinFeeMsgTypes []string
}

func (options HandlerOptions) Validate() error {
	if options.AccountKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.EvmKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "evm keeper is required for AnteHandler")
	}
	return nil
}

func newEthV0AnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethv0.NewEthSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewMempoolFeeDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ethv0.NewEthValidateBasicDecorator(options.EvmKeeperV0),
		ethv0.NewEthSigVerificationDecorator(options.EvmKeeperV0),
		ethv0.NewEthAccountVerificationDecorator(options.AccountKeeper, options.BankKeeper, options.EvmKeeperV0),
		ethv0.NewEthNonceVerificationDecorator(options.AccountKeeper),
		ethv0.NewEthGasConsumeDecorator(options.EvmKeeperV0),
		ethv0.NewCanTransferDecorator(options.EvmKeeperV0, options.FeeMarketKeeperV0),
		ethv0.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper), // innermost AnteDecorator.
	)
}

func newEthV1AnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethv1.NewEthSetUpContextDecorator(options.EvmKeeper), // outermost AnteDecorator. SetUpContext must be called first
		ethv1.NewEthMempoolFeeDecorator(options.EvmKeeper),
		ethv1.NewEthValidateBasicDecorator(options.EvmKeeper),
		ethv1.NewEthSigVerificationDecorator(options.EvmKeeper),
		ethv1.NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper),
		ethv1.NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted),
		ethv1.NewCanTransferDecorator(options.EvmKeeper),
		ethv1.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper), // innermost AnteDecorator.
	)
}

func newNormalTxAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		NewMempoolFeeDecorator(options.BypassMinFeeMsgTypes),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewRejectFeeGranterDecorator(),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
	)
}

func NewNormalTxAnteHandlerEip712(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethv1.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		// NOTE: extensions option decorator removed
		//NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewRejectFeeGranterDecorator(),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		// Note: signature verification uses EIP instead of the cosmos signature validator
		NewEip712SigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
	)
}
