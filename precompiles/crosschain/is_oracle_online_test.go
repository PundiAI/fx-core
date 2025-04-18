package crosschain_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/crosschain"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func TestCrosschainIsOracleOnlineABI(t *testing.T) {
	isOracleOnlineABI := crosschain.NewIsOracleOnlineABI()
	require.Len(t, isOracleOnlineABI.Method.Inputs, 2)
	require.Len(t, isOracleOnlineABI.Method.Outputs, 1)
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
				oracle := suite.SetOracle(true)
				return contract.IsOracleOnlineArgs{
					Chain:           contract.MustStrToByte32(suite.chainName),
					ExternalAddress: types.ExternalAddrToHexAddr(suite.chainName, oracle.ExternalAddress),
				}, nil
			},
			result: true,
		},
		{
			name: "oracle offline",
			malleate: func() (contract.IsOracleOnlineArgs, error) {
				oracle := suite.SetOracle(false)
				return contract.IsOracleOnlineArgs{
					Chain:           contract.MustStrToByte32(suite.chainName),
					ExternalAddress: types.ExternalAddrToHexAddr(suite.chainName, oracle.ExternalAddress),
				}, nil
			},
			result: false,
		},
		{
			name: "oracle not found",
			malleate: func() (contract.IsOracleOnlineArgs, error) {
				return contract.IsOracleOnlineArgs{
					Chain:           contract.MustStrToByte32(ethtypes.ModuleName),
					ExternalAddress: helpers.GenHexAddress(),
				}, nil
			},
			result: false,
		},
		{
			name: "invalid chain",
			malleate: func() (contract.IsOracleOnlineArgs, error) {
				return contract.IsOracleOnlineArgs{
					Chain:           contract.MustStrToByte32("invalid_chain"),
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
			suite.Require().Equal(tc.result, result, suite.chainName)
		})
	}
}
