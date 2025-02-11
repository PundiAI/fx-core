package keeper_test

import (
	"encoding/hex"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestBridgeCallHandler() {
	testCases := []struct {
		name    string
		initMsg func(msg *types.MsgBridgeCallClaim)
		assert  func(msg types.MsgBridgeCallClaim, err error)
	}{
		{
			name: "success - bridge token",
			initMsg: func(msg *types.MsgBridgeCallClaim) {
				msg.To = helpers.GenExternalAddr(suite.chainName)
				msg.Memo = types.MemoSendCallTo
			},
			assert: func(msg types.MsgBridgeCallClaim, resErr error) {
				suite.Require().NoError(resErr)
				for i, tokenContract := range msg.TokenContracts {
					erc20Token := suite.GetERC20TokenByBridgeContract(tokenContract)

					suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract())
					balanceOf := suite.erc20TokenSuite.BalanceOf(suite.Ctx, msg.GetReceiverAddr())
					suite.Equal(msg.Amounts[i].BigInt().String(), balanceOf.String())
				}
			},
		},
		{
			name: "refund - bridge token",
			initMsg: func(msg *types.MsgBridgeCallClaim) {
				erc20Token := suite.GetERC20TokenByBridgeContract(msg.TokenContracts[0])
				msg.To = fxtypes.ExternalAddrToStr(suite.chainName, erc20Token.GetERC20Contract().Bytes())
				data := helpers.PackERC20Mint(helpers.GenHexAddress(), big.NewInt(100))
				msg.Data = hex.EncodeToString(data)
				msg.Memo = types.MemoSendCallTo
			},
			assert: func(msg types.MsgBridgeCallClaim, resErr error) {
				suite.Require().NoError(resErr)

				outgoingBridgeCall, found := suite.Keeper().GetOutgoingBridgeCallByNonce(suite.Ctx, 1)
				suite.True(found)
				suite.Equal(hex.EncodeToString([]byte("Ownable: caller is not the owner: evm transaction execution failed")), outgoingBridgeCall.Data)
				suite.Equal(msg.EventNonce, outgoingBridgeCall.EventNonce)
				suite.Equal(
					fxtypes.ExternalAddrToStr(suite.chainName, suite.Keeper().GetCallbackFrom().Bytes()),
					outgoingBridgeCall.Sender,
				)

				for i, tokenContract := range msg.TokenContracts {
					erc20Token := suite.GetERC20TokenByBridgeContract(tokenContract)
					suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract())
					balanceOf := suite.erc20TokenSuite.BalanceOf(suite.Ctx, msg.GetRefundAddr())
					suite.Equal(msg.Amounts[i].BigInt().String(), balanceOf.String())
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := types.MsgBridgeCallClaim{
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
				TxOrigin: helpers.GenExternalAddr(suite.chainName),
			}
			for i, tokenContract := range msg.TokenContracts {
				denom := helpers.NewRandDenom()
				err := suite.Keeper().AddBridgeTokenExecuted(suite.Ctx,
					&types.MsgBridgeTokenClaim{
						TokenContract: tokenContract,
						Name:          denom,
						Symbol:        strings.ToUpper(denom),
						Decimals:      18,
					})
				suite.Require().NoError(err)

				if i != 0 {
					continue
				}
				quoteInfo := suite.bridgeFeeSuite.MockQuote(suite.Ctx, suite.chainName, denom)
				msg.QuoteId = sdkmath.NewIntFromBigInt(quoteInfo.Id)
				msg.GasLimit = sdkmath.NewInt(int64(quoteInfo.GasLimit))

				erc20Token, err := suite.App.Erc20Keeper.GetERC20Token(suite.Ctx, denom)
				suite.Require().NoError(err)
				suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract())
				suite.erc20TokenSuite.MintFromERC20Module(suite.Ctx, msg.GetSenderAddr(), big.NewInt(1))
			}
			tc.initMsg(&msg)
			err := suite.Keeper().BridgeCallExecuted(suite.Ctx, suite.App.EvmKeeper, &msg)
			tc.assert(msg, err)
		})
	}
}
