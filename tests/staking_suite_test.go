package tests

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func (suite *IntegrationTest) StakingTest() {
	var (
		delAddr     = suite.staking.AccAddress()
		valAddr     = suite.staking.GetFirstValAddr()
		initBalance = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance  = sdkmath.NewInt(1000).MulRaw(1e18)
	)
	suite.Send(delAddr, sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// delegate
	suite.staking.Delegate(suite.staking.privKey, valAddr.String(), delBalance.BigInt())

	// query delegate
	share, delegateAmount := suite.staking.Delegation(valAddr.String(), suite.staking.Address())
	query, err := suite.staking.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{DelegatorAddr: delAddr.String(), ValidatorAddr: valAddr.String()})
	suite.Require().NoError(err)
	suite.Require().EqualValues(delegateAmount.String(), delegateAmount.String())
	suite.Require().EqualValues(share.String(), query.DelegationResponse.Delegation.GetShares().TruncateInt().BigInt().String())
	suite.Require().EqualValues(delegateAmount.String(), query.DelegationResponse.GetBalance().Amount.String())

	// set WithdrawAddress
	rewardAddress := sdk.AccAddress(helpers.NewEthPrivKey().Bytes())
	suite.staking.SetWithdrawAddress(delAddr, rewardAddress)

	// delegation rewards
	rewards := suite.staking.Rewards(valAddr.String(), suite.staking.Address())
	delegationRewards := suite.staking.DelegationRewards(delAddr.String(), valAddr.String())
	suite.Require().EqualValues(rewards.String(), delegationRewards.AmountOf(fxtypes.DefaultDenom).TruncateInt().BigInt().String())
	beforeBalance := suite.QueryBalances(rewardAddress)
	suite.Require().True(beforeBalance.IsZero())

	// withdrawReward
	suite.staking.WithdrawReward(suite.staking.privKey, valAddr.String())
	afterBalance := suite.QueryBalances(rewardAddress)
	suite.Require().True(afterBalance.IsAllGTE(beforeBalance))

	// undelegate
	suite.staking.UnDelegate(suite.staking.privKey, valAddr.String(), delBalance.BigInt())
	rewards = suite.staking.Rewards(valAddr.String(), suite.staking.Address())
	suite.Require().EqualValues(rewards.String(), sdkmath.NewInt(0).String())
	afterBalance2 := suite.QueryBalances(rewardAddress)
	suite.Require().True(afterBalance2.IsAllGTE(beforeBalance))
}
