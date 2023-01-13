package tests

import (
	"context"
	"math/big"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/client"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/types/contract"
)

type EvmTestSuite struct {
	*TestSuite
	privKey cryptotypes.PrivKey
}

func NewEvmTestSuite(ts *TestSuite) EvmTestSuite {
	return EvmTestSuite{
		TestSuite: ts,
		privKey:   helpers.NewEthPrivKey(),
	}
}

func (suite *EvmTestSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()

	// transfer to eth private key
	suite.Send(suite.AccAddress(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))
}

func (suite *EvmTestSuite) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.privKey.PubKey().Address())
}

func (suite *EvmTestSuite) HexAddress() common.Address {
	return common.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *EvmTestSuite) EthClient() *ethclient.Client {
	return suite.GetFirstValidator().JSONRPCClient
}

func (suite *EvmTestSuite) TransactOpts() *bind.TransactOpts {
	ecdsa, err := crypto.ToECDSA(suite.privKey.Bytes())
	suite.Require().NoError(err)

	transactOpts, err := bind.NewKeyedTransactorWithChainID(ecdsa, fxtypes.EIP155ChainID())
	suite.Require().NoError(err)

	return transactOpts
}

func (suite *EvmTestSuite) Balance(addr common.Address) *big.Int {
	at, err := suite.EthClient().BalanceAt(suite.ctx, addr, nil)
	suite.Require().NoError(err)
	return at
}

func (suite *EvmTestSuite) CodeAt(addr common.Address) []byte {
	at, err := suite.EthClient().CodeAt(suite.ctx, addr, nil)
	suite.Require().NoError(err)
	return at
}

func (suite *EvmTestSuite) TotalSupply(contractAddr common.Address) *big.Int {
	caller, err := contract.NewFIP20(contractAddr, suite.EthClient())
	suite.NoError(err)
	totalSupply, err := caller.TotalSupply(nil)
	suite.NoError(err)
	return totalSupply
}

func (suite *EvmTestSuite) BalanceOf(contractAddr, address common.Address) *big.Int {
	caller, err := contract.NewFIP20(contractAddr, suite.EthClient())
	suite.NoError(err)
	balance, err := caller.BalanceOf(nil, address)
	suite.NoError(err)
	return balance
}

func (suite *EvmTestSuite) CheckBalanceOf(contractAddr, address common.Address, value *big.Int) bool {
	return suite.BalanceOf(contractAddr, address).Cmp(value) == 0
}

func (suite *EvmTestSuite) Allowance(contractAddr, owner, spender common.Address) *big.Int {
	caller, err := contract.NewFIP20(contractAddr, suite.EthClient())
	suite.NoError(err)
	allowance, err := caller.Allowance(nil, owner, spender)
	suite.NoError(err)
	return allowance
}

func (suite *EvmTestSuite) CheckAllowance(contractAddr, owner, spender common.Address, value *big.Int) bool {
	return suite.Allowance(contractAddr, owner, spender).Cmp(value) == 0
}

func (suite *EvmTestSuite) BlockHeight() uint64 {
	number, err := suite.EthClient().BlockNumber(suite.ctx)
	suite.Require().NoError(err)
	return number
}

func (suite *EvmTestSuite) Transfer(privateKey cryptotypes.PrivKey, recipient common.Address, value *big.Int) *ethtypes.Transaction {
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &recipient, value, nil)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *EvmTestSuite) WFXDeposit(privateKey cryptotypes.PrivKey, wfx common.Address, value *big.Int) *ethtypes.Transaction {
	testAddress := common.BytesToAddress(privateKey.PubKey().Address().Bytes())
	expectBalance := new(big.Int).Add(suite.BalanceOf(wfx, testAddress), value)
	expectTotalSupply := new(big.Int).Add(suite.TotalSupply(wfx), value)

	suite.True(suite.Balance(testAddress).Cmp(value) >= 0)
	pack, err := fxtypes.GetWFX().ABI.Pack("deposit")
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &wfx, value, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)

	suite.CheckBalanceOf(wfx, testAddress, expectBalance)
	suite.True(suite.TotalSupply(wfx).Cmp(expectTotalSupply) == 0)
	return ethTx
}

func (suite *EvmTestSuite) WFXWithdraw(privateKey cryptotypes.PrivKey, wfx, recipient common.Address, value *big.Int) *ethtypes.Transaction {
	testAddress := common.BytesToAddress(privateKey.PubKey().Address().Bytes())
	expectWfxBalance := new(big.Int).Sub(suite.BalanceOf(wfx, testAddress), value)
	expectTotalSupply := new(big.Int).Sub(suite.TotalSupply(wfx), value)
	expectFxbalance := new(big.Int).Add(suite.Balance(recipient), value)

	suite.True(suite.BalanceOf(wfx, testAddress).Cmp(value) >= 0)
	pack, err := fxtypes.GetWFX().ABI.Pack("withdraw", recipient, value)
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &wfx, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)

	suite.CheckBalanceOf(wfx, testAddress, expectWfxBalance)
	suite.CheckBalance(recipient.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(expectFxbalance)))
	suite.True(suite.TotalSupply(wfx).Cmp(expectTotalSupply) == 0)
	return ethTx
}

func (suite *EvmTestSuite) TransferERC20(privateKey cryptotypes.PrivKey, token, recipient common.Address, value *big.Int) *ethtypes.Transaction {
	testAddress := common.BytesToAddress(privateKey.PubKey().Address().Bytes())
	expectFromBalance := new(big.Int).Sub(suite.BalanceOf(token, testAddress), value)
	expectToBalance := new(big.Int).Add(suite.BalanceOf(token, recipient), value)

	pack, err := fxtypes.GetERC20().ABI.Pack("transfer", recipient, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	suite.CheckBalanceOf(token, testAddress, expectFromBalance)
	suite.CheckBalanceOf(token, recipient, expectToBalance)
	return ethTx
}

func (suite *EvmTestSuite) ApproveERC20(privateKey cryptotypes.PrivKey, token, spender common.Address, value *big.Int) *ethtypes.Transaction {
	testAddress := common.BytesToAddress(privateKey.PubKey().Address().Bytes())
	expectAllowance := value

	pack, err := fxtypes.GetERC20().ABI.Pack("approve", spender, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)

	suite.True(suite.Allowance(token, testAddress, spender).Cmp(expectAllowance) == 0)
	return ethTx
}

func (suite *EvmTestSuite) TransferFromERC20(privateKey cryptotypes.PrivKey, token, sender, recipient common.Address, value *big.Int) *ethtypes.Transaction {
	testAddress := common.BytesToAddress(privateKey.PubKey().Address().Bytes())
	expectAllowance := new(big.Int).Sub(suite.Allowance(token, sender, testAddress), value)
	expectSenderBalanceOf := new(big.Int).Sub(suite.BalanceOf(token, sender), value)
	expectRecipientBalanceOf := new(big.Int).Add(suite.BalanceOf(token, recipient), value)
	pack, err := fxtypes.GetERC20().ABI.Pack("transferFrom", sender, recipient, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)

	suite.True(suite.Allowance(token, sender, testAddress).Cmp(expectAllowance) == 0)
	suite.CheckBalanceOf(token, sender, expectSenderBalanceOf)
	suite.CheckBalanceOf(token, recipient, expectRecipientBalanceOf)
	return ethTx
}

func (suite *EvmTestSuite) WFXTransferCrossChain(privateKey cryptotypes.PrivKey, wfx common.Address, recipient string, totalAmount, fxAmount, fee *big.Int, target string) {
	sender := common.BytesToAddress(privateKey.PubKey().Address().Bytes())
	suite.True(suite.Balance(sender).Cmp(fxAmount) > 0)
	suite.True(new(big.Int).Add(suite.Balance(sender), suite.BalanceOf(wfx, sender)).Cmp(new(big.Int).Add(totalAmount, fee)) > 0)
	var expectBalance *big.Int
	if totalAmount.Cmp(fxAmount) > 0 {
		suite.True(suite.BalanceOf(wfx, sender).Cmp(new(big.Int).Sub(new(big.Int).Add(totalAmount, fee), fxAmount)) >= 0)
		expectBalance = new(big.Int).Sub(suite.BalanceOf(wfx, sender), new(big.Int).Sub(new(big.Int).Add(totalAmount, fee), fxAmount))
	} else if totalAmount.Cmp(fxAmount) == 0 {
		expectBalance = new(big.Int).Sub(suite.BalanceOf(wfx, sender), fee)
	} else {
		if new(big.Int).Sub(fxAmount, totalAmount).Cmp(fee) >= 0 {
			expectBalance = new(big.Int).Add(suite.BalanceOf(wfx, sender), new(big.Int).Sub(fxAmount, new(big.Int).Add(totalAmount, fee)))
		} else {
			expectBalance = new(big.Int).Sub(suite.BalanceOf(wfx, sender), new(big.Int).Sub(new(big.Int).Add(totalAmount, fee), fxAmount))
		}
	}
	pack, err := fxtypes.GetWFX().ABI.Pack("transferCrossChain", recipient, totalAmount, fee, fxtypes.MustStrToByte32(target))
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &wfx, fxAmount, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	suite.CheckBalanceOf(wfx, sender, expectBalance)
}

func (suite *EvmTestSuite) MintERC20(privateKey cryptotypes.PrivKey, token, account common.Address, value *big.Int) *ethtypes.Transaction {
	expectBalanceOf := new(big.Int).Add(suite.BalanceOf(token, account), value)
	expectTotalSupply := new(big.Int).Add(suite.TotalSupply(token), value)

	pack, err := fxtypes.GetERC20().ABI.Pack("mint", account, value)
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)

	suite.CheckBalanceOf(token, account, expectBalanceOf)
	suite.True(suite.TotalSupply(token).Cmp(expectTotalSupply) == 0)
	return ethTx
}

func (suite *EvmTestSuite) BurnERC20(privateKey cryptotypes.PrivKey, token, account common.Address, value *big.Int) *ethtypes.Transaction {
	expectBalanceOf := new(big.Int).Sub(suite.BalanceOf(token, account), value)
	expectTotalSupply := new(big.Int).Sub(suite.TotalSupply(token), value)

	pack, err := fxtypes.GetERC20().ABI.Pack("burn", account, value)
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)

	suite.CheckBalanceOf(account, token, expectBalanceOf)
	suite.True(suite.TotalSupply(token).Cmp(expectTotalSupply) == 0)
	return ethTx
}

func (suite *EvmTestSuite) BalanceOfERC721(contractAddr, account common.Address) *big.Int {
	caller, err := contract.NewERC721Token(contractAddr, suite.EthClient())
	suite.NoError(err)
	balanceOf, err := caller.BalanceOf(nil, account)
	suite.NoError(err)
	return balanceOf
}

func (suite *EvmTestSuite) CheckBalanceOfERC721(contractAddr, account common.Address, value *big.Int) bool {
	return suite.BalanceOfERC721(contractAddr, account).Cmp(value) == 0
}

func (suite *EvmTestSuite) TokenURI(contractAddr common.Address, id *big.Int) string {
	caller, err := contract.NewERC721Token(contractAddr, suite.EthClient())
	suite.NoError(err)
	uri, err := caller.TokenURI(nil, id)
	suite.NoError(err)
	return uri
}

func (suite *EvmTestSuite) IsApprovedForAll(contractAddr, owner, operator common.Address) bool {
	caller, err := contract.NewERC721Token(contractAddr, suite.EthClient())
	suite.NoError(err)
	isApproved, err := caller.IsApprovedForAll(nil, owner, operator)
	suite.NoError(err)
	return isApproved
}

func (suite *EvmTestSuite) GetApproved(contractAddr common.Address, id *big.Int) common.Address {
	caller, err := contract.NewERC721Token(contractAddr, suite.EthClient())
	suite.NoError(err)
	approved, err := caller.GetApproved(nil, id)
	suite.NoError(err)
	return approved
}

func (suite *EvmTestSuite) SafeMintERC721(privateKey cryptotypes.PrivKey, contractAddr, account common.Address) *ethtypes.Transaction {
	expectBalanceOf := new(big.Int).Add(suite.BalanceOf(contractAddr, account), big.NewInt(1))
	pack, err := GetERC721().ABI.Pack("safeMint", account)
	suite.NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &contractAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	suite.CheckBalanceOfERC721(contractAddr, account, expectBalanceOf)
	return ethTx
}

func (suite *EvmTestSuite) ApproveERC721(privateKey cryptotypes.PrivKey, contractAddr, operator common.Address, id *big.Int) *ethtypes.Transaction {
	pack, err := GetERC721().ABI.Pack("approve", operator, id)
	suite.NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &contractAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)

	suite.True(suite.GetApproved(contractAddr, id).Hex() == operator.Hex())
	return ethTx
}

func (suite *EvmTestSuite) SetApprovalForAll(privateKey cryptotypes.PrivKey, contractAddr, operator common.Address, approved bool) *ethtypes.Transaction {
	pack, err := GetERC721().ABI.Pack("setApprovalForAll", operator, approved)
	suite.NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &contractAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	suite.True(suite.IsApprovedForAll(contractAddr, common.BytesToAddress(privateKey.PubKey().Address().Bytes()), operator))
	return ethTx
}

func (suite *EvmTestSuite) SafeTransferFrom(privateKey cryptotypes.PrivKey, contractAddr, from, to common.Address, id *big.Int) *ethtypes.Transaction {
	expectFromBalanceOf := new(big.Int).Sub(suite.BalanceOf(contractAddr, from), big.NewInt(1))
	expectToBalanceOf := new(big.Int).Add(suite.BalanceOf(contractAddr, to), big.NewInt(1))

	pack, err := GetERC721().ABI.Pack("safeTransferFrom", from, to, id)
	suite.NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &contractAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)

	suite.CheckBalanceOfERC721(contractAddr, from, expectFromBalanceOf)
	suite.CheckBalanceOfERC721(contractAddr, to, expectToBalanceOf)

	return ethTx
}

func (suite *EvmTestSuite) SendTransaction(tx *ethtypes.Transaction) {
	err := suite.EthClient().SendTransaction(suite.ctx, tx)
	suite.Require().NoError(err)

	ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.EthClient(), tx)
	suite.Require().NoError(err)
	suite.Require().Equal(receipt.Status, ethtypes.ReceiptStatusSuccessful)
}

func GetERC721() fxtypes.Contract {
	return fxtypes.Contract{
		ABI: fxtypes.MustABIJson(contract.ERC721TokenMetaData.ABI),
		Bin: fxtypes.MustDecodeHex(contract.ERC721TokenMetaData.Bin),
	}
}
