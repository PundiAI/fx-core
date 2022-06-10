package keeper_test

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/crosschain/types"
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
			"batch fee BaseFee is Negative",
			func() {
				request = &types.QueryBatchFeeRequest{
					ChainName: "tron",
					MinBatchFees: []types.MinBatchFee{
						{
							TokenContract: suite.bridgeTokens[0].token,
							BaseFee:       sdk.NewInt(-1),
						},
					},
				}
			},
			false,
		},
		{
			name: "batch fee normal",
			malleate: func() {
				minBatchFee := []types.MinBatchFee{
					{
						TokenContract: suite.bridgeTokens[0].token,
						BaseFee:       sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100))),
					},
				}
				for i := uint64(1); i <= 3; i++ {
					_, err := suite.app.TronKeeper.AddToOutgoingPool(
						suite.ctx,
						suite.bridgeAcc,
						GenTronContractAddress(),
						sdk.NewCoin(suite.bridgeTokens[0].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))),
						sdk.NewCoin(suite.bridgeTokens[0].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))))
					suite.Require().NoError(err)
				}
				for i := uint64(1); i <= 2; i++ {
					_, err := suite.app.TronKeeper.AddToOutgoingPool(
						suite.ctx,
						suite.bridgeAcc,
						GenTronContractAddress(),
						sdk.NewCoin(suite.bridgeTokens[0].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))),
						sdk.NewCoin(suite.bridgeTokens[0].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(10)))))
					suite.Require().NoError(err)
				}
				request = &types.QueryBatchFeeRequest{
					ChainName:    "tron",
					MinBatchFees: minBatchFee,
				}
				response = &types.QueryBatchFeeResponse{BatchFees: []*types.BatchFees{
					{
						TokenContract: suite.bridgeTokens[0].token,
						TotalFees:     sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(300))),
						TotalTxs:      3,
						TotalAmount:   sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(300))),
					},
				}}
			},
			expPass: true,
		},
		{
			name: "batch fee mul normal",
			malleate: func() {
				minBatchFee := []types.MinBatchFee{
					{
						TokenContract: suite.bridgeTokens[0].token,
						BaseFee:       sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(100), big.NewInt(1e6))),
					},
					{
						TokenContract: suite.bridgeTokens[1].token,
						BaseFee:       sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))),
					},
				}
				for i := uint64(1); i <= 2; i++ {
					_, err := suite.app.TronKeeper.AddToOutgoingPool(
						suite.ctx,
						suite.bridgeAcc,
						GenTronContractAddress(),
						sdk.NewCoin(suite.bridgeTokens[0].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))),
						sdk.NewCoin(suite.bridgeTokens[0].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(10)))))
					suite.Require().NoError(err)
				}
				_, err := suite.app.TronKeeper.AddToOutgoingPool(
					suite.ctx,
					suite.bridgeAcc,
					GenTronContractAddress(),
					sdk.NewCoin(suite.bridgeTokens[0].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))),
					sdk.NewCoin(suite.bridgeTokens[0].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100)))))
				suite.Require().NoError(err)
				for i := uint64(1); i <= 3; i++ {
					_, err := suite.app.TronKeeper.AddToOutgoingPool(
						suite.ctx,
						suite.bridgeAcc,
						GenTronContractAddress(),
						sdk.NewCoin(suite.bridgeTokens[1].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100)))),
						sdk.NewCoin(suite.bridgeTokens[1].denom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100)))))
					suite.Require().NoError(err)
				}
				request = &types.QueryBatchFeeRequest{
					ChainName:    "tron",
					MinBatchFees: minBatchFee,
				}
				response = &types.QueryBatchFeeResponse{BatchFees: []*types.BatchFees{
					{
						TokenContract: suite.bridgeTokens[0].token,
						TotalFees:     sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100))),
						TotalTxs:      1,
						TotalAmount:   sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(100))),
					},
					{
						TokenContract: suite.bridgeTokens[1].token,
						TotalFees:     sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(300))),
						TotalTxs:      3,
						TotalAmount:   sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(300))),
					},
				}}
			},
			expPass: true,
		},
		{
			"batch fee abnormal",
			func() {
				request = &types.QueryBatchFeeRequest{
					ChainName: "tron",
					MinBatchFees: []types.MinBatchFee{
						{
							TokenContract: "0xaD6D458402F60fD3Bd25163575031ACDce07538D",
						},
					},
				}
			},
			false,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()

			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()

			res, err := suite.queryClient.BatchFees(ctx, request)

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
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "tron",
					TokenContract: suite.bridgeTokens[0].token,
					Nonce:         3,
				}
				err := suite.app.TronKeeper.StoreBatch(suite.ctx, &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Token: types.ERC20Token{
								Contract: suite.bridgeTokens[0].token,
								Amount:   sdk.NewIntFromBigInt(big.NewInt(1e18)),
							},
							Fee: types.ERC20Token{
								Contract: suite.bridgeTokens[0].token,
								Amount:   sdk.NewIntFromBigInt(big.NewInt(1e18)),
							},
						},
					},
					TokenContract: suite.bridgeTokens[0].token,
					FeeReceive:    suite.bridgeTokens[1].token,
				})
				suite.Require().NoError(err)
				response = &types.QueryBatchRequestByNonceResponse{
					Batch: &types.OutgoingTxBatch{
						BatchNonce:   3,
						BatchTimeout: 10000,
						Transactions: []*types.OutgoingTransferTx{
							{
								Token: types.ERC20Token{
									Contract: suite.bridgeTokens[0].token,
									Amount:   sdk.NewIntFromBigInt(big.NewInt(1e18)),
								},
								Fee: types.ERC20Token{
									Contract: suite.bridgeTokens[0].token,
									Amount:   sdk.NewIntFromBigInt(big.NewInt(1e18)),
								},
							},
						},
						TokenContract: suite.bridgeTokens[0].token,
						FeeReceive:    suite.bridgeTokens[1].token,
					},
				}
			},
			expPass: true,
		},
		{
			name: "request error nonce",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "tron",
					TokenContract: suite.bridgeTokens[0].token,
					Nonce:         0,
				}
			},
			expPass: false,
		},
		{
			name: "request error token",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "tron",
					TokenContract: "0xaD6D458402F60fD3Bd25163575031ACDce07538D",
					Nonce:         8,
				}
			},
			expPass: false,
		},
		{
			name: "request nonexistent nonce",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "tron",
					TokenContract: suite.bridgeTokens[0].token,
					Nonce:         8,
				}
			},
			expPass: false,
		},
		{
			name: "request nonexistent token",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "tron",
					TokenContract: GenTronContractAddress(),
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
			res, err := suite.queryClient.BatchRequestByNonce(sdk.WrapSDKContext(suite.ctx), request)
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
					ChainName:     "tron",
					TokenContract: "0xaD6D458402F60fD3Bd25163575031ACDce07538D",
				}
			},
			false,
		},
		{
			"token nonce is zero",
			func() {
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     "tron",
					TokenContract: GenTronContractAddress(),
					Nonce:         0,
				}
			},
			false,
		},
		{
			"request confirm nonexistent nonce",
			func() {
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     "tron",
					TokenContract: suite.bridgeTokens[0].token,
					Nonce:         2,
				}
				suite.app.TronKeeper.SetBatchConfirm(suite.ctx, suite.bridgeAcc, &types.MsgConfirmBatch{
					Nonce:     1,
					ChainName: "tron",
				})
				response = &types.QueryBatchConfirmsResponse{}
			},
			true,
		},
		{
			"set correct batch confirm",
			func() {
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     "tron",
					TokenContract: suite.bridgeTokens[0].token,
					Nonce:         1,
				}
				suite.app.TronKeeper.SetBatchConfirm(suite.ctx, suite.bridgeAcc, &types.MsgConfirmBatch{
					Nonce:           1,
					TokenContract:   suite.bridgeTokens[0].token,
					BridgerAddress:  suite.bridgeAcc.String(),
					ExternalAddress: suite.bridgeTokens[1].token,
					Signature:       "0x1",
					ChainName:       "tron",
				})
				response = &types.QueryBatchConfirmsResponse{Confirms: []*types.MsgConfirmBatch{
					{
						Nonce:           1,
						TokenContract:   suite.bridgeTokens[0].token,
						BridgerAddress:  suite.bridgeAcc.String(),
						ExternalAddress: suite.bridgeTokens[1].token,
						Signature:       "0x1",
						ChainName:       "tron",
					},
				}}
			},
			true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()

			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()

			res, err := suite.queryClient.BatchConfirms(ctx, request)

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
					ChainName: "tron",
					Token:     "0xaD6D458402F60fD3Bd25163575031ACDce07538D",
				}
			},
			expPass: false,
		},
		{
			name: "token that does not exist",
			malleate: func() {
				request = &types.QueryTokenToDenomRequest{
					ChainName: "tron",
					Token:     GenTronContractAddress(),
				}
			},
			expPass: false,
		},
		{
			name: "token normal",
			malleate: func() {
				request = &types.QueryTokenToDenomRequest{
					ChainName: "tron",
					Token:     suite.bridgeTokens[0].token,
				}
				response = &types.QueryTokenToDenomResponse{
					Denom:      suite.bridgeTokens[0].denom,
					ChannelIbc: "",
				}
			},
			expPass: true,
		},
		{
			name: "token is channel ibc normal",
			malleate: func() {
				request = &types.QueryTokenToDenomRequest{
					ChainName: "tron",
					Token:     suite.bridgeTokens[2].token,
				}

				response = &types.QueryTokenToDenomResponse{
					Denom:      suite.bridgeTokens[2].denom,
					ChannelIbc: "transfer/channel-0",
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

			res, err := suite.queryClient.TokenToDenom(ctx, request)

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
					ChainName:       "tron",
					ExternalAddress: "0xaD6D458402F60fD3Bd25163575031ACDce07538D",
				}
			},
			expPass: false,
		},
		{
			name: "nonexistent external address",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       "tron",
					ExternalAddress: GenTronContractAddress(),
				}
			},
			expPass: false,
		},
		{
			name: "normal external address and oracle",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       "tron",
					ExternalAddress: suite.externalAccList[0].address,
				}
				response = &types.QueryOracleResponse{Oracle: &types.Oracle{
					OracleAddress:   suite.oracleAddressList[0].String(),
					BridgerAddress:  suite.orchestratorAddressList[0].String(),
					ExternalAddress: suite.externalAccList[0].address,
					DelegateAmount:  sdk.ZeroInt(),
					StartHeight:     3,
				}}
			},
			expPass: true,
		},
		{
			name: "nonexistent oracle",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       "tron",
					ExternalAddress: suite.externalAccList[1].address,
				}
			},
			expPass: false,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			testCase.malleate()
			res, err := suite.queryClient.GetOracleByExternalAddr(sdk.WrapSDKContext(suite.ctx), request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response.Oracle, res.Oracle)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
