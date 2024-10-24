package types

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
)

type BridgeCoinAmountArgs struct {
	Token  common.Address `abi:"_token"`
	Target [32]byte       `abi:"_target"`
}

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

func (args *CancelSendToExternalArgs) Validate() error {
	if err := ValidateModuleName(args.Chain); err != nil {
		return err
	}
	if args.TxID == nil || args.TxID.Sign() <= 0 {
		return errors.New("invalid tx id")
	}
	return nil
}

type CrosschainArgs struct {
	Token   common.Address `abi:"_token"`
	Receipt string         `abi:"_receipt"`
	Amount  *big.Int       `abi:"_amount"`
	Fee     *big.Int       `abi:"_fee"`
	Target  [32]byte       `abi:"_target"`
	Memo    string         `abi:"_memo"`
}

func (args *CrosschainArgs) Validate() error {
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

func (args *IncreaseBridgeFeeArgs) Validate() error {
	if err := ValidateModuleName(args.Chain); err != nil {
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
	Refund   common.Address   `abi:"_refund"`
	Tokens   []common.Address `abi:"_tokens"`
	Amounts  []*big.Int       `abi:"_amounts"`
	To       common.Address   `abi:"_to"`
	Data     []byte           `abi:"_data"`
	Value    *big.Int         `abi:"_value"`
	Memo     []byte           `abi:"_memo"`
}

func (args *BridgeCallArgs) Validate() error {
	if err := ValidateModuleName(args.DstChain); err != nil {
		return err
	}
	if args.Value.Sign() != 0 {
		return errors.New("value must be zero")
	}
	if len(args.Tokens) != len(args.Amounts) {
		return errors.New("tokens and amounts do not match")
	}
	if len(args.Amounts) > 0 && contract.IsZeroEthAddress(args.Refund) {
		return errors.New("refund cannot be empty")
	}
	return nil
}

type ExecuteClaimArgs struct {
	Chain      string   `abi:"_chain"`
	EventNonce *big.Int `abi:"_eventNonce"`
}

func (args *ExecuteClaimArgs) Validate() error {
	if err := ValidateModuleName(args.Chain); err != nil {
		return err
	}
	if args.EventNonce == nil || args.EventNonce.Sign() <= 0 {
		return errors.New("invalid event nonce")
	}
	return nil
}

type HasOracleArgs struct {
	Chain           string         `abi:"_chain"`
	ExternalAddress common.Address `abi:"_externalAddress"`
}

func (args *HasOracleArgs) Validate() error {
	if err := ValidateModuleName(args.Chain); err != nil {
		return err
	}
	if contract.IsZeroEthAddress(args.ExternalAddress) {
		return errors.New("invalid external address")
	}
	return nil
}

type IsOracleOnlineArgs struct {
	Chain           string         `abi:"_chain"`
	ExternalAddress common.Address `abi:"_externalAddress"`
}

func (args *IsOracleOnlineArgs) Validate() error {
	if err := ValidateModuleName(args.Chain); err != nil {
		return err
	}
	if contract.IsZeroEthAddress(args.ExternalAddress) {
		return errors.New("invalid external address")
	}
	return nil
}
