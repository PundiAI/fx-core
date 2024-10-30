package contract

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type ValidatorSortBy uint8

const (
	ValidatorSortByPower ValidatorSortBy = iota
	ValidatorSortByMissed
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

type DelegateV2Args struct {
	Validator string   `abi:"_val"`
	Amount    *big.Int `abi:"_amount"`
}

// Validate validates the args
func (args *DelegateV2Args) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	if args.Amount == nil || args.Amount.Sign() <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *DelegateV2Args) GetValidator() sdk.ValAddress {
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

type RedelegateV2Args struct {
	ValidatorSrc string   `abi:"_valSrc"`
	ValidatorDst string   `abi:"_valDst"`
	Amount       *big.Int `abi:"_amount"`
}

// Validate validates the args
func (args *RedelegateV2Args) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.ValidatorSrc); err != nil {
		return fmt.Errorf("invalid validator src address: %s", args.ValidatorSrc)
	}
	if _, err := sdk.ValAddressFromBech32(args.ValidatorDst); err != nil {
		return fmt.Errorf("invalid validator dst address: %s", args.ValidatorDst)
	}
	if args.Amount == nil || args.Amount.Sign() <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}

// GetValidatorSrc returns the validator src address, caller must ensure the validator address is valid
func (args *RedelegateV2Args) GetValidatorSrc() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.ValidatorSrc)
	return valAddr
}

// GetValidatorDst returns the validator dest address, caller must ensure the validator address is valid
func (args *RedelegateV2Args) GetValidatorDst() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.ValidatorDst)
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

type TransferSharesRet struct {
	Token  *big.Int
	Reward *big.Int
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

type TransferFromSharesRet struct {
	Token  *big.Int
	Reward *big.Int
}

type UndelegateV2Args struct {
	Validator string   `abi:"_val"`
	Amount    *big.Int `abi:"_amount"`
}

// Validate validates the args
func (args *UndelegateV2Args) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	if args.Amount == nil || args.Amount.Sign() <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *UndelegateV2Args) GetValidator() sdk.ValAddress {
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

type SlashingInfoArgs struct {
	Validator string `abi:"_val"`
}

// Validate validates the args
func (args *SlashingInfoArgs) Validate() error {
	if _, err := sdk.ValAddressFromBech32(args.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	return nil
}

// GetValidator returns the validator address, caller must ensure the validator address is valid
func (args *SlashingInfoArgs) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(args.Validator)
	return valAddr
}

type ValidatorListArgs struct {
	SortBy uint8 `abi:"_val"`
}

// Validate validates the args
func (args *ValidatorListArgs) Validate() error {
	if args.SortBy > uint8(ValidatorSortByMissed) {
		return fmt.Errorf("over the sort by limit")
	}
	return nil
}

func (args *ValidatorListArgs) GetSortBy() ValidatorSortBy {
	return ValidatorSortBy(args.SortBy)
}
