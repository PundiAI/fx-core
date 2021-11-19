package app

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/functionx/fx-core/crypto/ethsecp256k1"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	"github.com/palantir/stacktrace"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"runtime/debug"
)

const (
	secp256k1VerifyCost uint64 = 21000
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

func NewAnteHandlerWithEVM(
	ak evmtypes.AccountKeeper,
	bankKeeper evmtypes.BankKeeper,
	evmKeeper EVMKeeper,
	feeMarketKeeper evmtypes.FeeMarketKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler authsigning.SignModeHandler,
) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		defer Recover(ctx.Logger(), &err)

		txWithExtensions, ok := tx.(ante.HasExtensionOptionsTx)
		if ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 0 {
				switch typeURL := opts[0].GetTypeUrl(); typeURL {
				case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
					// handle as *evmtypes.MsgEthereumTx

					anteHandler = sdk.ChainAnteDecorators(
						NewEthSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
						ante.NewMempoolFeeDecorator(),
						ante.TxTimeoutHeightDecorator{},
						ante.NewValidateMemoDecorator(ak),
						NewEthValidateBasicDecorator(evmKeeper),
						NewEthSigVerificationDecorator(evmKeeper),
						NewEthAccountVerificationDecorator(ak, bankKeeper, evmKeeper),
						NewEthNonceVerificationDecorator(ak),
						NewEthGasConsumeDecorator(evmKeeper),
						NewCanTransferDecorator(evmKeeper, feeMarketKeeper),
						NewEthIncrementSenderSequenceDecorator(ak), // innermost AnteDecorator.
					)

				default:
					return ctx, stacktrace.Propagate(
						sdkerrors.Wrap(sdkerrors.ErrUnknownExtensionOptions, typeURL),
						"rejecting tx with unsupported extension option",
					)
				}

				return anteHandler(ctx, tx, sim)
			}
		}

		// handle as totally normal Cosmos SDK tx

		switch tx.(type) {
		case sdk.Tx:
			anteHandler = sdk.ChainAnteDecorators(
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
		default:
			return ctx, stacktrace.Propagate(
				sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx),
				"transaction is not an SDK tx",
			)
		}

		return anteHandler(ctx, tx, sim)
	}
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

func Recover(logger tmlog.Logger, err *error) {
	if r := recover(); r != nil {
		*err = sdkerrors.Wrapf(sdkerrors.ErrPanic, "%v", r)

		if e, ok := r.(error); ok {
			logger.Error(
				"ante handler panicked",
				"error", e,
				"stack trace", string(debug.Stack()),
			)
		} else {
			logger.Error(
				"ante handler panicked",
				"recover", fmt.Sprintf("%v", r),
			)
		}
	}
}

var _ ante.SignatureVerificationGasConsumer = DefaultSigVerificationGasConsumer

// DefaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func DefaultSigVerificationGasConsumer(
	meter sdk.GasMeter, sig txsigning.SignatureV2, params types.Params,
) error {
	// support for ethereum ECDSA secp256k1 keys
	_, ok := sig.PubKey.(*ethsecp256k1.PubKey)
	if ok {
		meter.ConsumeGas(secp256k1VerifyCost, "ante verify: eth_secp256k1")
		return nil
	}

	return ante.DefaultSigVerificationGasConsumer(meter, sig, params)
}
