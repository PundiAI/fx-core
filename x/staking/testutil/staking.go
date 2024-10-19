package testutil

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	fxstakingkeeper "github.com/functionx/fx-core/v8/x/staking/keeper"
)

type StakingSuite struct {
	*require.Assertions
	ctx           sdk.Context
	stakingKeeper *fxstakingkeeper.Keeper
}

func (s *StakingSuite) Init(ass *require.Assertions, ctx sdk.Context, stakingKeeper *fxstakingkeeper.Keeper) *StakingSuite {
	s.Assertions = ass
	s.ctx = ctx
	s.stakingKeeper = stakingKeeper
	return s
}

func (s *StakingSuite) GetFirstValidator() stakingtypes.Validator {
	validators, err := s.stakingKeeper.GetValidators(s.ctx, 1)
	s.NoError(err)
	s.NotEmpty(validators)
	return validators[0]
}

func (s *StakingSuite) GetValidators() []stakingtypes.Validator {
	validators, err := s.stakingKeeper.GetValidators(s.ctx, 10)
	s.NoError(err)
	return validators
}

func (s *StakingSuite) GetValidator(valAddr sdk.ValAddress) stakingtypes.Validator {
	validator, err := s.stakingKeeper.GetValidator(s.ctx, valAddr)
	s.NoError(err)
	return validator
}

func (s *StakingSuite) Delegate(delAddr sdk.AccAddress, delAmount sdkmath.Int, val sdk.ValAddress) {
	validator, err := s.stakingKeeper.GetValidator(s.ctx, val)
	s.NoError(err)
	_, err = s.stakingKeeper.Delegate(s.ctx, delAddr, delAmount, stakingtypes.Unbonded, validator, true)
	s.NoError(err)
}

func (s *StakingSuite) GetDelegation(delAddr sdk.AccAddress, val sdk.ValAddress) stakingtypes.Delegation {
	delegation, err := s.stakingKeeper.GetDelegation(s.ctx, delAddr, val)
	s.NoError(err)
	return delegation
}

func (s *StakingSuite) Undelegate(delAddr sdk.AccAddress, val sdk.ValAddress) {
	delegation, err := s.stakingKeeper.GetDelegation(s.ctx, delAddr, val)
	s.NoError(err)
	_, _, err = s.stakingKeeper.Undelegate(s.ctx, delAddr, val, delegation.Shares)
	s.NoError(err)
}

func (s *StakingSuite) Redelegate(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) {
	delegation, err := s.stakingKeeper.GetDelegation(s.ctx, delAddr, valSrcAddr)
	s.NoError(err)
	_, err = s.stakingKeeper.BeginRedelegation(s.ctx, delAddr, valSrcAddr, valDstAddr, delegation.Shares)
	s.NoError(err)
}
