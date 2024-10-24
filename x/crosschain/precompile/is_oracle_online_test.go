package precompile_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
)

func TestCrosschainIsOracleOnlineABI(t *testing.T) {
	method := precompile.NewIsOracleOnlineMethod(nil)
	require.Len(t, method.Method.Inputs, 2)
	require.Len(t, method.Method.Outputs, 1)
}

func (suite *PrecompileTestSuite) TestIsOracleOnline() {
	testCases := []struct {
		name     string
		malleate func() (crosschaintypes.IsOracleOnlineArgs, error)
		result   bool
	}{
		{
			name: "oracle online",
			malleate: func() (crosschaintypes.IsOracleOnlineArgs, error) {
				moduleName := suite.GenerateModuleName()
				oracle := suite.GenerateRandOracle(moduleName, true)
				return crosschaintypes.IsOracleOnlineArgs{
					Chain:           moduleName,
					ExternalAddress: oracle.GetExternalHexAddr(),
				}, nil
			},
			result: true,
		},
		{
			name: "oracle offline",
			malleate: func() (crosschaintypes.IsOracleOnlineArgs, error) {
				moduleName := suite.GenerateModuleName()
				oracle := suite.GenerateRandOracle(moduleName, false)
				return crosschaintypes.IsOracleOnlineArgs{
					Chain:           moduleName,
					ExternalAddress: oracle.GetExternalHexAddr(),
				}, nil
			},
			result: false,
		},
		{
			name: "oracle not found",
			malleate: func() (crosschaintypes.IsOracleOnlineArgs, error) {
				return crosschaintypes.IsOracleOnlineArgs{
					Chain:           ethtypes.ModuleName,
					ExternalAddress: helpers.GenHexAddress(),
				}, nil
			},
			result: false,
		},
		{
			name: "invalid chain",
			malleate: func() (crosschaintypes.IsOracleOnlineArgs, error) {
				return crosschaintypes.IsOracleOnlineArgs{
					Chain:           "invalid_chain",
					ExternalAddress: helpers.GenHexAddress(),
				}, fmt.Errorf("invalid module name: %s: evm transaction execution failed", "invalid_chain")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			signer := suite.RandSigner()

			args, expectErr := tc.malleate()

			hasOracle := precompile.NewIsOracleOnlineMethod(nil)
			packData, err := hasOracle.PackInput(args)
			suite.Require().NoError(err)

			res, err := suite.App.EvmKeeper.CallEVMWithoutGas(suite.Ctx, signer.Address(), &suite.crosschainAddr, nil, packData, false)
			if err != nil {
				suite.Require().EqualError(err, expectErr.Error())
				return
			}
			result, err := hasOracle.UnpackOutput(res.Ret)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.result, result)
		})
	}
}
