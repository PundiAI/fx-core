package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/functionx/fx-core/tests"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
)

func (suite *KeeperTestSuite) TestQueryFIP20() {
	var contract common.Address
	testCases := []struct {
		name     string
		malleate func()
		res      bool
	}{
		{
			"fip20 not deployed",
			func() { contract = common.Address{} },
			false,
		},
		{
			"ok",
			func() { contract = suite.DeployContract(types.ModuleAddress, "coin", "token", 18) },
			true,
		},
	}
	for _, tc := range testCases {
		suite.SetupTest() // reset

		tc.malleate()

		res, err := suite.app.IntrarelayerKeeper.QueryFIP20(suite.ctx, contract)
		if tc.res {
			suite.Require().NoError(err)
			suite.Require().Equal(
				types.FIP20Data{Name: "coin", Symbol: "token", Decimals: 0x12},
				res,
			)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestCallEVM() {
	testCases := []struct {
		name    string
		method  string
		expPass bool
	}{
		{
			"unknown method",
			"",
			false,
		},
		{
			"pass",
			"balanceOf",
			true,
		},
	}
	for _, tc := range testCases {
		suite.SetupTest() // reset

		fip20ABI := contracts.MustGetABI(suite.ctx.BlockHeight(), contracts.FIP20UpgradeType)
		contract := suite.DeployContract(types.ModuleAddress, "coin", "token", 18)
		account := tests.GenerateAddress()

		res, err := suite.app.IntrarelayerKeeper.CallEVM(suite.ctx, fip20ABI, types.ModuleAddress, contract, tc.method, account)
		if tc.expPass {
			suite.Require().IsTypef(&evmtypes.MsgEthereumTxResponse{}, res, tc.name)
			suite.Require().NoError(err)
		} else {
			suite.Require().Error(err)
		}
	}
}
