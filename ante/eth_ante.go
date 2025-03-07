package ante

import (
	"bytes"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	evmante "github.com/evmos/ethermint/app/ante"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// Just copy, because the pendingTxListener property of the TxListenerDecorator is internal
// https://github.com/pundiai/ethermint/blob/fxcore/v0.22.x/app/ante/tx_listener.go

type TxListenerDecorator struct {
	pendingTxListener evmante.PendingTxListener
}

// newTxListenerDecorator creates a new TxListenerDecorator with the provided PendingTxListener.
// CONTRACT: must be put at the last of the chained decorators
func newTxListenerDecorator(pendingTxListener evmante.PendingTxListener) TxListenerDecorator {
	return TxListenerDecorator{pendingTxListener}
}

func (d TxListenerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}
	if ctx.IsCheckTx() && !simulate && d.pendingTxListener != nil {
		for _, msg := range tx.GetMsgs() {
			if ethTx, ok := msg.(*evmtypes.MsgEthereumTx); ok {
				d.pendingTxListener(ethTx.Hash())
			}
		}
	}
	return next(ctx, tx, simulate)
}

// CheckAndSetEthSenderNonce handles incrementing the sequence of the signer (i.e sender). If the transaction is a
// contract creation, the nonce will be incremented during the transaction execution and not within
// this AnteHandler decorator.
func CheckAndSetEthSenderNonce(
	ctx sdk.Context, tx sdk.Tx, ak evmtypes.AccountKeeper, unsafeUnOrderedTx bool, accountGetter evmante.AccountGetter, signer ethtypes.Signer,
) error {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		ethTx := msgEthTx.AsTransaction()

		// increase sequence of sender
		from := msgEthTx.GetFrom()
		acc := accountGetter(from)
		if acc == nil {
			return errorsmod.Wrapf(
				errortypes.ErrUnknownAddress,
				"account %s is nil", common.BytesToAddress(from.Bytes()),
			)
		}
		nonce := acc.GetSequence()

		if !unsafeUnOrderedTx {
			// we merged the nonce verification to nonce increment, so when tx includes multiple messages
			// with same sender, they'll be accepted.
			if ethTx.Nonce() != nonce {
				return errorsmod.Wrapf(
					errortypes.ErrInvalidSequence,
					"invalid nonce; got %d, expected %d", ethTx.Nonce(), nonce,
				)
			}
		}

		if err := acc.SetSequence(nonce + 1); err != nil {
			return errorsmod.Wrapf(err, "failed to set sequence to %d", acc.GetSequence()+1)
		}

		if acc.GetPubKey() == nil {
			if pubKey, err := EthPubkeyParse(ethTx, signer, from); err == nil {
				if err = acc.SetPubKey(pubKey); err != nil {
					return errorsmod.Wrapf(err, "failed to set pubkey to %s", pubKey.Address().String())
				}
			}
		}

		ak.SetAccount(ctx, acc)
	}

	return nil
}

func EthPubkeyParse(tx *ethtypes.Transaction, signer ethtypes.Signer, accAddr sdk.AccAddress) (pub cryptotypes.PubKey, err error) {
	v, r, s := tx.RawSignatureValues()
	switch tx.Type() {
	case ethtypes.LegacyTxType:
		if tx.Protected() {
			v = new(big.Int).Sub(v, new(big.Int).Mul(signer.ChainID(), big.NewInt(2)))
			v.Sub(v, big.NewInt(8))
		}
	case ethtypes.AccessListTxType, ethtypes.DynamicFeeTxType:
		v = new(big.Int).Add(v, big.NewInt(27))
	default:
		return nil, ethtypes.ErrTxTypeNotSupported
	}
	if v.BitLen() > 8 {
		return nil, ethtypes.ErrInvalidSig
	}
	vb := byte(v.Uint64() - 27)
	if !crypto.ValidateSignatureValues(vb, r, s, true) {
		return nil, ethtypes.ErrInvalidSig
	}
	rb, sb := r.Bytes(), s.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(rb):32], rb)
	copy(sig[64-len(sb):64], sb)
	sig[64] = vb

	pubKey, err := crypto.SigToPub(signer.Hash(tx).Bytes(), sig)
	if err != nil {
		return nil, err
	}
	ethPubkey := &ethsecp256k1.PubKey{Key: crypto.CompressPubkey(pubKey)}
	if !bytes.Equal(ethPubkey.Address().Bytes(), accAddr.Bytes()) {
		return nil, errortypes.ErrInvalidAddress
	}
	return ethPubkey, nil
}
