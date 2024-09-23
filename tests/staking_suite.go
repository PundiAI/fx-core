package tests

import (
	"context"
	"math/big"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/v8/client"
	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	stakingprecompile "github.com/functionx/fx-core/v8/x/staking/precompile"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type StakingSuite struct {
	Erc20TestSuite
	abi      abi.ABI
	grantKey cryptotypes.PrivKey
}

func NewStakingSuite(ts *TestSuite) StakingSuite {
	key := helpers.NewEthPrivKey()
	return StakingSuite{
		Erc20TestSuite: NewErc20TestSuite(ts),
		abi:            stakingprecompile.GetABI(),
		grantKey:       key,
	}
}

func (suite *StakingSuite) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.privKey.PubKey().Address())
}

func (suite *StakingSuite) Address() common.Address {
	return common.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *StakingSuite) GrantPrivKey() cryptotypes.PrivKey {
	return suite.grantKey
}

func (suite *StakingSuite) GrantAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.grantKey.PubKey().Address())
}

func (suite *StakingSuite) StakingQuery() stakingtypes.QueryClient {
	return suite.GRPCClient().StakingQuery()
}

func (suite *StakingSuite) TransactionOpts(privateKey cryptotypes.PrivKey) *bind.TransactOpts {
	ecdsa, err := crypto.ToECDSA(privateKey.Bytes())
	suite.Require().NoError(err)

	chainId, err := suite.EthClient().ChainID(suite.ctx)
	suite.Require().NoError(err)

	auth, err := bind.NewKeyedTransactorWithChainID(ecdsa, chainId)
	suite.Require().NoError(err)

	auth.GasTipCap = big.NewInt(1e9)
	auth.GasFeeCap = big.NewInt(6e12)
	return auth
}

func (suite *StakingSuite) DeployStakingContract(privKey cryptotypes.PrivKey) (common.Address, common.Hash) {
	stakingBin := contract.MustDecodeHex(testscontract.StakingTestMetaData.Bin)
	return suite.DeployContract(privKey, stakingBin)
}

// DelegationRewards Get delegatorAddress rewards
func (suite *StakingSuite) DelegationRewards(delAddr, valAddr string) sdk.DecCoins {
	response, err := suite.GRPCClient().DistrQuery().DelegationRewards(suite.ctx, &distrtypes.QueryDelegationRewardsRequest{DelegatorAddress: delAddr, ValidatorAddress: valAddr})
	suite.Require().NoError(err)
	return response.Rewards
}

func (suite *StakingSuite) SetWithdrawAddress(withdrawAddr sdk.AccAddress) {
	suite.SetWithdrawAddressWithResponse(suite.privKey, withdrawAddr)
}

func (suite *StakingSuite) SetWithdrawAddressWithResponse(privKey cryptotypes.PrivKey, withdrawAddr sdk.AccAddress) *sdk.TxResponse {
	delAddr := sdk.AccAddress(privKey.PubKey().Address())
	setWithdrawAddress := distrtypes.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)
	txResponse := suite.BroadcastTx(privKey, setWithdrawAddress)
	suite.Require().True(txResponse.Code == 0)
	response, err := suite.GRPCClient().DistrQuery().DelegatorWithdrawAddress(suite.ctx,
		&distrtypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delAddr.String()})
	suite.Require().NoError(err)
	suite.Require().EqualValues(response.WithdrawAddress, withdrawAddr.String())
	return txResponse
}

func (suite *StakingSuite) send(privateKey cryptotypes.PrivKey, value *big.Int, data []byte) *ethtypes.Receipt {
	stakingContract := stakingprecompile.GetAddress()
	if value == nil {
		value = big.NewInt(0)
	}
	transaction, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &stakingContract, value, data)
	suite.Require().NoError(err)
	return suite.SendTransaction(transaction)
}

func (suite *StakingSuite) DelegateV2(privateKey cryptotypes.PrivKey, valAddr string, delAmount *big.Int) *ethtypes.Receipt {
	method := stakingprecompile.NewDelegateV2Method(nil)
	pack, err := method.PackInput(fxstakingtypes.DelegateV2Args{Validator: valAddr, Amount: delAmount})
	suite.Require().NoError(err)
	return suite.send(privateKey, big.NewInt(0), pack)
}

func (suite *StakingSuite) RedelegateV2(privateKey cryptotypes.PrivKey, valSrc, valDst string, amount *big.Int) *ethtypes.Receipt {
	method := stakingprecompile.NewRedelegateV2Method(nil)
	pack, err := method.PackInput(fxstakingtypes.RedelegateV2Args{ValidatorSrc: valSrc, ValidatorDst: valDst, Amount: amount})
	suite.Require().NoError(err)
	return suite.send(privateKey, nil, pack)
}

func (suite *StakingSuite) WithdrawByContract(privateKey cryptotypes.PrivKey, contract common.Address, valAddr string) *ethtypes.Receipt {
	stakingContract, err := testscontract.NewStakingTest(contract, suite.EthClient())
	suite.Require().NoError(err)

	auth := suite.TransactionOpts(privateKey)

	tx, err := stakingContract.Withdraw(auth, valAddr)
	suite.Require().NoError(err)

	ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.EthClient(), tx)
	suite.Require().NoError(err)
	suite.Require().Equal(receipt.Status, ethtypes.ReceiptStatusSuccessful)
	return receipt
}

func (suite *StakingSuite) UnDelegateV2(privateKey cryptotypes.PrivKey, valAddr string, amount *big.Int) *ethtypes.Receipt {
	method := stakingprecompile.NewUndelegateV2Method(nil)
	pack, err := method.PackInput(fxstakingtypes.UndelegateV2Args{Validator: valAddr, Amount: amount})
	suite.Require().NoError(err)
	return suite.send(privateKey, nil, pack)
}

func (suite *StakingSuite) WithdrawReward(privateKey cryptotypes.PrivKey, valAddr string) *ethtypes.Receipt {
	method := stakingprecompile.NewWithdrawMethod(nil)
	pack, err := method.PackInput(fxstakingtypes.WithdrawArgs{Validator: valAddr})
	suite.Require().NoError(err)
	return suite.send(privateKey, nil, pack)
}

func (suite *StakingSuite) Delegation(valAddr string, delAddr common.Address) (*big.Int, *big.Int) {
	stakingContract := stakingprecompile.GetAddress()
	method := stakingprecompile.NewDelegationMethod(nil)
	pack, err := method.PackInput(fxstakingtypes.DelegationArgs{Validator: valAddr, Delegator: delAddr})
	suite.Require().NoError(err)
	output, err := suite.EthClient().CallContract(suite.ctx, ethereum.CallMsg{To: &stakingContract, Data: pack}, nil)
	suite.Require().NoError(err)
	shares, amount, err := method.UnpackOutput(output)
	suite.Require().NoError(err)
	return shares, amount
}

func (suite *StakingSuite) Rewards(valAddr string, delAddr common.Address) *big.Int {
	stakingContract := stakingprecompile.GetAddress()
	method := stakingprecompile.NewDelegationRewardsMethod(nil)
	pack, err := method.PackInput(fxstakingtypes.DelegationRewardsArgs{Validator: valAddr, Delegator: delAddr})
	suite.Require().NoError(err)
	output, err := suite.EthClient().CallContract(suite.ctx, ethereum.CallMsg{To: &stakingContract, Data: pack}, nil)
	suite.Require().NoError(err)
	amount, err := method.UnpackOutput(output)
	suite.Require().NoError(err)
	return amount
}

func (suite *StakingSuite) TransferShares(privateKey cryptotypes.PrivKey, valAddr string, receipt common.Address, shares *big.Int) *ethtypes.Receipt {
	method := stakingprecompile.NewTransferSharesMethod(nil)
	pack, err := method.PackInput(fxstakingtypes.TransferSharesArgs{Validator: valAddr, To: receipt, Shares: shares})
	suite.Require().NoError(err)
	return suite.send(privateKey, nil, pack)
}

func (suite *StakingSuite) TransferSharesByContract(privateKey cryptotypes.PrivKey, valAddr string, contract, to common.Address, shares *big.Int) *ethtypes.Receipt {
	stakingContract, err := testscontract.NewStakingTest(contract, suite.EthClient())
	suite.Require().NoError(err)

	auth := suite.TransactionOpts(privateKey)

	tx, err := stakingContract.TransferShares(auth, valAddr, to, shares)
	suite.Require().NoError(err)

	ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.EthClient(), tx)
	suite.Require().NoError(err)
	suite.Require().Equal(receipt.Status, ethtypes.ReceiptStatusSuccessful)
	return receipt
}

func (suite *StakingSuite) TransferFromShares(privateKey cryptotypes.PrivKey, valAddr string, from, receipt common.Address, shares *big.Int) *ethtypes.Receipt {
	method := stakingprecompile.NewTransferFromSharesMethod(nil)
	pack, err := method.PackInput(fxstakingtypes.TransferFromSharesArgs{Validator: valAddr, From: from, To: receipt, Shares: shares})
	suite.Require().NoError(err)
	return suite.send(privateKey, nil, pack)
}

func (suite *StakingSuite) TransferFromSharesByContract(privateKey cryptotypes.PrivKey, valAddr string, contract, from, to common.Address, shares *big.Int) *ethtypes.Receipt {
	stakingContract, err := testscontract.NewStakingTest(contract, suite.EthClient())
	suite.Require().NoError(err)

	auth := suite.TransactionOpts(privateKey)

	tx, err := stakingContract.TransferFromShares(auth, valAddr, from, to, shares)
	suite.Require().NoError(err)

	ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.EthClient(), tx)
	suite.Require().NoError(err)
	suite.Require().Equal(receipt.Status, ethtypes.ReceiptStatusSuccessful)
	return receipt
}

func (suite *StakingSuite) ApproveShares(privateKey cryptotypes.PrivKey, valAddr string, spender common.Address, shares *big.Int) *ethtypes.Receipt {
	method := stakingprecompile.NewApproveSharesMethod(nil)
	pack, err := method.PackInput(fxstakingtypes.ApproveSharesArgs{Validator: valAddr, Spender: spender, Shares: shares})
	suite.Require().NoError(err)
	return suite.send(privateKey, nil, pack)
}

func (suite *StakingSuite) AllowanceShares(valAddr string, owner, spender common.Address) *big.Int {
	stakingContract := stakingprecompile.GetAddress()
	method := stakingprecompile.NewAllowanceSharesMethod(nil)
	pack, err := method.PackInput(fxstakingtypes.AllowanceSharesArgs{Validator: valAddr, Owner: owner, Spender: spender})
	suite.Require().NoError(err)
	output, err := suite.EthClient().CallContract(suite.ctx, ethereum.CallMsg{To: &stakingContract, Data: pack}, nil)
	suite.Require().NoError(err)
	amount, err := method.UnpackOutput(output)
	suite.Require().NoError(err)
	return amount
}

func (suite *StakingSuite) LogReward(logs []*ethtypes.Log, valAddr string, addr common.Address) *big.Int {
	method := stakingprecompile.NewWithdrawMethod(nil)
	for _, log := range logs {
		if log.Address == stakingprecompile.GetAddress() &&
			log.Topics[0] == method.Event.ID &&
			log.Topics[1] == addr.Hash() {
			unpack, err := method.Event.Inputs.NonIndexed().Unpack(log.Data)
			suite.Require().NoError(err)
			suite.Require().Equal(unpack[0].(string), valAddr)
			return unpack[1].(*big.Int)
		}
	}
	return big.NewInt(0)
}

type AuthzSuite struct {
	*TestSuite
}

func NewAuthzSuite(ts *TestSuite) AuthzSuite {
	return AuthzSuite{TestSuite: ts}
}

func (suite *AuthzSuite) AuthzQuery() authz.QueryClient {
	return suite.GRPCClient().AuthzQuery()
}

type SlashingSuite struct {
	*TestSuite
}

func NewSlashingSuite(ts *TestSuite) SlashingSuite {
	return SlashingSuite{TestSuite: ts}
}

func (suite *SlashingSuite) SlashingQuery() slashingtypes.QueryClient {
	return suite.GRPCClient().SlashingQuery()
}
