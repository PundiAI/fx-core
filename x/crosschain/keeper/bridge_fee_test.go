package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func (suite *KeeperTestSuite) TestBridgeFeeOracle() {
	defOracle, err := suite.bridgeFeeSuite.DefaultOracle(suite.Ctx)
	suite.Require().NoError(err)
	suite.True(suite.HasValAddr(defOracle.Bytes()))

	_, err = suite.bridgeFeeSuite.GetQuoteById(suite.Ctx, big.NewInt(0))
	suite.Require().Error(err)
	_, err = suite.bridgeFeeSuite.GetQuoteById(suite.Ctx, big.NewInt(2))
	suite.Require().Error(err)

	// Oracle information will be stored only after the quote
	suite.bridgeFeeSuite.MockQuote(suite.Ctx, suite.chainName, fxtypes.DefaultDenom)
	oracleList, err := suite.bridgeFeeSuite.GetOracleList(suite.Ctx, contract.MustStrToByte32(suite.chainName))
	suite.Require().NoError(err)
	suite.Len(oracleList, 1, suite.chainName)
	suite.Equal(defOracle.String(), oracleList[0].String())
}

func (suite *KeeperTestSuite) TestValidateQuote() {
	quote := suite.bridgeFeeSuite.MockQuote(suite.Ctx, suite.chainName, fxtypes.DefaultDenom)

	quoteInfo, err := suite.Keeper().ValidateQuote(suite.Ctx, suite.App.EvmKeeper, quote.Id, 21000)
	suite.NoError(err)
	suite.Equal(quote, quoteInfo)
}

func (suite *KeeperTestSuite) TestHandlerBridgeCallInFee() {
	type args struct {
		caller   contract.Caller
		from     common.Address
		quoteId  *big.Int
		gasLimit uint64
	}
	tests := []struct {
		name     string
		malleate func() args
		wantErr  bool
	}{
		{
			name: "skip bridge fee",
			malleate: func() args {
				return args{
					caller:   suite.App.EvmKeeper,
					from:     helpers.GenHexAddress(),
					quoteId:  big.NewInt(0),
					gasLimit: 1e6,
				}
			},
		},
		{
			name: "bridge fee with origin token",
			malleate: func() args {
				from := helpers.GenHexAddress()
				suite.MintToken(from.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1)))

				quoteInfo := suite.bridgeFeeSuite.MockQuote(suite.Ctx, suite.chainName, fxtypes.DefaultDenom)

				return args{
					caller:   suite.App.EvmKeeper,
					from:     from,
					quoteId:  quoteInfo.Id,
					gasLimit: quoteInfo.GasLimit,
				}
			},
		},
		{
			name: "bridge fee with erc20 token",
			malleate: func() args {
				symbol := helpers.NewRandSymbol()
				erc20Token, err := suite.App.Erc20Keeper.RegisterNativeCoin(suite.Ctx, symbol, symbol, 18)
				suite.Require().NoError(err)

				signer := suite.NewSigner()
				suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract())
				suite.erc20TokenSuite.MintFromERC20Module(suite.Ctx, signer.Address(), big.NewInt(1))

				quoteInfo := suite.bridgeFeeSuite.MockQuote(suite.Ctx, suite.chainName, erc20Token.Denom)

				return args{
					caller:   suite.App.EvmKeeper,
					from:     signer.Address(),
					quoteId:  quoteInfo.Id,
					gasLimit: quoteInfo.GasLimit,
				}
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testArgs := tt.malleate()

			err := suite.Keeper().HandlerBridgeCallInFee(suite.Ctx, testArgs.caller, testArgs.from, testArgs.quoteId, testArgs.gasLimit)
			if tt.wantErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
