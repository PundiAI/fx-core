package tests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	stakingprecompile "github.com/functionx/fx-core/v8/x/staking/precompile"
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
	suite.staking.DelegateV2(suite.staking.privKey, valAddr.String(), delBalance.BigInt())

	// query delegate
	share, delegateAmount := suite.staking.Delegation(valAddr.String(), suite.staking.Address())
	query, err := suite.staking.StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{DelegatorAddr: delAddr.String(), ValidatorAddr: valAddr.String()})
	suite.Require().NoError(err)
	suite.Require().EqualValues(delegateAmount.String(), delegateAmount.String())
	suite.Require().EqualValues(share.String(), query.DelegationResponse.Delegation.GetShares().TruncateInt().BigInt().String())
	suite.Require().EqualValues(delegateAmount.String(), query.DelegationResponse.GetBalance().Amount.String())

	// set WithdrawAddress
	rewardAddress := sdk.AccAddress(helpers.NewEthPrivKey().Bytes())
	suite.staking.SetWithdrawAddress(rewardAddress)

	// delegation rewards
	rewards := suite.staking.Rewards(valAddr.String(), suite.staking.Address())
	suite.Require().EqualValues(1, rewards.Cmp(big.NewInt(0)))
	delegationRewards := suite.staking.DelegationRewards(delAddr.String(), valAddr.String())
	suite.Require().EqualValues(1, delegationRewards.AmountOf(fxtypes.DefaultDenom).TruncateInt().BigInt().Cmp(big.NewInt(0)))
	beforeBalance := suite.QueryBalances(rewardAddress)
	suite.Require().True(beforeBalance.IsZero())

	// withdrawReward
	suite.staking.WithdrawReward(suite.staking.privKey, valAddr.String())
	afterBalance := suite.QueryBalances(rewardAddress)
	suite.Require().True(afterBalance.IsAllGTE(beforeBalance))

	// undelegate
	suite.staking.UnDelegateV2(suite.staking.privKey, valAddr.String(), delBalance.BigInt())
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
	receipt := suite.staking.DelegateV2(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
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

	withdrawMethod := stakingprecompile.NewWithdrawMethod(nil)
	for _, log := range receipt.Logs {
		if log.Address == stakingprecompile.GetAddress() && log.Topics[0] == withdrawMethod.Event.ID {
			unpack, err := withdrawMethod.Event.Inputs.NonIndexed().Unpack(log.Data)
			suite.Require().NoError(err)
			reward := unpack[1].(*big.Int)
			delBal = suite.QueryBalances(contract.Bytes())
			suite.Require().Equal(reward.String(), delBal.AmountOf(fxtypes.DefaultDenom).BigInt().String())
		}
	}

	delBal = suite.QueryBalances(contract.Bytes())
	suite.Require().True(delBal.AmountOf(fxtypes.DefaultDenom).GT(sdkmath.NewInt(0)))

	// undelegate by contract
	receipt = suite.staking.UnDelegateV2(delSigner.PrivKey(), valAddr.String(), shares)
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
		delSigner          = helpers.NewSigner(helpers.NewEthPrivKey())
		receiptSigner      = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr            = suite.staking.GetFirstValAddr()
		initBalance        = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance         = sdkmath.NewInt(1000).MulRaw(1e18)
		receiptInitBalance = sdkmath.NewInt(100).MulRaw(1e18)
	)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))
	receiptAccAddr := receiptSigner.AccAddress()
	receiptEthAddr := receiptSigner.Address()
	suite.Send(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptInitBalance))

	// delegate
	receipt := suite.staking.DelegateV2(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
	txFee1 := suite.evm.TxFee(receipt.TxHash)

	// check receipt delegate
	_, amount := suite.staking.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(big.NewInt(0).String(), amount.String())

	// check del delegate
	delSignerAddr := delSigner.Address()
	shares, amount := suite.staking.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(delBalance.BigInt().String(), amount.String())

	halfShares := big.NewInt(0).Div(shares, big.NewInt(2))

	// transfer shares
	receipt = suite.staking.TransferShares(delSigner.PrivKey(), valAddr.String(), receiptEthAddr, halfShares)
	txFee2 := suite.evm.TxFee(receipt.TxHash)

	reward1 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSignerAddr)

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(shares, halfShares)

	// check receipt delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(shares, halfShares)

	// check del balance
	delBal := suite.QueryBalances(delSigner.AccAddress())
	txFee := sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2))
	total := delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1))
	suite.Require().Equal(initBalance.String(), total.String())

	// transfer shares
	receipt = suite.staking.TransferShares(delSigner.PrivKey(), valAddr.String(), receiptEthAddr, halfShares)
	txFee3 := suite.evm.TxFee(receipt.TxHash)

	reward2 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSignerAddr)
	reward3 := suite.staking.LogReward(receipt.Logs, valAddr.String(), receiptEthAddr)

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check receipt delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3)))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1)).Sub(sdkmath.NewIntFromBigInt(reward2))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	receiptExpectBalance := receiptInitBalance.Add(sdkmath.NewIntFromBigInt(reward3))
	suite.CheckBalance(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptExpectBalance))

	// approve
	receipt = suite.staking.ApproveShares(receiptSigner.PrivKey(), valAddr.String(), delSignerAddr, big.NewInt(0).Mul(big.NewInt(3), halfShares))
	txFee4 := suite.evm.TxFee(receipt.TxHash)

	// check approve
	allowance := suite.staking.AllowanceShares(valAddr.String(), receiptEthAddr, delSignerAddr)
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(3), halfShares).String(), allowance.String())

	// check receipt balance
	receiptExpectBalance = receiptExpectBalance.Sub(sdkmath.NewIntFromBigInt(txFee4))
	suite.CheckBalance(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptExpectBalance))

	// transfer from
	receipt = suite.staking.TransferFromShares(delSigner.PrivKey(), valAddr.String(), receiptEthAddr, delSignerAddr, halfShares)
	txFee5 := suite.evm.TxFee(receipt.TxHash)

	reward4 := suite.staking.LogReward(receipt.Logs, valAddr.String(), receiptEthAddr)
	receiptExpectBalance = receiptExpectBalance.Add(sdkmath.NewIntFromBigInt(reward4))

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(shares.String(), halfShares.String())

	// check receipt delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(shares.String(), halfShares.String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee5))))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1)).Sub(sdkmath.NewIntFromBigInt(reward2))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	suite.CheckBalance(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptExpectBalance))

	// check approve
	allowance = suite.staking.AllowanceShares(valAddr.String(), receiptEthAddr, delSignerAddr)
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(2), halfShares).String(), allowance.String())

	// transfer from
	receipt = suite.staking.TransferFromShares(delSigner.PrivKey(), valAddr.String(), receiptEthAddr, delSignerAddr, halfShares)
	txFee6 := suite.evm.TxFee(receipt.TxHash)

	reward5 := suite.staking.LogReward(receipt.Logs, valAddr.String(), delSignerAddr)
	reward6 := suite.staking.LogReward(receipt.Logs, valAddr.String(), receiptEthAddr)
	receiptExpectBalance = receiptExpectBalance.Add(sdkmath.NewIntFromBigInt(reward6))

	// check del delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check receipt delegate
	shares, _ = suite.staking.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check del balance
	delBal = suite.QueryBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee5).Add(sdkmath.NewIntFromBigInt(txFee6)))))
	totalReward := big.NewInt(0).Add(big.NewInt(0).Add(reward1, reward2), reward5)
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(totalReward))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	suite.CheckBalance(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptExpectBalance))

	// check approve
	allowance = suite.staking.AllowanceShares(valAddr.String(), receiptEthAddr, delSignerAddr)
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
	receipt := suite.staking.DelegateV2(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
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
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1))
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

func (suite *IntegrationTest) StakingPrecompileRedelegateTest() {
	var (
		delSigner   = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr     = suite.staking.GetFirstValAddr()
		valNew      = helpers.NewSigner(helpers.NewEthPrivKey())
		initBalance = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance  = sdkmath.NewInt(1000).MulRaw(1e18)
	)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))
	suite.Send(valNew.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// delegate
	receipt := suite.staking.DelegateV2(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
	txFee1 := suite.evm.TxFee(receipt.TxHash)

	delBal := suite.QueryBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())

	// set WithdrawAddress
	rewardAddress := sdk.AccAddress(helpers.NewEthPrivKey().PubKey().Address().Bytes())
	txRsp := suite.staking.SetWithdrawAddressWithResponse(delSigner.PrivKey(), rewardAddress)
	gasPrice, err := sdk.ParseCoinNormalized(suite.network.Config.MinGasPrices)
	suite.Require().NoError(err)
	gasFee := gasPrice.Amount.Mul(sdkmath.NewInt(txRsp.GasWanted))

	hexAddr := common.BytesToAddress(delSigner.AccAddress().Bytes())
	// query delegate
	valAddrShares1, _ := suite.staking.Delegation(valAddr.String(), hexAddr)

	resp := suite.staking.CreateValidator(valNew.PrivKey(), false)
	suite.Equal(resp.Code, uint32(0))

	receipt = suite.staking.RedelegateV2(delSigner.PrivKey(), valAddr.String(), sdk.ValAddress(valNew.AccAddress()).String(), valAddrShares1)
	txFee2 := suite.evm.TxFee(receipt.TxHash)

	valAddrShares2, _ := suite.staking.Delegation(valAddr.String(), hexAddr)
	suite.Equal(big.NewInt(0).String(), valAddrShares2.String())

	valNewShares, _ := suite.staking.Delegation(sdk.ValAddress(valNew.AccAddress()).String(), hexAddr)
	suite.Equal(valAddrShares1, valNewShares)

	delBal = suite.QueryBalances(delSigner.AccAddress())
	total = delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).
		Add(gasFee).
		Add(sdkmath.NewIntFromBigInt(txFee2)).
		Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())
}

func (suite *IntegrationTest) StakingPrecompileRedelegateByContractTest() {
	var (
		delSigner   = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr     = suite.staking.GetFirstValAddr()
		valNew      = helpers.NewSigner(helpers.NewEthPrivKey())
		initBalance = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance  = sdkmath.NewInt(1000).MulRaw(1e18)
	)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))
	suite.Send(valNew.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// deploy contract to staking
	contract, txHash := suite.staking.DeployStakingContract(delSigner.PrivKey())
	txFee1 := suite.evm.TxFee(txHash)

	// delegate by contract
	receipt := suite.staking.DelegateV2(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
	txFee2 := suite.evm.TxFee(receipt.TxHash)

	delBal := suite.QueryBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(sdkmath.NewIntFromBigInt(txFee2)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())

	// query delegate
	valAddrShares1, _ := suite.staking.Delegation(valAddr.String(), contract)

	resp := suite.staking.CreateValidator(valNew.PrivKey(), false)
	suite.Equal(resp.Code, uint32(0))

	receipt = suite.staking.RedelegateV2(delSigner.PrivKey(), valAddr.String(), sdk.ValAddress(valNew.AccAddress()).String(), valAddrShares1)
	txFee3 := suite.evm.TxFee(receipt.TxHash)

	valAddrShares2, _ := suite.staking.Delegation(valAddr.String(), contract)
	suite.Equal(big.NewInt(0).String(), valAddrShares2.String())

	valNewShares, _ := suite.staking.Delegation(sdk.ValAddress(valNew.AccAddress()).String(), contract)
	suite.Equal(valAddrShares1, valNewShares)

	delBal = suite.QueryBalances(delSigner.AccAddress())
	total = delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).
		Add(sdkmath.NewIntFromBigInt(txFee2)).
		Add(sdkmath.NewIntFromBigInt(txFee3)).
		Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())
}

func (suite *IntegrationTest) StakingPrecompileV2() {
	// 1. create a new account, send some balance to it
	delSigner := helpers.NewSigner(helpers.NewEthPrivKey())
	initBalance := sdkmath.NewInt(2000).MulRaw(1e18)
	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// 2. delegate to first validator
	valAddr := suite.staking.GetFirstValAddr()
	delBalance := sdkmath.NewInt(1000)
	receipt := suite.staking.DelegateV2(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
	txFee1 := suite.evm.TxFee(receipt.TxHash)

	delBal := suite.QueryBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())

	hexAddr := common.BytesToAddress(delSigner.AccAddress().Bytes())
	_, delAmt1 := suite.staking.Delegation(valAddr.String(), hexAddr)
	suite.Equal(delAmt1.String(), delBalance.String())

	halfDelegateAmount := big.NewInt(0).Div(delBalance.BigInt(), big.NewInt(2))

	// 2. undelegate half of the first validator amount
	suite.staking.UnDelegateV2(delSigner.PrivKey(), valAddr.String(), halfDelegateAmount)
	_, delAmt2 := suite.staking.Delegation(valAddr.String(), hexAddr)
	suite.Equal(halfDelegateAmount.String(), delAmt2.String())

	// 3. create a new validator, and redelegate half of the first validator amount to it(new validator is not bonded)
	valNew := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(valNew.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))
	resp := suite.staking.CreateValidator(valNew.PrivKey(), false)
	suite.Equal(resp.Code, uint32(0))

	// 4. redelegate half of the first validator amount to new validator
	suite.staking.RedelegateV2(delSigner.PrivKey(), valAddr.String(), sdk.ValAddress(valNew.AccAddress()).String(), halfDelegateAmount)

	delShare, _ := suite.staking.Delegation(valAddr.String(), hexAddr)
	suite.Equal(int64(0), delShare.Int64())

	// 5. check new validator's delegation amount, expecting half of the first validator amount
	_, delAmt3 := suite.staking.Delegation(sdk.ValAddress(valNew.AccAddress()).String(), hexAddr)
	suite.Equal(halfDelegateAmount.String(), delAmt3.String())

	{
		// finally, clear the new validator
		suite.staking.UnDelegateV2(delSigner.PrivKey(), sdk.ValAddress(valNew.AccAddress()).String(), halfDelegateAmount)
		_, valSelfDelegation := suite.staking.Delegation(sdk.ValAddress(valNew.AccAddress()).String(), valNew.Address())
		suite.staking.UnDelegateV2(valNew.PrivKey(), sdk.ValAddress(valNew.AccAddress()).String(), valSelfDelegation)
	}
}
