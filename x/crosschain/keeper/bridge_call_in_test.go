package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestBridgeCallHandler() {
	suite.T().SkipNow() // todo: re-enable this test
	testCases := []struct {
		Name              string
		Msg               types.MsgBridgeCallClaim
		TokenIsNativeCoin []bool
		Success           bool
		CallContract      bool
	}{
		{
			Name: "success - token",
			Msg: types.MsgBridgeCallClaim{
				ChainName:      suite.chainName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    1,
				Sender:         helpers.GenExternalAddr(suite.chainName),
				Refund:         helpers.GenExternalAddr(suite.chainName),
				TokenContracts: []string{
					helpers.GenExternalAddr(suite.chainName),
					helpers.GenExternalAddr(suite.chainName),
				},
				Amounts: []sdkmath.Int{
					helpers.NewRandAmount(),
					helpers.NewRandAmount(),
				},
				To:       helpers.GenExternalAddr(suite.chainName),
				Data:     "",
				Value:    sdkmath.ZeroInt(),
				Memo:     "",
				TxOrigin: helpers.GenExternalAddr(suite.chainName),
			},
			TokenIsNativeCoin: []bool{true, true},
			Success:           true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			erc20Addrs := make([]common.Address, len(tc.Msg.TokenContracts))

			err := suite.Keeper().BridgeCallHandler(suite.Ctx, &tc.Msg)
			if tc.Success {
				suite.Require().NoError(err)
				if !tc.CallContract {
					for i, addr := range erc20Addrs {
						erc20Token := contract.NewERC20TokenKeeper(suite.App.EvmKeeper)
						balanceOf, err := erc20Token.BalanceOf(suite.Ctx, addr, tc.Msg.GetToAddr())
						suite.Require().NoError(err)
						suite.Equal(tc.Msg.Amounts[i].BigInt().String(), balanceOf.String())
					}
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
