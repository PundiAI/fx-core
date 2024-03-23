package staking

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
)

var (
	// AllowanceSharesMethod Query the amount of shares an owner allowed to a spender.
	AllowanceSharesMethod = abi.NewMethod(
		AllowanceSharesMethodName,
		AllowanceSharesMethodName,
		abi.Function, "view", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
			abi.Argument{Name: "_owner", Type: contract.TypeAddress},
			abi.Argument{Name: "_spender", Type: contract.TypeAddress},
		},
		abi.Arguments{
			abi.Argument{Name: "_shares", Type: contract.TypeUint256},
		},
	)

	// DelegationMethod Query the amount of shares in val.
	DelegationMethod = abi.NewMethod(
		DelegationMethodName,
		DelegationMethodName,
		abi.Function, "view", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
			abi.Argument{Name: "_del", Type: contract.TypeAddress},
		},
		abi.Arguments{
			abi.Argument{Name: "_shares", Type: contract.TypeUint256},
			abi.Argument{Name: "_delegateAmount", Type: contract.TypeUint256},
		},
	)

	// DelegationRewardsMethod Query the amount of rewards in val.
	DelegationRewardsMethod = abi.NewMethod(
		DelegationRewardsMethodName,
		DelegationRewardsMethodName,
		abi.Function, "view", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
			abi.Argument{Name: "_del", Type: contract.TypeAddress},
		},
		abi.Arguments{
			abi.Argument{Name: "_reward", Type: contract.TypeUint256},
		},
	)
)

var (
	// ApproveSharesMethod Approve shares to a spender.
	ApproveSharesMethod = abi.NewMethod(
		ApproveSharesMethodName,
		ApproveSharesMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
			abi.Argument{Name: "_spender", Type: contract.TypeAddress},
			abi.Argument{Name: "_shares", Type: contract.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: contract.TypeBool},
		},
	)

	// DelegateMethod Delegate token to a validator.
	DelegateMethod = abi.NewMethod(
		DelegateMethodName,
		DelegateMethodName,
		abi.Function, "payable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
		},
		abi.Arguments{
			abi.Argument{Name: "_shares", Type: contract.TypeUint256},
			abi.Argument{Name: "_reward", Type: contract.TypeUint256},
		},
	)

	// TransferSharesMethod Transfer shares to a recipient.
	TransferSharesMethod = abi.NewMethod(
		TransferSharesMethodName,
		TransferSharesMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
			abi.Argument{Name: "_to", Type: contract.TypeAddress},
			abi.Argument{Name: "_shares", Type: contract.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_token", Type: contract.TypeUint256},
			abi.Argument{Name: "_reward", Type: contract.TypeUint256},
		},
	)

	// TransferFromSharesMethod Transfer shares from a sender to a recipient.
	TransferFromSharesMethod = abi.NewMethod(
		TransferFromSharesMethodName,
		TransferFromSharesMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
			abi.Argument{Name: "_from", Type: contract.TypeAddress},
			abi.Argument{Name: "_to", Type: contract.TypeAddress},
			abi.Argument{Name: "_shares", Type: contract.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_token", Type: contract.TypeUint256},
			abi.Argument{Name: "_reward", Type: contract.TypeUint256},
		},
	)

	// UndelegateMethod Undelegate shares from a validator.
	UndelegateMethod = abi.NewMethod(
		UndelegateMethodName,
		UndelegateMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
			abi.Argument{Name: "_shares", Type: contract.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_amount", Type: contract.TypeUint256},
			abi.Argument{Name: "_reward", Type: contract.TypeUint256},
			abi.Argument{Name: "_completionTime", Type: contract.TypeUint256},
		},
	)

	// WithdrawMethod Withdraw rewards from a validator.
	WithdrawMethod = abi.NewMethod(
		WithdrawMethodName,
		WithdrawMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: contract.TypeString},
		},
		abi.Arguments{
			abi.Argument{Name: "_reward", Type: contract.TypeUint256},
		},
	)

	// RedelegateMethod Redelegate share from src validator to dest validator
	RedelegateMethod = abi.NewMethod(
		RedelegateMethodName,
		RedelegateMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_valSrc", Type: contract.TypeString},
			abi.Argument{Name: "_valDst", Type: contract.TypeString},
			abi.Argument{Name: "_shares", Type: contract.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_amount", Type: contract.TypeUint256},
			abi.Argument{Name: "_reward", Type: contract.TypeUint256},
			abi.Argument{Name: "_completionTime", Type: contract.TypeUint256},
		},
	)
)

type AllowanceSharesArgs struct {
	Validator string         `abi:"_val"`
	Owner     common.Address `abi:"_owner"`
	Spender   common.Address `abi:"_spender"`
}

// Validate validates the args
func (args *AllowanceSharesArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *AllowanceSharesArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type ApproveSharesArgs struct {
	Validator string         `abi:"_val"`
	Spender   common.Address `abi:"_spender"`
	Shares    *big.Int       `abi:"_shares"`
}

// Validate validates the args
func (args *ApproveSharesArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	if args.Shares == nil || args.Shares.Sign() < 0 {
		return errors.New("invalid shares")
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *ApproveSharesArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type DelegateArgs struct {
	Validator string `abi:"_val"`
}

// Validate validates the args
func (args *DelegateArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *DelegateArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type DelegationArgs struct {
	Validator string         `abi:"_val"`
	Delegator common.Address `abi:"_del"`
}

// Validate validates the args
func (args *DelegationArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *DelegationArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type DelegationRewardsArgs struct {
	Validator string         `abi:"_val"`
	Delegator common.Address `abi:"_del"`
}

// Validate validates the args
func (args *DelegationRewardsArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *DelegationRewardsArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type TransferSharesArgs struct {
	Validator string         `abi:"_val"`
	To        common.Address `abi:"_to"`
	Shares    *big.Int       `abi:"_shares"`
}

// Validate validates the args
func (args *TransferSharesArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	if args.Shares == nil || args.Shares.Sign() <= 0 {
		return errors.New("invalid shares")
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *TransferSharesArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type TransferFromSharesArgs struct {
	Validator string         `abi:"_val"`
	From      common.Address `abi:"_from"`
	To        common.Address `abi:"_to"`
	Shares    *big.Int       `abi:"_shares"`
}

// Validate validates the args
func (args *TransferFromSharesArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	if args.Shares == nil || args.Shares.Sign() <= 0 {
		return errors.New("invalid shares")
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *TransferFromSharesArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type UndelegateArgs struct {
	Validator string   `abi:"_val"`
	Shares    *big.Int `abi:"_shares"`
}

// Validate validates the args
func (args *UndelegateArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	if args.Shares == nil || args.Shares.Sign() <= 0 {
		return errors.New("invalid shares")
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *UndelegateArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type WithdrawArgs struct {
	Validator string `abi:"_val"`
}

// Validate validates the args
func (args *WithdrawArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *WithdrawArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type RedelegateArgs struct {
	ValidatorSrc string   `abi:"_valSrc"`
	ValidatorDst string   `abi:"_valDst"`
	Shares       *big.Int `abi:"_shares"`
}

// Validate validates the args
func (args *RedelegateArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.ValidatorSrc); err != nil {
		return fmt.Errorf("invalid validator src address: %s", args.ValidatorSrc)
	}
	if _, err := sdk.ValAddressFromBech32(args.ValidatorDst); err != nil {
		return fmt.Errorf("invalid validator dst address: %s", args.ValidatorDst)
	}
	if args.Shares == nil || args.Shares.Sign() <= 0 {
		return errors.New("invalid shares")
	}
	return nil
}

// GetValidatorSrc returns the validator src address, caller must ensure the validator address is valid
func (args *RedelegateArgs) GetValidatorSrc() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.ValidatorSrc)
	return valAddr
}

// GetValidatorDst returns the validator dest address, caller must ensure the validator address is valid
func (args *RedelegateArgs) GetValidatorDst() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.ValidatorDst)
	return valAddr
}
