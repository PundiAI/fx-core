package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_BatchFees() {
	var (
		request  *types.QueryBatchFeeRequest
		response *types.QueryBatchFeeResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"baseFee is negative",
			func() {
				request = &types.QueryBatchFeeRequest{
					MinBatchFees: []types.MinBatchFee{
						{
							TokenContract: helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
							BaseFee:       sdkmath.NewInt(-1),
						},
					},
				}
			},
			false,
		},
		{
			"validate tron address error",
			func() {
				request = &types.QueryBatchFeeRequest{
					MinBatchFees: []types.MinBatchFee{
						{
							TokenContract: helpers.GenerateAddress().Hex(),
						},
					},
				}
			},
			false,
		},
		{
			name: "baseFee normal",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenerateAddress().Bytes())
				minBatchFee := []types.MinBatchFee{
					{
						TokenContract: bridgeToken[0].Token,
						BaseFee:       sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100))),
					},
				}
				for i := uint64(1); i <= 3; i++ {
					_, err := suite.app.TronKeeper.AddToOutgoingPool(
						suite.ctx,
						suite.signer.AccAddress(),
						helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
						sdk.NewCoin(bridgeToken[0].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))),
						sdk.NewCoin(bridgeToken[0].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))))
					suite.Require().NoError(err)
				}
				for i := uint64(1); i <= 2; i++ {
					_, err := suite.app.TronKeeper.AddToOutgoingPool(
						suite.ctx,
						suite.signer.AccAddress(),
						helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
						sdk.NewCoin(bridgeToken[0].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))),
						sdk.NewCoin(bridgeToken[0].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(10)))))
					suite.Require().NoError(err)
				}
				request = &types.QueryBatchFeeRequest{
					MinBatchFees: minBatchFee,
				}
				response = &types.QueryBatchFeeResponse{BatchFees: []*types.BatchFees{
					{
						TokenContract: bridgeToken[0].Token,
						TotalFees:     sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(300))),
						TotalTxs:      3,
						TotalAmount:   sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(300))),
					},
				}}
			},
			expPass: true,
		},
		{
			name: "batch fee mul normal",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenerateAddress().Bytes())
				minBatchFee := []types.MinBatchFee{
					{
						TokenContract: bridgeToken[0].Token,
						BaseFee:       sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(100), big.NewInt(1e6))),
					},
					{
						TokenContract: bridgeToken[1].Token,
						BaseFee:       sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))),
					},
				}
				for i := 1; i <= 2; i++ {
					_, err := suite.app.TronKeeper.AddToOutgoingPool(
						suite.ctx,
						suite.signer.AccAddress(),
						helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
						sdk.NewCoin(bridgeToken[0].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))),
						sdk.NewCoin(bridgeToken[0].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(10)))))
					suite.Require().NoError(err)
				}
				_, err := suite.app.TronKeeper.AddToOutgoingPool(
					suite.ctx,
					suite.signer.AccAddress(),
					helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
					sdk.NewCoin(bridgeToken[0].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))),
					sdk.NewCoin(bridgeToken[0].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))))
				suite.Require().NoError(err)
				for i := 1; i <= 3; i++ {
					_, err = suite.app.TronKeeper.AddToOutgoingPool(
						suite.ctx,
						suite.signer.AccAddress(),
						helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
						sdk.NewCoin(bridgeToken[1].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100)))),
						sdk.NewCoin(bridgeToken[1].Denom, sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100)))))
					suite.Require().NoError(err)
				}
				request = &types.QueryBatchFeeRequest{
					MinBatchFees: minBatchFee,
				}
				response = &types.QueryBatchFeeResponse{BatchFees: []*types.BatchFees{
					{
						TokenContract: bridgeToken[0].Token,
						TotalFees:     sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100))),
						TotalTxs:      1,
						TotalAmount:   sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100))),
					},
					{
						TokenContract: bridgeToken[1].Token,
						TotalFees:     sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(300))),
						TotalTxs:      3,
						TotalAmount:   sdkmath.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(300))),
					},
				}}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()
			res, err := suite.app.TronKeeper.BatchFees(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().ElementsMatch(response.BatchFees, res.BatchFees)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_BatchRequestByNonce() {
	var (
		request  *types.QueryBatchRequestByNonceRequest
		response *types.QueryBatchRequestByNonceResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "store normal batch",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenerateAddress().Bytes())
				feeReceive := helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex())
				request = &types.QueryBatchRequestByNonceRequest{
					TokenContract: bridgeToken[0].Token,
					Nonce:         3,
				}
				err := suite.app.TronKeeper.StoreBatch(suite.ctx, &types.OutgoingTxBatch{
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
				response = &types.QueryBatchRequestByNonceResponse{
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
				request = &types.QueryBatchRequestByNonceRequest{
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
					Nonce:         0,
				}
			},
			expPass: false,
		},
		{
			name: "request error token",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					TokenContract: helpers.GenerateAddress().Hex(),
					Nonce:         8,
				}
			},
			expPass: false,
		},
		{
			name: "request nonexistent nonce",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenerateAddress().Bytes())
				request = &types.QueryBatchRequestByNonceRequest{
					TokenContract: bridgeToken[0].Token,
					Nonce:         8,
				}
			},
			expPass: false,
		},
		{
			name: "request nonexistent token",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
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
			res, err := suite.app.TronKeeper.BatchRequestByNonce(sdk.WrapSDKContext(suite.ctx), request)
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
					TokenContract: helpers.GenerateAddress().Hex(),
				}
			},
			false,
		},
		{
			"token nonce is zero",
			func() {
				request = &types.QueryBatchConfirmsRequest{
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
					Nonce:         0,
				}
			},
			false,
		},
		{
			name: "request confirm nonexistent nonce",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenerateAddress().Bytes())
				request = &types.QueryBatchConfirmsRequest{
					TokenContract: bridgeToken[0].Token,
					Nonce:         2,
				}
				suite.app.TronKeeper.SetBatchConfirm(suite.ctx, suite.signer.AccAddress(), &types.MsgConfirmBatch{
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
				bridgeToken := suite.NewBridgeToken(helpers.GenerateAddress().Bytes())
				request = &types.QueryBatchConfirmsRequest{
					TokenContract: bridgeToken[0].Token,
					Nonce:         1,
				}
				newConfirmBatch := &types.MsgConfirmBatch{
					Nonce:           1,
					TokenContract:   bridgeToken[0].Token,
					BridgerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       helpers.GenerateAddress().Hex(),
				}
				suite.app.TronKeeper.SetBatchConfirm(suite.ctx, suite.signer.AccAddress(), newConfirmBatch)
				response = &types.QueryBatchConfirmsResponse{Confirms: []*types.MsgConfirmBatch{newConfirmBatch}}
			},
			true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()

			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()

			res, err := suite.app.TronKeeper.BatchConfirms(ctx, request)

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
					Token: helpers.GenerateAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "token that does not exist",
			malleate: func() {
				request = &types.QueryTokenToDenomRequest{
					Token: helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "token normal",
			malleate: func() {
				bridgeToken := suite.NewBridgeToken(helpers.GenerateAddress().Bytes())
				request = &types.QueryTokenToDenomRequest{
					Token: bridgeToken[0].Token,
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
				bridgeToken := suite.NewBridgeToken(helpers.GenerateAddress().Bytes())
				request = &types.QueryTokenToDenomRequest{
					Token: bridgeToken[2].Token,
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

			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()

			res, err := suite.app.TronKeeper.TokenToDenom(ctx, request)

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
					ExternalAddress: helpers.GenerateAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "nonexistent external address",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ExternalAddress: helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "normal external address and oracle",
			malleate: func() {
				oracle, bridger, externalKey := suite.NewOracleByBridger()
				request = &types.QueryOracleByExternalAddrRequest{
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
			res, err := suite.app.TronKeeper.GetOracleByExternalAddr(sdk.WrapSDKContext(suite.ctx), request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response.Oracle, res.Oracle)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
