package types

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// NewGenesisState creates a new GenesisState instanc e
func NewGenesisState(params stakingtypes.Params, validators []stakingtypes.Validator, delegations []stakingtypes.Delegation) *GenesisState {
	return &GenesisState{
		Params:      params,
		Validators:  validators,
		Delegations: delegations,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: stakingtypes.DefaultParams(),
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.Validators {
		if err := g.Validators[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}

func (a Allowance) Validate() error {
	if len(a.ValidatorAddress) == 0 {
		return errors.New("validator address cannot be empty")
	}
	if _, err := sdk.ValAddressFromBech32(a.ValidatorAddress); err != nil {
		return fmt.Errorf("invalid validator address: %s", err.Error())
	}
	if len(a.OwnerAddress) == 0 {
		return errors.New("owner address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(a.OwnerAddress); err != nil {
		return fmt.Errorf("invalid owner address: %s", err.Error())
	}
	if len(a.SpenderAddress) == 0 {
		return errors.New("spender address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(a.SpenderAddress); err != nil {
		return fmt.Errorf("invalid spender address: %s", err.Error())
	}
	if a.Allowance.LTE(sdkmath.ZeroInt()) {
		return fmt.Errorf("allowance must be greater than 0, is %s", a.Allowance)
	}
	return nil
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data *GenesisState) error {
	if err := ValidateGenesisStateValidators(data.Validators); err != nil {
		return err
	}

	if err := ValidateGenesisStateAllowances(data.Allowances); err != nil {
		return err
	}

	return data.Params.Validate()
}

func ValidateGenesisStateAllowances(allowances []Allowance) error {
	for _, a := range allowances {
		if err := a.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func ValidateGenesisStateValidators(validators []stakingtypes.Validator) error {
	addrMap := make(map[string]bool, len(validators))

	for i := 0; i < len(validators); i++ {
		val := validators[i]
		consPk, err := val.ConsPubKey()
		if err != nil {
			return err
		}

		strKey := string(consPk.Bytes())

		if _, ok := addrMap[strKey]; ok {
			consAddr, err := val.GetConsAddr()
			if err != nil {
				return err
			}
			return fmt.Errorf("duplicate validator in genesis state: moniker %v, address %v", val.Description.Moniker, consAddr)
		}

		if val.Jailed && val.IsBonded() {
			consAddr, err := val.GetConsAddr()
			if err != nil {
				return err
			}
			return fmt.Errorf("validator is bonded and jailed in genesis state: moniker %v, address %v", val.Description.Moniker, consAddr)
		}

		if val.DelegatorShares.IsZero() && !val.IsUnbonding() {
			return fmt.Errorf("bonded/unbonded genesis validator cannot have zero delegator shares, validator: %v", val)
		}

		addrMap[strKey] = true
	}

	return nil
}

type Allowances []Allowance

func (a Allowances) Len() int { return len(a) }
func (a Allowances) Less(i, j int) bool {
	if a[i].ValidatorAddress != a[j].ValidatorAddress {
		return a[i].ValidatorAddress < a[j].ValidatorAddress
	}
	if a[i].OwnerAddress != a[j].OwnerAddress {
		return a[i].OwnerAddress < a[j].OwnerAddress
	}
	return a[i].SpenderAddress < a[j].SpenderAddress
}
func (a Allowances) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
