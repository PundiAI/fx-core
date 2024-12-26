package contract

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransferFromModuleToAccountArgs struct {
	Module  string         `abi:"_module"`
	Account common.Address `abi:"_account"`
	Token   common.Address `abi:"_token"`
	Amount  *big.Int       `abi:"_amount"`
}

func (args *TransferFromModuleToAccountArgs) Validate() error {
	if args.Module == "" {
		return errors.New("module address is required")
	}
	if args.Amount == nil || args.Amount.Sign() <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}

type TransferFromAccountToModuleArgs struct {
	Account common.Address `abi:"_account"`
	Module  string         `abi:"_module"`
	Token   common.Address `abi:"_token"`
	Amount  *big.Int       `abi:"_amount"`
}

func (args *TransferFromAccountToModuleArgs) Validate() error {
	if args.Module == "" {
		return errors.New("module address is required")
	}
	if args.Amount == nil || args.Amount.Sign() <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}
