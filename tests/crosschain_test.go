package tests

import (
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (suite *IntegrationTest) CrosschainTest() {
}

func (suite *IntegrationTest) UpdateParamsTest() {
	for _, chain := range suite.crosschain {
		chain.UpdateParams(func(params *crosschaintypes.Params) {
			params.DelegateMultiple = 100
		})
		params := chain.QueryParams()
		suite.Require().Equal(params.DelegateMultiple, int64(100))
	}
}
