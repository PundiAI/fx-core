package tests_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/crosschain"
)

func TestIncreaseBridgeFeeABI(t *testing.T) {
	crossChainABI := crosschain.GetABI()

	method := crossChainABI.Methods[crosschain.IncreaseBridgeFeeMethodName]
	require.Equal(t, method, crosschain.IncreaseBridgeFeeMethod)
	require.Equal(t, 4, len(method.Inputs))
	require.Equal(t, 1, len(method.Outputs))

	event := crossChainABI.Events[crosschain.IncreaseBridgeFeeEventName]
	require.Equal(t, event, crosschain.IncreaseBridgeFeeEvent)
	require.Equal(t, 5, len(event.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestIncreaseBridgeFee() {
	randBridgeFee := big.NewInt(int64(tmrand.Uint32() + 10))
	crossChainTxFunc := func(signer *helpers.Signer, contact common.Address, moduleName string, amount, fee, value *big.Int) {
		data, err := crosschain.GetABI().Pack(
			"crossChain",
			contact,
			helpers.GenExternalAddr(moduleName),
			amount,
			fee,
			fxtypes.MustStrToByte32(moduleName),
			"",
		)
		suite.Require().NoError(err)

		tx, err := suite.PackEthereumTx(signer, crosschain.GetAddress(), value, data)
		suite.Require().NoError(err)
		res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
		suite.Require().NoError(err)
		suite.Require().False(res.Failed(), res.VmError)
	}
	transferCrossChainTxFunc := func(signer *helpers.Signer, contactAddr common.Address, moduleName string, amount, fee, value *big.Int) {
		data, err := contract.GetFIP20().ABI.Pack(
			"transferCrossChain",
			helpers.GenExternalAddr(moduleName),
			amount,
			fee,
			fxtypes.MustStrToByte32(moduleName),
		)
		suite.Require().NoError(err)
		tx, err := suite.PackEthereumTx(signer, contactAddr, value, data)
		suite.Require().NoError(err)
		res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
		suite.Require().NoError(err)
		suite.Require().False(res.Failed(), res.VmError)
	}
	increaseBridgeFeeFunc := func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
		queryServer := crosschainkeeper.NewQueryServerImpl(suite.CrossChainKeepers()[moduleName])
		pendingTx, err := queryServer.GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
			&crosschaintypes.QueryPendingSendToExternalRequest{
				ChainName:     moduleName,
				SenderAddress: signer.AccAddress().String(),
			})
		suite.Require().NoError(err)
		suite.Require().Equal(1, len(pendingTx.UnbatchedTransfers))
		totalAmount := pendingTx.UnbatchedTransfers[0].Token.Amount.Add(pendingTx.UnbatchedTransfers[0].Fee.Amount)
		suite.Require().Equal(randMint.String(), totalAmount.String())

		coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randBridgeFee))
		helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
		_, err = suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
			Coin:     coin,
			Receiver: signer.Address().Hex(),
			Sender:   signer.AccAddress().String(),
		})
		suite.Require().NoError(err)

		suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschain.GetAddress(), randBridgeFee)

		data, err := crosschain.GetABI().Pack(
			"increaseBridgeFee",
			moduleName,
			big.NewInt(int64(pendingTx.UnbatchedTransfers[0].Id)),
			pair.GetERC20Contract(),
			randBridgeFee,
		)
		suite.Require().NoError(err)
		return data, nil
	}

	testCases := []struct {
		name     string
		prepare  func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string)
		malleate func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string)
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

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschain.GetAddress(), randMint)

				crossChainTxFunc(signer, pair.GetERC20Contract(), moduleName, randMint, big.NewInt(0), big.NewInt(0))
				return pair, moduleName, ""
			},
			malleate: increaseBridgeFeeFunc,
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
			malleate: increaseBridgeFeeFunc,
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

				return &pair, moduleName, fxtypes.DefaultDenom
			},
			malleate: increaseBridgeFeeFunc,
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
			malleate: increaseBridgeFeeFunc,
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

				return pair, moduleName, ""
			},
			malleate: increaseBridgeFeeFunc,
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

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschain.GetAddress(), randMint)

				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)

				crossChainTxFunc(signer, pair.GetERC20Contract(), moduleName, amount, fee, big.NewInt(0))

				return &pair, moduleName, fxtypes.DefaultDenom
			},
			malleate: increaseBridgeFeeFunc,
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

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschain.GetAddress(), randMint)

				crossChainTxFunc(signer, pair.GetERC20Contract(), bsctypes.ModuleName, randMint, big.NewInt(0), big.NewInt(0))

				return pair, bsctypes.ModuleName, ""
			},
			malleate: increaseBridgeFeeFunc,
			result:   true,
		},
		{
			name: "failed - invalid chain name",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				chain := "123"
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					chain,
					big.NewInt(1),
					pair.GetERC20Contract(),
					randBridgeFee,
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
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				txID := big.NewInt(0)
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					moduleName,
					txID,
					pair.GetERC20Contract(),
					randBridgeFee,
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
			name: "failed - invalid bridge fee",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				fee := big.NewInt(0)
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					moduleName,
					big.NewInt(1),
					pair.GetERC20Contract(),
					fee,
				)
				suite.Require().NoError(err)
				return data, []string{fee.String()}
			},
			error: func(args []string) string {
				return "invalid add bridge fee"
			},
			result: false,
		},
		{
			name: "failed - not approve token",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					moduleName,
					big.NewInt(1),
					pair.GetERC20Contract(),
					randBridgeFee,
				)
				suite.Require().NoError(err)
				return data, []string{}
			},
			error: func(args []string) string {
				return "call transferFrom: execution reverted"
			},
			result: false,
		},
		{
			name: "failed - tx id not found",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randBridgeFee))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschain.GetAddress(), randBridgeFee)

				txID := big.NewInt(10)
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					moduleName,
					txID,
					pair.GetERC20Contract(),
					randBridgeFee,
				)
				suite.Require().NoError(err)
				return data, []string{txID.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("txId %s not in unbatched index! Must be in a batch!: invalid", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// token pair
			md := suite.GenerateCrossChainDenoms()
			pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, md.GetMetadata())
			suite.Require().NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(big.NewInt(0).Add(randMint, randBridgeFee)))
			moduleName := md.RandModule()

			pair, moduleName, evmToken := tc.prepare(pair, moduleName, signer, randMint)

			// check init balance zero
			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
			suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

			// get total supply
			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			packData, errArgs := tc.malleate(moduleName, pair, md, signer, randMint)
			tx, err := suite.PackEthereumTx(signer, crosschain.GetAddress(), big.NewInt(0), packData)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}
			// check result
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				queryServer := crosschainkeeper.NewQueryServerImpl(suite.CrossChainKeepers()[moduleName])
				pendingTx, err := queryServer.GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
					&crosschaintypes.QueryPendingSendToExternalRequest{
						ChainName:     moduleName,
						SenderAddress: signer.AccAddress().String(),
					})
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(pendingTx.UnbatchedTransfers))
				totalAmount := pendingTx.UnbatchedTransfers[0].Token.Amount.Add(pendingTx.UnbatchedTransfers[0].Fee.Amount)
				suite.Require().Equal(big.NewInt(0).Add(randMint, randBridgeFee).String(), totalAmount.String())

				totalAfter, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)

				for _, coin := range totalBefore.Supply {
					if has, find := totalAfter.Supply.Find(coin.Denom); has {
						if len(evmToken) > 0 && coin.Denom == evmToken {
							suite.Require().Equal(big.NewInt(0).Add(randBridgeFee, coin.Amount.BigInt()).String(), find.Amount.String(), coin.Denom)
							continue
						}
						suite.Require().Equal(coin.String(), find.String(), coin.Denom)
						continue
					}
					suite.Require().Equal(coin.Amount.String(), randBridgeFee.String())
				}

				for _, log := range res.Logs {
					if log.Topics[0] == crosschain.IncreaseBridgeFeeEvent.ID.String() {
						suite.Require().Equal(log.Address, crosschain.GetAddress().String())
						suite.Require().Equal(log.Topics[1], signer.Address().Hash().String())
						suite.Require().Equal(log.Topics[2], pair.GetERC20Contract().Hash().String())
						unpack, err := crosschain.IncreaseBridgeFeeEvent.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						chain := unpack[0].(string)
						suite.Require().Equal(chain, moduleName)
						txID := unpack[1].(*big.Int)
						suite.Require().True(txID.Uint64() > 0)
						fee := unpack[2].(*big.Int)
						suite.Require().Equal(fee.String(), randBridgeFee.String())
					}
				}

			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestIncreaseBridgeFeeExternal() {
	randBridgeFee := big.NewInt(int64(tmrand.Uint32() + 10))
	crossChainTxFunc := func(signer *helpers.Signer, contact common.Address, moduleName string, amount, fee, value *big.Int) {
		data, err := crosschain.GetABI().Pack(
			"crossChain",
			contact,
			helpers.GenExternalAddr(moduleName),
			amount,
			fee,
			fxtypes.MustStrToByte32(moduleName),
			"",
		)
		suite.Require().NoError(err)

		tx, err := suite.PackEthereumTx(signer, crosschain.GetAddress(), value, data)
		suite.Require().NoError(err)
		res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
		suite.Require().NoError(err)
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
		tx, err := suite.PackEthereumTx(signer, contact, value, data)
		suite.Require().NoError(err)
		res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
		suite.Require().NoError(err)
		suite.Require().False(res.Failed(), res.VmError)
	}
	increaseBridgeFeeFunc := func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
		queryServer := crosschainkeeper.NewQueryServerImpl(suite.CrossChainKeepers()[moduleName])
		pendingTx, err := queryServer.GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
			&crosschaintypes.QueryPendingSendToExternalRequest{
				ChainName:     moduleName,
				SenderAddress: signer.AccAddress().String(),
			})
		suite.Require().NoError(err)
		suite.Require().Equal(1, len(pendingTx.UnbatchedTransfers))
		totalAmount := pendingTx.UnbatchedTransfers[0].Token.Amount.Add(pendingTx.UnbatchedTransfers[0].Fee.Amount)
		suite.Require().Equal(randMint.String(), totalAmount.String())

		coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randBridgeFee))
		helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
		suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randBridgeFee)

		_, err = suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
			Coin:     coin,
			Receiver: signer.Address().Hex(),
			Sender:   signer.AccAddress().String(),
		})
		suite.Require().NoError(err)

		suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschain.GetAddress(), randBridgeFee)

		data, err := crosschain.GetABI().Pack(
			"increaseBridgeFee",
			moduleName,
			big.NewInt(int64(pendingTx.UnbatchedTransfers[0].Id)),
			pair.GetERC20Contract(),
			randBridgeFee,
		)
		suite.Require().NoError(err)
		return data, nil
	}

	testCases := []struct {
		name     string
		prepare  func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string)
		malleate func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - address + erc20 token",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))

				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)

				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschain.GetAddress(), randMint)

				crossChainTxFunc(signer, pair.GetERC20Contract(), moduleName, randMint, big.NewInt(0), big.NewInt(0))
				return pair, moduleName, ""
			},
			malleate: increaseBridgeFeeFunc,
			result:   true,
		},
		{
			name: "ok - fip20 contract + erc20 token",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))

				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)

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
			malleate: increaseBridgeFeeFunc,
			result:   true,
		},
		{
			name: "failed - invalid chain name",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				chain := "123"
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					chain,
					big.NewInt(1),
					pair.GetERC20Contract(),
					randBridgeFee,
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
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				txID := big.NewInt(0)
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					moduleName,
					txID,
					pair.GetERC20Contract(),
					randBridgeFee,
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
			name: "failed - invalid bridge fee",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				fee := big.NewInt(0)
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					moduleName,
					big.NewInt(1),
					pair.GetERC20Contract(),
					fee,
				)
				suite.Require().NoError(err)
				return data, []string{fee.String()}
			},
			error: func(args []string) string {
				return "invalid add bridge fee"
			},
			result: false,
		},
		{
			name: "failed - not approve token",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					moduleName,
					big.NewInt(1),
					pair.GetERC20Contract(),
					randBridgeFee,
				)
				suite.Require().NoError(err)
				return data, []string{}
			},
			error: func(args []string) string {
				return "call transferFrom: execution reverted"
			},
			result: false,
		},
		{
			name: "failed - tx id not found",
			prepare: func(pair *types.TokenPair, moduleName string, signer *helpers.Signer, randMint *big.Int) (*types.TokenPair, string, string) {
				return pair, moduleName, ""
			},
			malleate: func(moduleName string, pair *types.TokenPair, md Metadata, signer *helpers.Signer, _ *big.Int) ([]byte, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randBridgeFee))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randBridgeFee)
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschain.GetAddress(), randBridgeFee)

				txID := big.NewInt(10)
				data, err := crosschain.GetABI().Pack(
					"increaseBridgeFee",
					moduleName,
					txID,
					pair.GetERC20Contract(),
					randBridgeFee,
				)
				suite.Require().NoError(err)
				return data, []string{txID.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("txId %s not in unbatched index! Must be in a batch!: invalid", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// token pair
			md := suite.GenerateCrossChainDenoms()

			// deploy fip20 external
			fip20External, err := suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, signer.Address(), "Test token", "TEST", 18)
			suite.Require().NoError(err)
			// token pair
			pair, err := suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, fip20External, md.GetMetadata().DenomUnits[0].Aliases...)
			suite.Require().NoError(err)

			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(big.NewInt(0).Add(randMint, randBridgeFee)))
			moduleName := md.RandModule()

			pair, moduleName, _ = tc.prepare(pair, moduleName, signer, randMint)

			// check init balance zero
			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
			suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

			// get total supply
			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			packData, errArgs := tc.malleate(moduleName, pair, md, signer, randMint)
			tx, err := suite.PackEthereumTx(signer, crosschain.GetAddress(), big.NewInt(0), packData)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}
			// check result
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				queryServer := crosschainkeeper.NewQueryServerImpl(suite.CrossChainKeepers()[moduleName])
				pendingTx, err := queryServer.GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
					&crosschaintypes.QueryPendingSendToExternalRequest{
						ChainName:     moduleName,
						SenderAddress: signer.AccAddress().String(),
					})
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(pendingTx.UnbatchedTransfers))
				totalAmount := pendingTx.UnbatchedTransfers[0].Token.Amount.Add(pendingTx.UnbatchedTransfers[0].Fee.Amount)
				suite.Require().Equal(big.NewInt(0).Add(randMint, randBridgeFee).String(), totalAmount.String())

				totalAfter, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)

				for _, coin := range totalBefore.Supply {
					if pair.GetDenom() == coin.Denom {
						find, expect := totalAfter.Supply.Find(coin.Denom)
						suite.True(find)
						suite.Require().Equal(coin.Amount.Add(sdkmath.NewIntFromBigInt(randBridgeFee)).String(), expect.Amount.String(), coin.Denom)
						continue
					}
					if pair.IsNativeERC20() && strings.HasPrefix(coin.Denom, moduleName) {
						suite.Equal(totalAfter.Supply.AmountOf(coin.Denom), coin.Amount.Add(sdkmath.NewIntFromBigInt(randBridgeFee)))
						continue
					}
					suite.Equal(coin.Amount.String(), totalAfter.Supply.AmountOf(coin.Denom).String(), coin.Denom)
				}

				for _, log := range res.Logs {
					if log.Topics[0] == crosschain.IncreaseBridgeFeeEvent.ID.String() {
						suite.Require().Equal(log.Address, crosschain.GetAddress().String())
						suite.Require().Equal(log.Topics[1], signer.Address().Hash().String())
						suite.Require().Equal(log.Topics[2], pair.GetERC20Contract().Hash().String())
						unpack, err := crosschain.IncreaseBridgeFeeEvent.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						chain := unpack[0].(string)
						suite.Require().Equal(chain, moduleName)
						txID := unpack[1].(*big.Int)
						suite.Require().True(txID.Uint64() > 0)
						fee := unpack[2].(*big.Int)
						suite.Require().Equal(fee.String(), randBridgeFee.String())
					}
				}

			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}
