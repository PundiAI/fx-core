package tests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ConvertDenomTestSuite struct {
	CrosschainERC20TestSuite
}

func TestConvertDenomTestSuite(t *testing.T) {
	testSuite := NewTestSuite()
	convertDenomTestSuite := &ConvertDenomTestSuite{
		CrosschainERC20TestSuite: NewCrosschainERC20TestSuite(testSuite),
	}
	suite.Run(t, convertDenomTestSuite)
}

func (suite *ConvertDenomTestSuite) TestERC20ConvertDenom() {
	suite.InitCrossChain()
	suite.InitRegisterCoinUSDT()

	bscUSDTDenom := fmt.Sprintf("%s%s", suite.BSCCrossChain.chainName, bscUSDToken)

	usdtTokenPair := suite.ERC20.TokenPair("usdt")
	suite.T().Log("token pair", usdtTokenPair.String())

	// bsc -> fx -> evm
	beforeSendToFx := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.BSCCrossChain.HexAddr())
	suite.BSCCrossChain.SendToFxClaim(bscUSDToken, sdk.NewInt(100).MulRaw(1e18), "module/evm")
	afterSendToFx := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.BSCCrossChain.HexAddr())
	suite.Equal(big.NewInt(0).Sub(afterSendToFx, beforeSendToFx), sdk.NewInt(100).MulRaw(1e18).BigInt())

	beforeBalances := suite.QueryBalances(suite.BSCCrossChain.AccAddr())
	suite.BSCCrossChain.SendToFxClaim(bscUSDToken, sdk.NewInt(100).MulRaw(1e18), "")
	afterSendToFxBalances := suite.QueryBalances(suite.BSCCrossChain.AccAddr())
	suite.Equal(afterSendToFxBalances.AmountOf("usdt").Sub(beforeBalances.AmountOf("usdt")), sdk.NewInt(100).MulRaw(1e18))

	suite.ERC20.ConvertDenom(suite.BSCCrossChain.privKey, suite.BSCCrossChain.AccAddr(), sdk.NewCoin("usdt", sdk.NewInt(100).MulRaw(1e18)), "bsc")
	afterConvertDenomUSDT := suite.QueryBalances(suite.BSCCrossChain.AccAddr())
	suite.Equal(afterConvertDenomUSDT.AmountOf(bscUSDTDenom).Sub(afterSendToFxBalances.AmountOf(bscUSDTDenom)), sdk.NewInt(100).MulRaw(1e18))

	suite.ERC20.ConvertDenom(suite.BSCCrossChain.privKey, suite.BSCCrossChain.AccAddr(), sdk.NewCoin(bscUSDTDenom, sdk.NewInt(100).MulRaw(1e18)), "")
	afterConvertDenomBscUSDT := suite.QueryBalances(suite.BSCCrossChain.AccAddr())
	suite.Equal(afterConvertDenomBscUSDT.AmountOf("usdt").Sub(afterConvertDenomUSDT.AmountOf("usdt")), sdk.NewInt(100).MulRaw(1e18))
}
