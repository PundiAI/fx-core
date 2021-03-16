package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak ante.AccountKeeper, bankKeeper types.BankKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewRejectFeeGranterDecorator(),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, bankKeeper),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

// RejectExtensionOptionsDecorator is an AnteDecorator that rejects all extension
// options which can optionally be included in protobuf transactions. Users that
// need extension options should create a custom AnteHandler chain that handles
// needed extension options properly and rejects unknown ones.
type RejectExtensionOptionsDecorator struct{}

// NewRejectExtensionOptionsDecorator creates a new RejectExtensionOptionsDecorator
func NewRejectExtensionOptionsDecorator() RejectExtensionOptionsDecorator {
	return RejectExtensionOptionsDecorator{}
}

var _ sdk.AnteDecorator = RejectExtensionOptionsDecorator{}

// AnteHandle implements the AnteDecorator.AnteHandle method
func (r RejectExtensionOptionsDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if hasExtOptsTx, ok := tx.(ante.HasExtensionOptionsTx); ok {
		if len(hasExtOptsTx.GetExtensionOptions()) != 0 {
			return ctx, sdkerrors.ErrUnknownExtensionOptions
		}
		if len(hasExtOptsTx.GetNonCriticalExtensionOptions()) != 0 {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown non critical extension options")
		}
	}
	return next(ctx, tx, simulate)
}
