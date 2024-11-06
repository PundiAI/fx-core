package integration

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func (suite *IntegrationTest) StakingTest() {
	signer := helpers.NewSigner(helpers.NewEthPrivKey())

	var (
		valAddr    = suite.GetValAddr()
		delBalance = sdkmath.NewInt(1000).MulRaw(1e18)
	)
	suite.Send(signer.AccAddress(), suite.NewStakingCoin(2000, 18))

	stakingSuite := NewStakingSuite(suite.EthSuite, common.HexToAddress(contract.StakingAddress), signer)
	stakingSuite.DelegateV2(valAddr.String(), delBalance.BigInt())

	share, delegateAmount := stakingSuite.Delegation(valAddr.String(), signer.Address())

	delegation := suite.GetDelegation(signer.AccAddress(), valAddr)
	suite.Require().EqualValues(delegateAmount.String(), delegateAmount.String())
	suite.Require().EqualValues(share.String(), delegation.Delegation.GetShares().TruncateInt().BigInt().String())
	suite.Require().EqualValues(delegateAmount.String(), delegation.GetBalance().Amount.String())

	rewardAddress := helpers.GenAccAddress()
	suite.SetWithdrawAddress(stakingSuite.signer, rewardAddress)

	// delegation rewards
	rewards := stakingSuite.Rewards(valAddr.String(), signer.Address())
	suite.Require().EqualValues(1, rewards.Cmp(big.NewInt(0)))
	delegationRewards := suite.DelegationRewards(signer.AccAddress().String(), valAddr.String())
	suite.Require().EqualValues(1, delegationRewards.AmountOf(fxtypes.DefaultDenom).TruncateInt().BigInt().Cmp(big.NewInt(0)))
	beforeBalance := suite.GetAllBalances(rewardAddress)
	suite.Require().True(beforeBalance.IsZero())

	suite.WithdrawReward(stakingSuite.signer, valAddr)
	afterBalance := suite.GetAllBalances(rewardAddress)
	suite.Require().True(afterBalance.IsAllGTE(beforeBalance))

	stakingSuite.UnDelegateV2(valAddr.String(), delBalance.BigInt())
	afterBalance2 := suite.GetAllBalances(rewardAddress)
	suite.Require().True(afterBalance2.IsAllGTE(beforeBalance))
}

func (suite *IntegrationTest) StakingContractTest() {
	var (
		delSigner   = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr     = suite.GetValAddr()
		initBalance = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance  = sdkmath.NewInt(1000).MulRaw(1e18)
	)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// deploy contract to staking
	contractAddr, txHash := suite.DeployStaking(delSigner)
	txFee1 := suite.TxFee(txHash)

	stakingSuite := NewStakingSuite(suite.EthSuite, contractAddr, delSigner)

	// delegate by contract
	receipt := stakingSuite.DelegateV2(valAddr.String(), delBalance.BigInt())
	txFee2 := suite.TxFee(receipt.TxHash)

	delBal := suite.GetAllBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(sdkmath.NewIntFromBigInt(txFee2)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Require().Equal(initBalance.String(), total.String())

	// query delegate by contract
	shares, amount := stakingSuite.Delegation(valAddr.String(), contractAddr)
	suite.Require().Equal(delBalance.BigInt().String(), amount.String())

	// withdraw by contract
	delBal = suite.GetAllBalances(contractAddr.Bytes())
	suite.Require().True(delBal.IsZero())
	receipt = stakingSuite.Withdraw(valAddr.String())
	txFee3 := suite.TxFee(receipt.TxHash)

	delBal = suite.GetAllBalances(contractAddr.Bytes())
	suite.Require().True(delBal.AmountOf(fxtypes.DefaultDenom).GT(sdkmath.NewInt(0)))

	// undelegate by contract
	receipt = stakingSuite.UnDelegateV2(valAddr.String(), shares)
	txFee4 := suite.TxFee(receipt.TxHash)

	delBal = suite.GetAllBalances(delSigner.AccAddress())
	txFee := sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(
		sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee4))))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Require().Equal(initBalance.String(), total.String())

	// query delegate by contract
	shares, amount = stakingSuite.Delegation(valAddr.String(), contractAddr)
	suite.Require().Equal(amount.String(), big.NewInt(0).String())
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	resp, err := suite.StakingQuery().UnbondingDelegation(suite.ctx, &stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: sdk.AccAddress(contractAddr.Bytes()).String(),
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
		valAddr            = suite.GetValAddr()
		initBalance        = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance         = sdkmath.NewInt(1000).MulRaw(1e18)
		receiptInitBalance = sdkmath.NewInt(100).MulRaw(1e18)
	)

	stakingSuite := NewStakingSuite(suite.EthSuite, common.HexToAddress(contract.StakingAddress), delSigner)
	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	receiptAccAddr := receiptSigner.AccAddress()
	receiptEthAddr := receiptSigner.Address()
	suite.Send(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptInitBalance))

	// delegate
	receipt := stakingSuite.DelegateV2(valAddr.String(), delBalance.BigInt())
	txFee1 := suite.TxFee(receipt.TxHash)

	// check receipt delegate
	_, amount := stakingSuite.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(big.NewInt(0).String(), amount.String())

	// check del delegate
	delSignerAddr := delSigner.Address()
	shares, amount := stakingSuite.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(delBalance.BigInt().String(), amount.String())

	halfShares := big.NewInt(0).Div(shares, big.NewInt(2))

	// transfer shares
	receipt = stakingSuite.TransferShares(valAddr.String(), receiptEthAddr, halfShares)
	txFee2 := suite.TxFee(receipt.TxHash)

	reward1 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), delSignerAddr)

	// check del delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(shares, halfShares)

	// check receipt delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(shares, halfShares)

	// check del balance
	delBal := suite.GetAllBalances(delSigner.AccAddress())
	txFee := sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2))
	total := delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1))
	suite.Require().Equal(initBalance.String(), total.String())

	// transfer shares
	receipt = stakingSuite.TransferShares(valAddr.String(), receiptEthAddr, halfShares)
	txFee3 := suite.TxFee(receipt.TxHash)

	reward2 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), delSignerAddr)
	reward3 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), receiptEthAddr)

	// check del delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check receipt delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check del balance
	delBal = suite.GetAllBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3)))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1)).Sub(sdkmath.NewIntFromBigInt(reward2))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	receiptExpectBalance := receiptInitBalance.Add(sdkmath.NewIntFromBigInt(reward3))
	suite.EqualBalance(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptExpectBalance))

	// approve
	receipt = stakingSuite.WithSigner(receiptSigner).ApproveShares(valAddr.String(), delSignerAddr, big.NewInt(0).Mul(big.NewInt(3), halfShares))
	txFee4 := suite.TxFee(receipt.TxHash)

	// check approve
	allowance := stakingSuite.AllowanceShares(valAddr.String(), receiptEthAddr, delSignerAddr)
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(3), halfShares).String(), allowance.String())

	// check receipt balance
	receiptExpectBalance = receiptExpectBalance.Sub(sdkmath.NewIntFromBigInt(txFee4))
	suite.EqualBalance(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptExpectBalance))

	// transfer from
	receipt = stakingSuite.TransferFromShares(valAddr.String(), receiptEthAddr, delSignerAddr, halfShares)
	txFee5 := suite.TxFee(receipt.TxHash)

	reward4 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), receiptEthAddr)
	receiptExpectBalance = receiptExpectBalance.Add(sdkmath.NewIntFromBigInt(reward4))

	// check del delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(shares.String(), halfShares.String())

	// check receipt delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(shares.String(), halfShares.String())

	// check del balance
	delBal = suite.GetAllBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee5))))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1)).Sub(sdkmath.NewIntFromBigInt(reward2))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	suite.EqualBalance(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptExpectBalance))

	// check approve
	allowance = stakingSuite.AllowanceShares(valAddr.String(), receiptEthAddr, delSignerAddr)
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(2), halfShares).String(), allowance.String())

	// transfer from
	receipt = stakingSuite.TransferFromShares(valAddr.String(), receiptEthAddr, delSignerAddr, halfShares)
	txFee6 := suite.TxFee(receipt.TxHash)

	reward5 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), delSignerAddr)
	reward6 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), receiptEthAddr)
	receiptExpectBalance = receiptExpectBalance.Add(sdkmath.NewIntFromBigInt(reward6))

	// check del delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), delSignerAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check receipt delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), receiptEthAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check del balance
	delBal = suite.GetAllBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee5).Add(sdkmath.NewIntFromBigInt(txFee6)))))
	totalReward := big.NewInt(0).Add(big.NewInt(0).Add(reward1, reward2), reward5)
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(totalReward))
	suite.Require().Equal(initBalance.String(), total.String())

	// check receipt balance
	suite.EqualBalance(receiptAccAddr, sdk.NewCoin(fxtypes.DefaultDenom, receiptExpectBalance))

	// check approve
	allowance = stakingSuite.AllowanceShares(valAddr.String(), receiptEthAddr, delSignerAddr)
	suite.Require().Equal(halfShares.String(), allowance.String())
}

func (suite *IntegrationTest) StakingSharesContractTest() {
	var (
		delSigner   = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr     = suite.GetValAddr()
		initBalance = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance  = sdkmath.NewInt(1000).MulRaw(1e18)
	)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// deploy contract to staking
	contractAddr, txHash := suite.DeployStaking(delSigner)
	txFee1 := suite.TxFee(txHash)

	stakingSuite := NewStakingSuite(suite.EthSuite, contractAddr, delSigner)

	// delegate
	receipt := stakingSuite.DelegateV2(valAddr.String(), delBalance.BigInt())
	txFee2 := suite.TxFee(receipt.TxHash)

	// check del balance
	delBal := suite.GetAllBalances(delSigner.AccAddress())
	txFee := sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2))
	total := delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract delegate
	_, amount := stakingSuite.Delegation(valAddr.String(), contractAddr)
	suite.Require().Equal(big.NewInt(0).String(), amount.String())

	// check del delegate
	shares, amount := stakingSuite.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(delBalance.BigInt().String(), amount.String())

	halfShares := big.NewInt(0).Div(shares, big.NewInt(2))

	// approve
	receipt = stakingSuite.ApproveShares(valAddr.String(), contractAddr, big.NewInt(0).Mul(big.NewInt(3), halfShares))
	txFee3 := suite.TxFee(receipt.TxHash)

	// check approve
	allowance := stakingSuite.AllowanceShares(valAddr.String(), delSigner.Address(), contractAddr)
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(3), halfShares).String(), allowance.String())

	// transferFromShares
	receipt = stakingSuite.TransferFromShares(valAddr.String(), delSigner.Address(), contractAddr, halfShares)
	txFee4 := suite.TxFee(receipt.TxHash)

	reward1 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())

	// check del delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), halfShares.String())

	// check contract delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), contractAddr)
	suite.Require().Equal(shares.String(), halfShares.String())

	// check del balance
	delBal = suite.GetAllBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee4))))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdkmath.NewIntFromBigInt(reward1))
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract balance
	contractBal := suite.GetAllBalances(contractAddr.Bytes())
	suite.Require().Equal(contractBal.AmountOf(fxtypes.DefaultDenom).String(), big.NewInt(0).String())

	// check approve
	allowance = stakingSuite.AllowanceShares(valAddr.String(), delSigner.Address(), contractAddr)
	suite.Require().Equal(big.NewInt(0).Mul(big.NewInt(2), halfShares).String(), allowance.String())

	// transferFromShares
	receipt = stakingSuite.TransferFromShares(valAddr.String(), delSigner.Address(), contractAddr, halfShares)
	txFee5 := suite.TxFee(receipt.TxHash)
	reward2 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())
	reward3 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), contractAddr)

	// check del delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check contract delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), contractAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check del balance
	delBal = suite.GetAllBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).Add(sdkmath.NewIntFromBigInt(txFee4).Add(sdkmath.NewIntFromBigInt(txFee5)))))
	totalReward := sdkmath.NewIntFromBigInt(reward1).Add(sdkmath.NewIntFromBigInt(reward2))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(totalReward)
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract balance
	contractBal = suite.GetAllBalances(contractAddr.Bytes())
	suite.Require().Equal(contractBal.AmountOf(fxtypes.DefaultDenom).String(), reward3.String())

	// check approve
	allowance = stakingSuite.AllowanceShares(valAddr.String(), delSigner.Address(), contractAddr)
	suite.Require().Equal(halfShares.String(), allowance.String())

	// contract transfer
	receipt = stakingSuite.TransferShares(valAddr.String(), delSigner.Address(), halfShares)
	txFee6 := suite.TxFee(receipt.TxHash)
	reward4 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), contractAddr)

	// check del delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), halfShares.String())

	// check contract delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), contractAddr)
	suite.Require().Equal(shares.String(), halfShares.String())

	// check del balance
	delBal = suite.GetAllBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).
		Add(sdkmath.NewIntFromBigInt(txFee4).Add(sdkmath.NewIntFromBigInt(txFee5).Add(sdkmath.NewIntFromBigInt(txFee6))))))
	totalReward = sdkmath.NewIntFromBigInt(reward1).Add(sdkmath.NewIntFromBigInt(reward2))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(totalReward)
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract balance
	contractBal = suite.GetAllBalances(contractAddr.Bytes())
	suite.Require().Equal(contractBal.AmountOf(fxtypes.DefaultDenom).String(), big.NewInt(0).Add(reward3, reward4).String())

	// contract transfer
	receipt = stakingSuite.TransferShares(valAddr.String(), delSigner.Address(), halfShares)
	txFee7 := suite.TxFee(receipt.TxHash)
	reward5 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), delSigner.Address())
	reward6 := stakingSuite.LogReward(receipt.Logs, valAddr.String(), contractAddr)

	// check del delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), delSigner.Address())
	suite.Require().Equal(shares.String(), big.NewInt(0).Mul(big.NewInt(2), halfShares).String())

	// check contract delegate
	shares, _ = stakingSuite.Delegation(valAddr.String(), contractAddr)
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	// check del balance
	delBal = suite.GetAllBalances(delSigner.AccAddress())
	txFee = sdkmath.NewIntFromBigInt(txFee1).Add(sdkmath.NewIntFromBigInt(txFee2).Add(sdkmath.NewIntFromBigInt(txFee3).
		Add(sdkmath.NewIntFromBigInt(txFee4).Add(sdkmath.NewIntFromBigInt(txFee5).Add(sdkmath.NewIntFromBigInt(txFee6).Add(sdkmath.NewIntFromBigInt(txFee7)))))))
	totalReward = sdkmath.NewIntFromBigInt(reward1).Add(sdkmath.NewIntFromBigInt(reward2).Add(sdkmath.NewIntFromBigInt(reward5)))
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(totalReward)
	suite.Require().Equal(initBalance.String(), total.String())

	// check contract balance
	contractBal = suite.GetAllBalances(contractAddr.Bytes())
	totalReward = sdkmath.NewIntFromBigInt(reward3).Add(sdkmath.NewIntFromBigInt(reward4).Add(sdkmath.NewIntFromBigInt(reward6)))
	suite.Require().Equal(contractBal.AmountOf(fxtypes.DefaultDenom).String(), totalReward.String())
}

func (suite *IntegrationTest) StakingPrecompileRedelegateTest() {
	var (
		delSigner    = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr      = suite.GetValAddr()
		valNewSigner = helpers.NewSigner(helpers.NewEthPrivKey())
		initBalance  = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance   = sdkmath.NewInt(1000).MulRaw(1e18)
	)

	stakingSuite := NewStakingSuite(suite.EthSuite, common.HexToAddress(contract.StakingAddress), delSigner)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))
	suite.Send(valNewSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// delegate
	receipt := stakingSuite.DelegateV2(valAddr.String(), delBalance.BigInt())
	txFee1 := suite.TxFee(receipt.TxHash)

	delBal := suite.GetAllBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())

	// set WithdrawAddress
	rewardAddress := sdk.AccAddress(helpers.NewEthPrivKey().PubKey().Address().Bytes())
	txRsp := suite.SetWithdrawAddress(stakingSuite.signer, rewardAddress)
	gasPrice, err := sdk.ParseCoinNormalized(suite.network.Config.MinGasPrices)
	suite.Require().NoError(err)
	gasFee := gasPrice.Amount.Mul(sdkmath.NewInt(txRsp.GasWanted))

	hexAddr := common.BytesToAddress(delSigner.AccAddress().Bytes())
	// query delegate
	valAddrShares1, _ := stakingSuite.Delegation(valAddr.String(), hexAddr)

	resp := suite.CreateValidator(valNewSigner, false)
	suite.Equal(resp.Code, uint32(0))

	receipt = stakingSuite.RedelegateV2(valAddr.String(), sdk.ValAddress(valNewSigner.AccAddress()).String(), valAddrShares1)
	txFee2 := suite.TxFee(receipt.TxHash)

	valAddrShares2, _ := stakingSuite.Delegation(valAddr.String(), hexAddr)
	suite.Equal(big.NewInt(0).String(), valAddrShares2.String())

	valNewShares, _ := stakingSuite.Delegation(sdk.ValAddress(valNewSigner.AccAddress()).String(), hexAddr)
	suite.Equal(valAddrShares1, valNewShares)

	delBal = suite.GetAllBalances(delSigner.AccAddress())
	total = delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).
		Add(gasFee).
		Add(sdkmath.NewIntFromBigInt(txFee2)).
		Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())
}

func (suite *IntegrationTest) StakingPrecompileRedelegateByContractTest() {
	var (
		delSigner   = helpers.NewSigner(helpers.NewEthPrivKey())
		valAddr     = suite.GetValAddr()
		valNew      = helpers.NewSigner(helpers.NewEthPrivKey())
		initBalance = sdkmath.NewInt(2000).MulRaw(1e18)
		delBalance  = sdkmath.NewInt(1000).MulRaw(1e18)
	)

	suite.Send(delSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))
	suite.Send(valNew.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// deploy contract to staking
	contractAddr, txHash := suite.DeployStaking(delSigner)
	txFee1 := suite.TxFee(txHash)

	stakingSuite := NewStakingSuite(suite.EthSuite, contractAddr, delSigner)

	// delegate by contract
	receipt := stakingSuite.DelegateV2(valAddr.String(), delBalance.BigInt())
	txFee2 := suite.TxFee(receipt.TxHash)

	delBal := suite.GetAllBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(sdkmath.NewIntFromBigInt(txFee2)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())

	// query delegate
	valAddrShares1, _ := stakingSuite.Delegation(valAddr.String(), contractAddr)

	resp := suite.CreateValidator(valNew, false)
	suite.Equal(resp.Code, uint32(0))

	receipt = stakingSuite.RedelegateV2(valAddr.String(), sdk.ValAddress(valNew.AccAddress()).String(), valAddrShares1)
	txFee3 := suite.TxFee(receipt.TxHash)

	valAddrShares2, _ := stakingSuite.Delegation(valAddr.String(), contractAddr)
	suite.Equal(big.NewInt(0).String(), valAddrShares2.String())

	valNewShares, _ := stakingSuite.Delegation(sdk.ValAddress(valNew.AccAddress()).String(), contractAddr)
	suite.Equal(valAddrShares1, valNewShares)

	delBal = suite.GetAllBalances(delSigner.AccAddress())
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

	stakingSuite := NewStakingSuite(suite.EthSuite, common.HexToAddress(contract.StakingAddress), delSigner)

	// 2. delegate to first validator
	valAddr := suite.GetValAddr()
	delBalance := sdkmath.NewInt(1000)
	receipt := stakingSuite.DelegateV2(valAddr.String(), delBalance.BigInt())
	txFee1 := suite.TxFee(receipt.TxHash)

	delBal := suite.GetAllBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())

	hexAddr := common.BytesToAddress(delSigner.AccAddress().Bytes())
	_, delAmt1 := stakingSuite.Delegation(valAddr.String(), hexAddr)
	suite.Equal(delAmt1.String(), delBalance.String())

	halfDelegateAmount := big.NewInt(0).Div(delBalance.BigInt(), big.NewInt(2))

	// 2. undelegate half of the first validator amount
	stakingSuite.UnDelegateV2(valAddr.String(), halfDelegateAmount)
	_, delAmt2 := stakingSuite.Delegation(valAddr.String(), hexAddr)
	suite.Equal(halfDelegateAmount.String(), delAmt2.String())

	// 3. create a new validator, and redelegate half of the first validator amount to it(new validator is not bonded)
	valNewSigner := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(valNewSigner.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))
	resp := suite.CreateValidator(valNewSigner, false)
	suite.Equal(resp.Code, uint32(0))

	// 4. redelegate half of the first validator amount to new validator
	stakingSuite.RedelegateV2(valAddr.String(), sdk.ValAddress(valNewSigner.AccAddress()).String(), halfDelegateAmount)

	delShare, _ := stakingSuite.Delegation(valAddr.String(), hexAddr)
	suite.Equal(int64(0), delShare.Int64())

	// 5. check new validator's delegation amount, expecting half of the first validator amount
	_, delAmt3 := stakingSuite.Delegation(sdk.ValAddress(valNewSigner.AccAddress()).String(), hexAddr)
	suite.Equal(halfDelegateAmount.String(), delAmt3.String())

	{
		// finally, clear the new validator
		stakingSuite.UnDelegateV2(sdk.ValAddress(valNewSigner.AccAddress()).String(), halfDelegateAmount)
		_, valSelfDelegation := stakingSuite.Delegation(sdk.ValAddress(valNewSigner.AccAddress()).String(), valNewSigner.Address())
		stakingSuite.WithSigner(valNewSigner).UnDelegateV2(sdk.ValAddress(valNewSigner.AccAddress()).String(), valSelfDelegation)
	}
}
