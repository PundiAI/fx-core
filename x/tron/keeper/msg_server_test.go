package keeper_test

import (
	"encoding/hex"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"
)

func (suite *KeeperTestSuite) Test_msgServer_ConfirmBatch() {
	var (
		msg *types.MsgConfirmBatch
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "confirm batch msg bridge address error",
			malleate: func() {
				msg = &types.MsgConfirmBatch{
					Nonce:          3,
					BridgerAddress: "fx1",
					ChainName:      "tron",
				}
			},
			expPass: false,
		},
		{
			name: "confirm batch nonexistent tx nonce",
			malleate: func() {
				newOutgoingTx := &types.OutgoingTxBatch{
					BatchNonce:   2,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Sender:      suite.bridgeAcc.String(),
							DestAddress: GenTronContractAddress(),
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
				}

				err := suite.app.TronKeeper.StoreBatch(suite.ctx, newOutgoingTx)
				suite.Require().NoError(err)
				msg = &types.MsgConfirmBatch{
					Nonce:          3,
					TokenContract:  suite.bridgeTokens[0].token,
					BridgerAddress: suite.orchestratorAddressList[0].String(),
					ChainName:      "tron",
				}
			},
			expPass: false,
		},
		{
			name: "confirm batch nonexistent tx token contract",
			malleate: func() {
				newOutgoingTx := &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Sender:      suite.bridgeAcc.String(),
							DestAddress: GenTronContractAddress(),
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
				}

				err := suite.app.TronKeeper.StoreBatch(suite.ctx, newOutgoingTx)
				suite.Require().NoError(err)
				msg = &types.MsgConfirmBatch{
					Nonce:          3,
					TokenContract:  GenTronContractAddress(),
					BridgerAddress: suite.orchestratorAddressList[0].String(),
					ChainName:      "tron",
				}
			},
			expPass: false,
		},
		{
			name: "confirm normal batch tx",
			malleate: func() {
				newOutgoingTx := &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Sender:      suite.bridgeAcc.String(),
							DestAddress: GenTronContractAddress(),
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
				}
				batchHash, err := trontypes.GetCheckpointConfirmBatch(newOutgoingTx, "tron")
				suite.Require().NoError(err)
				signature, err := trontypes.NewTronSignature(batchHash, suite.externalAccList[0].key)
				suite.Require().NoError(err)
				msg = &types.MsgConfirmBatch{
					Nonce:           3,
					TokenContract:   suite.bridgeTokens[0].token,
					BridgerAddress:  suite.orchestratorAddressList[0].String(),
					ExternalAddress: suite.externalAccList[0].address,
					Signature:       hex.EncodeToString(signature),
					ChainName:       "tron",
				}

				err = suite.app.TronKeeper.StoreBatch(suite.ctx, newOutgoingTx)
				suite.Require().NoError(err)

			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			testCase.malleate()
			suite.Require().Empty(suite.app.TronKeeper.GetBatchConfirm(suite.ctx, 3, suite.bridgeTokens[0].token, suite.oracleAddressList[0]))
			_, err := suite.msgServer.ConfirmBatch(sdk.WrapSDKContext(suite.ctx), msg)
			confirm := suite.app.TronKeeper.GetBatchConfirm(suite.ctx, 3, suite.bridgeTokens[0].token, suite.oracleAddressList[0])
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(msg, confirm)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_msgServer_OracleSetConfirm() {
	var (
		msg *types.MsgOracleSetConfirm
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "oracle set bridger Address error msg",
			malleate: func() {
				msg = &types.MsgOracleSetConfirm{
					Nonce:          3,
					BridgerAddress: "fx1",
					ChainName:      "tron",
				}
			},
			expPass: false,
		},
		{
			name: "oracle set nonexistent nonce msg",
			malleate: func() {
				newOracleSet := types.NewOracleSet(2, 10, types.BridgeValidators{
					{
						Power:           1000000,
						ExternalAddress: suite.externalAccList[0].address,
					},
				})
				suite.app.TronKeeper.StoreOracleSet(suite.ctx, newOracleSet)
				msg = &types.MsgOracleSetConfirm{
					Nonce:           3,
					BridgerAddress:  suite.orchestratorAddressList[0].String(),
					ExternalAddress: suite.externalAccList[0].address,
					ChainName:       "tron",
				}
			},
			expPass: false,
		},
		{
			name: "oracle set checkpoint error msg",
			malleate: func() {
				newOracleSet := types.NewOracleSet(3, 10, types.BridgeValidators{
					{
						Power:           1000000,
						ExternalAddress: suite.externalAccList[0].address,
					},
				})
				suite.app.TronKeeper.StoreOracleSet(suite.ctx, newOracleSet)

				msg = &types.MsgOracleSetConfirm{
					Nonce:           3,
					BridgerAddress:  suite.orchestratorAddressList[0].String(),
					ExternalAddress: suite.externalAccList[0].address,
					Signature:       "0x1",
					ChainName:       "tron",
				}
			},
			expPass: false,
		},
		{
			name: "oracle set normal batch tx",
			malleate: func() {
				newOracleSet := types.NewOracleSet(3, 10, types.BridgeValidators{
					{
						Power:           1000000,
						ExternalAddress: suite.externalAccList[0].address,
					},
				})
				suite.app.TronKeeper.StoreOracleSet(suite.ctx, newOracleSet)
				oracleSetHash, err := trontypes.GetCheckpointOracleSet(newOracleSet, "tron")
				suite.Require().NoError(err)
				signature, err := trontypes.NewTronSignature(oracleSetHash, suite.externalAccList[0].key)
				suite.Require().NoError(err)

				msg = &types.MsgOracleSetConfirm{
					Nonce:           3,
					BridgerAddress:  suite.orchestratorAddressList[0].String(),
					ExternalAddress: suite.externalAccList[0].address,
					Signature:       hex.EncodeToString(signature),
					ChainName:       "tron",
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			testCase.malleate()
			suite.Require().Empty(suite.app.TronKeeper.GetOracleSetConfirm(suite.ctx, 3, suite.oracleAddressList[0]))
			_, err := suite.msgServer.OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), msg)
			confirm := suite.app.TronKeeper.GetOracleSetConfirm(suite.ctx, 3, suite.oracleAddressList[0])
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(msg, confirm)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
