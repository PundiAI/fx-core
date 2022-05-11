package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	fxtypes "github.com/functionx/fx-core/types"
)

func (suite *KeeperTestSuite) TestEndBlock() {
	testCases := []struct {
		name       string
		NoBaseFee  bool
		malleate   func()
		expGasUsed uint64
	}{
		{
			"basFee nil",
			true,
			func() {},
			uint64(0),
		},
		{
			"Block gas meter is nil",
			false,
			func() {},
			uint64(0),
		},
		{
			"pass",
			false,
			func() {
				meter := sdk.NewGasMeter(uint64(1000000000))
				suite.ctx = suite.ctx.WithBlockGasMeter(meter)
				suite.ctx.BlockGasMeter().ConsumeGas(uint64(5000000), "consume gas")
			},
			uint64(5000000),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			params := suite.app.FeeMarketKeeperV0.GetParams(suite.ctx)
			params.NoBaseFee = tc.NoBaseFee
			suite.app.FeeMarketKeeperV0.SetParams(suite.ctx, params)

			tc.malleate()

			req := abci.RequestEndBlock{Height: 1}
			suite.ctx = suite.ctx.WithBlockHeight(fxtypes.EvmV0SupportBlock() + 1)
			suite.app.FeeMarketKeeperV0.EndBlock(suite.ctx, req)
			gasUsed := suite.app.FeeMarketKeeperV0.GetBlockGasUsed(suite.ctx)
			suite.Require().Equal(tc.expGasUsed, gasUsed, tc.name)
		})
	}
}
