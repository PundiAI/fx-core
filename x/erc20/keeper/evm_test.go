package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestQueryERC20() {
	var contract common.Address
	testCases := []struct {
		name     string
		malleate func()
		res      bool
	}{
		{
			"erc20 not deployed",
			func() { contract = common.Address{} },
			false,
		},
		{
			"ok",
			func() { contract, _ = suite.DeployContract(suite.address, "coin", "token", erc20Decimals) },
			true,
		},
	}
	for _, tc := range testCases {
		suite.SetupTest() // reset

		tc.malleate()

		res, err := suite.app.Erc20Keeper.QueryERC20(suite.ctx, contract)
		if tc.res {
			suite.Require().NoError(err)
			suite.Require().Equal(
				types.ERC20Data{Name: "coin", Symbol: "token", Decimals: erc20Decimals},
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

		erc20Config := fxtypes.GetERC20()
		contract, err := suite.DeployContract(suite.address, "coin", "token", erc20Decimals)
		suite.Require().NoError(err)
		account := helpers.GenerateAddress()

		res, err := suite.app.Erc20Keeper.CallEVM(suite.ctx, erc20Config.ABI, types.ModuleAddress, contract, true, tc.method, account)
		if tc.expPass {
			suite.Require().IsTypef(&evmtypes.MsgEthereumTxResponse{}, res, tc.name)
			suite.Require().NoError(err)
		} else {
			suite.Require().Error(err)
		}
	}
}
