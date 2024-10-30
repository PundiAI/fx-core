package precompile_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
)

func TestCrosschainIsOracleOnlineABI(t *testing.T) {
	method := precompile.NewIsOracleOnlineMethod(nil)
	require.Len(t, method.Method.Inputs, 2)
	require.Len(t, method.Method.Outputs, 1)
}

func (suite *CrosschainPrecompileTestSuite) TestIsOracleOnline() {
	testCases := []struct {
		name     string
		malleate func() (contract.IsOracleOnlineArgs, error)
		result   bool
	}{
		{
			name: "oracle online",
			malleate: func() (contract.IsOracleOnlineArgs, error) {
				moduleName := suite.GenerateModuleName()
				oracle := suite.GenerateRandOracle(moduleName, true)
				return contract.IsOracleOnlineArgs{
					Chain:           moduleName,
					ExternalAddress: oracle.GetExternalHexAddr(),
				}, nil
			},
			result: true,
		},
		{
			name: "oracle offline",
			malleate: func() (contract.IsOracleOnlineArgs, error) {
				moduleName := suite.GenerateModuleName()
				oracle := suite.GenerateRandOracle(moduleName, false)
				return contract.IsOracleOnlineArgs{
					Chain:           moduleName,
					ExternalAddress: oracle.GetExternalHexAddr(),
				}, nil
			},
			result: false,
		},
		{
			name: "oracle not found",
			malleate: func() (contract.IsOracleOnlineArgs, error) {
				return contract.IsOracleOnlineArgs{
					Chain:           ethtypes.ModuleName,
					ExternalAddress: helpers.GenHexAddress(),
				}, nil
			},
			result: false,
		},
		{
			name: "invalid chain",
			malleate: func() (contract.IsOracleOnlineArgs, error) {
				return contract.IsOracleOnlineArgs{
					Chain:           "invalid_chain",
					ExternalAddress: helpers.GenHexAddress(),
				}, fmt.Errorf("invalid module name: %s: evm transaction execution failed", "invalid_chain")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			args, expectErr := tc.malleate()

			result := suite.WithError(expectErr).IsOracleOnline(suite.Ctx, args)
			suite.Require().Equal(tc.result, result)
		})
	}
}
