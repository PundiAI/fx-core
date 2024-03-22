package tests

import (
	"context"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/client"
	"github.com/functionx/fx-core/v7/contract"
	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
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
	suite.Send(suite.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
}

func (suite *EvmTestSuite) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.privKey.PubKey().Address())
}

func (suite *EvmTestSuite) HexAddress() common.Address {
	return common.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *EvmTestSuite) EthClient() *ethclient.Client {
	if suite.IsUseLocalNetwork() {
		cli, err := ethclient.Dial("http://localhost:8545")
		suite.Require().NoError(err)
		return cli
	}
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

func (suite *EvmTestSuite) TotalSupply(contractAddr common.Address) *big.Int {
	caller, err := contract.NewFIP20Upgradable(contractAddr, suite.EthClient())
	suite.NoError(err)
	totalSupply, err := caller.TotalSupply(nil)
	suite.NoError(err)
	return totalSupply
}

func (suite *EvmTestSuite) BalanceOf(contractAddr, address common.Address) *big.Int {
	caller, err := contract.NewFIP20Upgradable(contractAddr, suite.EthClient())
	suite.NoError(err)
	balance, err := caller.BalanceOf(nil, address)
	suite.NoError(err)
	return balance
}

func (suite *EvmTestSuite) Symbol(contractAddr common.Address) string {
	caller, err := contract.NewFIP20Upgradable(contractAddr, suite.EthClient())
	suite.NoError(err)
	symbol, err := caller.Symbol(nil)
	suite.NoError(err)
	return symbol
}

func (suite *EvmTestSuite) CheckBalanceOf(contractAddr, address common.Address, value *big.Int) bool {
	return suite.BalanceOf(contractAddr, address).Cmp(value) == 0
}

func (suite *EvmTestSuite) HandleWithCheckBalance(contractAddr, address common.Address, addValue *big.Int, h func()) {
	value := suite.BalanceOf(contractAddr, address)
	h()
	newValue := suite.BalanceOf(contractAddr, address)
	suite.Equal(big.NewInt(0).Add(value, addValue).String(), newValue.String())
}

func (suite *EvmTestSuite) Allowance(contractAddr, owner, spender common.Address) *big.Int {
	caller, err := contract.NewFIP20Upgradable(contractAddr, suite.EthClient())
	suite.NoError(err)
	allowance, err := caller.Allowance(nil, owner, spender)
	suite.NoError(err)
	return allowance
}

func (suite *EvmTestSuite) Owner(contractAddr common.Address) common.Address {
	caller, err := contract.NewFIP20Upgradable(contractAddr, suite.EthClient())
	suite.NoError(err)
	owner, err := caller.Owner(nil)
	suite.NoError(err)
	return owner
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

func (suite *EvmTestSuite) WFXDeposit(privateKey cryptotypes.PrivKey, address common.Address, value *big.Int) *ethtypes.Transaction {
	suite.True(suite.Balance(common.BytesToAddress(privateKey.PubKey().Address().Bytes())).Cmp(value) >= 0)
	pack, err := contract.GetWFX().ABI.Pack("deposit")
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &address, value, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	suite.True(suite.BalanceOf(address, common.BytesToAddress(privateKey.PubKey().Address().Bytes())).Cmp(value) >= 0)
	suite.True(suite.TotalSupply(address).Cmp(value) >= 0)
	return ethTx
}

func (suite *EvmTestSuite) WFXWithdraw(privateKey cryptotypes.PrivKey, address, recipient common.Address, value *big.Int) *ethtypes.Transaction {
	suite.True(suite.TotalSupply(address).Cmp(value) >= 0)
	suite.True(suite.BalanceOf(address, common.BytesToAddress(privateKey.PubKey().Address().Bytes())).Cmp(value) >= 0)
	pack, err := contract.GetWFX().ABI.Pack("withdraw", recipient, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &address, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)

	suite.True(suite.Balance(recipient).Cmp(value) >= 0)
	return ethTx
}

func (suite *EvmTestSuite) TransferERC20(privateKey cryptotypes.PrivKey, token, recipient common.Address, value *big.Int) *ethtypes.Transaction {
	suite.True(suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes())).Cmp(value) >= 0)

	pack, err := contract.GetFIP20().ABI.Pack("transfer", recipient, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	suite.True(suite.BalanceOf(token, recipient).Cmp(value) >= 0)
	return ethTx
}

func (suite *EvmTestSuite) ApproveERC20(privateKey cryptotypes.PrivKey, token, spender common.Address, value *big.Int) *ethtypes.Transaction {
	pack, err := contract.GetFIP20().ABI.Pack("approve", spender, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	suite.True(suite.Allowance(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()), spender).Cmp(value) >= 0)
	return ethTx
}

func (suite *EvmTestSuite) TransferOwnership(privateKey cryptotypes.PrivKey, token, newOwner common.Address) *ethtypes.Transaction {
	pack, err := contract.GetFIP20().ABI.Pack("transferOwnership", newOwner)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	suite.Require().Equal(suite.Owner(token).String(), newOwner.String())
	return ethTx
}

func (suite *EvmTestSuite) TransferFromERC20(privateKey cryptotypes.PrivKey, token, sender, recipient common.Address, value *big.Int) *ethtypes.Transaction {
	pack, err := contract.GetFIP20().ABI.Pack("transferFrom", sender, recipient, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	suite.True(suite.BalanceOf(token, recipient).Cmp(value) >= 0)
	return ethTx
}

func (suite *EvmTestSuite) MintERC20(privateKey cryptotypes.PrivKey, token, account common.Address, value *big.Int) *ethtypes.Transaction {
	pack, err := contract.GetFIP20().ABI.Pack("mint", account, value)
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	suite.True(suite.BalanceOf(token, account).Cmp(value) >= 0)
	suite.True(suite.TotalSupply(token).Cmp(value) >= 0)
	return ethTx
}

func (suite *EvmTestSuite) BurnERC20(privateKey cryptotypes.PrivKey, token, account common.Address, value *big.Int) *ethtypes.Transaction {
	beforeBalance := suite.BalanceOf(token, account)
	suite.True(beforeBalance.Cmp(value) >= 0)
	beforeTotalSupply := suite.TotalSupply(token)
	suite.True(beforeTotalSupply.Cmp(value) >= 0)
	pack, err := contract.GetFIP20().ABI.Pack("burn", account, value)
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	suite.True(new(big.Int).Sub(beforeBalance, suite.BalanceOf(token, account)).Cmp(value) == 0)
	suite.True(new(big.Int).Sub(beforeTotalSupply, suite.TotalSupply(token)).Cmp(value) == 0)
	return ethTx
}

func (suite *EvmTestSuite) BalanceOfERC721(contractAddr, account common.Address) *big.Int {
	caller, err := testscontract.NewERC721TokenTest(contractAddr, suite.EthClient())
	suite.NoError(err)
	balanceOf, err := caller.BalanceOf(nil, account)
	suite.NoError(err)
	return balanceOf
}

func (suite *EvmTestSuite) CheckBalanceOfERC721(contractAddr, account common.Address, value *big.Int) bool {
	return suite.BalanceOfERC721(contractAddr, account).Cmp(value) == 0
}

func (suite *EvmTestSuite) TokenURI(contractAddr common.Address, id *big.Int) string {
	caller, err := testscontract.NewERC721TokenTest(contractAddr, suite.EthClient())
	suite.NoError(err)
	uri, err := caller.TokenURI(nil, id)
	suite.NoError(err)
	return uri
}

func (suite *EvmTestSuite) IsApprovedForAll(contractAddr, owner, operator common.Address) bool {
	caller, err := testscontract.NewERC721TokenTest(contractAddr, suite.EthClient())
	suite.NoError(err)
	isApproved, err := caller.IsApprovedForAll(nil, owner, operator)
	suite.NoError(err)
	return isApproved
}

func (suite *EvmTestSuite) SafeMintERC721(privateKey cryptotypes.PrivKey, contractAddr, account common.Address) *ethtypes.Transaction {
	pack, err := GetERC721().ABI.Pack("safeMint", account, "ipfs://test-url")
	suite.NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &contractAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *EvmTestSuite) ApproveERC721(privateKey cryptotypes.PrivKey, contractAddr, operator common.Address, id *big.Int) *ethtypes.Transaction {
	pack, err := GetERC721().ABI.Pack("approve", operator, id)
	suite.NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &contractAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *EvmTestSuite) SetApprovalForAll(privateKey cryptotypes.PrivKey, contractAddr, operator common.Address, approved bool) *ethtypes.Transaction {
	pack, err := GetERC721().ABI.Pack("setApprovalForAll", operator, approved)
	suite.NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &contractAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *EvmTestSuite) SafeTransferFrom(privateKey cryptotypes.PrivKey, contractAddr, from, to common.Address, id *big.Int) *ethtypes.Transaction {
	pack, err := GetERC721().ABI.Pack("safeTransferFrom", from, to, id)
	suite.NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &contractAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *EvmTestSuite) SendTransaction(tx *ethtypes.Transaction) *ethtypes.Receipt {
	err := suite.EthClient().SendTransaction(suite.ctx, tx)
	suite.Require().NoError(err)

	ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.EthClient(), tx)
	suite.Require().NoError(err)
	suite.Require().Equal(receipt.Status, ethtypes.ReceiptStatusSuccessful)
	return receipt
}

func (suite *EvmTestSuite) DeployContract(privKey cryptotypes.PrivKey, contractBin []byte) (common.Address, common.Hash) {
	tx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privKey, nil, nil, contractBin)
	suite.Require().NoError(err)
	receipt := suite.SendTransaction(tx)

	suite.Require().NotEqualf(contract.EmptyEvmAddress, receipt.ContractAddress.String(), "contract address is empty")
	return receipt.ContractAddress, receipt.TxHash
}

func (suite *EvmTestSuite) DeployProxy(privateKey cryptotypes.PrivKey, logic common.Address, initData []byte) common.Address {
	input, err := contract.GetERC1967Proxy().ABI.Pack("", logic, initData)
	suite.NoError(err)
	tx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, nil, nil, append(contract.GetERC1967Proxy().Bin, input...))
	suite.NoError(err)
	suite.SendTransaction(tx)
	return crypto.CreateAddress(common.BytesToAddress(privateKey.PubKey().Address().Bytes()), tx.Nonce())
}

func (suite *EvmTestSuite) TxFee(hash common.Hash) *big.Int {
	receipt, err := suite.EthClient().TransactionReceipt(suite.ctx, hash)
	suite.Require().NoError(err)

	tx, pending, err := suite.EthClient().TransactionByHash(suite.ctx, hash)
	suite.Require().NoError(err)
	suite.Require().False(pending)

	block, err := suite.EthClient().BlockByNumber(suite.ctx, receipt.BlockNumber)
	suite.Require().NoError(err)
	baseFee := block.BaseFee()

	txData, err := evmtypes.NewTxDataFromTx(tx)
	suite.Require().NoError(err)
	effectiveGasPrice := txData.EffectiveGasPrice(baseFee)
	return big.NewInt(0).Mul(effectiveGasPrice, big.NewInt(0).SetUint64(receipt.GasUsed))
}

func (suite *EvmTestSuite) DeployERC20Contract(privKey cryptotypes.PrivKey) common.Address {
	tx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privKey, nil, nil, contract.GetFIP20().Bin)
	suite.NoError(err)
	suite.SendTransaction(tx)
	logic := crypto.CreateAddress(common.BytesToAddress(privKey.PubKey().Address().Bytes()), tx.Nonce())
	proxy := suite.DeployProxy(privKey, logic, []byte{})
	pack, err := contract.GetFIP20().ABI.Pack("initialize", "Test ERC20", helpers.NewRandSymbol(), uint8(18), common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes()))
	suite.NoError(err)
	tx, err = client.BuildEthTransaction(suite.ctx, suite.EthClient(), privKey, &proxy, nil, pack)
	suite.NoError(err)
	suite.SendTransaction(tx)
	return proxy
}

func GetERC721() contract.Contract {
	return contract.Contract{
		ABI: contract.MustABIJson(testscontract.ERC721TokenTestMetaData.ABI),
		Bin: contract.MustDecodeHex(testscontract.ERC721TokenTestMetaData.Bin),
	}
}
