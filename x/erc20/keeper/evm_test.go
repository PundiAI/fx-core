package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/x/erc20/types"
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
