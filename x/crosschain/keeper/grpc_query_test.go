package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/helpers"
	"github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/x/crosschain/types"
)

type CrossChainGrpcTestSuite struct {
	suite.Suite

	app *app.App
	ctx sdk.Context

	oracles  []sdk.AccAddress
	bridgers []sdk.AccAddress

	queryClient types.QueryClient
}

func TestCrossChainGrpcTestSuite(t *testing.T) {
	suite.Run(t, new(CrossChainGrpcTestSuite))
}

func (suite *CrossChainGrpcTestSuite) SetupTest() {
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(types.MaxOracleSize, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.CrosschainKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	suite.oracles = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.bridgers = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
}

func (suite *CrossChainGrpcTestSuite) Keeper() keeper.Keeper {
	return suite.app.BscKeeper
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_OracleSetRequest() {
	var (
		request       *types.QueryOracleSetRequestRequest
		response      *types.QueryCurrentOracleSetResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "oracle set nonce does not exist",
			malleate: func() {
				request = &types.QueryOracleSetRequestRequest{
					Nonce: 1,
				}
				response = &types.QueryCurrentOracleSetResponse{OracleSet: nil}
			},
			expPass: true,
		},
		{
			name: "oracle set nonce is zero",
			malleate: func() {
				request = &types.QueryOracleSetRequestRequest{
					Nonce: 0,
				}
				expectedError = sdkerrors.Wrapf(types.ErrUnknown, "nonce")
			},
			expPass: false,
		},
		{
			name: "normal oracle set",
			malleate: func() {
				members := []types.BridgeValidator{
					{
						Power:           10000,
						ExternalAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
					},
				}
				request = &types.QueryOracleSetRequestRequest{
					Nonce: 3,
				}
				suite.Keeper().StoreOracleSet(suite.ctx, &types.OracleSet{
					Nonce:   3,
					Members: members,
					Height:  100,
				})
				response = &types.QueryCurrentOracleSetResponse{
					OracleSet: types.NewOracleSet(3, 100, members),
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
			res, err := suite.Keeper().OracleSetRequest(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response.OracleSet, res.OracleSet)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_OracleSetConfirm() {
	var (
		request       *types.QueryOracleSetConfirmRequest
		response      *types.QueryOracleSetConfirmResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "oracle set bridger address error",
			malleate: func() {
				request = &types.QueryOracleSetConfirmRequest{
					ChainName:      "bsc",
					BridgerAddress: "fx1",
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "bridger address")
			},
			expPass: false,
		},
		{
			name: "oracle set nonce error",
			malleate: func() {
				request = &types.QueryOracleSetConfirmRequest{
					ChainName:      "bsc",
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
					Nonce:          0,
				}
				expectedError = sdkerrors.Wrap(types.ErrUnknown, "nonce")
			},
			expPass: false,
		},
		{
			name: "oracle set bridger address does not exist",
			malleate: func() {
				request = &types.QueryOracleSetConfirmRequest{
					ChainName:      "bsc",
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
					Nonce:          3,
				}
				expectedError = types.ErrNoFoundOracle
			},
			expPass: false,
		},
		{
			"oracle set normal",
			func() {
				bridger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				oracle := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())

				request = &types.QueryOracleSetConfirmRequest{
					ChainName:      "bsc",
					BridgerAddress: bridger.String(),
					Nonce:          3,
				}
				suite.Keeper().SetOracleByBridger(suite.ctx, oracle, bridger)
				suite.Keeper().SetOracleSetConfirm(suite.ctx, oracle, &types.MsgOracleSetConfirm{
					Nonce:          3,
					BridgerAddress: bridger.String(),
					ChainName:      "bsc",
				})
				response = &types.QueryOracleSetConfirmResponse{
					Confirm: &types.MsgOracleSetConfirm{
						Nonce:          3,
						BridgerAddress: bridger.String(),
						ChainName:      "bsc",
					},
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()
			res, err := suite.Keeper().OracleSetConfirm(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response, res)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_OracleSetConfirmsByNonce() {
	var (
		request       *types.QueryOracleSetConfirmsByNonceRequest
		response      *types.QueryOracleSetConfirmsByNonceResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "query nonce is zero",
			malleate: func() {
				request = &types.QueryOracleSetConfirmsByNonceRequest{
					ChainName: "bsc",
					Nonce:     0,
				}
				expectedError = sdkerrors.Wrapf(types.ErrUnknown, "nonce")
			},
			expPass: false,
		},
		{
			name: "query nonce does not exist",
			malleate: func() {
				request = &types.QueryOracleSetConfirmsByNonceRequest{
					ChainName: "bsc",
					Nonce:     5,
				}
				response = &types.QueryOracleSetConfirmsByNonceResponse{}
			},
			expPass: true,
		},
		{
			name: "query nonce normal",
			malleate: func() {
				bridger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				oracle := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())

				suite.Keeper().SetOracleByBridger(suite.ctx, oracle, bridger)
				suite.Keeper().SetOracleSetConfirm(suite.ctx, oracle, &types.MsgOracleSetConfirm{
					Nonce:          3,
					BridgerAddress: bridger.String(),
					ChainName:      "bsc",
				})
				request = &types.QueryOracleSetConfirmsByNonceRequest{
					ChainName: "bsc",
					Nonce:     3,
				}
				response = &types.QueryOracleSetConfirmsByNonceResponse{Confirms: []*types.MsgOracleSetConfirm{
					{
						Nonce:          3,
						BridgerAddress: bridger.String(),
						ChainName:      "bsc",
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
			res, err := suite.Keeper().OracleSetConfirmsByNonce(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response, res)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_LastPendingOracleSetRequestByAddr() {
	var (
		request       *types.QueryLastPendingOracleSetRequestByAddrRequest
		response      *types.QueryLastPendingOracleSetRequestByAddrResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "query oracla set address error",
			malleate: func() {
				request = &types.QueryLastPendingOracleSetRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: "fx1",
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "bridger address")
			},
			expPass: false,
		},
		{
			name: "not found oracle address by bridger",
			malleate: func() {
				request = &types.QueryLastPendingOracleSetRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
				}
				expectedError = types.ErrNoFoundOracle
			},
			expPass: false,
		},
		{
			name: "not found oracle by oracle address",
			malleate: func() {
				bridger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				oracle := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				suite.Keeper().SetOracleByBridger(suite.ctx, oracle, bridger)

				request = &types.QueryLastPendingOracleSetRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: bridger.String(),
				}
				expectedError = types.ErrNoFoundOracle
			},
			expPass: false,
		},
		{
			name: "not found oracle by oracle address",
			malleate: func() {
				bridger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				oracle := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				key, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(key.PubKey().Address().Bytes())
				suite.ctx = suite.ctx.WithBlockHeight(100)

				oracleSet := &types.OracleSet{
					Nonce: 3,
					Members: []types.BridgeValidator{
						{
							Power:           10000,
							ExternalAddress: externalAcc.String(),
						},
					},
					Height: 100,
				}

				suite.Keeper().SetOracleByBridger(suite.ctx, oracle, bridger)
				suite.Keeper().SetOracle(suite.ctx, types.Oracle{
					OracleAddress:   oracle.String(),
					BridgerAddress:  bridger.String(),
					ExternalAddress: externalAcc.String(),
					StartHeight:     0,
				})
				suite.Keeper().StoreOracleSet(suite.ctx, oracleSet)
				request = &types.QueryLastPendingOracleSetRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: bridger.String(),
				}

				response = &types.QueryLastPendingOracleSetRequestByAddrResponse{
					OracleSets: []*types.OracleSet{oracleSet},
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
			res, err := suite.Keeper().LastPendingOracleSetRequestByAddr(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response, res)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_LastPendingBatchRequestByAddr() {
	var (
		request       *types.QueryLastPendingBatchRequestByAddrRequest
		response      *types.QueryLastPendingBatchRequestByAddrResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "bridger address error",
			malleate: func() {
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: "fx1",
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "bridger address")
			},
			expPass: false,
		},
		{
			name: "not found oracle by bridger",
			malleate: func() {
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
				}
				expectedError = types.ErrNoFoundOracle
			},
			expPass: false,
		},
		{
			name: "not found oracle",
			malleate: func() {
				suite.Keeper().SetOracleByBridger(suite.ctx, suite.oracles[0], suite.bridgers[0])
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: suite.bridgers[0].String(),
				}
				expectedError = types.ErrNoFoundOracle
			},
			expPass: false,
		},
		{
			name: "normal test",
			malleate: func() {
				externalKey, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(externalKey.PubKey().Address().Bytes())
				externalToken := crypto.CreateAddress(common.BytesToAddress(externalKey.PubKey().Address().Bytes()), 0)

				suite.Keeper().SetOracleByBridger(suite.ctx, suite.oracles[0], suite.bridgers[0])
				suite.Keeper().SetOracle(suite.ctx, types.Oracle{
					OracleAddress:   suite.oracles[0].String(),
					BridgerAddress:  suite.bridgers[0].String(),
					ExternalAddress: externalAcc.String(),
					StartHeight:     10,
				})
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: suite.bridgers[0].String(),
				}
				suite.ctx = suite.ctx.WithBlockHeight(100)
				err = suite.Keeper().StoreBatch(suite.ctx, &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Id:          0,
							Sender:      sdk.AccAddress(externalKey.PubKey().Bytes()).String(),
							DestAddress: externalAcc.String(),
							Token:       types.NewERC20Token(sdk.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
							Fee:         types.NewERC20Token(sdk.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
						},
					},
					TokenContract: externalToken.String(),
					FeeReceive:    externalAcc.String(),
				})
				suite.Require().NoError(err)
				response = &types.QueryLastPendingBatchRequestByAddrResponse{Batch: &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Id:          0,
							Sender:      sdk.AccAddress(externalKey.PubKey().Bytes()).String(),
							DestAddress: externalAcc.String(),
							Token:       types.NewERC20Token(sdk.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
							Fee:         types.NewERC20Token(sdk.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
						},
					},
					TokenContract: externalToken.String(),
					Block:         100,
					FeeReceive:    externalAcc.String(),
				}}
			},
			expPass: true,
		},
		{
			name: "test batch confirm tx",
			malleate: func() {
				externalKey, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(externalKey.PubKey().Address().Bytes())
				externalToken := crypto.CreateAddress(common.BytesToAddress(externalKey.PubKey().Address().Bytes()), 0)
				suite.Keeper().SetOracleByBridger(suite.ctx, suite.oracles[0], suite.bridgers[0])
				suite.Keeper().SetOracle(suite.ctx, types.Oracle{
					OracleAddress:   suite.oracles[0].String(),
					BridgerAddress:  suite.bridgers[0].String(),
					ExternalAddress: externalAcc.String(),
					StartHeight:     10,
				})
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: suite.bridgers[0].String(),
				}
				suite.ctx = suite.ctx.WithBlockHeight(100)
				err = suite.Keeper().StoreBatch(suite.ctx, &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Id:          0,
							Sender:      sdk.AccAddress(externalKey.PubKey().Bytes()).String(),
							DestAddress: externalAcc.String(),
							Token:       types.NewERC20Token(sdk.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
							Fee:         types.NewERC20Token(sdk.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
						},
					},
					TokenContract: externalToken.String(),
					FeeReceive:    externalAcc.String(),
				})
				suite.Require().NoError(err)
				suite.Keeper().SetBatchConfirm(suite.ctx, suite.oracles[0], &types.MsgConfirmBatch{
					Nonce:           3,
					TokenContract:   externalToken.String(),
					BridgerAddress:  suite.bridgers[0].String(),
					ExternalAddress: externalAcc.String(),
					Signature:       "0x1",
					ChainName:       "bsc",
				})
				response = &types.QueryLastPendingBatchRequestByAddrResponse{}
			},
			expPass: true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()
			res, err := suite.Keeper().LastPendingBatchRequestByAddr(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response.Batch, res.Batch)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_BatchRequestByNonce() {
	var (
		request       *types.QueryBatchRequestByNonceRequest
		response      *types.QueryBatchRequestByNonceResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "query token contract error",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "bsc",
					TokenContract: "0x1",
					Nonce:         3,
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "token contract address")
			},
			expPass: false,
		},
		{
			name: "query token contract error",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "bsc",
					TokenContract: crypto.CreateAddress(common.BytesToAddress(key.PubKey().Bytes()), 0).String(),
					Nonce:         0,
				}
				expectedError = sdkerrors.Wrap(types.ErrUnknown, "nonce")
			},
			expPass: false,
		},
		{
			name: "query does not exist tx batch",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "bsc",
					TokenContract: crypto.CreateAddress(common.BytesToAddress(key.PubKey().Bytes()), 0).String(),
					Nonce:         3,
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "can not find tx batch")
			},
			expPass: false,
		},
		{
			name: "query tx batch normal",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				token := crypto.CreateAddress(common.BytesToAddress(key.PubKey().Bytes()), 0)

				newBatch := &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Id:    0,
							Token: types.NewERC20Token(sdk.NewIntFromBigInt(big.NewInt(1e18)), token.String()),
							Fee:   types.NewERC20Token(sdk.NewIntFromBigInt(big.NewInt(1e18)), token.String()),
						},
					},
					TokenContract: token.String(),
					Block:         100,
				}
				err := suite.Keeper().StoreBatch(suite.ctx, newBatch)
				suite.Require().NoError(err)
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     "bsc",
					TokenContract: token.String(),
					Nonce:         3,
				}
				response = &types.QueryBatchRequestByNonceResponse{Batch: newBatch}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			testCase.malleate()
			res, err := suite.Keeper().BatchRequestByNonce(sdk.WrapSDKContext(suite.ctx), request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response, res)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_BatchConfirm() {
	var (
		request       *types.QueryBatchConfirmRequest
		response      *types.QueryBatchConfirmResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "bridger address error",
			malleate: func() {
				request = &types.QueryBatchConfirmRequest{
					ChainName:      "bsc",
					BridgerAddress: "fx1",
					Nonce:          3,
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "bridger address")
			},
			expPass: false,
		},
		{
			name: "query nonce error",
			malleate: func() {
				request = &types.QueryBatchConfirmRequest{
					ChainName:      "bsc",
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
					Nonce:          0,
				}
				expectedError = sdkerrors.Wrap(types.ErrUnknown, "nonce")
			},
			expPass: false,
		},
		{
			name: "query oracle not found",
			malleate: func() {
				request = &types.QueryBatchConfirmRequest{
					ChainName:      "bsc",
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
					Nonce:          3,
				}
				expectedError = types.ErrNoFoundOracle
			},
			expPass: false,
		},
		{
			name: "query batch confirm normal",
			malleate: func() {
				bridger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				oracle := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				suite.Keeper().SetOracleByBridger(suite.ctx, oracle, bridger)

				suite.Keeper().SetBatchConfirm(suite.ctx, oracle, &types.MsgConfirmBatch{
					Nonce:          3,
					BridgerAddress: bridger.String(),
					ChainName:      "bsc",
				})
				request = &types.QueryBatchConfirmRequest{
					ChainName:      "bsc",
					BridgerAddress: bridger.String(),
					Nonce:          3,
				}
				response = &types.QueryBatchConfirmResponse{Confirm: &types.MsgConfirmBatch{
					Nonce:          3,
					BridgerAddress: bridger.String(),
					ChainName:      "bsc",
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
			res, err := suite.Keeper().BatchConfirm(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response, res)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_BatchConfirms() {
	var (
		request       *types.QueryBatchConfirmsRequest
		response      *types.QueryBatchConfirmsResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "query token address error",
			malleate: func() {
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     "bsc",
					TokenContract: "0x11",
					Nonce:         3,
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "token contract address")
			},
			expPass: false,
		},
		{
			name: "query nonce error",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				token := crypto.CreateAddress(common.BytesToAddress(key.PubKey().Bytes()), 0)

				request = &types.QueryBatchConfirmsRequest{
					ChainName:     "bsc",
					TokenContract: token.String(),
					Nonce:         0,
				}
				expectedError = sdkerrors.Wrap(types.ErrUnknown, "nonce")
			},
			expPass: false,
		},
		{
			name: "batch confirms normal",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				token := crypto.CreateAddress(common.BytesToAddress(key.PubKey().Bytes()), 0)
				confirms := make([]*types.MsgConfirmBatch, 0)

				for i := 0; i < 3; i++ {
					oracle := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
					bridger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
					newMsg := &types.MsgConfirmBatch{
						Nonce:          3,
						TokenContract:  token.String(),
						BridgerAddress: bridger.String(),
						ChainName:      "bsc",
					}
					suite.Keeper().SetBatchConfirm(suite.ctx, oracle, newMsg)
					confirms = append(confirms, newMsg)
				}

				request = &types.QueryBatchConfirmsRequest{
					ChainName:     "bsc",
					TokenContract: token.String(),
					Nonce:         3,
				}
				response = &types.QueryBatchConfirmsResponse{Confirms: confirms}
			},
			expPass: true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()
			res, err := suite.Keeper().BatchConfirms(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().ElementsMatch(response.Confirms, res.Confirms)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_LastEventNonceByAddr() {
	var (
		request       *types.QueryLastEventNonceByAddrRequest
		response      *types.QueryLastEventNonceByAddrResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "query bridger address ",
			malleate: func() {
				request = &types.QueryLastEventNonceByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: "fx1",
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "bridger address")
			},
			expPass: false,
		},
		{
			name: "query not found oracle by bridger",
			malleate: func() {
				request = &types.QueryLastEventNonceByAddrRequest{
					ChainName:      "bsc",
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
				}
				expectedError = types.ErrNoFoundOracle
			},
			expPass: false,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()
			res, err := suite.Keeper().LastEventNonceByAddr(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response, res)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_GetOracleByAddr() {
	var (
		request       *types.QueryOracleByAddrRequest
		response      *types.QueryOracleResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "query oracle address error",
			malleate: func() {
				request = &types.QueryOracleByAddrRequest{
					ChainName:     "bsc",
					OracleAddress: "fx1",
				}
				expectedError = sdkerrors.Wrap(types.ErrInvalid, "oracle address")
			},
			expPass: false,
		},
		{
			name: "query oracle does not exist",
			malleate: func() {
				request = &types.QueryOracleByAddrRequest{
					ChainName:     "bsc",
					OracleAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
				}
				expectedError = types.ErrNoFoundOracle
			},
			expPass: false,
		},
		{
			name: "query oracle normal",
			malleate: func() {
				bridger := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				oracle := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes())
				key, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(key.PubKey().Address().Bytes())
				suite.ctx = suite.ctx.WithBlockHeight(100)
				newOralce := types.Oracle{
					OracleAddress:   oracle.String(),
					BridgerAddress:  bridger.String(),
					ExternalAddress: externalAcc.String(),
					DelegateAmount:  sdk.NewIntFromBigInt(big.NewInt(10000)),
					StartHeight:     0,
				}
				suite.Keeper().SetOracle(suite.ctx, newOralce)

				request = &types.QueryOracleByAddrRequest{
					ChainName:     "bsc",
					OracleAddress: oracle.String(),
				}
				response = &types.QueryOracleResponse{Oracle: &newOralce}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			ctx := sdk.WrapSDKContext(suite.ctx)
			testCase.malleate()
			res, err := suite.Keeper().GetOracleByAddr(ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response, res)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}
