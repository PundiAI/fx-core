package tests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	precompilesstaking "github.com/functionx/fx-core/v3/x/evm/precompiles/staking"
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

func (suite *IntegrationTest) StakingContractTest() {
	var (
		delSigner   = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr     = suite.staking.GetFirstValAddr()
		initBalance = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance  = sdkmath.NewInt(1000).MulRaw(1e18)
	)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// deploy contract to staking
	contract, txHash := suite.staking.DeployStakingContract(delSigner.PrivKey())
	txFee1 := suite.evm.TxFee(txHash)

	// delegate by contract
	receipt := suite.staking.DelegateByContract(delSigner.PrivKey(), contract, valAddr.String(), delBalance.BigInt())
	txFee2 := suite.evm.TxFee(receipt.TxHash)

	delBal := suite.QueryBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(sdkmath.NewIntFromBigInt(txFee2)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Require().Equal(initBalance.String(), total.String())

	// query delegate by contract
	shares, amount := suite.staking.Delegation(valAddr.String(), contract)
	suite.Require().Equal(delBalance.BigInt().String(), amount.String())

	// withdraw by contract
	delBal = suite.QueryBalances(contract.Bytes())
	suite.Require().True(delBal.IsZero())
	receipt = suite.staking.WithdrawByContract(delSigner.PrivKey(), contract, valAddr.String())
	txFee3 := suite.evm.TxFee(receipt.TxHash)

	for _, log := range receipt.Logs {
		if log.Address == precompilesstaking.GetAddress() && log.Topics[0] == precompilesstaking.WithdrawEvent.ID {
			unpack, err := precompilesstaking.WithdrawEvent.Inputs.NonIndexed().Unpack(log.Data)
			suite.Require().NoError(err)
			reward := unpack[1].(*big.Int)
			delBal = suite.QueryBalances(contract.Bytes())
			suite.Require().Equal(reward.String(), delBal.AmountOf(fxtypes.DefaultDenom).BigInt().String())
		}
	}

	delBal = suite.QueryBalances(contract.Bytes())
	suite.Require().True(delBal.AmountOf(fxtypes.DefaultDenom).GT(sdkmath.NewInt(0)))

	// undelegate by contract
	receipt = suite.staking.UndelegateByContract(delSigner.PrivKey(), contract, valAddr.String(), shares)
	txFee4 := suite.evm.TxFee(receipt.TxHash)

	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee := sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(
		sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee4))))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Require().Equal(initBalance.String(), total.String())

	// query delegate by contract
	shares, amount = suite.staking.Delegation(valAddr.String(), contract)
	suite.Require().Equal(amount.String(), big.NewInt(0).String())
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	resp, err := suite.GRPCClient().StakingQuery().UnbondingDelegation(suite.ctx, &stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: sdk.AccAddress(contract.Bytes()).String(),
		ValidatorAddr: valAddr.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(resp.Unbond.Entries))
	suite.Require().Equal(delBalance.BigInt().String(), resp.Unbond.Entries[0].Balance.BigInt().String())
}
