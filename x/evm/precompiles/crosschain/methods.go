package crosschain

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
)

// BridgeCoinAmountMethod query the amount of bridge coin
var BridgeCoinAmountMethod = abi.NewMethod(
	BridgeCoinAmountMethodName,
	BridgeCoinAmountMethodName,
	abi.Function, "view", false, false,
	abi.Arguments{
		abi.Argument{Name: "_token", Type: contract.TypeAddress},
		abi.Argument{Name: "_target", Type: contract.TypeBytes32},
	},
	abi.Arguments{
		abi.Argument{Name: "_amount", Type: contract.TypeUint256},
	},
)

var (
	// CancelSendToExternalMethod cancel send to external tx
	CancelSendToExternalMethod = abi.NewMethod(
		CancelSendToExternalMethodName,
		CancelSendToExternalMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_chain", Type: contract.TypeString},
			abi.Argument{Name: "_txID", Type: contract.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: contract.TypeBool},
		},
	)

	// FIP20CrossChainMethod cross chain with FIP20 token, only for FIP20 token
	// Deprecated: use CrossChainMethod instead
	FIP20CrossChainMethod = abi.NewMethod(
		FIP20CrossChainMethodName,
		FIP20CrossChainMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_sender", Type: contract.TypeAddress},
			abi.Argument{Name: "_receipt", Type: contract.TypeString},
			abi.Argument{Name: "_amount", Type: contract.TypeUint256},
			abi.Argument{Name: "_fee", Type: contract.TypeUint256},
			abi.Argument{Name: "_target", Type: contract.TypeBytes32},
			abi.Argument{Name: "_memo", Type: contract.TypeString},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: contract.TypeBool},
		},
	)
	// CrossChainMethod cross chain with FIP20 token
	CrossChainMethod = abi.NewMethod(
		CrossChainMethodName,
		CrossChainMethodName,
		abi.Function, "payable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_token", Type: contract.TypeAddress},
			abi.Argument{Name: "_receipt", Type: contract.TypeString},
			abi.Argument{Name: "_amount", Type: contract.TypeUint256},
			abi.Argument{Name: "_fee", Type: contract.TypeUint256},
			abi.Argument{Name: "_target", Type: contract.TypeBytes32},
			abi.Argument{Name: "_memo", Type: contract.TypeString},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: contract.TypeBool},
		},
	)

	// IncreaseBridgeFeeMethod increase bridge fee
	IncreaseBridgeFeeMethod = abi.NewMethod(
		IncreaseBridgeFeeMethodName,
		IncreaseBridgeFeeMethodName,
		abi.Function, "payable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_chain", Type: contract.TypeString},
			abi.Argument{Name: "_txID", Type: contract.TypeUint256},
			abi.Argument{Name: "_token", Type: contract.TypeAddress},
			abi.Argument{Name: "_fee", Type: contract.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: contract.TypeBool},
		},
	)

	// BridgeCallMethod bridge call other chain
	BridgeCallMethod = abi.NewMethod(
		BridgeCallMethodName,
		BridgeCallMethodName,
		abi.Function, "payable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_dstChainId", Type: contract.TypeString},
			abi.Argument{Name: "_gasLimit", Type: contract.TypeUint256},
			abi.Argument{Name: "_receiver", Type: contract.TypeAddress},
			abi.Argument{Name: "_to", Type: contract.TypeAddress},
			abi.Argument{Name: "_tokens", Type: contract.TypeAddressArray},
			abi.Argument{Name: "_amounts", Type: contract.TypeUint256Array},
			abi.Argument{Name: "_message", Type: contract.TypeBytes},
			abi.Argument{Name: "_value", Type: contract.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: contract.TypeBool},
		},
	)
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
	DstChainId string           `abi:"_dstChainId"`
	GasLimit   *big.Int         `abi:"_gasLimit"`
	Receiver   common.Address   `abi:"_receiver"`
	To         common.Address   `abi:"_to"`
	Tokens     []common.Address `abi:"_tokens"`
	Amounts    []*big.Int       `abi:"_amounts"`
	Message    []byte           `abi:"_message"`
	Value      *big.Int         `abi:"_value"`
}

// Validate validates the args
func (args *BridgeCallArgs) Validate() error {
	if len(args.DstChainId) == 0 {
		return errors.New("empty dst chain id")
	}
	if args.GasLimit.Cmp(big.NewInt(0)) == -1 || !args.GasLimit.IsUint64() {
		return errors.New("invalid gas limit")
	}
	if args.Value.Cmp(big.NewInt(0)) == -1 {
		return errors.New("invalid value")
	}
	if len(args.Tokens) != len(args.Amounts) {
		return errors.New("token not match amount")
	}
	return nil
}
