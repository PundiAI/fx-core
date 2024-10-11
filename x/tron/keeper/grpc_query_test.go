package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
)

func (suite *KeeperTestSuite) TestKeeper_OutgoingTxBatch() {
	var (
		request  *types.QueryOutgoingTxBatchRequest
		response *types.QueryOutgoingTxBatchResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "store normal batch",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenHexAddress().Bytes())
				feeReceive := helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex())
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: bridgeToken[0].Token,
					Nonce:         3,
				}
				err := suite.App.TronKeeper.StoreBatch(suite.Ctx, &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Token: types.ERC20Token{
								Contract: bridgeToken[0].Token,
								Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
							},
							Fee: types.ERC20Token{
								Contract: bridgeToken[0].Token,
								Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
							},
						},
					},
					TokenContract: bridgeToken[0].Token,
					FeeReceive:    feeReceive,
				})
				suite.Require().NoError(err)
				response = &types.QueryOutgoingTxBatchResponse{
					Batch: &types.OutgoingTxBatch{
						BatchNonce:   3,
						BatchTimeout: 10000,
						Transactions: []*types.OutgoingTransferTx{
							{
								Token: types.ERC20Token{
									Contract: bridgeToken[0].Token,
									Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
								},
								Fee: types.ERC20Token{
									Contract: bridgeToken[0].Token,
									Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
								},
							},
						},
						TokenContract: bridgeToken[0].Token,
						FeeReceive:    feeReceive,
					},
				}
			},
			expPass: true,
		},
		{
			name: "request error nonce",
			malleate: func() {
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					Nonce:         0,
				}
			},
			expPass: false,
		},
		{
			name: "request error token",
			malleate: func() {
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.GenHexAddress().Hex(),
					Nonce:         8,
				}
			},
			expPass: false,
		},
		{
			name: "request nonexistent nonce",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: bridgeToken[0].Token,
					Nonce:         8,
				}
			},
			expPass: false,
		},
		{
			name: "request nonexistent token",
			malleate: func() {
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					Nonce:         8,
				}
			},
			expPass: false,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			testCase.malleate()
			res, err := suite.queryServer.OutgoingTxBatch(suite.Ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response.Batch, res.Batch)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_BatchConfirms() {
	var (
		request  *types.QueryBatchConfirmsRequest
		response *types.QueryBatchConfirmsResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"token address error",
			func() {
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.GenHexAddress().Hex(),
				}
			},
			false,
		},
		{
			"token nonce is zero",
			func() {
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					Nonce:         0,
				}
			},
			false,
		},
		{
			name: "request confirm nonexistent nonce",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: bridgeToken[0].Token,
					Nonce:         2,
				}
				suite.App.TronKeeper.SetBatchConfirm(suite.Ctx, suite.signer.AccAddress(), &types.MsgConfirmBatch{
					Nonce: 1,
				})
				response = &types.QueryBatchConfirmsResponse{}
			},
			expPass: true,
		},
		{
			"set correct batch confirm",
			func() {
				_, bridger, externalKey := suite.NewOracleByBridger()
				bridgeToken := suite.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: bridgeToken[0].Token,
					Nonce:         1,
				}
				newConfirmBatch := &types.MsgConfirmBatch{
					ChainName:       trontypes.ModuleName,
					Nonce:           1,
					TokenContract:   bridgeToken[0].Token,
					BridgerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       helpers.GenHexAddress().Hex(),
				}
				suite.App.TronKeeper.SetBatchConfirm(suite.Ctx, suite.signer.AccAddress(), newConfirmBatch)
				response = &types.QueryBatchConfirmsResponse{Confirms: []*types.MsgConfirmBatch{newConfirmBatch}}
			},
			true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()

			ctx := suite.Ctx
			testCase.malleate()

			res, err := suite.queryServer.BatchConfirms(ctx, request)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().ElementsMatch(response.Confirms, res.Confirms)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_TokenToDenom() {
	var (
		request  *types.QueryTokenToDenomRequest
		response *types.QueryTokenToDenomResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "token address error",
			malleate: func() {
				request = &types.QueryTokenToDenomRequest{
					ChainName: trontypes.ModuleName,
					Token:     helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "token that does not exist",
			malleate: func() {
				request = &types.QueryTokenToDenomRequest{
					ChainName: trontypes.ModuleName,
					Token:     helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "token normal",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryTokenToDenomRequest{
					ChainName: trontypes.ModuleName,
					Token:     bridgeToken[0].Token,
				}
				response = &types.QueryTokenToDenomResponse{
					Denom: bridgeToken[0].Denom,
				}
			},
			expPass: true,
		},
		{
			name: "token is channel ibc normal",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryTokenToDenomRequest{
					ChainName: trontypes.ModuleName,
					Token:     bridgeToken[2].Token,
				}
				response = &types.QueryTokenToDenomResponse{
					Denom: bridgeToken[2].Denom,
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()

			ctx := suite.Ctx
			testCase.malleate()

			res, err := suite.queryServer.TokenToDenom(ctx, request)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetOracleByExternalAddr() {
	var (
		request  *types.QueryOracleByExternalAddrRequest
		response *types.QueryOracleResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "external address is error",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "nonexistent external address",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "normal external address and oracle",
			malleate: func() {
				oracle, bridger, externalKey := suite.NewOracleByBridger()
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
				}
				response = &types.QueryOracleResponse{Oracle: &types.Oracle{
					OracleAddress:   oracle.String(),
					BridgerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					DelegateAmount:  sdkmath.ZeroInt(),
				}}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			testCase.malleate()
			res, err := suite.queryServer.GetOracleByExternalAddr(suite.Ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response.Oracle, res.Oracle)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
