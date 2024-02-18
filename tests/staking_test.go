package tests

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/metadata"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	precompilesstaking "github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
	fxstakingtypes "github.com/functionx/fx-core/v7/x/staking/types"
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
	suite.staking.UnDelegate(suite.staking.privKey, valAddr.String(), delBalance.BigInt())
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
	receipt := suite.staking.Delegate(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
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
	total = delBalance.Add(txFee).Add(delBal.AmountOf(fxtypes.DefaultDenom)).Sub(sdk.NewIntFromBigInt(totalReward))
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
	receipt := suite.staking.Delegate(delSigner.PrivKey(), valAddr.String(), delBalance.BigInt())
	txFee1 := suite.evm.TxFee(receipt.TxHash)

	delBal := suite.QueryBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())

	// set WithdrawAddress
	rewardAddress := sdk.AccAddress(helpers.NewEthPrivKey().PubKey().Address().Bytes())
	txRsp := suite.staking.SetWithdrawAddressWithResponse(delSigner.PrivKey(), rewardAddress)
	gasPrice, err := sdk.ParseCoinNormalized(suite.network.Config.MinGasPrices)
	suite.NoError(err)
	gasFee := gasPrice.Amount.Mul(sdkmath.NewInt(txRsp.GasWanted))

	hexAddr := common.BytesToAddress(delSigner.AccAddress().Bytes())
	// query delegate
	valAddrShares1, _ := suite.staking.Delegation(valAddr.String(), hexAddr)

	resp := suite.staking.CreateValidator(valNew.PrivKey())
	suite.Equal(resp.Code, uint32(0))

	receipt = suite.staking.Redelegate(delSigner.PrivKey(), valAddr.String(), sdk.ValAddress(valNew.AccAddress()).String(), valAddrShares1)
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
	receipt := suite.staking.DelegateByContract(delSigner.PrivKey(), contract, valAddr.String(), delBalance.BigInt())
	txFee2 := suite.evm.TxFee(receipt.TxHash)

	delBal := suite.QueryBalances(delSigner.AccAddress())
	total := delBalance.Add(sdkmath.NewIntFromBigInt(txFee1)).Add(sdkmath.NewIntFromBigInt(txFee2)).Add(delBal.AmountOf(fxtypes.DefaultDenom))
	suite.Equal(initBalance.String(), total.String())

	// query delegate
	valAddrShares1, _ := suite.staking.Delegation(valAddr.String(), contract)

	resp := suite.staking.CreateValidator(valNew.PrivKey())
	suite.Equal(resp.Code, uint32(0))

	receipt = suite.staking.RedelegateByContract(delSigner.PrivKey(), contract, valAddr.String(), sdk.ValAddress(valNew.AccAddress()).String(), valAddrShares1)
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

func (suite *IntegrationMultiNodeTest) StakingGrantPrivilege() {
	initBalance := sdkmath.NewInt(2000).MulRaw(1e18)
	suite.Send(suite.staking.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	allVals := suite.GetAllValidators()
	valAddr := allVals[1].ValAddress

	from := sdk.AccAddress(valAddr)
	to := sdk.AccAddress(suite.staking.privKey.PubKey().Address())

	to2 := suite.staking.GrantAddress()
	suite.Send(to2, sdk.NewCoin(fxtypes.DefaultDenom, initBalance))

	// check val with to grants
	grantResp, err := suite.authz.AuthzQuery().Grants(suite.ctx, &authz.QueryGrantsRequest{Granter: from.String(), Grantee: to.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(0, len(grantResp.Grants))

	// grant privilege val with to address
	sign, err := suite.staking.privKey.Sign(fxstakingtypes.GrantPrivilegeSignatureData(valAddr, from, to))
	suite.Require().NoError(err)
	msg, err := fxstakingtypes.NewMsgGrantPrivilege(valAddr, sdk.AccAddress(valAddr), suite.staking.privKey.PubKey(), hex.EncodeToString(sign))
	suite.Require().NoError(err)
	suite.BroadcastTx(suite.GetValidatorPrivKeys(from), msg)

	// check val with to grants
	grantResp, err = suite.authz.AuthzQuery().Grants(suite.ctx, &authz.QueryGrantsRequest{Granter: from.String(), Grantee: to.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(grantResp.Grants))

	// val send coins error, not have privilege
	msgSend := banktypes.NewMsgSend(from, to, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100))))
	_, err = suite.GRPCClient().BuildTxV1(suite.GetValidatorPrivKeys(from), []sdk.Msg{msgSend}, 500000, "", 0)
	suite.Require().Error(err)

	// grant MsgSend privilege with to address
	exp := time.Now().Add(time.Hour)
	msgGrant, err := authz.NewMsgGrant(from, to, authz.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgSend{})), &exp)
	suite.Require().NoError(err)
	msgExec := authz.NewMsgExec(to, []sdk.Msg{msgGrant, msgSend})

	_, err = suite.GRPCClient().BuildTxV1(suite.staking.privKey, []sdk.Msg{&msgExec}, 500000, "", 0)
	suite.Require().NoError(err)
	tx := suite.BroadcastTx(suite.staking.privKey, &msgExec)
	suite.Require().Equal(uint32(0), tx.Code)

	// check val with to grants
	grantResp, err = suite.authz.AuthzQuery().Grants(suite.ctx, &authz.QueryGrantsRequest{Granter: from.String(), Grantee: to.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(2, len(grantResp.Grants))

	// val grant to to2 address error, val not have privilege
	sign, err = suite.staking.GrantPrivKey().Sign(fxstakingtypes.GrantPrivilegeSignatureData(valAddr, from, to2))
	suite.Require().NoError(err)
	msg, err = fxstakingtypes.NewMsgGrantPrivilege(valAddr, from, suite.staking.GrantPrivKey().PubKey(), hex.EncodeToString(sign))
	suite.Require().NoError(err)
	_, err = suite.GRPCClient().BuildTxV1(suite.staking.privKey, []sdk.Msg{msg}, 500000, "", 0)
	suite.Require().Error(err)

	// to grant to to2 address
	sign, err = suite.staking.GrantPrivKey().Sign(fxstakingtypes.GrantPrivilegeSignatureData(valAddr, to, to2))
	suite.Require().NoError(err)
	msg, err = fxstakingtypes.NewMsgGrantPrivilege(valAddr, to, suite.staking.GrantPrivKey().PubKey(), hex.EncodeToString(sign))
	suite.Require().NoError(err)
	_, err = suite.GRPCClient().BuildTxV1(suite.staking.privKey, []sdk.Msg{msg}, 500000, "", 0)
	suite.NoError(err)
	suite.BroadcastTx(suite.staking.privKey, msg)

	// check val with to grants
	grantResp, err = suite.authz.AuthzQuery().Grants(suite.ctx, &authz.QueryGrantsRequest{Granter: from.String(), Grantee: to.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(0, len(grantResp.Grants))

	// check to with to2 grants
	grantResp, err = suite.authz.AuthzQuery().Grants(suite.ctx, &authz.QueryGrantsRequest{Granter: from.String(), Grantee: to2.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(grantResp.Grants))

	//  grant privilege edit consenus pubkey
	msgGrant, err = authz.NewMsgGrant(from, suite.staking.GrantAddress(), authz.NewGenericAuthorization(sdk.MsgTypeURL(&fxstakingtypes.MsgEditConsensusPubKey{})), &exp)
	suite.Require().NoError(err)

	// edit consenus pubkey
	newPriKey := ed25519.GenPrivKey()
	msgEdit, err := fxstakingtypes.NewMsgEditConsensusPubKey(valAddr, suite.staking.GrantAddress(), newPriKey.PubKey())
	suite.Require().NoError(err)

	msgExec = authz.NewMsgExec(suite.staking.GrantAddress(), []sdk.Msg{msgGrant, msgEdit})
	tx = suite.BroadcastTx(suite.staking.GrantPrivKey(), &msgExec)
	suite.Require().Equal(uint32(0), tx.Code)
}

func (suite *IntegrationMultiNodeTest) StakingEditPubKey() {
	if suite.QueryBalances(suite.staking.AccAddress()).IsZero() {
		suite.Send(suite.staking.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(2000).MulRaw(1e18)))
	}

	// query all validator
	allVals := suite.GetAllValidators()
	valAddr := allVals[2].ValAddress
	valResp, err := suite.staking.StakingQuery().Validator(suite.ctx, &stakingtypes.QueryValidatorRequest{ValidatorAddr: valAddr.String()})
	suite.Require().NoError(err)
	//	val addr
	validator := valResp.Validator
	// val consensus pubkey
	var oldPubKey cryptotypes.PubKey
	err = app.MakeEncodingConfig().InterfaceRegistry.UnpackAny(validator.ConsensusPubkey, &oldPubKey)
	suite.Require().NoError(err)
	valFrom := sdk.AccAddress(valAddr)
	// new consenus pubkey
	newPriKey := ed25519.GenPrivKey()

	// check signing info
	infos, err := suite.slasing.SlashingQuery().SigningInfos(suite.ctx, &slashingtypes.QuerySigningInfosRequest{})
	suite.Require().NoError(err)
	for _, info := range infos.Info {
		suite.True(info.MissedBlocksCounter == 0)
	}

	// edit consensus pubkey to new consenus pubkey
	msg, err := fxstakingtypes.NewMsgEditConsensusPubKey(valAddr, valFrom, newPriKey.PubKey())
	suite.Require().NoError(err)
	suite.BroadcastTx(suite.GetValidatorPrivKeys(valFrom), msg)

	// wait 15 block
	_, err = suite.network.WaitNumberBlock(15)
	suite.Require().NoError(err)

	// check signing info
	infos, err = suite.slasing.SlashingQuery().SigningInfos(suite.ctx, &slashingtypes.QuerySigningInfosRequest{})
	suite.Require().NoError(err)
	for _, info := range infos.Info {
		if info.Address == sdk.ConsAddress(newPriKey.PubKey().Address()).String() {
			suite.True(info.MissedBlocksCounter > 0)
			continue
		}
		suite.True(info.MissedBlocksCounter == 0)
	}

	// validator jailed
	valResp, err = suite.staking.StakingQuery().Validator(suite.ctx, &stakingtypes.QueryValidatorRequest{ValidatorAddr: valAddr.String()})
	suite.Require().NoError(err)
	suite.Require().True(valResp.Validator.Jailed)

	// edit consensus pubkey to old consenus pubkey
	msg, err = fxstakingtypes.NewMsgEditConsensusPubKey(valAddr, valFrom, oldPubKey)
	suite.Require().NoError(err)
	suite.BroadcastTx(suite.GetValidatorPrivKeys(valFrom), msg)

	// check signing info
	infos, err = suite.slasing.SlashingQuery().SigningInfos(suite.ctx, &slashingtypes.QuerySigningInfosRequest{})
	suite.Require().NoError(err)
	for _, info := range infos.Info {
		if info.Address == sdk.ConsAddress(oldPubKey.Address()).String() || info.Address == sdk.ConsAddress(newPriKey.PubKey().Address()).String() {
			suite.True(info.MissedBlocksCounter > 0)
			continue
		}
		suite.True(info.MissedBlocksCounter == 0)
	}

	// unjail validator
	msgUnjail := slashingtypes.NewMsgUnjail(valAddr)
	suite.BroadcastTx(suite.GetValidatorPrivKeys(valFrom), msgUnjail)

	// check validator unjailed
	valResp, err = suite.staking.StakingQuery().Validator(suite.ctx, &stakingtypes.QueryValidatorRequest{ValidatorAddr: valAddr.String()})
	suite.Require().NoError(err)
	suite.Require().False(valResp.Validator.Jailed)

	// check signing info
	infoResp, err := suite.slasing.SlashingQuery().SigningInfo(suite.ctx, &slashingtypes.QuerySigningInfoRequest{ConsAddress: sdk.ConsAddress(oldPubKey.Address()).String()})
	suite.Require().NoError(err)
	missBlock := infoResp.ValSigningInfo.MissedBlocksCounter

	// wait 15 block
	_, err = suite.network.WaitNumberBlock(15)
	suite.Require().NoError(err)

	// check signing info
	infoResp, err = suite.slasing.SlashingQuery().SigningInfo(suite.ctx, &slashingtypes.QuerySigningInfoRequest{ConsAddress: sdk.ConsAddress(oldPubKey.Address()).String()})
	suite.Require().NoError(err)
	suite.Require().Equal(missBlock, infoResp.ValSigningInfo.MissedBlocksCounter)
}

func (suite *IntegrationMultiNodeTest) StakingEditPubKeyJailBlock() {
	if suite.QueryBalances(suite.staking.AccAddress()).IsZero() {
		suite.Send(suite.staking.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(2000).MulRaw(1e18)))
	}

	// query all validator
	allVals := suite.GetAllValidators()
	valAddr := allVals[3].ValAddress
	valResp, err := suite.staking.StakingQuery().Validator(suite.ctx, &stakingtypes.QueryValidatorRequest{ValidatorAddr: valAddr.String()})
	suite.Require().NoError(err)
	validator := valResp.Validator

	// val consensus pubkey
	var oldPubKey cryptotypes.PubKey
	err = app.MakeEncodingConfig().InterfaceRegistry.UnpackAny(validator.ConsensusPubkey, &oldPubKey)
	suite.Require().NoError(err)
	oldConsAddr := sdk.ConsAddress(oldPubKey.Address())

	valFrom := sdk.AccAddress(valAddr)
	// new consenus pubkey
	newPriKey := ed25519.GenPrivKey()
	newPubKey := newPriKey.PubKey()
	newConsAddr := sdk.ConsAddress(newPubKey.Address())

	// edit consensus pubkey to new consenus pubkey and undelegate all
	editPubkeyMsg, err := fxstakingtypes.NewMsgEditConsensusPubKey(valAddr, valFrom, newPubKey)
	suite.Require().NoError(err)
	undelegateMsg := stakingtypes.NewMsgUndelegate(sdk.AccAddress(valAddr), valAddr, sdk.NewCoin(fxtypes.DefaultDenom, validator.Tokens))
	txResp := suite.BroadcastTx(suite.GetValidatorPrivKeys(valFrom), editPubkeyMsg, undelegateMsg)
	suite.Require().Equal(uint32(0), txResp.Code)
	suite.Require().Greater(txResp.Height, int64(0))

	_, _ = suite.network.WaitNumberBlock(10)

	valResp, err = suite.staking.StakingQuery().Validator(suite.ctx, &stakingtypes.QueryValidatorRequest{ValidatorAddr: valAddr.String()})
	suite.Require().NoError(err)
	suite.Require().True(valResp.Validator.Jailed)

	height := txResp.Height + 1 // val update, edit skip to next block

	// block tx process
	ctx := metadata.AppendToOutgoingContext(suite.ctx, grpctypes.GRPCBlockHeightHeader, fmt.Sprintf("%d", height))
	info1, err := suite.slasing.SlashingQuery().SigningInfo(ctx, &slashingtypes.QuerySigningInfoRequest{ConsAddress: oldConsAddr.String()})
	suite.NoError(err)
	info2, err := suite.slasing.SlashingQuery().SigningInfo(ctx, &slashingtypes.QuerySigningInfoRequest{ConsAddress: newConsAddr.String()})
	suite.NoError(err)
	suite.Equal(info1.ValSigningInfo.IndexOffset, info2.ValSigningInfo.IndexOffset)
	suite.Equal(info1.ValSigningInfo.MissedBlocksCounter, info2.ValSigningInfo.MissedBlocksCounter)

	// block tx process +1
	ctx = metadata.AppendToOutgoingContext(suite.ctx, grpctypes.GRPCBlockHeightHeader, fmt.Sprintf("%d", height+1))
	info1, err = suite.slasing.SlashingQuery().SigningInfo(ctx, &slashingtypes.QuerySigningInfoRequest{ConsAddress: oldConsAddr.String()})
	suite.NoError(err)
	info2, err = suite.slasing.SlashingQuery().SigningInfo(ctx, &slashingtypes.QuerySigningInfoRequest{ConsAddress: newConsAddr.String()})
	suite.NoError(err)
	suite.Equal(info1.ValSigningInfo.IndexOffset, info2.ValSigningInfo.IndexOffset)
	suite.Equal(info1.ValSigningInfo.MissedBlocksCounter, info2.ValSigningInfo.MissedBlocksCounter)

	// block tx process +2
	ctx = metadata.AppendToOutgoingContext(suite.ctx, grpctypes.GRPCBlockHeightHeader, fmt.Sprintf("%d", height+2))
	_, err = suite.slasing.SlashingQuery().SigningInfo(ctx, &slashingtypes.QuerySigningInfoRequest{ConsAddress: oldConsAddr.String()})
	suite.Error(err)
	info2, err = suite.slasing.SlashingQuery().SigningInfo(ctx, &slashingtypes.QuerySigningInfoRequest{ConsAddress: newConsAddr.String()})
	suite.NoError(err)
	suite.Equal(info1.ValSigningInfo.IndexOffset, info2.ValSigningInfo.IndexOffset)

	// delegate and unjail
	delegateMsg := stakingtypes.NewMsgDelegate(sdk.AccAddress(valAddr), valAddr, sdk.NewCoin(fxtypes.DefaultDenom, validator.Tokens))
	unjailMsg := slashingtypes.NewMsgUnjail(valAddr)
	editPubkeyMsg, err = fxstakingtypes.NewMsgEditConsensusPubKey(valAddr, valFrom, oldPubKey)
	suite.Require().NoError(err)
	_ = suite.BroadcastTx(suite.GetValidatorPrivKeys(valFrom), delegateMsg, unjailMsg, editPubkeyMsg)

	_, _ = suite.network.WaitNumberBlock(3)

	valResp, err = suite.staking.StakingQuery().Validator(suite.ctx, &stakingtypes.QueryValidatorRequest{ValidatorAddr: valAddr.String()})
	suite.Require().NoError(err)
	suite.Require().False(valResp.Validator.Jailed)
}
