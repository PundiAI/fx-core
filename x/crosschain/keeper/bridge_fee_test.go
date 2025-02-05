package keeper_test

import (
	"math/big"
	"time"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func (suite *KeeperTestSuite) MockQuote() contract.IBridgeFeeQuoteQuoteInfo {
	oracleKeeper := contract.NewBridgeFeeOracleKeeper(suite.App.EvmKeeper)
	defOracle, err := oracleKeeper.DefaultOracle(suite.Ctx)
	suite.Require().NoError(err)
	suite.True(suite.HasValAddr(defOracle.Bytes()))

	keeper := contract.NewBridgeFeeQuoteKeeper(suite.App.EvmKeeper)
	input := contract.IBridgeFeeQuoteQuoteInput{
		Cap:       0,
		GasLimit:  21000,
		Expiry:    uint64(time.Now().Add(time.Hour).Unix()),
		ChainName: contract.MustStrToByte32(ethtypes.ModuleName),
		TokenName: contract.MustStrToByte32(fxtypes.DefaultDenom),
		Amount:    big.NewInt(100),
	}
	_, err = keeper.Quote(suite.Ctx, defOracle, []contract.IBridgeFeeQuoteQuoteInput{input})
	suite.Require().NoError(err)

	quote, err := keeper.GetQuoteById(suite.Ctx, big.NewInt(1))
	suite.Require().NoError(err)
	suite.Equal(input.ChainName, quote.ChainName)
	suite.Equal(input.TokenName, quote.TokenName)
	suite.Equal(input.Amount, quote.Amount)
	suite.Equal(input.GasLimit, quote.GasLimit)
	suite.Equal(input.Expiry, quote.Expiry)

	_, err = keeper.GetQuoteById(suite.Ctx, big.NewInt(0))
	suite.Require().Error(err)
	_, err = keeper.GetQuoteById(suite.Ctx, big.NewInt(2))
	suite.Require().Error(err)

	oracleList, err := oracleKeeper.GetOracleList(suite.Ctx, contract.MustStrToByte32(ethtypes.ModuleName))
	suite.Require().NoError(err)
	suite.Len(oracleList, 1)
	suite.Equal(defOracle.String(), oracleList[0].String())
	return quote
}

func (suite *KeeperTestSuite) TestValidateQuote() {
	quote := suite.MockQuote()

	quoteInfo, err := suite.Keeper().ValidateQuote(suite.Ctx, suite.App.EvmKeeper, quote.Id, 21000)
	suite.NoError(err)
	suite.Equal(quote, quoteInfo)
}
