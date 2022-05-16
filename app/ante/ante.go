package ante

import (
	"fmt"
	"runtime/debug"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	tmlog "github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/crypto/ethsecp256k1"
)

func NewAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		defer Recover(ctx.Logger(), &err)

		txWithExtensions, ok := tx.(ante.HasExtensionOptionsTx)
		if ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 0 {
				typeURL := opts[0].GetTypeUrl()
				if ctx.BlockHeight() >= fxtypes.EvmV0SupportBlock() &&
					ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() &&
					typeURL == "/ethermint.evm.v1.ExtensionOptionsEthereumTx" {
					//evm v0
					anteHandler = newEthV0AnteHandler(options)
				} else if ctx.BlockHeight() >= fxtypes.EvmV1SupportBlock() &&
					typeURL == "/fx.ethereum.evm.v1.ExtensionOptionsEthereumTx" {
					//evm v1
					anteHandler = newEthV1AnteHandler(options)
				} else {
					//unsupported
					return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownExtensionOptions, "rejecting tx with unsupported extension option: %s", typeURL)
				}

				return anteHandler(ctx, tx, sim)
			}
		}

		// handle as totally normal Cosmos SDK tx

		switch tx.(type) {
		case sdk.Tx:
			anteHandler = newNormalTxAnteHandler(options)
		default:
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}
}

func Recover(logger tmlog.Logger, err *error) {
	if r := recover(); r != nil {
		//*err = sdkerrors.Wrapf(sdkerrors.ErrPanic, "%v", r)

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

const (
	secp256k1VerifyCost uint64 = 21000
)

var _ SignatureVerificationGasConsumer = DefaultSigVerificationGasConsumer

// DefaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func DefaultSigVerificationGasConsumer(ctx sdk.Context, sig txsigning.SignatureV2, params types.Params) error {
	pubkey := sig.PubKey
	meter := ctx.GasMeter()
	switch pubkey := pubkey.(type) {
	case *ethsecp256k1.PubKey: // support for ethereum ECDSA secp256k1 keys
		ethsecp256k1VerifyCost := params.SigVerifyCostSecp256k1
		if ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() {
			ethsecp256k1VerifyCost = secp256k1VerifyCost
		}
		meter.ConsumeGas(ethsecp256k1VerifyCost, "ante verify: eth_secp256k1")
		return nil

	case *ed25519.PubKey:
		meter.ConsumeGas(params.SigVerifyCostED25519, "ante verify: ed25519")
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "ED25519 public keys are unsupported")

	case *secp256k1.PubKey:
		meter.ConsumeGas(params.SigVerifyCostSecp256k1, "ante verify: secp256k1")
		return nil

	case multisig.PubKey:
		multisignature, ok := sig.Data.(*txsigning.MultiSignatureData)
		if !ok {
			return fmt.Errorf("expected %T, got, %T", &txsigning.MultiSignatureData{}, sig.Data)
		}
		err := ante.ConsumeMultisignatureVerificationGas(meter, multisignature, pubkey, params, sig.Sequence)
		if err != nil {
			return err
		}
		return nil

	default:
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "unrecognized public key type: %T", pubkey)
	}
}
