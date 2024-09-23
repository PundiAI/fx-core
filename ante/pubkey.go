package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type PubKeyDecorator struct {
	ak ante.AccountKeeper
}

func NewPubKeyDecorator(ak ante.AccountKeeper) PubKeyDecorator {
	return PubKeyDecorator{ak: ak}
}

func (pkd PubKeyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.ErrTxDecode.Wrapf("invalid tx type")
	}

	pubkeys, err := sigTx.GetPubKeys()
	if err != nil {
		return ctx, err
	}
	signers, err := sigTx.GetSigners()
	if err != nil {
		return ctx, err
	}
	for i := range pubkeys {
		if err = checkPubKeyDisabled(ctx, pkd.ak, signers[i]); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

type EthPubKeyDecorator struct {
	ak ante.AccountKeeper
}

func NewEthPubKeyDecorator(ak ante.AccountKeeper) EthPubKeyDecorator {
	return EthPubKeyDecorator{ak: ak}
}

func (epkd EthPubKeyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, sdkerrors.ErrUnknownRequest.Wrapf("invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		from := common.BytesToAddress(msgEthTx.From)
		if err := checkPubKeyDisabled(ctx, epkd.ak, from.Bytes()); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

func checkPubKeyDisabled(ctx sdk.Context, ak ante.AccountKeeper, address sdk.AccAddress) error {
	signerAcc, err := ante.GetSignerAcc(ctx, ak, address)
	if err != nil {
		return err
	}
	pubKey := signerAcc.GetPubKey()
	if pubKey == nil {
		return nil
	}
	for _, b := range pubKey.Bytes() {
		if b != 0xff {
			return nil
		}
	}
	return sdkerrors.ErrInvalidAddress.Wrap("account disabled")
}
