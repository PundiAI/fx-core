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

func TestCrosschainHasOracleABI(t *testing.T) {
	method := precompile.NewHasOracleMethod(nil)
	require.Len(t, method.Method.Inputs, 2)
	require.Len(t, method.Method.Outputs, 1)
}

func (suite *PrecompileTestSuite) TestHasOracle() {
	testCases := []struct {
		name     string
		malleate func() (crosschaintypes.HasOracleArgs, error)
		result   bool
	}{
		{
			name: "has oracle",
			malleate: func() (crosschaintypes.HasOracleArgs, error) {
				moduleName := suite.GenerateModuleName()
				oracle := suite.GenerateRandOracle(moduleName, true)
				return crosschaintypes.HasOracleArgs{
					Chain:           moduleName,
					ExternalAddress: oracle.GetExternalHexAddr(),
				}, nil
			},
			result: true,
		},
		{
			name: "not has oracle",
			malleate: func() (crosschaintypes.HasOracleArgs, error) {
				return crosschaintypes.HasOracleArgs{
					Chain:           ethtypes.ModuleName,
					ExternalAddress: helpers.GenHexAddress(),
				}, nil
			},
			result: false,
		},
		{
			name: "invalid chain",
			malleate: func() (crosschaintypes.HasOracleArgs, error) {
				return crosschaintypes.HasOracleArgs{
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

			hasOracle := precompile.NewHasOracleMethod(nil)
			packData, err := hasOracle.PackInput(args)
			suite.Require().NoError(err)

			contractAddr := crosschaintypes.GetAddress()

			res, err := suite.App.EvmKeeper.CallEVMWithoutGas(suite.Ctx, signer.Address(), &contractAddr, nil, packData, false)
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
