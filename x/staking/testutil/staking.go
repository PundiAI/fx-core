package testutil

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	fxstakingkeeper "github.com/functionx/fx-core/v7/x/staking/keeper"
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
	validators := s.stakingKeeper.GetValidators(s.ctx, 1)
	s.True(len(validators) > 0)
	return validators[0]
}

func (s *StakingSuite) GetValidators() []stakingtypes.Validator {
	return s.stakingKeeper.GetValidators(s.ctx, 10)
}

func (s *StakingSuite) GetValidator(valAddr sdk.ValAddress) stakingtypes.Validator {
	validator, found := s.stakingKeeper.GetValidator(s.ctx, valAddr)
	s.True(found)
	return validator
}

func (s *StakingSuite) Delegate(delAddr sdk.AccAddress, delAmount sdkmath.Int, val sdk.ValAddress) {
	validator, found := s.stakingKeeper.GetValidator(s.ctx, val)
	s.True(found)
	_, err := s.stakingKeeper.Delegate(s.ctx, delAddr, delAmount, stakingtypes.Unbonded, validator, true)
	s.NoError(err)
}

func (s *StakingSuite) GetDelegation(delAddr sdk.AccAddress, val sdk.ValAddress) stakingtypes.Delegation {
	delegation, found := s.stakingKeeper.GetDelegation(s.ctx, delAddr, val)
	s.True(found)
	return delegation
}

func (s *StakingSuite) Undelegate(delAddr sdk.AccAddress, val sdk.ValAddress) {
	delegation, found := s.stakingKeeper.GetDelegation(s.ctx, delAddr, val)
	s.True(found)
	_, err := s.stakingKeeper.Undelegate(s.ctx, delAddr, val, delegation.Shares)
	s.NoError(err)
}

func (s *StakingSuite) Redelegate(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) {
	delegation, found := s.stakingKeeper.GetDelegation(s.ctx, delAddr, valSrcAddr)
	s.True(found)
	_, err := s.stakingKeeper.BeginRedelegation(s.ctx, delAddr, valSrcAddr, valDstAddr, delegation.Shares)
	s.NoError(err)
}
