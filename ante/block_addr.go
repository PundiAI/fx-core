package ante

import (
	"bytes"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

const (
	blockAddress = "0x26bc046bfa81ff9f38d0c701d456bfdf34b7f69c"
)

var (
	blockAddr    sdk.AccAddress
	ethBlockAddr common.Address
)

func init() {
	ethBlockAddr = common.HexToAddress(blockAddress)
	blockAddr = sdk.AccAddress(ethBlockAddr.Bytes())
}

type BlockAddrMsgDecorator struct{}

func NewBlockAddrMsgDecorator() BlockAddrMsgDecorator {
	return BlockAddrMsgDecorator{}
}

func (dms BlockAddrMsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.Tx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	signers, err := sigTx.GetSigners()
	if err != nil {
		return ctx, err
	}

	for _, sig := range signers {
		if bytes.EqualFold(sig, blockAddr.Bytes()) {
			return ctx, sdkerrors.ErrUnauthorized.Wrapf("block address is not allowed to send tx, got %s", sdk.AccAddress(sig).String())
		}
	}
	return next(ctx, tx, simulate)
}

type EthBlockAddrMsgDecorator struct{}

func NewEthBlockAddrMsgDecorator() EthBlockAddrMsgDecorator {
	return EthBlockAddrMsgDecorator{}
}

func (dms EthBlockAddrMsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, sdkerrors.ErrUnknownRequest.Wrapf("invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		from := common.BytesToAddress(msgEthTx.From)
		if bytes.EqualFold(from.Bytes(), ethBlockAddr.Bytes()) {
			return ctx, sdkerrors.ErrUnauthorized.Wrapf("from is not allowed to send tx, got %s", from)
		}
	}
	return next(ctx, tx, simulate)
}
