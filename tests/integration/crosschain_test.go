package integration

import (
	"github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *IntegrationTest) CrosschainTest() {
	suite.updateParamsTest()
}

func (suite *IntegrationTest) updateParamsTest() {
	chains := types.GetSupportChains()
	for _, chain := range chains {
		crosschainSuite := NewCrosschainSuite(chain, suite.FxCoreSuite)
		crosschainSuite.UpdateParams(func(params *crosschaintypes.Params) {
			params.DelegateMultiple = 100
		})
		params := crosschainSuite.QueryParams()
		suite.Require().Equal(params.DelegateMultiple, int64(100))
	}
}
