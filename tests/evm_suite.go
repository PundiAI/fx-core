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
	return suite.GetFirstValidtor().JSONRPCClient
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

func (suite *EvmTestSuite) BalanceOf(contractAddr, address common.Address) *big.Int {
	caller, err := contract.NewFIP20(contractAddr, suite.EthClient())
	suite.NoError(err)
	balance, err := caller.BalanceOf(nil, address)
	suite.NoError(err)
	return balance
}

func (suite *EvmTestSuite) BlockHeight() uint64 {
	number, err := suite.EthClient().BlockNumber(suite.ctx)
	suite.Require().NoError(err)
	return number
}

func (suite *EvmTestSuite) Transfer(privateKey cryptotypes.PrivKey, recipient common.Address, value *big.Int) common.Hash {
	suite.T().Logf("transfer to %s value %s\n", recipient.String(), value.String())
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &recipient, value, nil)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	return ethTx.Hash()
}

func (suite *EvmTestSuite) WFXDeposit(address common.Address, amount *big.Int) common.Hash {
	pack, err := fxtypes.GetWFX().ABI.Pack("deposit")
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), suite.privKey, &address, amount, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)

	return ethTx.Hash()
}

func (suite *EvmTestSuite) WFXWithdraw(address, recipient common.Address, value *big.Int) common.Hash {
	pack, err := fxtypes.GetWFX().ABI.Pack("withdraw", recipient, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), suite.privKey, &address, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)

	return ethTx.Hash()
}

func (suite *EvmTestSuite) SendTransaction(tx *ethtypes.Transaction) {
	err := suite.EthClient().SendTransaction(suite.ctx, tx)
	suite.Require().NoError(err)

	suite.T().Log("pending tx hash", tx.Hash())
	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.EthClient(), tx)
	suite.Require().NoError(err)
	suite.Require().Equal(receipt.Status, ethtypes.ReceiptStatusSuccessful)
}
