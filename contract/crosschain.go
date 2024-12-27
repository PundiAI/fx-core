package contract

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

// Deprecated: After the upgrade to v8
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

type BridgeCallArgs struct {
	DstChain string           `abi:"_dstChain"`
	Refund   common.Address   `abi:"_refund"`
	Tokens   []common.Address `abi:"_tokens"`
	Amounts  []*big.Int       `abi:"_amounts"`
	To       common.Address   `abi:"_to"`
	Data     []byte           `abi:"_data"`
	QuoteId  *big.Int         `abi:"_quoteId"`
	GasLimit *big.Int         `abi:"_gasLimit"`
	Memo     []byte           `abi:"_memo"`
}

func (args *BridgeCallArgs) Validate() error {
	if args.DstChain == "" {
		return errors.New("empty chain")
	}
	if len(args.Tokens) != len(args.Amounts) {
		return errors.New("tokens and amounts do not match")
	}
	if len(args.Amounts) > 0 && IsZeroEthAddress(args.Refund) {
		return errors.New("refund cannot be empty")
	}
	if args.QuoteId.Sign() < 0 {
		return errors.New("quoteId cannot be negative")
	}
	return nil
}

type ExecuteClaimArgs struct {
	Chain      string   `abi:"_chain"`
	EventNonce *big.Int `abi:"_eventNonce"`
}

func (args *ExecuteClaimArgs) Validate() error {
	if args.Chain == "" {
		return errors.New("empty chain")
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
	if args.Chain == "" {
		return errors.New("empty chain")
	}
	if IsZeroEthAddress(args.ExternalAddress) {
		return errors.New("invalid external address")
	}
	return nil
}

type IsOracleOnlineArgs struct {
	Chain           string         `abi:"_chain"`
	ExternalAddress common.Address `abi:"_externalAddress"`
}

func (args *IsOracleOnlineArgs) Validate() error {
	if args.Chain == "" {
		return errors.New("empty chain")
	}
	if IsZeroEthAddress(args.ExternalAddress) {
		return errors.New("invalid external address")
	}
	return nil
}
