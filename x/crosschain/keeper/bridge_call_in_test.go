package keeper_test

import (
	"strings"

	sdkmath "cosmossdk.io/math"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (suite *KeeperTestSuite) TestBridgeCallHandler() {
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
				QuoteId:  sdkmath.ZeroInt(),
				GasLimit: sdkmath.ZeroInt(),
				Memo:     "",
				TxOrigin: helpers.GenExternalAddr(suite.chainName),
			},
			TokenIsNativeCoin: []bool{true, true},
			Success:           true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			erc20Tokens := make([]erc20types.ERC20Token, 0, len(tc.Msg.TokenContracts))
			for _, tokenContract := range tc.Msg.TokenContracts {
				denom := helpers.NewRandDenom()
				err := suite.App.Erc20Keeper.AddBridgeToken(suite.Ctx, denom, suite.chainName, tokenContract, false)
				suite.Require().NoError(err)

				erc20Token, err := suite.App.Erc20Keeper.RegisterNativeCoin(suite.Ctx, denom, strings.ToUpper(denom), 18)
				suite.Require().NoError(err)
				erc20Tokens = append(erc20Tokens, erc20Token)
			}

			err := suite.Keeper().BridgeCallHandler(suite.Ctx, &tc.Msg)
			if tc.Success {
				suite.Require().NoError(err)
				if !tc.CallContract {
					for i, erc20Token := range erc20Tokens {
						erc20TokenKeeper := contract.NewERC20TokenKeeper(suite.App.EvmKeeper)
						balanceOf, err := erc20TokenKeeper.BalanceOf(suite.Ctx, erc20Token.GetERC20Contract(), tc.Msg.GetToAddr())
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
