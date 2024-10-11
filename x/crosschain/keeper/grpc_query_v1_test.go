package keeper_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	bsctypes "github.com/functionx/fx-core/v8/x/bsc/types"
	"github.com/functionx/fx-core/v8/x/crosschain/keeper"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
)

type CrossChainGrpcTestSuite struct {
	KeeperTestSuite
	chainName    string
	oracleAddrs  []sdk.AccAddress
	bridgerAddrs []sdk.AccAddress

	msgServer   types.MsgServer
	queryClient types.QueryClient
}

// func TestCrossChainGrpcTestSuite_bsc(t *testing.T) {
//	suite.Run(t, &CrossChainGrpcTestSuite{KeeperTestSuite: KeeperTestSuite{chainName: bsctypes.ModuleName}, chainName: bsctypes.ModuleName})
// }

func TestCrossChainGrpcTestSuite_eth(t *testing.T) {
	suite.Run(t, &CrossChainGrpcTestSuite{KeeperTestSuite: KeeperTestSuite{chainName: ethtypes.ModuleName}, chainName: ethtypes.ModuleName})
}

func (suite *CrossChainGrpcTestSuite) SetupTest() {
	suite.KeeperTestSuite.SetupTest()
	suite.Commit(10)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.App.CrosschainRouterKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	suite.oracleAddrs = suite.AddTestAddress(types.MaxOracleSize, types.NewDelegateAmount(sdkmath.NewInt(300*1e3).MulRaw(1e18)))
	suite.bridgerAddrs = suite.AddTestAddress(types.MaxOracleSize, sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(300*1e3).MulRaw(1e18)))
	suite.msgServer = keeper.NewMsgServerImpl(suite.Keeper())
}

func (suite *CrossChainGrpcTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *CrossChainGrpcTestSuite) Keeper() keeper.Keeper {
	switch suite.chainName {
	case bsctypes.ModuleName:
		return suite.App.BscKeeper
	case ethtypes.ModuleName:
		return suite.App.EthKeeper
	default:
		panic("invalid chain name")
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_CurrentOracleSet() {
	testCases := []struct {
		name          string
		malleate      func() *types.QueryCurrentOracleSetResponse
		expectedError error
		expPass       bool
	}{
		{
			name: "no oracle set",
			malleate: func() *types.QueryCurrentOracleSetResponse {
				return &types.QueryCurrentOracleSetResponse{OracleSet: types.NewOracleSet(1, 11, nil)}
			},
			expPass: true,
		},
		{
			name: "query current oracle set",
			malleate: func() *types.QueryCurrentOracleSetResponse {
				newOracleSet := &types.OracleSet{
					Members: make([]types.BridgeValidator, 0),
				}
				for i := 0; i < 6; i++ {
					key, _ := ethsecp256k1.GenerateKey()
					externalAcc := common.BytesToAddress(key.PubKey().Address())
					delegateAmount := sdk.DefaultPowerReduction.Mul(sdkmath.NewInt(100))
					if i == 5 {
						delegateAmount = sdkmath.ZeroInt()
					}
					suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
						OracleAddress:   suite.oracleAddrs[i].String(),
						BridgerAddress:  suite.bridgerAddrs[i].String(),
						ExternalAddress: externalAcc.String(),
						DelegateAmount:  delegateAmount,
						Online:          true,
						StartHeight:     int64(10 + i),
					})
					if i != 5 {
						newOracleSet.Members = append(newOracleSet.Members, types.BridgeValidator{
							Power:           858993459,
							ExternalAddress: externalAcc.String(),
						})
					}
				}
				suite.Ctx = suite.Ctx.WithBlockHeight(100)
				newOracleSet.Height = 11
				suite.Keeper().SetLatestOracleSetNonce(suite.Ctx, 10)
				newOracleSet.Nonce = 11
				return &types.QueryCurrentOracleSetResponse{OracleSet: newOracleSet}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			response := testCase.malleate()
			res, err := suite.queryClient.CurrentOracleSet(
				suite.Ctx,
				&types.QueryCurrentOracleSetRequest{ChainName: suite.chainName},
			)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().ElementsMatch(response.OracleSet.Members, res.OracleSet.Members)
				suite.Require().Equal(response.OracleSet.Nonce, res.OracleSet.Nonce)
				suite.Require().Equal(response.OracleSet.Height, res.OracleSet.Height)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, testCase.expectedError)
			}
		})
	}
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
					ChainName: suite.chainName,
					Nonce:     1,
				}
				response = &types.QueryCurrentOracleSetResponse{OracleSet: nil}
			},
			expPass: true,
		},
		{
			name: "normal oracle set",
			malleate: func() {
				members := []types.BridgeValidator{
					{
						Power:           10000,
						ExternalAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
					},
				}
				request = &types.QueryOracleSetRequestRequest{
					ChainName: suite.chainName,
					Nonce:     3,
				}
				suite.Keeper().StoreOracleSet(suite.Ctx, &types.OracleSet{
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
			testCase.malleate()
			res, err := suite.queryClient.OracleSetRequest(suite.Ctx, request)
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
					ChainName:      suite.chainName,
					BridgerAddress: "fx1",
					Nonce:          tmrand.Uint64() + 1,
				}
				expectedError = status.Error(codes.InvalidArgument, "bridger address")
			},
			expPass: false,
		},
		{
			name: "oracle set nonce error",
			malleate: func() {
				request = &types.QueryOracleSetConfirmRequest{
					ChainName:      suite.chainName,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
					Nonce:          0,
				}
				expectedError = status.Error(codes.InvalidArgument, "nonce")
			},
			expPass: false,
		},
		{
			name: "oracle set bridger address does not exist",
			malleate: func() {
				request = &types.QueryOracleSetConfirmRequest{
					ChainName:      suite.chainName,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
					Nonce:          3,
				}
				expectedError = status.Error(codes.NotFound, "oracle not found by bridger address")
			},
			expPass: false,
		},
		{
			"oracle set normal",
			func() {
				request = &types.QueryOracleSetConfirmRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
					Nonce:          3,
				}
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				suite.Keeper().SetOracleSetConfirm(suite.Ctx, suite.oracleAddrs[0], &types.MsgOracleSetConfirm{
					Nonce:          3,
					BridgerAddress: suite.bridgerAddrs[0].String(),
					ChainName:      suite.chainName,
				})
				response = &types.QueryOracleSetConfirmResponse{
					Confirm: &types.MsgOracleSetConfirm{
						Nonce:          3,
						BridgerAddress: suite.bridgerAddrs[0].String(),
						ChainName:      suite.chainName,
					},
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.OracleSetConfirm(suite.Ctx, request)
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
					ChainName: suite.chainName,
					Nonce:     0,
				}
				expectedError = status.Error(codes.InvalidArgument, "nonce")
			},
			expPass: false,
		},
		{
			name: "query nonce does not exist",
			malleate: func() {
				request = &types.QueryOracleSetConfirmsByNonceRequest{
					ChainName: suite.chainName,
					Nonce:     5,
				}
				response = &types.QueryOracleSetConfirmsByNonceResponse{}
			},
			expPass: true,
		},
		{
			name: "query nonce normal",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				suite.Keeper().SetOracleSetConfirm(suite.Ctx, suite.oracleAddrs[0], &types.MsgOracleSetConfirm{
					Nonce:          3,
					BridgerAddress: suite.bridgerAddrs[0].String(),
					ChainName:      suite.chainName,
				})
				request = &types.QueryOracleSetConfirmsByNonceRequest{
					ChainName: suite.chainName,
					Nonce:     3,
				}
				response = &types.QueryOracleSetConfirmsByNonceResponse{Confirms: []*types.MsgOracleSetConfirm{
					{
						Nonce:          3,
						BridgerAddress: suite.bridgerAddrs[0].String(),
						ChainName:      suite.chainName,
					},
				}}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.OracleSetConfirmsByNonce(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_LastOracleSetRequest() {
	testCases := []struct {
		name          string
		malleate      func() *types.QueryLastOracleSetRequestsResponse
		expectedError error
		expPass       bool
	}{
		{
			name: "query params",
			malleate: func() *types.QueryLastOracleSetRequestsResponse {
				oracleSetList := make([]*types.OracleSet, 0)
				for i := 0; i < 10; i++ {
					key, _ := ethsecp256k1.GenerateKey()
					newOracleSet := &types.OracleSet{
						Nonce: uint64(i),
						Members: []types.BridgeValidator{
							{
								Power:           100000,
								ExternalAddress: common.BytesToAddress(key.PubKey().Address().Bytes()).String(),
							},
						},
						Height: uint64((i + 1) * 33),
					}
					suite.Keeper().StoreOracleSet(suite.Ctx, newOracleSet)
					oracleSetList = append(oracleSetList, newOracleSet)
				}
				return &types.QueryLastOracleSetRequestsResponse{
					OracleSets: oracleSetList[len(oracleSetList)-5:],
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			response := testCase.malleate()
			res, err := suite.queryClient.LastOracleSetRequests(
				suite.Ctx,
				&types.QueryLastOracleSetRequestsRequest{ChainName: suite.chainName},
			)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().ElementsMatch(response.OracleSets, res.OracleSets)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, testCase.expectedError)
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
			name: "query oracle set address error",
			malleate: func() {
				request = &types.QueryLastPendingOracleSetRequestByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: "fx1",
				}
				expectedError = status.Error(codes.InvalidArgument, "bridger address")
			},
			expPass: false,
		},
		{
			name: "not found oracle address by bridger",
			malleate: func() {
				request = &types.QueryLastPendingOracleSetRequestByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found by bridger address")
			},
			expPass: false,
		},
		{
			name: "not found oracle by oracle address",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				request = &types.QueryLastPendingOracleSetRequestByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found")
			},
			expPass: false,
		},
		{
			name: "ok",
			malleate: func() {
				key, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(key.PubKey().Address().Bytes())
				suite.Ctx = suite.Ctx.WithBlockHeight(100)

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

				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])

				suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
					OracleAddress:   suite.oracleAddrs[0].String(),
					BridgerAddress:  suite.bridgerAddrs[0].String(),
					ExternalAddress: externalAcc.String(),
					StartHeight:     0,
				})
				suite.Keeper().StoreOracleSet(suite.Ctx, oracleSet)
				request = &types.QueryLastPendingOracleSetRequestByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
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
			testCase.malleate()
			res, err := suite.queryClient.LastPendingOracleSetRequestByAddr(suite.Ctx, request)
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
					ChainName:      suite.chainName,
					BridgerAddress: "fx1",
				}
				expectedError = status.Error(codes.InvalidArgument, "bridger address")
			},
			expPass: false,
		},
		{
			name: "not found oracle by bridger",
			malleate: func() {
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found by bridger address")
			},
			expPass: false,
		},
		{
			name: "not found oracle",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found")
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

				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
					OracleAddress:   suite.oracleAddrs[0].String(),
					BridgerAddress:  suite.bridgerAddrs[0].String(),
					ExternalAddress: externalAcc.String(),
					StartHeight:     10,
				})
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				suite.Ctx = suite.Ctx.WithBlockHeight(100)
				err = suite.Keeper().StoreBatch(suite.Ctx, &types.OutgoingTxBatch{
					Block:        uint64(suite.Ctx.BlockHeight()),
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Id:          0,
							Sender:      sdk.AccAddress(externalKey.PubKey().Address()).String(),
							DestAddress: externalAcc.String(),
							Token:       types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
							Fee:         types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
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
							Sender:      sdk.AccAddress(externalKey.PubKey().Address()).String(),
							DestAddress: externalAcc.String(),
							Token:       types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
							Fee:         types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
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
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
					OracleAddress:   suite.oracleAddrs[0].String(),
					BridgerAddress:  suite.bridgerAddrs[0].String(),
					ExternalAddress: externalAcc.String(),
					StartHeight:     10,
				})
				request = &types.QueryLastPendingBatchRequestByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				suite.Ctx = suite.Ctx.WithBlockHeight(100)
				err = suite.Keeper().StoreBatch(suite.Ctx, &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Id:          0,
							Sender:      sdk.AccAddress(externalKey.PubKey().Address()).String(),
							DestAddress: externalAcc.String(),
							Token:       types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
							Fee:         types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), externalToken.String()),
						},
					},
					TokenContract: externalToken.String(),
					FeeReceive:    externalAcc.String(),
				})
				suite.Require().NoError(err)
				suite.Keeper().SetBatchConfirm(suite.Ctx, suite.oracleAddrs[0], &types.MsgConfirmBatch{
					Nonce:           3,
					TokenContract:   externalToken.String(),
					BridgerAddress:  suite.bridgerAddrs[0].String(),
					ExternalAddress: externalAcc.String(),
					Signature:       "0x1",
					ChainName:       suite.chainName,
				})
				response = &types.QueryLastPendingBatchRequestByAddrResponse{}
			},
			expPass: true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.LastPendingBatchRequestByAddr(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_OutgoingTxBatches() {
	testCases := []struct {
		name          string
		malleate      func() *types.QueryOutgoingTxBatchesResponse
		expectedError error
		expPass       bool
	}{
		{
			name: "query outgoing tx batches",
			malleate: func() *types.QueryOutgoingTxBatchesResponse {
				newBatchList := make([]*types.OutgoingTxBatch, 0)
				for i := 0; i < 10; i++ {
					suite.Ctx = suite.Ctx.WithBlockHeight(int64(i + 3))
					token := helpers.GenHexAddress().String()
					newOutgoingTx := &types.OutgoingTxBatch{
						BatchNonce:   uint64(i + 3),
						BatchTimeout: uint64(1000),
						Transactions: []*types.OutgoingTransferTx{
							{
								Id:          uint64(i),
								Sender:      sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
								DestAddress: helpers.GenHexAddress().String(),
								Token:       types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), token),
								Fee:         types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), token),
							},
						},
						TokenContract: token,
						Block:         uint64(i + 3),
						FeeReceive:    helpers.GenHexAddress().String(),
					}
					err := suite.Keeper().StoreBatch(suite.Ctx, newOutgoingTx)
					suite.Require().NoError(err)
					newBatchList = append(newBatchList, newOutgoingTx)
				}
				return &types.QueryOutgoingTxBatchesResponse{Batches: newBatchList}
			},
			expPass: true,
		},
		{
			name: "query outgoing tx batches more than 100",
			malleate: func() *types.QueryOutgoingTxBatchesResponse {
				for i := 1; i < 110; i++ {
					suite.Ctx = suite.Ctx.WithBlockHeight(int64(i))
					token := helpers.GenHexAddress().String()
					newOutgoingTx := &types.OutgoingTxBatch{
						BatchNonce:   uint64(i),
						BatchTimeout: uint64(1000 + i),
						Transactions: []*types.OutgoingTransferTx{
							{
								Id:          uint64(i),
								Sender:      sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
								DestAddress: helpers.GenHexAddress().String(),
								Token:       types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), token),
								Fee:         types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), token),
							},
						},
						TokenContract: token,
						Block:         uint64(i),
						FeeReceive:    helpers.GenHexAddress().String(),
					}
					err := suite.Keeper().StoreBatch(suite.Ctx, newOutgoingTx)
					suite.Require().NoError(err)
				}
				return &types.QueryOutgoingTxBatchesResponse{}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			response := testCase.malleate()
			res, err := suite.queryClient.OutgoingTxBatches(suite.Ctx, &types.QueryOutgoingTxBatchesRequest{ChainName: suite.chainName})
			suite.Require().NoError(err)
			if testCase.expPass {
				suite.Require().True(len(res.Batches) <= 100)
				if len(res.Batches) < 100 {
					suite.Require().ElementsMatch(response.Batches, res.Batches)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, testCase.expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_OutgoingTxBatch() {
	var (
		request       *types.QueryOutgoingTxBatchRequest
		response      *types.QueryOutgoingTxBatchResponse
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
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     suite.chainName,
					TokenContract: "0x1",
					Nonce:         3,
				}
				expectedError = status.Error(codes.InvalidArgument, "token contract address")
			},
			expPass: false,
		},
		{
			name: "query token contract error",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     suite.chainName,
					TokenContract: crypto.CreateAddress(common.BytesToAddress(key.PubKey().Address().Bytes()), 0).String(),
					Nonce:         0,
				}
				expectedError = status.Error(codes.InvalidArgument, "nonce")
			},
			expPass: false,
		},
		{
			name: "query does not exist tx batch",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     suite.chainName,
					TokenContract: crypto.CreateAddress(common.BytesToAddress(key.PubKey().Address().Bytes()), 0).String(),
					Nonce:         3,
				}
				expectedError = status.Error(codes.NotFound, "tx batch")
			},
			expPass: false,
		},
		{
			name: "query tx batch normal",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				token := crypto.CreateAddress(common.BytesToAddress(key.PubKey().Address().Bytes()), 0)

				newBatch := &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Id:    0,
							Token: types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), token.String()),
							Fee:   types.NewERC20Token(sdkmath.NewIntFromBigInt(big.NewInt(1e18)), token.String()),
						},
					},
					TokenContract: token.String(),
					Block:         100,
				}
				err := suite.Keeper().StoreBatch(suite.Ctx, newBatch)
				suite.Require().NoError(err)
				request = &types.QueryOutgoingTxBatchRequest{
					ChainName:     suite.chainName,
					TokenContract: token.String(),
					Nonce:         3,
				}
				response = &types.QueryOutgoingTxBatchResponse{Batch: newBatch}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.OutgoingTxBatch(suite.Ctx, request)
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
					ChainName:      suite.chainName,
					BridgerAddress: "fx1",
					Nonce:          3,
				}
				expectedError = status.Error(codes.InvalidArgument, "bridger address")
			},
			expPass: false,
		},
		{
			name: "query nonce error",
			malleate: func() {
				request = &types.QueryBatchConfirmRequest{
					ChainName:      suite.chainName,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
					Nonce:          0,
				}
				expectedError = status.Error(codes.InvalidArgument, "nonce")
			},
			expPass: false,
		},
		{
			name: "query oracle not found",
			malleate: func() {
				request = &types.QueryBatchConfirmRequest{
					ChainName:      suite.chainName,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
					Nonce:          3,
				}
				expectedError = status.Error(codes.NotFound, "oracle not found by bridger address")
			},
			expPass: false,
		},
		{
			name: "query batch confirm normal",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])

				suite.Keeper().SetBatchConfirm(suite.Ctx, suite.oracleAddrs[0], &types.MsgConfirmBatch{
					Nonce:          3,
					BridgerAddress: suite.bridgerAddrs[0].String(),
					ChainName:      suite.chainName,
				})
				request = &types.QueryBatchConfirmRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
					Nonce:          3,
				}
				response = &types.QueryBatchConfirmResponse{Confirm: &types.MsgConfirmBatch{
					Nonce:          3,
					BridgerAddress: suite.bridgerAddrs[0].String(),
					ChainName:      suite.chainName,
				}}
			},
			expPass: true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.BatchConfirm(suite.Ctx, request)
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
					ChainName:     suite.chainName,
					TokenContract: "0x11",
					Nonce:         3,
				}
				expectedError = status.Error(codes.InvalidArgument, "token contract address")
			},
			expPass: false,
		},
		{
			name: "query nonce error",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				token := crypto.CreateAddress(common.BytesToAddress(key.PubKey().Address()), 0)

				request = &types.QueryBatchConfirmsRequest{
					ChainName:     suite.chainName,
					TokenContract: token.String(),
					Nonce:         0,
				}
				expectedError = status.Error(codes.InvalidArgument, "nonce")
			},
			expPass: false,
		},
		{
			name: "batch confirms normal",
			malleate: func() {
				key, _ := ethsecp256k1.GenerateKey()
				token := crypto.CreateAddress(common.BytesToAddress(key.PubKey().Address()), 0)
				confirms := make([]*types.MsgConfirmBatch, 0)

				for i := 0; i < 3; i++ {
					newMsg := &types.MsgConfirmBatch{
						Nonce:          3,
						TokenContract:  token.String(),
						BridgerAddress: suite.bridgerAddrs[i].String(),
						ChainName:      suite.chainName,
					}
					suite.Keeper().SetBatchConfirm(suite.Ctx, suite.oracleAddrs[i], newMsg)
					confirms = append(confirms, newMsg)
				}

				request = &types.QueryBatchConfirmsRequest{
					ChainName:     suite.chainName,
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
			testCase.malleate()
			res, err := suite.queryClient.BatchConfirms(suite.Ctx, request)
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
					ChainName:      suite.chainName,
					BridgerAddress: "fx1",
				}
				expectedError = status.Error(codes.InvalidArgument, "bridger address")
			},
			expPass: false,
		},
		{
			name: "query not found oracle by bridger",
			malleate: func() {
				request = &types.QueryLastEventNonceByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found by bridger address")
			},
			expPass: false,
		},
		{
			name: "query last event nonce from lastObservedEventNonce",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				suite.Keeper().SetLastObservedEventNonce(suite.Ctx, 5)

				request = &types.QueryLastEventNonceByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				response = &types.QueryLastEventNonceByAddrResponse{EventNonce: 4}
			},
			expPass: true,
		},
		{
			name: "query last event nonce not found",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				request = &types.QueryLastEventNonceByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				response = &types.QueryLastEventNonceByAddrResponse{EventNonce: 0}
			},
			expPass: true,
		},
		{
			name: "query last event nonce normal",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				suite.Keeper().SetLastEventNonceByOracle(suite.Ctx, suite.oracleAddrs[0], 3)

				request = &types.QueryLastEventNonceByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				response = &types.QueryLastEventNonceByAddrResponse{EventNonce: 3}
			},
			expPass: true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.LastEventNonceByAddr(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_DenomToToken() {
	var (
		request       *types.QueryDenomToTokenRequest
		response      *types.QueryDenomToTokenResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"denom is nil",
			func() {
				request = &types.QueryDenomToTokenRequest{
					ChainName: suite.chainName,
				}
				expectedError = status.Error(codes.InvalidArgument, "denom")
			},
			false,
		},
		{
			"bridge token not exist",
			func() {
				request = &types.QueryDenomToTokenRequest{
					ChainName: suite.chainName,
					Denom:     "bsc0xfbbbb4f7b1e5bcb0345c5a5a61584b2547d5d582",
				}
				expectedError = status.Error(codes.NotFound, "bridge token")
			},
			false,
		},
		{
			"bridge token and ChannelIbc is exist and true",
			func() {
				key, _ := ethsecp256k1.GenerateKey()
				token := common.BytesToAddress(key.PubKey().Address()).String()

				err := suite.Keeper().AttestationHandler(suite.Ctx, &types.MsgBridgeTokenClaim{
					TokenContract: token,
					ChannelIbc:    hex.EncodeToString([]byte("transfer/channel-0")),
					Symbol:        "fxcoin",
				})
				suite.Require().NoError(err)
				bridgeDenom, found := suite.Keeper().GetBridgeDenomByContract(suite.Ctx, token)
				suite.Require().True(found)
				request = &types.QueryDenomToTokenRequest{
					ChainName: suite.chainName,
					Denom:     bridgeDenom,
				}
				response = &types.QueryDenomToTokenResponse{
					Token: token,
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.DenomToToken(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_TokenToDenom() {
	var (
		request       *types.QueryTokenToDenomRequest
		response      *types.QueryTokenToDenomResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"token address is error or null",
			func() {
				request = &types.QueryTokenToDenomRequest{
					ChainName: suite.chainName,
				}
				expectedError = status.Error(codes.InvalidArgument, "token address")
			},
			false,
		},
		{
			"bridge token is not exist",
			func() {
				key, _ := ethsecp256k1.GenerateKey()
				request = &types.QueryTokenToDenomRequest{
					ChainName: suite.chainName,
					Token:     common.BytesToAddress(key.PubKey().Address()).String(),
				}
				expectedError = status.Error(codes.NotFound, "bridge token")
			},
			false,
		},
		{
			"token normal",
			func() {
				key, _ := ethsecp256k1.GenerateKey()
				token := common.BytesToAddress(key.PubKey().Address()).String()
				err := suite.Keeper().AttestationHandler(suite.Ctx, &types.MsgBridgeTokenClaim{
					TokenContract: token,
					ChannelIbc:    hex.EncodeToString([]byte("transfer/channel-0")),
					Symbol:        "fxcoin",
				})
				suite.Require().NoError(err)
				request = &types.QueryTokenToDenomRequest{
					ChainName: suite.chainName,
					Token:     token,
				}
				bridgeDenom, found := suite.Keeper().GetBridgeDenomByContract(suite.Ctx, token)
				suite.Require().True(found)
				response = &types.QueryTokenToDenomResponse{
					Denom: bridgeDenom,
				}
				expectedError = types.ErrInvalid.Wrapf("bridge token is not exist")
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.TokenToDenom(suite.Ctx, request)
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
					ChainName:     suite.chainName,
					OracleAddress: "fx1",
				}
				expectedError = status.Error(codes.InvalidArgument, "oracle address")
			},
			expPass: false,
		},
		{
			name: "query oracle does not exist",
			malleate: func() {
				request = &types.QueryOracleByAddrRequest{
					ChainName:     suite.chainName,
					OracleAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found")
			},
			expPass: false,
		},
		{
			name: "query oracle normal",
			malleate: func() {
				key, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(key.PubKey().Address().Bytes())
				suite.Ctx = suite.Ctx.WithBlockHeight(100)
				newOracle := types.Oracle{
					OracleAddress:   suite.oracleAddrs[0].String(),
					BridgerAddress:  suite.bridgerAddrs[0].String(),
					ExternalAddress: externalAcc.String(),
					DelegateAmount:  sdkmath.NewIntFromBigInt(big.NewInt(10000)),
					StartHeight:     0,
				}
				suite.Keeper().SetOracle(suite.Ctx, newOracle)
				request = &types.QueryOracleByAddrRequest{
					ChainName:     suite.chainName,
					OracleAddress: suite.oracleAddrs[0].String(),
				}
				response = &types.QueryOracleResponse{Oracle: &newOracle}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.GetOracleByAddr(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_GetOracleByBridgerAddr() {
	var (
		request       *types.QueryOracleByBridgerAddrRequest
		response      *types.QueryOracleResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "query bridger address error",
			malleate: func() {
				request = &types.QueryOracleByBridgerAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: "fx1",
				}
				expectedError = status.Error(codes.InvalidArgument, "bridger address")
			},

			expPass: false,
		},
		{
			name: "query oracle by bridger address does not exist",
			malleate: func() {
				request = &types.QueryOracleByBridgerAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found by bridger address")
			},
			expPass: false,
		},
		{
			name: "query oracle by oracle address does not exist",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				request = &types.QueryOracleByBridgerAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found")
			},
			expPass: false,
		},
		{
			name: "query oracle by oracle address normal",
			malleate: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])
				key, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(key.PubKey().Address().Bytes())
				suite.Ctx = suite.Ctx.WithBlockHeight(100)

				newOracle := types.Oracle{
					OracleAddress:   suite.oracleAddrs[0].String(),
					BridgerAddress:  suite.bridgerAddrs[0].String(),
					ExternalAddress: externalAcc.String(),
					DelegateAmount:  sdkmath.NewIntFromBigInt(big.NewInt(10000)),
					StartHeight:     0,
				}

				suite.Keeper().SetOracle(suite.Ctx, newOracle)
				request = &types.QueryOracleByBridgerAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				response = &types.QueryOracleResponse{Oracle: &newOracle}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.GetOracleByBridgerAddr(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_GetOracleByExternalAddr() {
	var (
		request       *types.QueryOracleByExternalAddrRequest
		response      *types.QueryOracleResponse
		expectedError error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "query external address error",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       suite.chainName,
					ExternalAddress: "0x123",
				}
				expectedError = status.Error(codes.InvalidArgument, "external address")
			},
			expPass: false,
		},
		{
			name: "query oracle by external address does not exist",
			malleate: func() {
				key, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(key.PubKey().Address().Bytes())
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       suite.chainName,
					ExternalAddress: externalAcc.String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle")
			},
			expPass: false,
		},
		{
			name: "query oracle does not exist",
			malleate: func() {
				key, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(key.PubKey().Address().Bytes())
				suite.Keeper().SetOracleAddrByExternalAddr(suite.Ctx, externalAcc.String(), suite.oracleAddrs[0])
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       suite.chainName,
					ExternalAddress: externalAcc.String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle")
			},
			expPass: false,
		},
		{
			name: "query oracle normal",
			malleate: func() {
				key, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(key.PubKey().Address().Bytes())
				suite.Keeper().SetOracleAddrByExternalAddr(suite.Ctx, externalAcc.String(), suite.oracleAddrs[0])
				newOracle := types.Oracle{
					OracleAddress:   suite.oracleAddrs[0].String(),
					BridgerAddress:  suite.bridgerAddrs[0].String(),
					ExternalAddress: externalAcc.String(),
					DelegateAmount:  sdkmath.NewIntFromBigInt(big.NewInt(10000)),
					StartHeight:     0,
				}

				suite.Keeper().SetOracle(suite.Ctx, newOracle)
				suite.Ctx = suite.Ctx.WithBlockHeight(100)
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       suite.chainName,
					ExternalAddress: externalAcc.String(),
				}
				response = &types.QueryOracleResponse{Oracle: &newOracle}
			},
			expPass: true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.GetOracleByExternalAddr(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_LastEventBlockHeightByAddr() {
	var (
		request       *types.QueryLastEventBlockHeightByAddrRequest
		response      *types.QueryLastEventBlockHeightByAddrResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"BridgerAddress is error",
			func() {
				request = &types.QueryLastEventBlockHeightByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: "fx1",
				}
				expectedError = status.Error(codes.InvalidArgument, "bridger address")
			},
			false,
		},
		{
			"BridgerAddress exist oracle is nil",
			func() {
				request = &types.QueryLastEventBlockHeightByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				expectedError = status.Error(codes.NotFound, "oracle not found by bridger address")
			},
			false,
		},
		{
			"BridgerAddress exist oracle is not nil",
			func() {
				request = &types.QueryLastEventBlockHeightByAddrRequest{
					ChainName:      suite.chainName,
					BridgerAddress: suite.bridgerAddrs[0].String(),
				}
				suite.Ctx = suite.Ctx.WithBlockHeight(100)
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, suite.bridgerAddrs[0], suite.oracleAddrs[0])

				suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
					OracleAddress:  suite.oracleAddrs[0].String(),
					BridgerAddress: suite.bridgerAddrs[0].String(),
					StartHeight:    100,
					Online:         true,
				})
				claimMsg := &types.MsgBridgeTokenClaim{
					EventNonce:     1,
					BlockHeight:    100,
					TokenContract:  helpers.GenHexAddress().String(),
					Name:           "test token",
					Symbol:         "tt",
					Decimals:       18,
					BridgerAddress: suite.bridgerAddrs[0].String(),
					ChainName:      suite.chainName,
				}
				anyWithValue, err := codectypes.NewAnyWithValue(claimMsg)
				suite.Require().NoError(err)
				_, err = suite.msgServer.Claim(suite.Ctx, &types.MsgClaim{Claim: anyWithValue})
				suite.Require().NoError(err)
				response = &types.QueryLastEventBlockHeightByAddrResponse{
					BlockHeight: uint64(100),
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.LastEventBlockHeightByAddr(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_LastObservedBlockHeight() {
	var (
		request       *types.QueryLastObservedBlockHeightRequest
		response      *types.QueryLastObservedBlockHeightResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"ExternalBlockHeight not exist",
			func() {
				request = &types.QueryLastObservedBlockHeightRequest{
					ChainName: suite.chainName,
				}
				response = &types.QueryLastObservedBlockHeightResponse{
					ExternalBlockHeight: 0,
					BlockHeight:         0,
				}
			},
			true,
		},
		{
			"ExternalBlockHeight exist",
			func() {
				suite.Ctx = suite.Ctx.WithBlockHeight(100)
				suite.Keeper().SetLastObservedBlockHeight(suite.Ctx, uint64(30), uint64(suite.Ctx.BlockHeight()))

				request = &types.QueryLastObservedBlockHeightRequest{
					ChainName: suite.chainName,
				}

				response = &types.QueryLastObservedBlockHeightResponse{
					ExternalBlockHeight: uint64(30),
					BlockHeight:         uint64(100),
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.LastObservedBlockHeight(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_Oracles() {
	var (
		request       *types.QueryOraclesRequest
		response      *types.QueryOraclesResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"Oracles exist and online is false",
			func() {
				externalKey, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(externalKey.PubKey().Address())
				suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
					OracleAddress:   suite.oracleAddrs[0].String(),
					BridgerAddress:  suite.bridgerAddrs[0].String(),
					ExternalAddress: externalAcc.String(),
					DelegateAmount:  sdkmath.ZeroInt(),
					StartHeight:     10,
					Online:          false,
				})
				request = &types.QueryOraclesRequest{
					ChainName: suite.chainName,
				}
				response = &types.QueryOraclesResponse{
					Oracles: []types.Oracle{
						{
							OracleAddress:   suite.oracleAddrs[0].String(),
							BridgerAddress:  suite.bridgerAddrs[0].String(),
							ExternalAddress: externalAcc.String(),
							DelegateAmount:  sdkmath.ZeroInt(),
							StartHeight:     10,
							Online:          false,
						},
					},
				}
			},
			true,
		},
		{
			"Oracles  exist and online is true",
			func() {
				externalKey, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				externalAcc := common.BytesToAddress(externalKey.PubKey().Address().Bytes())
				for i := 1; i < 4; i++ {
					online := true
					if i == 2 {
						online = false
					}
					suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
						OracleAddress:   suite.oracleAddrs[i].String(),
						BridgerAddress:  suite.bridgerAddrs[i].String(),
						ExternalAddress: externalAcc.String(),
						DelegateAmount:  sdkmath.ZeroInt(),
						StartHeight:     int64(i),
						Online:          online,
					})
				}
				request = &types.QueryOraclesRequest{
					ChainName: suite.chainName,
				}
				response = &types.QueryOraclesResponse{
					Oracles: []types.Oracle{
						{
							OracleAddress:   suite.oracleAddrs[1].String(),
							BridgerAddress:  suite.bridgerAddrs[1].String(),
							ExternalAddress: externalAcc.String(),
							DelegateAmount:  sdkmath.ZeroInt(),
							StartHeight:     int64(1),
							Online:          true,
						},
						{
							OracleAddress:   suite.oracleAddrs[2].String(),
							BridgerAddress:  suite.bridgerAddrs[2].String(),
							ExternalAddress: externalAcc.String(),
							DelegateAmount:  sdkmath.ZeroInt(),
							StartHeight:     int64(2),
							Online:          false,
						},
						{
							OracleAddress:   suite.oracleAddrs[3].String(),
							BridgerAddress:  suite.bridgerAddrs[3].String(),
							ExternalAddress: externalAcc.String(),
							DelegateAmount:  sdkmath.ZeroInt(),
							StartHeight:     int64(3),
							Online:          true,
						},
					},
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.Oracles(suite.Ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().ElementsMatch(response.Oracles, res.Oracles)

			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_ProjectedBatchTimeoutHeight() {
	var (
		request       *types.QueryProjectedBatchTimeoutHeightRequest
		response      *types.QueryProjectedBatchTimeoutHeightResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"ExternalBlockHeight is 0",
			func() {
				request = &types.QueryProjectedBatchTimeoutHeightRequest{
					ChainName: suite.chainName,
				}
				suite.Require().Equal(uint64(0), suite.Keeper().GetLastObservedBlockHeight(suite.Ctx).ExternalBlockHeight)
				suite.Require().Equal(uint64(0), suite.Keeper().GetLastObservedBlockHeight(suite.Ctx).BlockHeight)
				response = &types.QueryProjectedBatchTimeoutHeightResponse{
					TimeoutHeight: 0,
				}
			},
			true,
		},
		{
			name: "ProjectedBatchTimeoutHeight exist",
			malleate: func() {
				suite.Ctx = suite.Ctx.WithBlockHeight(5)
				suite.Keeper().SetLastObservedBlockHeight(suite.Ctx, 99, uint64(suite.Ctx.BlockHeight()))
				heights := suite.Keeper().GetLastObservedBlockHeight(suite.Ctx)
				suite.Assert().Equal(uint64(99), heights.ExternalBlockHeight)
				suite.Assert().Equal(uint64(5), heights.BlockHeight)

				request = &types.QueryProjectedBatchTimeoutHeightRequest{
					ChainName: suite.chainName,
				}
				var timeoutHeight uint64
				switch suite.chainName {
				case ethtypes.ModuleName:
					timeoutHeight = 2981
				case bsctypes.ModuleName:
					timeoutHeight = 14513
				}
				response = &types.QueryProjectedBatchTimeoutHeightResponse{
					TimeoutHeight: timeoutHeight,
				}
			},
			expPass: true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.ProjectedBatchTimeoutHeight(suite.Ctx, request)
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

func (suite *CrossChainGrpcTestSuite) TestKeeper_BridgeTokens() {
	testCases := []struct {
		name          string
		malleate      func() *types.QueryBridgeTokensResponse
		expectedError error
		expPass       bool
	}{
		{
			name: "query bridge tokens",
			malleate: func() *types.QueryBridgeTokensResponse {
				newBridgeTokens := make([]*types.BridgeToken, 3)

				for i := 0; i < 3; i++ {
					key, _ := ethsecp256k1.GenerateKey()
					channelIbc := ""
					if i == 2 {
						channelIbc = "transfer/channel-0"
					}
					tokenContract := common.BytesToAddress(key.PubKey().Address()).String()
					err := suite.Keeper().AttestationHandler(suite.Ctx, &types.MsgBridgeTokenClaim{
						TokenContract:  tokenContract,
						BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
						ChannelIbc:     hex.EncodeToString([]byte(channelIbc)),
					})

					suite.Require().NoError(err)
					bridgeDenom, found := suite.Keeper().GetBridgeDenomByContract(suite.Ctx, tokenContract)
					suite.Require().True(found)

					tokenContractByBridgeDenom, found := suite.Keeper().GetContractByBridgeDenom(suite.Ctx, bridgeDenom)
					suite.Require().True(found)
					suite.Require().Equal(tokenContract, tokenContractByBridgeDenom)

					newBridgeTokens[i] = &types.BridgeToken{
						Token: tokenContract,
						Denom: bridgeDenom,
					}
				}
				return &types.QueryBridgeTokensResponse{BridgeTokens: newBridgeTokens}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			response := testCase.malleate()
			res, err := suite.queryClient.BridgeTokens(
				suite.Ctx,
				&types.QueryBridgeTokensRequest{ChainName: suite.chainName},
			)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().ElementsMatch(response.BridgeTokens, res.BridgeTokens)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, testCase.expectedError)
			}
		})
	}
}

func (suite *CrossChainGrpcTestSuite) TestKeeper_BridgeCoinByToken() {
	var (
		request       *types.QueryBridgeCoinByDenomRequest
		response      *types.QueryBridgeCoinByDenomResponse
		expectedError error
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "bridge token not exist",
			malleate: func() {
				denom := helpers.GenHexAddress().Hex()
				request = &types.QueryBridgeCoinByDenomRequest{
					ChainName: suite.chainName,
					Denom:     denom,
				}
				expectedError = status.Error(codes.NotFound, "denom")
			},
			expPass: false,
		},
		{
			name: "bridge token exist",
			malleate: func() {
				token := helpers.GenHexAddress().Hex()
				suite.App.BankKeeper.SetDenomMetaData(suite.Ctx, banktypes.Metadata{
					Description: "The cross chain token of the Function X",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "usdt",
							Exponent: 0,
							Aliases: []string{
								types.NewBridgeDenom(ethtypes.ModuleName, token),
								types.NewBridgeDenom(bsctypes.ModuleName, token),
							},
						},
						{
							Denom:    "USDT",
							Exponent: 18,
						},
					},
					Base:    "usdt",
					Display: "usdt",
					Name:    "Tether USD",
					Symbol:  "USDT",
				})
				suite.AddTokenPair("usdt", true)
				err := suite.Keeper().AttestationHandler(suite.Ctx, &types.MsgBridgeTokenClaim{
					ChainName:      suite.chainName,
					TokenContract:  token,
					BridgerAddress: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
					ChannelIbc:     hex.EncodeToString([]byte("")),
				})
				suite.Require().NoError(err)
				denom, err := suite.queryClient.TokenToDenom(suite.Ctx, &types.QueryTokenToDenomRequest{
					ChainName: suite.chainName,
					Token:     token,
				})
				suite.Require().NoError(err)
				request = &types.QueryBridgeCoinByDenomRequest{
					ChainName: suite.chainName,
					Denom:     denom.Denom,
				}
				amount := sdkmath.NewInt(int64(tmrand.Uint32() + 1))
				claim := &types.MsgSendToFxClaim{
					Sender:        helpers.GenExternalAddr(suite.chainName),
					ChainName:     suite.chainName,
					TokenContract: token,
					Amount:        amount,
					Receiver:      helpers.GenAccAddress().String(),
					TargetIbc:     hex.EncodeToString([]byte("")),
					EventNonce:    1,
				}
				err = suite.Keeper().AttestationHandler(suite.Ctx, claim)
				suite.Require().NoError(err)
				err = suite.Keeper().ExecuteClaim(suite.Ctx, claim.EventNonce)
				suite.Require().NoError(err)
				response = &types.QueryBridgeCoinByDenomResponse{
					Coin: sdk.Coin{
						Denom:  denom.GetDenom(),
						Amount: amount,
					},
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			testCase.malleate()
			res, err := suite.queryClient.BridgeCoinByDenom(suite.Ctx, request)
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
