package ante

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	tmstrings "github.com/tendermint/tendermint/libs/strings"
)

var (
	// simulation signature values used to estimate gas consumption
	key                   = make([]byte, secp256k1.PubKeySize)
	simSecp256k1Pubkey    = &secp256k1.PubKey{Key: key}
	simEthSecp256k1Pubkey = &ethsecp256k1.PubKey{Key: key}

	_ authsigning.SigVerifiableTx = (*legacytx.StdTx)(nil) // assert StdTx implements SigVerifiableTx
)

func init() {
	// This decodes a valid hex string into a sepc256k1Pubkey for use in transaction simulation
	bz, _ := hex.DecodeString("035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A")
	copy(key, bz)
	simSecp256k1Pubkey.Key = key
	simEthSecp256k1Pubkey.Key = key
}

// SetPubKeyDecorator sets PubKeys in context for any signer which does not already have pubkey set
// PubKeys must be set in context for all signers before any other sigverify decorators run
// CONTRACT: Tx must implement SigVerifiableTx interface
type SetPubKeyDecorator struct {
	ak ante.AccountKeeper
}

func NewSetPubKeyDecorator(ak ante.AccountKeeper) SetPubKeyDecorator {
	return SetPubKeyDecorator{
		ak: ak,
	}
}

func (spkd SetPubKeyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	pubkeys, err := sigTx.GetPubKeys()
	if err != nil {
		return ctx, err
	}
	signers := sigTx.GetSigners()

	for i, pk := range pubkeys {
		// PublicKey was omitted from slice since it has already been set in context
		if pk == nil {
			if !simulate {
				continue
			}
			pk = simEthSecp256k1Pubkey
		}
		// Only make check if simulate=false
		if !simulate && !bytes.Equal(pk.Address(), signers[i]) {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey,
				"pubKey does not match signer address %s with signer index: %d", signers[i], i)
		}

		acc, err := ante.GetSignerAcc(ctx, spkd.ak, signers[i])
		if err != nil {
			return ctx, err
		}
		// account already has pubkey set,no need to reset
		if acc.GetPubKey() != nil {
			continue
		}

		err = acc.SetPubKey(pk)
		if err != nil {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
		}
		spkd.ak.SetAccount(ctx, acc)
	}

	// Also emit the following events, so that txs can be indexed by these
	// indices:
	// - signature (via `tx.signature='<sig_as_base64>'`),
	// - concat(address,"/",sequence) (via `tx.acc_seq='cosmos1abc...def/42'`).
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}

	var events sdk.Events
	for i, sig := range sigs {
		events = append(events, sdk.NewEvent(sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyAccountSequence, fmt.Sprintf("%s/%d", signers[i], sig.Sequence)),
		))

		sigBzs, err := signatureDataToBz(sig.Data)
		if err != nil {
			return ctx, err
		}
		for _, sigBz := range sigBzs {
			events = append(events, sdk.NewEvent(sdk.EventTypeTx,
				sdk.NewAttribute(sdk.AttributeKeySignature, base64.StdEncoding.EncodeToString(sigBz)),
			))
		}
	}

	ctx.EventManager().EmitEvents(events)

	return next(ctx, tx, simulate)
}

// signatureDataToBz converts a SignatureData into raw bytes signature.
// For SingleSignatureData, it returns the signature raw bytes.
// For MultiSignatureData, it returns an array of all individual signatures,
// as well as the aggregated signature.
func signatureDataToBz(data signing.SignatureData) ([][]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("got empty SignatureData")
	}

	switch data := data.(type) {
	case *signing.SingleSignatureData:
		return [][]byte{data.Signature}, nil
	case *signing.MultiSignatureData:
		sigs := [][]byte{}
		var err error

		for _, d := range data.Signatures {
			nestedSigs, err := signatureDataToBz(d)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, nestedSigs...)
		}

		multisig := cryptotypes.MultiSignature{
			Signatures: sigs,
		}
		aggregatedSig, err := multisig.Marshal()
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, aggregatedSig)

		return sigs, nil
	default:
		return nil, sdkerrors.ErrInvalidType.Wrapf("unexpected signature data type %T", data)
	}
}

// SigGasConsumeDecorator Consume parameter-defined amount of gas for each signature according to the passed-in SignatureVerificationGasConsumer function
// before calling the next AnteHandler
// CONTRACT: Pubkeys are set in context for all signers before this decorator runs
// CONTRACT: Tx must implement SigVerifiableTx interface
type SigGasConsumeDecorator struct {
	ak             ante.AccountKeeper
	sigGasConsumer ante.SignatureVerificationGasConsumer
}

func NewSigGasConsumeDecorator(ak ante.AccountKeeper, sigGasConsumer ante.SignatureVerificationGasConsumer) SigGasConsumeDecorator {
	return SigGasConsumeDecorator{
		ak:             ak,
		sigGasConsumer: sigGasConsumer,
	}
}

func (sgcd SigGasConsumeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	params := sgcd.ak.GetParams(ctx)
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	signerAddrs := sigTx.GetSigners()

	for i, sig := range sigs {
		signerAcc, err := ante.GetSignerAcc(ctx, sgcd.ak, signerAddrs[i])
		if err != nil {
			return ctx, err
		}

		pubKey := signerAcc.GetPubKey()

		// In simulate mode the transaction comes with no signatures, thus if the
		// account's pubkey is nil, both signature verification and gasKVStore.Set()
		// shall consume the largest amount, i.e. it takes more gas to verify
		// secp256k1 keys than ed25519 ones.
		if simulate && pubKey == nil {
			pubKey = simSecp256k1Pubkey
		}

		// make a SignatureV2 with PubKey filled in from above
		sig = signing.SignatureV2{
			PubKey:   pubKey,
			Data:     sig.Data,
			Sequence: sig.Sequence,
		}

		if err = sgcd.sigGasConsumer(ctx.GasMeter(), sig, params); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

type MsgInterceptDecorator struct {
	InterceptMsgTypes map[int64][]string
}

func NewMsgInterceptDecorator(interceptMsgTypes map[int64][]string) MsgInterceptDecorator {
	return MsgInterceptDecorator{InterceptMsgTypes: interceptMsgTypes}
}

func (mid MsgInterceptDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	currentHeight := ctx.BlockHeight()

	validateTypes := make([]string, 0, 10)
	for supportHeight, types := range mid.InterceptMsgTypes {
		//before support height, intercept msg
		if currentHeight < supportHeight {
			validateTypes = append(validateTypes, types...)
		}
	}

	for _, m := range tx.GetMsgs() {
		var msgTypeURL string
		switch msg := m.(type) {
		case *govtypes.MsgSubmitProposal:
			msgTypeURL = msg.Content.GetTypeUrl()
		default:
			msgTypeURL = sdk.MsgTypeURL(m)
		}
		if tmstrings.StringInSlice(msgTypeURL, validateTypes) {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "msg type %s not support", msgTypeURL)
		}
	}

	return next(ctx, tx, simulate)
}
