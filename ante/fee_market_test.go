package ante_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/ante"
	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func (suite *AnteTestSuite) TestGasWantedDecorator() {
	testCases := []struct {
		name              string
		expectedGasWanted uint64
		malleate          func() sdk.Tx
	}{
		{
			"Cosmos Tx",
			TestGasLimit,
			func() sdk.Tx {
				denom := evmtypes.DefaultEVMDenom
				testMsg := banktypes.MsgSend{
					FromAddress: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
					ToAddress:   sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: denom}},
				}
				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), "stake", &testMsg)
				return txBuilder.GetTx()
			},
		},
		{
			"Ethereum Legacy Tx",
			TestGasLimit,
			func() sdk.Tx {
				signer := helpers.NewSigner(helpers.NewEthPrivKey())
				to := helpers.GenerateAddress()
				msg := suite.BuildTestEthTx(signer.Address(), to, nil, make([]byte, 0), big.NewInt(0), nil, nil, nil)
				return suite.CreateTestTx(msg, signer.PrivKey(), 1, false)
			},
		},
		{
			"Ethereum Access List Tx",
			TestGasLimit,
			func() sdk.Tx {
				signer := helpers.NewSigner(helpers.NewEthPrivKey())
				to := helpers.GenerateAddress()
				emptyAccessList := ethtypes.AccessList{}
				msg := suite.BuildTestEthTx(signer.Address(), to, nil, make([]byte, 0), big.NewInt(0), nil, nil, &emptyAccessList)
				return suite.CreateTestTx(msg, signer.PrivKey(), 1, false)
			},
		},
		{
			"Ethereum Dynamic Fee Tx (EIP1559)",
			TestGasLimit,
			func() sdk.Tx {
				signer := helpers.NewSigner(helpers.NewEthPrivKey())
				to := helpers.GenerateAddress()
				emptyAccessList := ethtypes.AccessList{}
				msg := suite.BuildTestEthTx(signer.Address(), to, nil, make([]byte, 0), big.NewInt(0), big.NewInt(100), big.NewInt(50), &emptyAccessList)
				return suite.CreateTestTx(msg, signer.PrivKey(), 1, false)
			},
		},
	}
	suite.SetupTest()
	suite.ctx = suite.ctx.WithBlockHeight(1)

	params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
	params.EnableHeight = 1
	params.NoBaseFee = false
	err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
	suite.NoError(err)
	dec := ante.NewGasWantedDecorator(suite.app.EvmKeeper, suite.app.FeeMarketKeeper)

	// cumulative gas wanted from all test transactions in the same block
	var expectedGasWanted uint64

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := dec.AnteHandle(suite.ctx, tc.malleate(), false, NextFn)
			suite.Require().NoError(err)

			gasWanted := suite.app.FeeMarketKeeper.GetTransientGasWanted(suite.ctx)
			expectedGasWanted += tc.expectedGasWanted
			suite.Require().Equal(expectedGasWanted, gasWanted)
		})
	}
}
