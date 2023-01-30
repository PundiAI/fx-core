package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestQueryERC20() {
	testCases := []struct {
		name     string
		malleate func() common.Address
		res      bool
	}{
		{
			"erc20 not deployed",
			func() common.Address { return common.Address{} },
			false,
		},
		{
			"ok",
			func() common.Address {
				contract, err := suite.DeployContract(suite.signer.Address())
				suite.NoError(err)
				return contract
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			contract := tc.malleate()
			res, err := suite.app.Erc20Keeper.QueryERC20(suite.ctx, contract)
			if tc.res {
				suite.Require().NoError(err)
				suite.Require().Equal(types.ERC20Data{Name: "Test token", Symbol: "TEST", Decimals: 18}, res)
			} else {
				suite.Require().Error(err)
			}
		})
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
		suite.Run(tc.name, func() {
			contract, err := suite.DeployContract(suite.signer.Address())
			suite.Require().NoError(err)

			account := helpers.GenerateAddress()
			erc20Config := fxtypes.GetERC20()
			res, err := suite.app.Erc20Keeper.CallEVM(suite.ctx, erc20Config.ABI, contract, contract, true, tc.method, account)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}
