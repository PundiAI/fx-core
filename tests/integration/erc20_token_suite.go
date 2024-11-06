package integration

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
)

type ERC20TokenSuite struct {
	*EthSuite

	erc20Token *contract.WFXUpgradable
	signer     *helpers.Signer
}

func NewERC20TokenSuite(suite *EthSuite, token common.Address, signer *helpers.Signer) *ERC20TokenSuite {
	erc20Token, err := contract.NewWFXUpgradable(token, suite.ethCli)
	suite.Require().NoError(err)
	return &ERC20TokenSuite{
		EthSuite:   suite,
		erc20Token: erc20Token,
		signer:     signer,
	}
}

func (suite *ERC20TokenSuite) WithSigner(signer *helpers.Signer) *ERC20TokenSuite {
	return &ERC20TokenSuite{
		EthSuite:   suite.EthSuite,
		erc20Token: suite.erc20Token,
		signer:     signer,
	}
}

func (suite *ERC20TokenSuite) Symbol() string {
	symbol, err := suite.erc20Token.Symbol(nil)
	suite.Require().NoError(err)
	return symbol
}

func (suite *ERC20TokenSuite) Name() string {
	name, err := suite.erc20Token.Name(nil)
	suite.Require().NoError(err)
	return name
}

func (suite *ERC20TokenSuite) Decimals() uint8 {
	decimals, err := suite.erc20Token.Decimals(nil)
	suite.Require().NoError(err)
	return decimals
}

func (suite *ERC20TokenSuite) BalanceOf(account common.Address) *big.Int {
	balance, err := suite.erc20Token.BalanceOf(nil, account)
	suite.Require().NoError(err)
	return balance
}

func (suite *ERC20TokenSuite) EqualBalanceOf(address common.Address, balance *big.Int) {
	balAfter := suite.BalanceOf(address)
	suite.Require().Equal(balAfter.String(), balance.String())
}

func (suite *ERC20TokenSuite) CheckBalance(address common.Address, addValue *big.Int, f func()) {
	value := suite.BalanceOf(address)
	f()
	newValue := suite.BalanceOf(address)
	suite.Require().Equal(big.NewInt(0).Add(value, addValue).String(), newValue.String())
}

func (suite *ERC20TokenSuite) TotalSupply() *big.Int {
	totalSupply, err := suite.erc20Token.TotalSupply(nil)
	suite.Require().NoError(err)
	return totalSupply
}

func (suite *ERC20TokenSuite) Allowance(owner, spender common.Address) *big.Int {
	allowance, err := suite.erc20Token.Allowance(nil, owner, spender)
	suite.Require().NoError(err)
	return allowance
}

func (suite *ERC20TokenSuite) EqualAllowance(owner, spender common.Address, value *big.Int) {
	suite.Require().Equal(suite.Allowance(owner, spender).Cmp(value), 0)
}

func (suite *ERC20TokenSuite) Owner() common.Address {
	owner, err := suite.erc20Token.Owner(nil)
	suite.Require().NoError(err)
	return owner
}

func (suite *ERC20TokenSuite) Deposit(value *big.Int) *ethtypes.Transaction {
	suite.Require().True(suite.Balance(suite.signer.Address()).Cmp(value) >= 0)

	opts := suite.TransactOpts(suite.signer)
	opts.Value = value
	ethTx, err := suite.erc20Token.Deposit(opts)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	suite.Require().True(suite.BalanceOf(suite.signer.Address()).Cmp(value) >= 0)
	suite.Require().True(suite.TotalSupply().Cmp(value) >= 0)
	return ethTx
}

func (suite *ERC20TokenSuite) Withdraw(recipient common.Address, value *big.Int) *ethtypes.Transaction {
	suite.Require().True(suite.TotalSupply().Cmp(value) >= 0)
	suite.Require().True(suite.BalanceOf(suite.signer.Address()).Cmp(value) >= 0)

	ethTx, err := suite.erc20Token.Withdraw0(suite.TransactOpts(suite.signer), recipient, value)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	suite.Require().True(suite.Balance(recipient).Cmp(value) >= 0)
	return ethTx
}

func (suite *ERC20TokenSuite) Transfer(recipient common.Address, value *big.Int) *ethtypes.Transaction {
	suite.Require().True(suite.BalanceOf(suite.signer.Address()).Cmp(value) >= 0)

	ethTx, err := suite.erc20Token.Transfer(suite.TransactOpts(suite.signer), recipient, value)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	suite.Require().True(suite.BalanceOf(recipient).Cmp(value) >= 0)
	return ethTx
}

func (suite *ERC20TokenSuite) Approve(spender common.Address, value *big.Int) *ethtypes.Transaction {
	ethTx, err := suite.erc20Token.Approve(suite.TransactOpts(suite.signer), spender, value)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	suite.Require().True(suite.Allowance(suite.signer.Address(), spender).Cmp(value) >= 0)
	return ethTx
}

func (suite *ERC20TokenSuite) TransferOwnership(newOwner common.Address) *ethtypes.Transaction {
	ethTx, err := suite.erc20Token.TransferOwnership(suite.TransactOpts(suite.signer), newOwner)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	suite.Require().Equal(suite.Owner().String(), newOwner.String())
	return ethTx
}

func (suite *ERC20TokenSuite) TransferFrom(sender, recipient common.Address, value *big.Int) *ethtypes.Transaction {
	ethTx, err := suite.erc20Token.TransferFrom(suite.TransactOpts(suite.signer), sender, recipient, value)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	suite.Require().True(suite.BalanceOf(recipient).Cmp(value) >= 0)
	return ethTx
}

func (suite *ERC20TokenSuite) Mint(account common.Address, value *big.Int) *ethtypes.Transaction {
	ethTx, err := suite.erc20Token.Mint(suite.TransactOpts(suite.signer), account, value)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	suite.Require().True(suite.BalanceOf(account).Cmp(value) >= 0)
	suite.Require().True(suite.TotalSupply().Cmp(value) >= 0)
	return ethTx
}

func (suite *ERC20TokenSuite) Burn(account common.Address, value *big.Int) *ethtypes.Transaction {
	beforeBalance := suite.BalanceOf(account)
	suite.Require().True(beforeBalance.Cmp(value) >= 0)
	beforeTotalSupply := suite.TotalSupply()
	suite.Require().True(beforeTotalSupply.Cmp(value) >= 0)

	ethTx, err := suite.erc20Token.Burn(suite.TransactOpts(suite.signer), account, value)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	suite.Require().True(new(big.Int).Sub(beforeBalance, suite.BalanceOf(account)).Cmp(value) == 0)
	suite.Require().True(new(big.Int).Sub(beforeTotalSupply, suite.TotalSupply()).Cmp(value) == 0)
	return ethTx
}
