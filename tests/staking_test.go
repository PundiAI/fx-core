package tests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v4/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v4/types"
	precompilesstaking "github.com/functionx/fx-core/v4/x/evm/precompiles/staking"
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

func (suite *IntegrationTest) StakingSharesTest() {
	var (
		delSigner     = helpers.NewSigner(helpers.NewEthPrivKey())
		receiptSigner = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr       = suite.staking.GetFirstValAddr()
		initBalance   = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance    = sdkmath.NewInt(1000).MulRaw(1e18)
	)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// delegate
	receipt := suite.staking.Delegate(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
	txFee1 := suite.evm.TxFee(receipt.TxHash)

	// check receipt delegate
	_, amount := suite.staking.Delegation(valAddr.String(), receiptSigner.Address())
	suite.Require().Equal(big.NewInt(0).String(), amount.String())

	// check del delegate
	shares, amount := suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(delBalance.BigInt().String(), amount.String())

	halfShares := big.NewInt(0).Div(shares, big.NewInt(2))

	// transfer shares
	receipt = suite.staking.TransferShares(delSigner.PrivKey(), valAddr.String(), receiptSigner.Address(), halfShares)
	txFee2 := suite.evm.TxFee(receipt.TxHash)

	reward1 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares, halfShares)

	// check receipt delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), receiptSigner.Address())
	suite.Require().Equal(shares, halfShares)

	// check del balance
	delBal := suite.QueryBalances(delSigner.AccAddress())
	txFee := sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2))
	total := delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1))
	suite.Require().Equal(initBalance.String(), total.String())

	// transfer shares
	receipt = suite.staking.TransferShares(delSigner.PrivKey(), valAddr.String(), receiptSigner.Address(), halfShares)
	txFee3 := suite.evm.TxFee(receipt.TxHash)

	reward2 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())
	reward3 := suite.staking.LogReward(receipt.Logs, valAddr.String(), receiptSigner.Address())

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check receipt delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), receiptSigner.Address())
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3)))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1)).Sub(sdkmath.NewIntFromBigInt(reward2))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	delBal = suite.QueryBalances(receiptSigner.AccAddress())
	suite.Require().Equal(sdkmath.NewIntFromBigInt(reward3).String(), delBal.AmountOf(fxtypes.DefaultDenom).String())

	// approve
	receipt = suite.staking.ApproveShares(receiptSigner.PrivKey(), valAddr.String(), delSigner.Address(), big.NewInt(0).Mul(big.NewInt(3), halfShares))
	txFee4 := suite.evm.TxFee(receipt.TxHash)

	// check approve
	allowance := suite.staking.AllowanceShares(valAddr.String(), receiptSigner.Address(), delSigner.Address())
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(3), halfShares).String(), allowance.String())

	// check receipt balance
	delBal = suite.QueryBalances(receiptSigner.AccAddress())
	suite.Require().Equal(delBal.AmountOf(fxtypes.DefaultDenom).String(), big.NewInt(0).Sub(reward3, txFee4).String())

	// transfer from
	receipt = suite.staking.TransferFromShares(delSigner.PrivKey(), valAddr.String(), receiptSigner.Address(), delSigner.Address(), halfShares)
	txFee5 := suite.evm.TxFee(receipt.TxHash)

	reward4 := suite.staking.LogReward(receipt.Logs, valAddr.String(), receiptSigner.Address())

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), halfShares.String())

	// check receipt delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), receiptSigner.Address())
	suite.Require().Equal(shares.String(), halfShares.String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee5))))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1)).Sub(sdkmath.NewIntFromBigInt(reward2))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	delBal = suite.QueryBalances(receiptSigner.AccAddress())
	suite.Require().Equal(delBal.AmountOf(fxtypes.DefaultDenom).String(), big.NewInt(0).Sub(big.NewInt(0).Add(reward3, reward4), txFee4).String())

	// check approve
	allowance = suite.staking.AllowanceShares(valAddr.String(), receiptSigner.Address(), delSigner.Address())
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(2), halfShares).String(), allowance.String())

	// transfer from
	receipt = suite.staking.TransferFromShares(delSigner.PrivKey(), valAddr.String(), receiptSigner.Address(), delSigner.Address(), halfShares)
	txFee6 := suite.evm.TxFee(receipt.TxHash)

	reward5 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())
	reward6 := suite.staking.LogReward(receipt.Logs, valAddr.String(), receiptSigner.Address())

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check receipt delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), receiptSigner.Address())
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee5).Add(sdkmath.NewIntFromBigInt(txFee6)))))
	totalReward := big.NewInt(0).Add(big.NewInt(0).Add(reward1, reward2), reward5)
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdk.NewIntFromBigInt(totalReward))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	delBal = suite.QueryBalances(receiptSigner.AccAddress())
	total = sdkmath.NewIntFromBigInt(reward3).Add(sdkmath.NewIntFromBigInt(reward4)).Add(sdkmath.NewIntFromBigInt(reward6)).Sub(sdkmath.NewIntFromBigInt(txFee4))
	suite.Require().Equal(delBal.AmountOf(fxtypes.DefaultDenom).String(), total.String())

	// check approve
	allowance = suite.staking.AllowanceShares(valAddr.String(), receiptSigner.Address(), delSigner.Address())
	suite.Require().Equal(halfShares.String(), allowance.String())
}

func (suite *IntegrationTest) StakingSharesContractTest() {
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

	// delegate
	receipt := suite.staking.Delegate(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
	txFee2 := suite.evm.TxFee(receipt.TxHash)

	// check del balance
	delBal := suite.QueryBalances(delSigner.AccAddress())
	txFee := sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2))
	total := delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract delegate
	_, amount := suite.staking.Delegation(valAddr.String(), contract)
	suite.Require().Equal(big.NewInt(0).String(), amount.String())

	// check del delegate
	shares, amount := suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(delBalance.BigInt().String(), amount.String())

	halfShares := big.NewInt(0).Div(shares, big.NewInt(2))

	// approve
	receipt = suite.staking.ApproveShares(delSigner.PrivKey(), valAddr.String(), contract, big.NewInt(0).Mul(big.NewInt(3), halfShares))
	txFee3 := suite.evm.TxFee(receipt.TxHash)

	// check approve
	allowance := suite.staking.AllowanceShares(valAddr.String(), delSigner.Address(), contract)
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(3), halfShares).String(), allowance.String())

	// transferFromShares
	receipt = suite.staking.TransferFromSharesByContract(delSigner.PrivKey(), valAddr.String(), contract, delSigner.Address(), contract, halfShares)
	txFee4 := suite.evm.TxFee(receipt.TxHash)

	reward1 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), halfShares.String())

	// check contract delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), contract)
	suite.Require().Equal(shares.String(), halfShares.String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee4))))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdk.NewIntFromBigInt(reward1))
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract balance
	contractBal := suite.QueryBalances(contract.Bytes())
	suite.Require().Equal(contractBal.AmountOf(fxtypes.DefaultDenom).String(), big.NewInt(0).String())

	// check approve
	allowance = suite.staking.AllowanceShares(valAddr.String(), delSigner.Address(), contract)
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(2), halfShares).String(), allowance.String())

	// transferFromShares
	receipt = suite.staking.TransferFromSharesByContract(delSigner.PrivKey(), valAddr.String(), contract, delSigner.Address(), contract, halfShares)
	txFee5 := suite.evm.TxFee(receipt.TxHash)
	reward2 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())
	reward3 := suite.staking.LogReward(receipt.Logs, valAddr.String(), contract)

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check contract delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), contract)
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee4).Add(sdkmath.NewIntFromBigInt(txFee5)))))
	totalReward := sdkmath.NewIntFromBigInt(reward1).Add(sdkmath.NewIntFromBigInt(reward2))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(totalReward)
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract balance
	contractBal = suite.QueryBalances(contract.Bytes())
	suite.Require().Equal(contractBal.AmountOf(fxtypes.DefaultDenom).String(), reward3.String())

	// check approve
	allowance = suite.staking.AllowanceShares(valAddr.String(), delSigner.Address(), contract)
	suite.Require().Equal(halfShares.String(), allowance.String())

	// contract transfer
	receipt = suite.staking.TransferSharesByContract(delSigner.PrivKey(), valAddr.String(), contract, delSigner.Address(), halfShares)
	txFee6 := suite.evm.TxFee(receipt.TxHash)
	reward4 := suite.staking.LogReward(receipt.Logs, valAddr.String(), contract)

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), halfShares.String())

	// check contract delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), contract)
	suite.Require().Equal(shares.String(), halfShares.String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).
		Add(sdkmath.NewIntFromBigInt(txFee4).Add(sdkmath.NewIntFromBigInt(txFee5).Add(sdkmath.NewIntFromBigInt(txFee6))))))
	totalReward = sdkmath.NewIntFromBigInt(reward1).Add(sdkmath.NewIntFromBigInt(reward2))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(totalReward)
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract balance
	contractBal = suite.QueryBalances(contract.Bytes())
	suite.Require().Equal(contractBal.AmountOf(fxtypes.DefaultDenom).String(), big.NewInt(0).Add(reward3, reward4).String())

	// contract transfer
	receipt = suite.staking.TransferSharesByContract(delSigner.PrivKey(), valAddr.String(), contract, delSigner.Address(), halfShares)
	txFee7 := suite.evm.TxFee(receipt.TxHash)
	reward5 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())
	reward6 := suite.staking.LogReward(receipt.Logs, valAddr.String(), contract)

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check contract delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), contract)
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).
		Add(sdkmath.NewIntFromBigInt(txFee4).Add(sdkmath.NewIntFromBigInt(txFee5).Add(sdkmath.NewIntFromBigInt(txFee6).Add(sdkmath.NewIntFromBigInt(txFee7)))))))
	totalReward = sdkmath.NewIntFromBigInt(reward1).Add(sdkmath.NewIntFromBigInt(reward2).Add(sdkmath.NewIntFromBigInt(reward5)))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(totalReward)
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract balance
	contractBal = suite.QueryBalances(contract.Bytes())
	totalReward = sdkmath.NewIntFromBigInt(reward3).Add(sdkmath.NewIntFromBigInt(reward4).Add(sdkmath.NewIntFromBigInt(reward6)))
	suite.Require().Equal(contractBal.AmountOf(fxtypes.DefaultDenom).String(), totalReward.String())
}
