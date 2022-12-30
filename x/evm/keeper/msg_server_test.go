package keeper_test

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func (suite *KeeperTestSuite) TestKeeper_EthereumTx() {
	contract := suite.DeployERC20Contract()

	recipient := helpers.GenerateAddress()
	amount := big.NewInt(10)
	data, err := fxtypes.GetERC20().ABI.Pack("transfer", recipient, amount)
	suite.Require().NoError(err)

	chanId := suite.app.EvmKeeper.ChainID()
	suite.Equal(fxtypes.EIP155ChainID(), chanId)

	tx := types.NewTx(
		chanId,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		&contract,
		big.NewInt(0),
		71000,
		big.NewInt(500*1e9),
		nil,
		nil,
		data,
		nil,
	)
	tx.From = suite.signer.Address().String()
	suite.NoError(tx.Sign(ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID()), suite.signer))

	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res)

	suite.Equal(amount, suite.BalanceOf(contract, recipient))
}

func (suite *KeeperTestSuite) TestKeeper_EthereumTx2() {
	recipient := helpers.GenerateAddress()
	amount := big.NewInt(10)

	chanId := suite.app.EvmKeeper.ChainID()
	suite.Equal(fxtypes.EIP155ChainID(), chanId)

	tx := types.NewTx(
		chanId,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		&recipient,
		amount,
		71000,
		big.NewInt(500*1e9),
		nil,
		nil,
		nil,
		nil,
	)
	tx.From = suite.signer.Address().String()
	suite.NoError(tx.Sign(ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID()), suite.signer))

	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res)

	balance := suite.app.EvmKeeper.GetBalance(suite.ctx, recipient)
	suite.Equal(balance, amount)

	balance = suite.app.BankKeeper.GetBalance(suite.ctx, recipient.Bytes(), fxtypes.DefaultDenom).Amount.BigInt()
	suite.Equal(balance, amount)
}
