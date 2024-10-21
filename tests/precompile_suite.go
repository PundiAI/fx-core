package tests

import (
	"math/big"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v8/client"
	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
)

type PrecompileTestSuite struct {
	EvmTestSuite
	privKey cryptotypes.PrivKey
}

func NewPrecompileTestSuite(ts *TestSuite) PrecompileTestSuite {
	return PrecompileTestSuite{
		EvmTestSuite: NewEvmTestSuite(ts),
		privKey:      helpers.NewEthPrivKey(),
	}
}

func (suite *PrecompileTestSuite) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.privKey.PubKey().Address())
}

func (suite *PrecompileTestSuite) HexAddress() common.Address {
	return common.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *PrecompileTestSuite) TransferCrosschain(token common.Address, recipient string, amount, fee *big.Int, target string) *ethtypes.Transaction {
	privateKey := suite.privKey
	beforeBalanceOf := suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()))
	pack, err := contract.GetFIP20().ABI.Pack("transferCrossChain", recipient, amount, fee, fxtypes.MustStrToByte32(target))
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err, target)
	suite.SendTransaction(ethTx)
	afterBalanceOf := suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()))
	suite.Require().True(new(big.Int).Sub(beforeBalanceOf, afterBalanceOf).Cmp(new(big.Int).Add(amount, fee)) == 0)
	return ethTx
}

func (suite *PrecompileTestSuite) CrosschainAndResponse(token common.Address, recipient string, amount, fee *big.Int, target string) *ethtypes.Transaction {
	privateKey := suite.privKey
	crosschainContract := crosschaintypes.GetAddress()
	suite.ApproveERC20(privateKey, token, crosschainContract, big.NewInt(0).Add(amount, fee))

	beforeBalanceOf := suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()))
	pack, err := crosschaintypes.GetABI().Pack("crossChain", token, recipient, amount, fee, fxtypes.MustStrToByte32(target), "")
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &crosschainContract, nil, pack)
	suite.Require().NoError(err, target)
	suite.SendTransaction(ethTx)
	afterBalanceOf := suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()))
	suite.Require().True(new(big.Int).Sub(beforeBalanceOf, afterBalanceOf).Cmp(new(big.Int).Add(amount, fee)) == 0)
	return ethTx
}
