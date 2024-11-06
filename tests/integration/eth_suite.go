package integration

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/client"
	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
)

type EthSuite struct {
	suite.Suite
	ctx    context.Context
	ethCli *ethclient.Client
}

func (suite *EthSuite) TransactOpts(signer *helpers.Signer) *bind.TransactOpts {
	ecdsa, err := crypto.ToECDSA(signer.PrivKey().Bytes())
	suite.Require().NoError(err)

	chainId, err := suite.ethCli.ChainID(suite.ctx)
	suite.Require().NoError(err)

	transactOpts, err := bind.NewKeyedTransactorWithChainID(ecdsa, chainId)
	suite.Require().NoError(err)

	transactOpts.GasTipCap = big.NewInt(1e9)
	transactOpts.GasFeeCap = big.NewInt(6e12)
	transactOpts.GasLimit = 200_000
	return transactOpts
}

func (suite *EthSuite) Balance(addr common.Address) *big.Int {
	at, err := suite.ethCli.BalanceAt(suite.ctx, addr, nil)
	suite.Require().NoError(err)
	return at
}

func (suite *EthSuite) BlockHeight() uint64 {
	number, err := suite.ethCli.BlockNumber(suite.ctx)
	suite.Require().NoError(err)
	return number
}

func (suite *EthSuite) SendTransaction(tx *ethtypes.Transaction) *ethtypes.Receipt {
	err := suite.ethCli.SendTransaction(suite.ctx, tx)
	suite.Require().NoError(err)

	return suite.WaitMined(tx)
}

func (suite *EthSuite) WaitMined(tx *ethtypes.Transaction) *ethtypes.Receipt {
	ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.ethCli, tx)
	suite.Require().NoError(err)
	suite.T().Log("broadcast tx", "msg:", "ethermint.evm.v1.MsgEthereumTx", "height:", receipt.BlockNumber, "txHash:", receipt.TxHash)
	suite.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)
	return receipt
}

func (suite *EthSuite) DeployERC20(signer *helpers.Signer, symbol string) common.Address {
	erc20 := contract.GetFIP20()
	tx, err := client.BuildEthTransaction(suite.ctx, suite.ethCli, signer.PrivKey(), nil, nil, erc20.Bin)
	suite.Require().NoError(err)
	suite.SendTransaction(tx)

	logicAddr := crypto.CreateAddress(signer.Address(), tx.Nonce())
	proxyAddr := suite.DeployProxy(signer, logicAddr, []byte{})

	pack, err := erc20.ABI.Pack("initialize", "Test ERC20", symbol, uint8(18), signer.Address())
	suite.Require().NoError(err)
	tx, err = client.BuildEthTransaction(suite.ctx, suite.ethCli, signer.PrivKey(), &proxyAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(tx)
	return proxyAddr
}

func (suite *EthSuite) DeployERC721(signer *helpers.Signer) common.Address {
	erc721ABI := contract.MustABIJson(testscontract.ERC721TokenTestMetaData.ABI)
	erc721Bin := contract.MustDecodeHex(testscontract.ERC721TokenTestMetaData.Bin)

	tx, err := client.BuildEthTransaction(suite.ctx, suite.ethCli, signer.PrivKey(), nil, nil, erc721Bin)
	suite.Require().NoError(err)
	suite.SendTransaction(tx)

	logicAddr := crypto.CreateAddress(signer.Address(), tx.Nonce())
	proxyAddr := suite.DeployProxy(signer, logicAddr, []byte{})
	pack, err := erc721ABI.Pack("initialize")
	suite.Require().NoError(err)
	tx, err = client.BuildEthTransaction(suite.ctx, suite.ethCli, signer.PrivKey(), &proxyAddr, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(tx)
	return proxyAddr
}

func (suite *EthSuite) DeployStaking(signer *helpers.Signer) (common.Address, common.Hash) {
	stakingBin := contract.MustDecodeHex(testscontract.StakingTestMetaData.Bin)
	return suite.DeployContract(signer, stakingBin)
}

func (suite *EthSuite) DeployCrosschain(signer *helpers.Signer) (common.Address, common.Hash) {
	crosschainBin := contract.MustDecodeHex(testscontract.CrosschainTestMetaData.Bin)
	return suite.DeployContract(signer, crosschainBin)
}

func (suite *EthSuite) DeployContract(signer *helpers.Signer, contractBin []byte) (common.Address, common.Hash) {
	tx, err := client.BuildEthTransaction(suite.ctx, suite.ethCli, signer.PrivKey(), nil, nil, contractBin)
	suite.Require().NoError(err)
	receipt := suite.SendTransaction(tx)

	suite.Require().False(contract.IsZeroEthAddress(receipt.ContractAddress))
	return receipt.ContractAddress, receipt.TxHash
}

func (suite *EthSuite) DeployProxy(signer *helpers.Signer, logic common.Address, initData []byte) common.Address {
	erc1967Proxy := contract.GetERC1967Proxy()
	input, err := erc1967Proxy.ABI.Pack("", logic, initData)
	suite.Require().NoError(err)
	tx, err := client.BuildEthTransaction(suite.ctx, suite.ethCli, signer.PrivKey(), nil, nil, append(erc1967Proxy.Bin, input...))
	suite.Require().NoError(err)
	suite.SendTransaction(tx)
	return crypto.CreateAddress(signer.Address(), tx.Nonce())
}

func (suite *EthSuite) TxFee(hash common.Hash) *big.Int {
	receipt, err := suite.ethCli.TransactionReceipt(suite.ctx, hash)
	suite.Require().NoError(err)

	tx, pending, err := suite.ethCli.TransactionByHash(suite.ctx, hash)
	suite.Require().NoError(err)
	suite.Require().False(pending)

	block, err := suite.ethCli.BlockByNumber(suite.ctx, receipt.BlockNumber)
	suite.Require().NoError(err)
	baseFee := block.BaseFee()

	txData, err := evmtypes.NewTxDataFromTx(tx)
	suite.Require().NoError(err)
	effectiveGasPrice := txData.EffectiveGasPrice(baseFee)
	return big.NewInt(0).Mul(effectiveGasPrice, big.NewInt(0).SetUint64(receipt.GasUsed))
}
