package tests

import (
	"math/big"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/client"
	precompilesstaking "github.com/functionx/fx-core/v3/x/evm/precompiles/staking"
)

type StakingSuite struct {
	Erc20TestSuite
	abi abi.ABI
}

func NewStakingSuite(ts *TestSuite) StakingSuite {
	return StakingSuite{
		Erc20TestSuite: NewErc20TestSuite(ts),
		abi:            precompilesstaking.GetABI(),
	}
}

func (suite *StakingSuite) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.privKey.PubKey().Address())
}

func (suite *StakingSuite) Address() common.Address {
	return common.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *StakingSuite) StakingQuery() stakingtypes.QueryClient {
	return suite.GRPCClient().StakingQuery()
}

// DelegationRewards Get delegatorAddress rewards
func (suite *StakingSuite) DelegationRewards(delAddr, valAddr string) sdk.DecCoins {
	response, err := suite.GRPCClient().DistrQuery().DelegationRewards(suite.ctx, &distrtypes.QueryDelegationRewardsRequest{DelegatorAddress: delAddr, ValidatorAddress: valAddr})
	suite.Require().NoError(err)
	return response.Rewards
}

func (suite *StakingSuite) SetWithdrawAddress(delAddr, withdrawAddr sdk.AccAddress) {
	setWithdrawAddress := distrtypes.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)
	txResponse := suite.BroadcastTx(suite.privKey, setWithdrawAddress)
	suite.Require().True(txResponse.Code == 0)
	response, err := suite.GRPCClient().DistrQuery().DelegatorWithdrawAddress(suite.ctx, &distrtypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delAddr.String()})
	suite.Require().NoError(err)
	suite.Require().EqualValues(response.WithdrawAddress, withdrawAddr.String())
}

func (suite *StakingSuite) Delegate(privateKey cryptotypes.PrivKey, valAddr string, delAmount *big.Int) {
	stakingContract := precompilesstaking.GetAddress()
	pack, err := precompilesstaking.GetABI().Pack("delegate", valAddr)
	suite.Require().NoError(err)
	transaction, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &stakingContract, delAmount, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(transaction)
}

func (suite *StakingSuite) UnDelegate(privateKey cryptotypes.PrivKey, valAddr string, shares *big.Int) {
	stakingContract := precompilesstaking.GetAddress()
	pack, err := precompilesstaking.GetABI().Pack(precompilesstaking.UndelegateMethodName, valAddr, shares)
	suite.Require().NoError(err)
	transaction, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &stakingContract, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(transaction)
}

func (suite *StakingSuite) WithdrawReward(privateKey cryptotypes.PrivKey, valAddr string) {
	stakingContract := precompilesstaking.GetAddress()
	pack, err := precompilesstaking.GetABI().Pack(precompilesstaking.WithdrawMethodName, valAddr)
	suite.Require().NoError(err)
	transaction, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &stakingContract, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(transaction)
}

func (suite *StakingSuite) Delegation(valAddr string, delAddr common.Address) (*big.Int, *big.Int) {
	stakingContract := precompilesstaking.GetAddress()
	pack, err := precompilesstaking.GetABI().Pack(precompilesstaking.DelegationMethodName, valAddr, delAddr)
	suite.Require().NoError(err)
	output, err := suite.EthClient().CallContract(suite.ctx, ethereum.CallMsg{To: &stakingContract, Data: pack}, nil)
	suite.Require().NoError(err)
	var out []interface{}
	res, err := suite.abi.Unpack(precompilesstaking.DelegationMethodName, output)
	suite.Require().NoError(err)
	out = res
	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return out0, out1
}

func (suite *StakingSuite) Rewards(valAddr string, delAddr common.Address) *big.Int {
	stakingContract := precompilesstaking.GetAddress()
	pack, err := precompilesstaking.GetABI().Pack(precompilesstaking.DelegationRewardsMethodName, valAddr, delAddr)
	suite.Require().NoError(err)
	output, err := suite.EthClient().CallContract(suite.ctx, ethereum.CallMsg{To: &stakingContract, Data: pack}, nil)
	suite.Require().NoError(err)
	var out []interface{}
	res, err := suite.abi.Unpack(precompilesstaking.DelegationRewardsMethodName, output)
	suite.Require().NoError(err)
	out = res
	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0
}
