package crosschain_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/crosschain"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func TestCrosschainGetERC20Tokwn(t *testing.T) {
	getERC20TokenABI := crosschain.NewGetERC20TokenMethod(nil)
	require.Len(t, getERC20TokenABI.Method.Inputs, 1)
	require.Len(t, getERC20TokenABI.Method.Outputs, 2)
}

func (suite *CrosschainPrecompileTestSuite) TestGetERC20Token() {
	testCases := []struct {
		name     string
		malleate func() (contract.GetERC20TokenArgs, error)
		result   bool
	}{
		{
			name: "erc20 token",
			malleate: func() (contract.GetERC20TokenArgs, error) {
				denom := "usdt"
				err := suite.App.Erc20Keeper.ERC20Token.Set(suite.Ctx, denom, erc20types.ERC20Token{
					Erc20Address:  helpers.GenHexAddress().Hex(),
					Denom:         denom,
					Enabled:       false,
					ContractOwner: erc20types.OWNER_MODULE,
				})
				suite.Require().NoError(err)
				return contract.GetERC20TokenArgs{
					Denom: contract.MustStrToByte32(denom),
				}, nil
			},
			result: true,
		},
		{
			name: "denom not found",
			malleate: func() (contract.GetERC20TokenArgs, error) {
				return contract.GetERC20TokenArgs{
					Denom: contract.MustStrToByte32("test"),
				}, fmt.Errorf("collections: not found: key 'test' of type github.com/cosmos/gogoproto/fx.erc20.v1.ERC20Token")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			args, expectErr := tc.malleate()
			token, enable := suite.WithError(expectErr).GetERC20Token(suite.Ctx, args)
			if tc.result {
				erc20Token, err := suite.App.Erc20Keeper.GetERC20Token(suite.Ctx, contract.Byte32ToString(args.Denom))
				suite.Require().NoError(err)
				suite.Require().Equal(common.HexToAddress(erc20Token.Erc20Address), token)
				suite.Require().Equal(erc20Token.Enabled, enable)
			}
		})
	}
}
