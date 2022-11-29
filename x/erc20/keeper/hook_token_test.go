package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestHookToken() {
	suite.mintToken(bsctypes.ModuleName, sdk.NewCoin(purseDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))

	signer1 := helpers.NewSigner(helpers.NewEthPrivKey())
	signer2 := helpers.NewSigner(helpers.NewEthPrivKey())

	purseId := suite.app.Erc20Keeper.GetDenomMap(suite.ctx, purseDenom)
	suite.Require().NotEmpty(purseId)

	purseTokenPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, purseId)
	suite.Require().True(found)
	suite.Require().NotNil(purseTokenPair)
	suite.Require().NotEmpty(purseTokenPair.GetERC20Contract())

	fip20, err := suite.app.Erc20Keeper.QueryERC20(suite.ctx, purseTokenPair.GetERC20Contract())
	suite.Require().NoError(err)
	suite.Require().Equal("PURSE", fip20.Symbol)

	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), signer1.Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt), sdk.NewCoin(purseDenom, amt)))
	suite.Require().NoError(err)

	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), signer2.Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt)))
	suite.Require().NoError(err)

	//convert coin purse
	transferAmount := sdk.NewInt(10).Mul(sdk.NewInt(1e18))
	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err = suite.app.Erc20Keeper.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
		Coin:     sdk.NewCoin(purseDenom, transferAmount),
		Receiver: signer1.Address().Hex(),
		Sender:   sdk.AccAddress(signer1.Address().Bytes()).String(),
	})
	suite.Require().NoError(err)

	//check contract addr1 balance
	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, purseTokenPair.GetERC20Contract(), signer1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(balanceOf, transferAmount.BigInt())

	//check addr2 balance
	bankBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer2.Address().Bytes())
	suite.Require().True(bankBalance.AmountOf(purseDenom).IsZero())

	balanceOf, err = suite.app.Erc20Keeper.BalanceOf(suite.ctx, purseTokenPair.GetERC20Contract(), signer2.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(balanceOf.Int64(), int64(0))

	data := packTransferData(suite.T(), signer2.Address(), big.NewInt(100))
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, signer1.Address(), purseTokenPair.GetERC20Contract(), data)

	//check addr2 contract balance
	balanceOf, err = suite.app.Erc20Keeper.BalanceOf(suite.ctx, purseTokenPair.GetERC20Contract(), signer2.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(balanceOf.Int64(), int64(100))

	//transfer hook
	data = packTransferData(suite.T(), erc20types.ModuleAddress, big.NewInt(100))
	sendEthTx(suite.T(), suite.ctx, suite.app, signer2, signer2.Address(), purseTokenPair.GetERC20Contract(), data)

	//check addr2 balance
	bankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, signer2.Address().Bytes())
	suite.Require().Equal(bankBalance.AmountOf(purseDenom).Int64(), int64(100))

	balanceOf, err = suite.app.Erc20Keeper.BalanceOf(suite.ctx, purseTokenPair.GetERC20Contract(), signer2.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(balanceOf.Int64(), int64(0))

}

func packTransferData(t *testing.T, to common.Address, amount *big.Int) []byte {
	fip20 := fxtypes.GetERC20()
	pack, err := fip20.ABI.Pack("transfer", to, amount)
	require.NoError(t, err)
	return pack
}
