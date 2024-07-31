package precompile_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	"github.com/functionx/fx-core/v7/x/crosschain/precompile"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
)

func TestCancelSendToExternalABI(t *testing.T) {
	cancelSendToExternal := precompile.NewCancelSendToExternalMethod(nil)

	require.Equal(t, 2, len(cancelSendToExternal.Method.Inputs))
	require.Equal(t, 1, len(cancelSendToExternal.Method.Outputs))

	require.Equal(t, 3, len(cancelSendToExternal.Event.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestCancelSendToExternal() {
	crossChainTxFunc := func(signer *helpers.Signer, contact common.Address, moduleName string, amount, fee, value *big.Int) {
		data, err := crosschaintypes.GetABI().Pack(
			"crossChain",
			contact,
			helpers.GenExternalAddr(moduleName),
			amount,
			fee,
			fxtypes.MustStrToByte32(moduleName),
			"",
		)
		suite.Require().NoError(err)

		res := suite.EthereumTx(signer, crosschaintypes.GetAddress(), value, data)
		suite.Require().False(res.Failed(), res.VmError)
	}
	transferCrossChainTxFunc := func(signer *helpers.Signer, contact common.Address, moduleName string, amount, fee, value *big.Int) {
		data, err := contract.GetFIP20().ABI.Pack(
			"transferCrossChain",
			helpers.GenExternalAddr(moduleName),
			amount,
			fee,
			fxtypes.MustStrToByte32(moduleName),
		)
		suite.Require().NoError(err)
		res := suite.EthereumTx(signer, contact, value, data)
		suite.Require().False(res.Failed(), res.VmError)
	}
	refundPackFunc := func(moduleName string, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
		queryServer := crosschainkeeper.NewQueryServerImpl(suite.CrossChainKeepers()[moduleName])
		externalTx, err := queryServer.GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
			&crosschaintypes.QueryPendingSendToExternalRequest{
				ChainName:     moduleName,
				SenderAddress: signer.AccAddress().String(),
			})
		suite.Require().NoError(err)
		suite.Require().Equal(1, len(externalTx.UnbatchedTransfers))

		data, err := crosschaintypes.GetABI().Pack(
			"cancelSendToExternal",
			moduleName,
			big.NewInt(int64(externalTx.UnbatchedTransfers[0].Id)),
		)
		suite.Require().NoError(err)
		return data, nil
	}

	testCases := []struct {
		name     string
		prepare  func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string)
		malleate func(moduleName string, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - address + erc20 token",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), randMint)

				crossChainTxFunc(signer, pair.GetERC20Contract(), moduleName, randMint, big.NewInt(0), big.NewInt(0))
				return pair, moduleName, ""
			},
			malleate: refundPackFunc,
			result:   true,
		},
		{
			name: "ok - fip20 contract + erc20 token",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)
				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				transferCrossChainTxFunc(signer, pair.GetERC20Contract(), moduleName, amount, fee, big.NewInt(0))

				return pair, moduleName, ""
			},
			malleate: refundPackFunc,
			result:   true,
		},
		{
			name: "ok - fip20 contract + evm token",
			prepare: func(_ *types.TokenPair, _ string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				moduleName := ethtypes.ModuleName
				pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)
				suite.CrossChainKeepers()[moduleName].AddBridgeToken(suite.ctx, helpers.GenHexAddress().String(), pair.GetDenom())

				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))

				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				transferCrossChainTxFunc(signer, pair.GetERC20Contract(), moduleName, amount, fee, randMint)

				return &pair, moduleName, ""
			},
			malleate: refundPackFunc,
			result:   true,
		},
		{
			name: "ok - address + evm token",
			prepare: func(_ *types.TokenPair, _ string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				moduleName := ethtypes.ModuleName

				suite.CrossChainKeepers()[moduleName].AddBridgeToken(suite.ctx, helpers.GenHexAddress().String(), fxtypes.DefaultDenom)

				coin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))

				pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)

				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)

				crossChainTxFunc(signer, common.Address{}, moduleName, amount, fee, randMint)

				return &pair, moduleName, fxtypes.DefaultDenom
			},
			malleate: refundPackFunc,
			result:   true,
		},
		{
			name: "ok - address + origin token",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				suite.CrossChainKeepers()[moduleName].AddBridgeToken(suite.ctx, helpers.GenHexAddress().String(), pair.GetDenom())

				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))

				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)

				// convert denom to many
				fxTarget := fxtypes.ParseFxTarget(moduleName)
				targetCoin, err := suite.app.Erc20Keeper.ConvertDenomToTarget(suite.ctx, signer.AccAddress(),
					sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint)), fxTarget)
				suite.Require().NoError(err)

				_, err = suite.CrossChainKeepers()[moduleName].AddToOutgoingPool(suite.ctx,
					signer.AccAddress(), signer.Address().String(),
					sdk.NewCoin(targetCoin.Denom, sdkmath.NewIntFromBigInt(amount)),
					sdk.NewCoin(targetCoin.Denom, sdkmath.NewIntFromBigInt(fee)))
				suite.Require().NoError(err)

				return pair, moduleName, pair.GetDenom()
			},
			malleate: refundPackFunc,
			result:   true,
		},
		{
			name: "ok - address + wrapper origin token",
			prepare: func(_ *types.TokenPair, _ string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				moduleName := ethtypes.ModuleName
				pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)

				suite.CrossChainKeepers()[moduleName].AddBridgeToken(suite.ctx, helpers.GenHexAddress().String(), fxtypes.DefaultDenom)

				coin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))

				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), randMint)

				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)

				crossChainTxFunc(signer, pair.GetERC20Contract(), moduleName, amount, fee, big.NewInt(0))

				return &pair, moduleName, ""
			},
			malleate: refundPackFunc,
			result:   true,
		},
		{
			name: "ok - address + ibc token",
			prepare: func(_ *types.TokenPair, _ string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				sourcePort, sourceChannel := suite.RandTransferChannel()
				tokenAddress := helpers.GenHexAddress()
				denom, err := suite.CrossChainKeepers()[bsctypes.ModuleName].SetIbcDenomTrace(suite.ctx,
					tokenAddress.Hex(), hex.EncodeToString([]byte(fmt.Sprintf("%s/%s", sourcePort, sourceChannel))))
				suite.Require().NoError(err)
				suite.CrossChainKeepers()[bsctypes.ModuleName].AddBridgeToken(suite.ctx, tokenAddress.Hex(), denom)

				symbol := helpers.NewRandSymbol()
				ibcMD := banktypes.Metadata{
					Description: "The cross chain token of the Function X",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    denom,
							Exponent: 0,
						},
						{
							Denom:    symbol,
							Exponent: 18,
						},
					},
					Base:    denom,
					Display: denom,
					Name:    fmt.Sprintf("%s Token", symbol),
					Symbol:  symbol,
				}
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, ibcMD)
				suite.Require().NoError(err)

				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				_, err = suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
					&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), randMint)

				crossChainTxFunc(signer, pair.GetERC20Contract(), bsctypes.ModuleName, randMint, big.NewInt(0), big.NewInt(0))

				return pair, bsctypes.ModuleName, ""
			},
			malleate: refundPackFunc,
			result:   true,
		},
		{
			name: "failed - invalid chain name",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				chain := "123"
				data, err := crosschaintypes.GetABI().Pack(
					"cancelSendToExternal",
					chain,
					big.NewInt(1),
				)
				suite.Require().NoError(err)
				return data, []string{chain}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid module name: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed - invalid tx id",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				txID := big.NewInt(0)
				data, err := crosschaintypes.GetABI().Pack(
					"cancelSendToExternal",
					moduleName,
					txID,
				)
				suite.Require().NoError(err)
				return data, []string{txID.String()}
			},
			error: func(args []string) string {
				return "invalid tx id"
			},
			result: false,
		},
		{
			name: "failed - tx id not found",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				txID := big.NewInt(10)
				data, err := crosschaintypes.GetABI().Pack(
					"cancelSendToExternal",
					moduleName,
					txID,
				)
				suite.Require().NoError(err)
				return data, []string{txID.String()}
			},
			error: func(args []string) string {
				return "pool transaction: unknown"
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			signer := suite.RandSigner()
			// token pair
			md := suite.GenerateCrossChainDenoms()
			pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, md.GetMetadata())
			suite.Require().NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))
			moduleName := md.RandModule()

			pair, moduleName, originToken := tc.prepare(pair, moduleName, signer, randMint)

			if len(originToken) > 0 && originToken != fxtypes.DefaultDenom {
				queryServer := crosschainkeeper.NewQueryServerImpl(suite.CrossChainKeepers()[moduleName])
				petxs, err := queryServer.GetPendingSendToExternal(suite.ctx, &crosschaintypes.QueryPendingSendToExternalRequest{
					ChainName:     moduleName,
					SenderAddress: signer.AccAddress().String(),
				})
				suite.Require().NoError(err)
				if len(petxs.UnbatchedTransfers) > 0 && !strings.Contains(tc.name, "ok - address + origin token") { // send by chain, not add relation
					relation := suite.app.Erc20Keeper.HasOutgoingTransferRelation(suite.ctx, moduleName, petxs.UnbatchedTransfers[0].Id)
					suite.Require().True(relation)
				}
			}

			packData, errArgs := tc.malleate(moduleName, md, signer, randMint)

			// check init balance zero
			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
			suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

			// get total supply
			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			res := suite.EthereumTx(signer, crosschaintypes.GetAddress(), big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)
				// check balance after tx
				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
				balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
				if len(originToken) > 0 {
					suite.Require().True(chainBalances.AmountOf(originToken).Equal(sdkmath.NewIntFromBigInt(randMint)), chainBalances.String())
					suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())
					chainDenom := md.GetDenom(originToken)
					if len(chainDenom) > 0 {
						suite.Require().True(chainBalances.AmountOf(chainDenom).Equal(sdkmath.NewIntFromBigInt(randMint)), chainBalances.String())
					}
				} else {
					suite.Require().True(chainBalances.IsZero(), chainBalances.String())
					suite.Require().True(balance.Cmp(randMint) == 0, balance.String())
				}

				// check total supply equal
				manyToOne := make(map[string]bool)
				suite.app.BankKeeper.IterateAllDenomMetaData(suite.ctx, func(md banktypes.Metadata) bool {
					if len(md.DenomUnits) > 0 && len(md.DenomUnits[0].Aliases) > 0 {
						manyToOne[md.Base] = true
					}
					return false
				})
				totalAfter, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)

				for _, coin := range totalAfter.Supply {
					if manyToOne[coin.Denom] {
						continue
					}
					expect := totalBefore.Supply.AmountOf(coin.Denom)

					md, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, pair.GetDenom())
					suite.Require().True(found)

					has := false
					if len(md.DenomUnits) > 0 && len(md.DenomUnits[0].Aliases) > 0 {
						for _, alias := range md.DenomUnits[0].Aliases {
							if strings.HasPrefix(alias, moduleName) && alias == coin.GetDenom() {
								has = true
								break
							}
						}
					}
					if has || strings.HasPrefix(coin.GetDenom(), "ibc/") {
						expect = expect.Add(sdkmath.NewIntFromBigInt(randMint))
					}
					suite.Require().Equal(coin.Amount.String(), expect.String(), coin.Denom)
				}

				for _, log := range res.Logs {
					event := crosschaintypes.GetABI().Events["CancelSendToExternal"]
					if log.Topics[0] == event.ID.String() {
						suite.Require().Equal(log.Address, crosschaintypes.GetAddress().String())
						suite.Require().Equal(log.Topics[1], signer.Address().Hash().String())
						unpack, err := event.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						chain := unpack[0].(string)
						suite.Require().Equal(chain, moduleName)
						txID := unpack[1].(*big.Int)
						suite.Require().True(txID.Uint64() > 0)
					}
				}

			} else {
				suite.Error(res, errors.New(tc.error(errArgs)))
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestDeleteOutgoingTransferRelation() {
	signer := suite.RandSigner()
	// token pair
	md := suite.GenerateCrossChainDenoms()
	pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, md.GetMetadata())
	suite.Require().NoError(err)
	randMint := big.NewInt(int64(tmrand.Uint32() + 10))
	suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))
	moduleName := md.RandModule()

	coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
	helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
	_, err = suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
		Coin:     coin,
		Receiver: signer.Address().Hex(),
		Sender:   signer.AccAddress().String(),
	})
	suite.Require().NoError(err)

	suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), randMint)

	fee := big.NewInt(1)
	amount := big.NewInt(0).Sub(randMint, fee)
	data, err := crosschaintypes.GetABI().Pack("crossChain", pair.GetERC20Contract(),
		helpers.GenExternalAddr(moduleName), amount, fee, fxtypes.MustStrToByte32(moduleName), "")
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, crosschaintypes.GetAddress(), big.NewInt(0), data)
	suite.Require().False(res.Failed(), res.VmError)

	// get crosschain pending tx
	queryServer := crosschainkeeper.NewQueryServerImpl(suite.CrossChainKeepers()[moduleName])
	petxs, err := queryServer.GetPendingSendToExternal(suite.ctx, &crosschaintypes.QueryPendingSendToExternalRequest{
		ChainName:     moduleName,
		SenderAddress: signer.AccAddress().String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(petxs.UnbatchedTransfers, 1)

	txId := petxs.UnbatchedTransfers[0].Id
	txContract := petxs.UnbatchedTransfers[0].Token.Contract

	suite.CrossChainKeepers()[moduleName].SetLastObservedBlockHeight(suite.ctx, 100, uint64(suite.ctx.BlockHeight()))

	batch, err := suite.CrossChainKeepers()[moduleName].BuildOutgoingTxBatch(suite.ctx, txContract,
		signer.Address().String(), 100, sdkmath.NewInt(0), sdkmath.NewInt(1))
	suite.Require().NoError(err)
	batchNonce := batch.BatchNonce

	relation := suite.app.Erc20Keeper.HasOutgoingTransferRelation(suite.ctx, moduleName, txId)
	suite.Require().True(relation)

	suite.CrossChainKeepers()[moduleName].OutgoingTxBatchExecuted(suite.ctx, txContract, batchNonce)

	relation = suite.app.Erc20Keeper.HasOutgoingTransferRelation(suite.ctx, moduleName, txId)
	suite.Require().False(relation)
}
