package crosschain

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
)

var (
	// BridgeCoinAmountMethod query the amount of bridge coin
	BridgeCoinAmountMethod = GetABI().Methods[BridgeCoinAmountMethodName]

	// CancelSendToExternalMethod cancel send to external tx
	CancelSendToExternalMethod = GetABI().Methods[CancelSendToExternalMethodName]

	// FIP20CrossChainMethod cross chain with FIP20 token, only for FIP20 token
	// Deprecated: use CrossChainMethod instead
	FIP20CrossChainMethod = GetABI().Methods[FIP20CrossChainMethodName]

	// CrossChainMethod cross chain with FIP20 token
	CrossChainMethod = GetABI().Methods[CrossChainMethodName]

	// IncreaseBridgeFeeMethod increase bridge fee
	IncreaseBridgeFeeMethod = GetABI().Methods[IncreaseBridgeFeeMethodName]

	// BridgeCallMethod bridge call other chain
	BridgeCallMethod = GetABI().Methods[BridgeCallMethodName]
)

type BridgeCoinAmountArgs struct {
	Token  common.Address `abi:"_token"`
	Target [32]byte       `abi:"_target"`
}

// Validate validates the args
func (args *BridgeCoinAmountArgs) Validate() error {
	if args.Target == [32]byte{} {
		return errors.New("empty target")
	}
	return nil
}

type CancelSendToExternalArgs struct {
	Chain string   `abi:"_chain"`
	TxID  *big.Int `abi:"_txID"`
}

// Validate validates the args
func (args *CancelSendToExternalArgs) Validate() error {
	if err := crosschaintypes.ValidateModuleName(args.Chain); err != nil {
		return err
	}
	if args.TxID == nil || args.TxID.Sign() <= 0 {
		return errors.New("invalid tx id")
	}
	return nil
}

type FIP20CrossChainArgs struct {
	Sender  common.Address `abi:"_sender"`
	Receipt string         `abi:"_receipt"`
	Amount  *big.Int       `abi:"_amount"`
	Fee     *big.Int       `abi:"_fee"`
	Target  [32]byte       `abi:"_target"`
	Memo    string         `abi:"_memo"`
}

// Validate validates the args
func (args *FIP20CrossChainArgs) Validate() error {
	if args.Receipt == "" {
		return errors.New("empty receipt")
	}
	if args.Amount == nil || args.Amount.Sign() <= 0 {
		return errors.New("invalid amount")
	}
	if args.Fee == nil || args.Fee.Sign() < 0 {
		return errors.New("invalid fee")
	}
	if args.Target == [32]byte{} {
		return errors.New("empty target")
	}
	return nil
}

type CrossChainArgs struct {
	Token   common.Address `abi:"_token"`
	Receipt string         `abi:"_receipt"`
	Amount  *big.Int       `abi:"_amount"`
	Fee     *big.Int       `abi:"_fee"`
	Target  [32]byte       `abi:"_target"`
	Memo    string         `abi:"_memo"`
}

// Validate validates the args
func (args *CrossChainArgs) Validate() error {
	if args.Receipt == "" {
		return errors.New("empty receipt")
	}
	if args.Amount == nil || args.Amount.Sign() <= 0 {
		return errors.New("invalid amount")
	}
	if args.Fee == nil || args.Fee.Sign() < 0 {
		return errors.New("invalid fee")
	}
	if args.Target == [32]byte{} {
		return errors.New("empty target")
	}
	return nil
}

type IncreaseBridgeFeeArgs struct {
	Chain string         `abi:"_chain"`
	TxID  *big.Int       `abi:"_txID"`
	Token common.Address `abi:"_token"`
	Fee   *big.Int       `abi:"_fee"`
}

// Validate validates the args
func (args *IncreaseBridgeFeeArgs) Validate() error {
	if err := crosschaintypes.ValidateModuleName(args.Chain); err != nil {
		return err
	}

	if args.TxID == nil || args.TxID.Sign() <= 0 {
		return errors.New("invalid tx id")
	}
	if args.Fee == nil || args.Fee.Sign() <= 0 {
		return errors.New("invalid add bridge fee")
	}
	return nil
}

type BridgeCallArgs struct {
	DstChain string           `abi:"_dstChain"`
	Receiver common.Address   `abi:"_receiver"`
	Tokens   []common.Address `abi:"_tokens"`
	Amounts  []*big.Int       `abi:"_amounts"`
	To       common.Address   `abi:"_to"`
	Data     []byte           `abi:"_data"`
	Value    *big.Int         `abi:"_value"`
	Memo     []byte           `abi:"_memo"`
}

// Validate validates the args
func (args *BridgeCallArgs) Validate() error {
	if err := crosschaintypes.ValidateModuleName(args.DstChain); err != nil {
		return err
	}
	if args.Value.Sign() != 0 {
		return errors.New("value must be zero")
	}
	if len(args.Tokens) != len(args.Amounts) {
		return errors.New("tokens and amounts do not match")
	}
	if len(args.Amounts) > 0 {
		if bytes.Equal(args.Receiver.Bytes(), common.Address{}.Bytes()) {
			return errors.New("receiver cannot be empty")
		}
	}
	if len(args.Tokens) == 0 && len(args.Data) == 0 {
		return errors.New("tokens and data cannot be empty at the same time")
	}
	return nil
}
