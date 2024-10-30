package helpers

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (s *BaseSuite) GetFirstValidator() stakingtypes.Validator {
	validators, err := s.App.StakingKeeper.GetValidators(s.Ctx, 1)
	s.NoError(err)
	s.NotEmpty(validators)
	return validators[0]
}

func (s *BaseSuite) GetFirstValAddr() sdk.ValAddress {
	val := s.GetFirstValidator()
	operator, err := s.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	s.NoError(err)
	return operator
}

func (s *BaseSuite) GetSecondValidator() stakingtypes.Validator {
	validators, err := s.App.StakingKeeper.GetValidators(s.Ctx, 2)
	s.NoError(err)
	s.NotEmpty(validators)
	return validators[1]
}

func (s *BaseSuite) GetSecondValAddr() sdk.ValAddress {
	val := s.GetSecondValidator()
	operator, err := s.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	s.NoError(err)
	return operator
}

func (s *BaseSuite) GetValidators() []stakingtypes.Validator {
	validators, err := s.App.StakingKeeper.GetValidators(s.Ctx, 10)
	s.NoError(err)
	return validators
}

func (s *BaseSuite) GetValidator(valAddr sdk.ValAddress) stakingtypes.Validator {
	validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
	s.NoError(err)
	return validator
}

func (s *BaseSuite) Delegate(delAddr sdk.AccAddress, delAmount sdkmath.Int, val sdk.ValAddress) {
	validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, val)
	s.NoError(err)
	_, err = s.App.StakingKeeper.Delegate(s.Ctx, delAddr, delAmount, stakingtypes.Unbonded, validator, true)
	s.NoError(err)
}

func (s *BaseSuite) GetDelegation(delAddr sdk.AccAddress, val sdk.ValAddress) stakingtypes.Delegation {
	delegation, err := s.App.StakingKeeper.GetDelegation(s.Ctx, delAddr, val)
	s.NoError(err)
	return delegation
}

func (s *BaseSuite) Undelegate(delAddr sdk.AccAddress, val sdk.ValAddress) {
	delegation, err := s.App.StakingKeeper.GetDelegation(s.Ctx, delAddr, val)
	s.NoError(err)
	_, _, err = s.App.StakingKeeper.Undelegate(s.Ctx, delAddr, val, delegation.Shares)
	s.NoError(err)
}

func (s *BaseSuite) Redelegate(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) {
	delegation, err := s.App.StakingKeeper.GetDelegation(s.Ctx, delAddr, valSrcAddr)
	s.NoError(err)
	_, err = s.App.StakingKeeper.BeginRedelegation(s.Ctx, delAddr, valSrcAddr, valDstAddr, delegation.Shares)
	s.NoError(err)
}

func (s *BaseSuite) SetAllowance(valAddr sdk.ValAddress, owner, spender sdk.AccAddress, shares *big.Int) {
	s.App.StakingKeeper.SetAllowance(s.Ctx, valAddr, owner, spender, shares)
}
