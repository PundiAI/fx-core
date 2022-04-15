package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func (suite *KeeperTestSuite) TestCalculateBaseFee() {
	fxtypes.ChangeNetworkForTest(fxtypes.NetworkDevnet())
	testCases := []struct {
		name      string
		NoBaseFee bool
		malleate  func()
		expFee    *big.Int
	}{
		{
			"without BaseFee",
			true,
			func() {},
			nil,
		},
		{
			"with BaseFee - initial EIP-1559 block",
			false,
			func() {},
			suite.app.FeeMarketKeeper.GetParams(suite.ctx).BaseFee.BigInt(),
		},
		{
			"with BaseFee - parent block used the same gas as its target",
			false,
			func() {
				// non initial block

				// Set gas used
				suite.app.FeeMarketKeeper.SetBlockGasUsed(suite.ctx, 100)

				// Set target/gasLimit through Consensus Param MaxGas
				blockParams := abci.BlockParams{
					MaxGas:   100,
					MaxBytes: 10,
				}
				consParams := abci.ConsensusParams{Block: &blockParams}
				suite.ctx = suite.ctx.WithConsensusParams(&consParams)

				// set ElasticityMultiplier
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.ElasticityMultiplier = 1
				params.MaxGas = sdk.NewInt(100)
				suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
			},
			suite.app.FeeMarketKeeper.GetParams(suite.ctx).BaseFee.BigInt(),
		},
		{
			"with BaseFee - parent block used more gas than its target",
			false,
			func() {

				suite.app.FeeMarketKeeper.SetBlockGasUsed(suite.ctx, 200)

				blockParams := abci.BlockParams{
					MaxGas:   100,
					MaxBytes: 10,
				}
				consParams := abci.ConsensusParams{Block: &blockParams}
				suite.ctx = suite.ctx.WithConsensusParams(&consParams)

				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.ElasticityMultiplier = 1
				params.MaxGas = sdk.NewInt(100)
				suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)

				suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
			},
			big.NewInt(1125000000),
		},
		{
			"with BaseFee - Parent gas used smaller than parent gas target",
			false,
			func() {

				suite.app.FeeMarketKeeper.SetBlockGasUsed(suite.ctx, 50)

				blockParams := abci.BlockParams{
					MaxGas:   100,
					MaxBytes: 10,
				}
				consParams := abci.ConsensusParams{Block: &blockParams}
				suite.ctx = suite.ctx.WithConsensusParams(&consParams)

				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.ElasticityMultiplier = 1
				params.MaxGas = sdk.NewInt(100)
				suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)

				suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
			},
			big.NewInt(937500000),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
			params.NoBaseFee = tc.NoBaseFee
			suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)

			tc.malleate()

			fee := suite.app.FeeMarketKeeper.CalculateBaseFee(suite.ctx)
			if tc.NoBaseFee {
				suite.Require().Nil(fee, tc.name)
			} else {
				suite.Require().Equal(tc.expFee, fee, tc.name)
			}
		})
	}
}
