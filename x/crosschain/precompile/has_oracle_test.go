package precompile_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
)

func TestCrosschainHasOracleABI(t *testing.T) {
	hasOracleABI := precompile.NewHasOracleABI()
	require.Len(t, hasOracleABI.Method.Inputs, 2)
	require.Len(t, hasOracleABI.Method.Outputs, 1)
}

func (suite *CrosschainPrecompileTestSuite) TestHasOracle() {
	testCases := []struct {
		name     string
		malleate func() (contract.HasOracleArgs, error)
		result   bool
	}{
		{
			name: "has oracle",
			malleate: func() (contract.HasOracleArgs, error) {
				oracle := suite.SetOracle(true)
				return contract.HasOracleArgs{
					Chain:           suite.chainName,
					ExternalAddress: types.ExternalAddrToHexAddr(suite.chainName, oracle.ExternalAddress),
				}, nil
			},
			result: true,
		},
		{
			name: "not has oracle",
			malleate: func() (contract.HasOracleArgs, error) {
				return contract.HasOracleArgs{
					Chain:           ethtypes.ModuleName,
					ExternalAddress: helpers.GenHexAddress(),
				}, nil
			},
			result: false,
		},
		{
			name: "invalid chain",
			malleate: func() (contract.HasOracleArgs, error) {
				return contract.HasOracleArgs{
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

			result := suite.WithError(expectErr).HasOracle(suite.Ctx, args)
			suite.Require().Equal(tc.result, result)
		})
	}
}
