package keeper_test

import (
	"encoding/hex"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	trontypes "github.com/pundiai/fx-core/v8/x/tron/types"
)

func (suite *KeeperTestSuite) Test_msgServer_ConfirmBatch() {
	var msg *crosschaintypes.MsgConfirmBatch
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "couldn't find batch",
			malleate: func() {
				msg = &crosschaintypes.MsgConfirmBatch{
					Nonce:          tmrand.Uint64(),
					TokenContract:  helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					BridgerAddress: helpers.GenAccAddress().String(),
				}
			},
			expPass: false,
		},
		{
			name: "no found oracle",
			malleate: func() {
				newOutgoingTx := suite.NewOutgoingTxBatch()
				msg = &crosschaintypes.MsgConfirmBatch{
					Nonce:          newOutgoingTx.BatchNonce,
					TokenContract:  newOutgoingTx.TokenContract,
					BridgerAddress: helpers.GenAccAddress().String(),
				}
			},
			expPass: false,
		},
		{
			name: "signature decoding",
			malleate: func() {
				newOutgoingTx := suite.NewOutgoingTxBatch()
				_, bridger, externalKey := suite.NewOracleByBridger()
				msg = &crosschaintypes.MsgConfirmBatch{
					Nonce:           newOutgoingTx.BatchNonce,
					TokenContract:   newOutgoingTx.TokenContract,
					BridgerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "confirm batch",
			malleate: func() {
				newOutgoingTx := suite.NewOutgoingTxBatch()
				_, bridger, externalKey := suite.NewOracleByBridger()
				params, err := suite.queryServer.Params(suite.Ctx, &crosschaintypes.QueryParamsRequest{ChainName: trontypes.ModuleName})
				suite.Require().NoError(err)
				checkpoint, err := newOutgoingTx.GetCheckpoint(params.Params.GravityId)
				suite.Require().NoError(err)
				key, err := externalKey.(*ethsecp256k1.PrivKey).ToECDSA()
				suite.Require().NoError(err)
				signature, err := trontypes.NewTronSignature(checkpoint, key)
				suite.Require().NoError(err)
				msg = &crosschaintypes.MsgConfirmBatch{
					Nonce:           newOutgoingTx.BatchNonce,
					TokenContract:   newOutgoingTx.TokenContract,
					BridgerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       hex.EncodeToString(signature),
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			_, err := suite.msgServer.ConfirmBatch(suite.Ctx, msg)
			if testCase.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().ErrorContains(err, testCase.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_msgServer_OracleSetConfirm() {
	var msg *crosschaintypes.MsgOracleSetConfirm
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "couldn't find oracleSet",
			malleate: func() {
				msg = &crosschaintypes.MsgOracleSetConfirm{
					Nonce:          tmrand.Uint64(),
					BridgerAddress: helpers.GenAccAddress().String(),
				}
			},
			expPass: false,
		},
		{
			name: "no found oracle",
			malleate: func() {
				newOracleSet := suite.NewOracleSet(helpers.NewEthPrivKey())
				msg = &crosschaintypes.MsgOracleSetConfirm{
					Nonce:           newOracleSet.Nonce,
					BridgerAddress:  helpers.GenAccAddress().String(),
					ExternalAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "signature decoding",
			malleate: func() {
				_, bridger, externalKey := suite.NewOracleByBridger()
				newOracleSet := suite.NewOracleSet(externalKey)
				msg = &crosschaintypes.MsgOracleSetConfirm{
					Nonce:           newOracleSet.Nonce,
					BridgerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					Signature:       helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "oracle set confirm",
			malleate: func() {
				_, bridger, externalKey := suite.NewOracleByBridger()
				newOracleSet := suite.NewOracleSet(externalKey)
				key, err := externalKey.(*ethsecp256k1.PrivKey).ToECDSA()
				suite.Require().NoError(err)
				params, err := suite.queryServer.Params(suite.Ctx, &crosschaintypes.QueryParamsRequest{ChainName: trontypes.ModuleName})
				suite.Require().NoError(err)
				checkpoint, err := newOracleSet.GetCheckpoint(params.Params.GravityId)
				suite.Require().NoError(err)
				signature, err := trontypes.NewTronSignature(checkpoint, key)
				suite.Require().NoError(err)
				msg = &crosschaintypes.MsgOracleSetConfirm{
					Nonce:           newOracleSet.Nonce,
					BridgerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       hex.EncodeToString(signature),
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			testCase.malleate()
			_, err := suite.msgServer.OracleSetConfirm(suite.Ctx, msg)
			if testCase.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().ErrorContains(err, testCase.name)
			}
		})
	}
}
