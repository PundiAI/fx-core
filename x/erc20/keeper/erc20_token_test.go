package keeper_test

import (
	"fmt"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (suite *KeeperTestSuite) TestKeeper_ToggleTokenConvert() {
	erc20Token, err := suite.GetKeeper().ToggleTokenConvert(suite.Ctx, "test")
	suite.Require().EqualError(err, "collections: not found: key 'test' of type github.com/cosmos/gogoproto/fx.erc20.v1.ERC20Token")
	suite.Require().Empty(erc20Token)
}

func (suite *KeeperTestSuite) TestKeeper_AddERC20Token() {
	testCases := []struct {
		name      string
		metadata  func() banktypes.Metadata
		erc20Addr common.Address
		Owner     types.Owner
		expErr    string
	}{
		{
			name: "success",
			metadata: func() banktypes.Metadata {
				return fxtypes.NewMetadata("test", helpers.NewRandSymbol(), 18)
			},
			erc20Addr: helpers.GenHexAddress(),
			Owner:     types.OWNER_MODULE,
		},
		{
			name: "failed - already exists",
			metadata: func() banktypes.Metadata {
				md := fxtypes.NewMetadata("test", helpers.NewRandSymbol(), 18)
				_, err := suite.GetKeeper().AddERC20Token(suite.Ctx, md, helpers.GenHexAddress(), types.OWNER_MODULE)
				suite.Require().NoError(err)
				return md
			},
			erc20Addr: helpers.GenHexAddress(),
			Owner:     types.OWNER_MODULE,
			expErr:    "denom %s is already registered: invalid request",
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			md := tc.metadata()
			erc20Token, err := suite.GetKeeper().AddERC20Token(suite.Ctx, md, tc.erc20Addr, tc.Owner)
			if tc.expErr != "" {
				suite.Require().Error(err)
				suite.Require().Equal(fmt.Sprintf(tc.expErr, md.Base), err.Error())
				return
			}
			suite.Require().NoError(err)
			suite.Equal(md.Base, erc20Token.Denom)
			suite.Equal(tc.erc20Addr.String(), erc20Token.Erc20Address)
			suite.Equal(tc.Owner, erc20Token.ContractOwner)
			suite.True(erc20Token.Enabled)
		})
	}
}
