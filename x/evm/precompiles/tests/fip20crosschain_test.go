package tests_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/crosschain"
)

func TestFIP20CrossChainABI(t *testing.T) {
	crossChainABI := crosschain.GetABI()

	method := crossChainABI.Methods[crosschain.FIP20CrossChainMethod.Name]
	require.Equal(t, method, crosschain.FIP20CrossChainMethod)
	require.Equal(t, 6, len(method.Inputs))
	require.Equal(t, 1, len(method.Outputs))
}

func (suite *PrecompileTestSuite) TestFIP20CrossChain() {
	testCases := []struct {
		name          string
		malleate      func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string)
		error         func(args []string) string
		result        bool
		isPendingPool bool
	}{
		{
			name: "ok",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				moduleName := md.RandModule()
				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					amount,
					fee,
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), moduleName, nil
			},
			result: true,
		},
		{
			name: "ok - msg.value",
			malleate: func(_ *types.TokenPair, _ Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)

				coin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))

				moduleName := ethtypes.ModuleName
				suite.CrossChainKeepers()[moduleName].AddBridgeToken(suite.ctx, helpers.GenerateAddress().String(), fxtypes.DefaultDenom)

				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				data, err := contract.GetWFX().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					amount,
					fee,
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				return data, &pair, randMint, moduleName, nil
			},
			result: true,
		},
		{
			name: "ok - ibc token",
			malleate: func(_ *types.TokenPair, _ Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				sourcePort, sourceChannel := suite.RandTransferChannel()
				tokenAddress := helpers.GenerateAddress()
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

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(bsctypes.ModuleName),
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32("chain/"+bsctypes.ModuleName),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), bsctypes.ModuleName, nil
			},
			result: true,
		},
		{
			name: "failed - from address insufficient funds",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
					&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				moduleName := md.RandModule()
				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					randMint,
					big.NewInt(1),
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), moduleName, nil
			},
			error: func(args []string) string {
				return "execution reverted: transfer amount exceeds balance"
			},
			result: false,
		},
		{
			name: "success - module insufficient funds - add to pending pool",

			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				// add relay token
				addAmount := big.NewInt(0).Add(randMint, big.NewInt(1))
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(addAmount))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
					&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				moduleName := md.RandModule()
				expectedModuleName := moduleName
				if moduleName == "gravity" {
					expectedModuleName = "eth"
				}

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					addAmount,
					big.NewInt(0),
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), moduleName, []string{
					fmt.Sprintf("%s%s", big.NewInt(0).Sub(addAmount, big.NewInt(1)).String(), md.GetDenom(expectedModuleName)),
					fmt.Sprintf("%s%s", addAmount.String(), md.GetDenom(expectedModuleName)),
				}
			},
			error: func(args []string) string {
				return ""
			},
			result:        true,
			isPendingPool: true,
		},
		{
			name: "failed - target not support",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				// add relay token
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
					&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				unknownChain := "chainabc"
				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddress().String(),
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(unknownChain),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), "", nil
			},
			error: func(args []string) string {
				return "execution reverted: fip-cross-chain failed: invalid target"
			},
			result: false,
		},
		{
			name: "failed - bridge token is not exist",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				// _, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
				// 	&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				// suite.Require().NoError(err)

				randDenom := helpers.NewRandDenom()
				moduleName := md.RandModule()
				aliasDenom := crosschaintypes.NewBridgeDenom(moduleName, helpers.GenerateAddressByModule(moduleName))
				newPair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, banktypes.Metadata{
					Description: "New Token",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    randDenom,
							Exponent: 0,
							Aliases:  []string{aliasDenom},
						},
						{
							Denom:    strings.ToUpper(randDenom),
							Exponent: 18,
						},
					},
					Base:    randDenom,
					Display: randDenom,
					Name:    randDenom,
					Symbol:  strings.ToUpper(randDenom),
				})
				suite.Require().NoError(err)

				randCoin := sdk.NewCoin(randDenom, sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(randCoin))

				aliasCoin := sdk.NewCoin(aliasDenom, sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(aliasCoin))
				err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, signer.AccAddress(), types.ModuleName, sdk.NewCoins(aliasCoin))
				suite.Require().NoError(err)

				_, err = suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
					&types.MsgConvertCoin{Coin: randCoin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				return data, newPair, big.NewInt(0), moduleName, []string{}
			},
			error: func(args []string) string {
				return "execution reverted: fip-cross-chain failed: cross chain error: bridge token is not exist: invalid"
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
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))

			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
			suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

			packData, newPair, value, moduleName, errArgs := tc.malleate(pair, md, signer, randMint)

			if len(moduleName) > 0 {
				resp, err := suite.CrossChainKeepers()[moduleName].GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
					&crosschaintypes.QueryPendingSendToExternalRequest{
						ChainName:     moduleName,
						SenderAddress: signer.AccAddress().String(),
					})
				suite.Require().NoError(err)
				suite.Require().Equal(0, len(resp.UnbatchedTransfers))
				suite.Require().Equal(0, len(resp.TransfersInBatches))
			}

			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			tx, err := suite.PackEthereumTx(signer, newPair.GetERC20Contract(), value, packData)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			// check result
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
				suite.Require().True(chainBalances.IsZero(), chainBalances.String())
				balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
				suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

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
					suite.Require().Equal(coin.Amount.String(), expect.String(), coin.Denom)
				}

				if tc.isPendingPool {
					resp, err := suite.CrossChainKeepers()[moduleName].GetPendingPoolSendToExternal(sdk.WrapSDKContext(suite.ctx), &crosschaintypes.QueryPendingPoolSendToExternalRequest{
						ChainName:     moduleName,
						SenderAddress: signer.AccAddress().String(),
					})
					suite.Require().NoError(err)
					suite.Require().Equal(1, len(resp.Txs))
				} else {
					resp, err := suite.CrossChainKeepers()[moduleName].GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
						&crosschaintypes.QueryPendingSendToExternalRequest{
							ChainName:     moduleName,
							SenderAddress: signer.AccAddress().String(),
						})
					suite.Require().NoError(err)
					suite.Require().Equal(1, len(resp.UnbatchedTransfers))
					suite.Require().Equal(0, len(resp.TransfersInBatches))
					suite.Require().Equal(signer.AccAddress().String(), resp.UnbatchedTransfers[0].Sender)
					// NOTE: fee + amount == randMint
					suite.Require().Equal(randMint.String(), resp.UnbatchedTransfers[0].Fee.Amount.Add(resp.UnbatchedTransfers[0].Token.Amount).BigInt().String())
					if !strings.EqualFold(resp.UnbatchedTransfers[0].Token.Contract, strings.TrimPrefix(md.GetDenom(moduleName), moduleName)) {
						bridgeToken := suite.CrossChainKeepers()[moduleName].GetDenomBridgeToken(suite.ctx, newPair.Denom)
						suite.Require().Equal(resp.UnbatchedTransfers[0].Token.Contract, bridgeToken.Token, moduleName)
					}
				}
			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestFIP20CrossChainExternal() {
	testCases := []struct {
		name     string
		malleate func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))

				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)

				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				moduleName := md.RandModule()
				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					amount,
					fee,
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), moduleName, nil
			},
			result: true,
		},
		{
			name: "failed - from address insufficient funds",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
					&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				moduleName := md.RandModule()
				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					randMint,
					big.NewInt(1),
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), moduleName, nil
			},
			error: func(args []string) string {
				return "execution reverted: transfer amount exceeds balance"
			},
			result: false,
		},
		{
			name: "failed - target not support",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				// add relay token
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
					&types.MsgConvertCoin{Coin: coin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				unknownChain := "chainabc"
				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddress().String(),
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(unknownChain),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), "", nil
			},
			error: func(args []string) string {
				return "execution reverted: fip-cross-chain failed: invalid target"
			},
			result: false,
		},
		{
			name: "failed - bridge token is not exist",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, *types.TokenPair, *big.Int, string, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)

				randDenom := helpers.NewRandDenom()
				moduleName := md.RandModule()
				aliasDenom := crosschaintypes.NewBridgeDenom(moduleName, helpers.GenerateAddressByModule(moduleName))
				newPair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, banktypes.Metadata{
					Description: "New Token",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    randDenom,
							Exponent: 0,
							Aliases:  []string{aliasDenom},
						},
						{
							Denom:    strings.ToUpper(randDenom),
							Exponent: 18,
						},
					},
					Base:    randDenom,
					Display: randDenom,
					Name:    randDenom,
					Symbol:  strings.ToUpper(randDenom),
				})
				suite.Require().NoError(err)

				randCoin := sdk.NewCoin(randDenom, sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(randCoin))

				aliasCoin := sdk.NewCoin(aliasDenom, sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(aliasCoin))
				err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, signer.AccAddress(), types.ModuleName, sdk.NewCoins(aliasCoin))
				suite.Require().NoError(err)

				_, err = suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
					&types.MsgConvertCoin{Coin: randCoin, Receiver: signer.Address().Hex(), Sender: signer.AccAddress().String()})
				suite.Require().NoError(err)

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				return data, newPair, big.NewInt(0), moduleName, []string{}
			},
			error: func(args []string) string {
				return "execution reverted: fip-cross-chain failed: cross chain error: bridge token is not exist: invalid"
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()

			md := suite.GenerateCrossChainDenoms()
			// deploy fip20 external
			fip20External, err := suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, signer.Address(), "Test token", "TEST", 18)
			suite.Require().NoError(err)
			// token pair
			pair, err := suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, fip20External, md.GetMetadata().DenomUnits[0].Aliases...)
			suite.Require().NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))

			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
			suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

			packData, newPair, value, moduleName, errArgs := tc.malleate(pair, md, signer, randMint)

			if len(moduleName) > 0 {
				resp, err := suite.CrossChainKeepers()[moduleName].GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
					&crosschaintypes.QueryPendingSendToExternalRequest{
						ChainName:     moduleName,
						SenderAddress: signer.AccAddress().String(),
					})
				suite.Require().NoError(err)
				suite.Require().Equal(0, len(resp.UnbatchedTransfers))
				suite.Require().Equal(0, len(resp.TransfersInBatches))
			}

			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			tx, err := suite.PackEthereumTx(signer, newPair.GetERC20Contract(), value, packData)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			// check result
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
				suite.Require().True(chainBalances.IsZero(), chainBalances.String())
				balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
				suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

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

					moduleDenom := md.GetDenom(moduleName)
					suite.Require().NotEmpty(moduleDenom)
					if coin.Denom == moduleDenom {
						// crosschain: erc20 token --> base denom --> module denom
						suite.Require().Equal(coin.Amount.String(), expect.Add(sdkmath.NewIntFromBigInt(randMint)).String(), coin.Denom)
					} else {
						suite.Require().Equal(coin.Amount.String(), expect.String(), coin.Denom)
					}
				}

				resp, err := suite.CrossChainKeepers()[moduleName].GetPendingSendToExternal(sdk.WrapSDKContext(suite.ctx),
					&crosschaintypes.QueryPendingSendToExternalRequest{
						ChainName:     moduleName,
						SenderAddress: signer.AccAddress().String(),
					})
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(resp.UnbatchedTransfers))
				suite.Require().Equal(0, len(resp.TransfersInBatches))
				suite.Require().Equal(signer.AccAddress().String(), resp.UnbatchedTransfers[0].Sender)
				// NOTE: fee + amount == randMint
				suite.Require().Equal(randMint.String(), resp.UnbatchedTransfers[0].Fee.Amount.Add(resp.UnbatchedTransfers[0].Token.Amount).BigInt().String())
				if !strings.EqualFold(resp.UnbatchedTransfers[0].Token.Contract, strings.TrimPrefix(md.GetDenom(moduleName), moduleName)) {
					bridgeToken := suite.CrossChainKeepers()[moduleName].GetDenomBridgeToken(suite.ctx, newPair.Denom)
					suite.Require().Equal(resp.UnbatchedTransfers[0].Token.Contract, bridgeToken.Token, moduleName)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestFIP20CrossChainIBC() {
	testCases := []struct {
		name     string
		malleate func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - ibc token",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				// add relay token
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				prefix, recipient := suite.RandPrefixAndAddress()

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					recipient,
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), nil
			},
			result: true,
		},
		{
			name: "ok - base token",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				// add relay token
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				sourcePort1, sourceChannel1 := suite.RandTransferChannel()
				prefix, recipient := suite.RandPrefixAndAddress()

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					recipient,
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(fmt.Sprintf("ibc/%s/%s/%s", prefix, sourcePort1, sourceChannel1)),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), nil
			},
			result: true,
		},
		{
			name: "ok - msg.value",
			malleate: func(_ *types.TokenPair, _ Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				coin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))

				moduleName := ethtypes.ModuleName
				suite.CrossChainKeepers()[moduleName].AddBridgeToken(suite.ctx, helpers.GenerateAddress().String(), fxtypes.DefaultDenom)

				fee := big.NewInt(1)
				amount := big.NewInt(0).Sub(randMint, fee)
				data, err := contract.GetWFX().ABI.Pack(
					"transferCrossChain",
					helpers.GenerateAddressByModule(moduleName),
					amount,
					fee,
					fxtypes.MustStrToByte32(moduleName),
				)
				suite.Require().NoError(err)

				pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)

				return data, &pair, randMint, nil
			},
			result: true,
		},
		{
			name: "failed - no zero fee",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)
				prefix, recipient := suite.RandPrefixAndAddress()

				fee := big.NewInt(int64(tmrand.Intn(1000) + 10))
				relayAmount := big.NewInt(0).Sub(randMint, fee)

				ibcTarget := fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					recipient,
					relayAmount,
					fee,
					fxtypes.MustStrToByte32(ibcTarget),
				)
				suite.Require().NoError(err)

				aliasDenom := ""
				for _, alias := range md.metadata.DenomUnits[0].Aliases {
					if strings.Contains(alias, ibctransfertypes.DenomPrefix+"/") {
						hexHash := strings.TrimPrefix(alias, ibctransfertypes.DenomPrefix+"/")
						hash, err := ibctransfertypes.ParseHexHash(hexHash)
						suite.Require().NoError(err)

						denomTrace, found := suite.app.IBCTransferKeeper.GetDenomTrace(suite.ctx, hash)
						if !found {
							continue
						}
						if !strings.HasPrefix(denomTrace.GetPath(), fmt.Sprintf("%s/%s", ibctransfertypes.PortID, sourceChannel)) {
							continue
						}
						aliasDenom = alias
						break
					}
				}
				return data, pair, big.NewInt(0), []string{sdk.NewCoin(aliasDenom, sdkmath.NewIntFromBigInt(fee)).String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: fip-cross-chain failed: ibc transfer fee must be zero: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed - invalid recipient address - hex",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				prefix, recipient := "0x", suite.RandSigner().AccAddress().String()
				ibcTarget := fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					recipient,
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(ibcTarget),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), []string{recipient}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: fip-cross-chain failed: invalid to address: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed - invalid recipient address - bench32",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				prefix, recipient := "px", helpers.GenerateAddress().Hex()
				ibcTarget := fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					recipient,
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(ibcTarget),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), []string{recipient}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: fip-cross-chain failed: invalid to address: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// set port channel
			sourcePort, sourceChannel := suite.RandTransferChannel()
			// add ibc token
			ibcToken := suite.AddIBCToken(sourcePort, sourceChannel)
			// token pair
			md := suite.GenerateCrossChainDenoms(ibcToken)
			pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, md.GetMetadata())
			suite.Require().NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 100000))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))

			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
			suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

			commitments := suite.app.IBCKeeper.ChannelKeeper.GetAllPacketCommitmentsAtChannel(suite.ctx, sourcePort, sourceChannel)
			ibcTxs := make(map[string]bool, len(commitments))
			for _, commitment := range commitments {
				ibcTxs[fmt.Sprintf("%s/%s/%d", commitment.PortId, commitment.ChannelId, commitment.Sequence)] = true
			}

			packData, newPair, value, errArgs := tc.malleate(pair, md, signer, randMint, sourcePort, sourceChannel)

			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			tx, err := suite.PackEthereumTx(signer, newPair.GetERC20Contract(), value, packData)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			// check result
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
				suite.Require().True(chainBalances.IsZero(), chainBalances.String())
				balance := suite.BalanceOf(newPair.GetERC20Contract(), signer.Address())
				suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

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
					suite.Require().Equal(coin.Amount.String(), expect.String(), coin.Denom)
				}

				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type != ibcchanneltypes.EventTypeSendPacket {
						continue
					}
					var eventPortId, eventChannelId string
					var sequence string
					var data []byte

					for _, attr := range event.Attributes {
						attrKey, attrValue := string(attr.Key), string(attr.Value)
						if attrKey == ibcchanneltypes.AttributeKeyDataHex {
							data, err = hex.DecodeString(attrValue)
							suite.Require().NoError(err)
						}
						if attrKey == ibcchanneltypes.AttributeKeySequence {
							sequence = attrValue
						}
						if attrKey == ibcchanneltypes.AttributeKeySrcPort {
							eventPortId = attrValue
						}
						if attrKey == ibcchanneltypes.AttributeKeySrcChannel {
							eventChannelId = attrValue
						}
					}
					if eventPortId != sourcePort || eventChannelId != sourceChannel {
						continue
					}
					txKey := fmt.Sprintf("%s/%s/%s", sourcePort, sourceChannel, sequence)
					if ibcTxs[txKey] {
						continue
					}
					var packet ibctransfertypes.FungibleTokenPacketData
					err = types.ModuleCdc.UnmarshalJSON(data, &packet)
					suite.Require().NoError(err)
					suite.Require().Equal(signer.AccAddress().String(), packet.Sender)
					suite.Require().Equal(randMint.String(), packet.Amount)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestFIP20CrossChainIBCExternal() {
	testCases := []struct {
		name     string
		malleate func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - ibc token",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				// add relay token
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				prefix, recipient := suite.RandPrefixAndAddress()

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					recipient,
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), nil
			},
			result: true,
		},
		{
			name: "ok - base token",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				// add relay token
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)

				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)

				sourcePort1, sourceChannel1 := suite.RandTransferChannel()
				prefix, recipient := suite.RandPrefixAndAddress()

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					recipient,
					randMint,
					big.NewInt(0),
					fxtypes.MustStrToByte32(fmt.Sprintf("ibc/%s/%s/%s", prefix, sourcePort1, sourceChannel1)),
				)
				suite.Require().NoError(err)

				return data, pair, big.NewInt(0), nil
			},
			result: true,
		},
		{
			name: "failed - no zero fee",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int, sourcePort, sourceChannel string) ([]byte, *types.TokenPair, *big.Int, []string) {
				coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(coin))
				suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)

				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)
				prefix, recipient := suite.RandPrefixAndAddress()

				fee := big.NewInt(int64(tmrand.Intn(1000) + 10))
				relayAmount := big.NewInt(0).Sub(randMint, fee)

				ibcTarget := fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(sourceChannel, ibcchanneltypes.ChannelPrefix), prefix)

				data, err := contract.GetFIP20().ABI.Pack(
					"transferCrossChain",
					recipient,
					relayAmount,
					fee,
					fxtypes.MustStrToByte32(ibcTarget),
				)
				suite.Require().NoError(err)

				aliasDenom := ""
				for _, alias := range md.metadata.DenomUnits[0].Aliases {
					if strings.Contains(alias, ibctransfertypes.DenomPrefix+"/") {
						hexHash := strings.TrimPrefix(alias, ibctransfertypes.DenomPrefix+"/")
						hash, err := ibctransfertypes.ParseHexHash(hexHash)
						suite.Require().NoError(err)

						denomTrace, found := suite.app.IBCTransferKeeper.GetDenomTrace(suite.ctx, hash)
						if !found {
							continue
						}
						if !strings.HasPrefix(denomTrace.GetPath(), fmt.Sprintf("%s/%s", ibctransfertypes.PortID, sourceChannel)) {
							continue
						}
						aliasDenom = alias
						break
					}
				}
				return data, pair, big.NewInt(0), []string{sdk.NewCoin(aliasDenom, sdkmath.NewIntFromBigInt(fee)).String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: fip-cross-chain failed: ibc transfer fee must be zero: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// set port channel
			sourcePort, sourceChannel := suite.RandTransferChannel()
			// add ibc token
			ibcToken := suite.AddIBCToken(sourcePort, sourceChannel)
			// token pair
			md := suite.GenerateCrossChainDenoms(ibcToken)

			// deploy fip20 external
			fip20External, err := suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, signer.Address(), "Test token", "TEST", 18)
			suite.Require().NoError(err)
			// token pair
			pair, err := suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, fip20External, md.GetMetadata().DenomUnits[0].Aliases...)
			suite.Require().NoError(err)

			suite.Require().NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 100000))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))

			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			balance := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
			suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

			commitments := suite.app.IBCKeeper.ChannelKeeper.GetAllPacketCommitmentsAtChannel(suite.ctx, sourcePort, sourceChannel)
			ibcTxs := make(map[string]bool, len(commitments))
			for _, commitment := range commitments {
				ibcTxs[fmt.Sprintf("%s/%s/%d", commitment.PortId, commitment.ChannelId, commitment.Sequence)] = true
			}

			packData, newPair, value, errArgs := tc.malleate(pair, md, signer, randMint, sourcePort, sourceChannel)

			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			tx, err := suite.PackEthereumTx(signer, newPair.GetERC20Contract(), value, packData)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			// check result
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
				suite.Require().True(chainBalances.IsZero(), chainBalances.String())
				balance := suite.BalanceOf(newPair.GetERC20Contract(), signer.Address())
				suite.Require().True(balance.Cmp(big.NewInt(0)) == 0, balance.String())

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
					// ibc transfer: erc20 token --> base denom
					// ibc base transfer: erc20 token --> base denom --> ibc token --> burn
					expect := totalBefore.Supply.AmountOf(coin.Denom)
					suite.Require().Equal(coin.Amount.String(), expect.String(), coin.Denom)
				}

				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type != ibcchanneltypes.EventTypeSendPacket {
						continue
					}
					var eventPortId, eventChannelId string
					var sequence string
					var data []byte

					for _, attr := range event.Attributes {
						attrKey, attrValue := string(attr.Key), string(attr.Value)
						if attrKey == ibcchanneltypes.AttributeKeyDataHex {
							data, err = hex.DecodeString(attrValue)
							suite.Require().NoError(err)
						}
						if attrKey == ibcchanneltypes.AttributeKeySequence {
							sequence = attrValue
						}
						if attrKey == ibcchanneltypes.AttributeKeySrcPort {
							eventPortId = attrValue
						}
						if attrKey == ibcchanneltypes.AttributeKeySrcChannel {
							eventChannelId = attrValue
						}
					}
					if eventPortId != sourcePort || eventChannelId != sourceChannel {
						continue
					}
					txKey := fmt.Sprintf("%s/%s/%s", sourcePort, sourceChannel, sequence)
					if ibcTxs[txKey] {
						continue
					}
					var packet ibctransfertypes.FungibleTokenPacketData
					err = types.ModuleCdc.UnmarshalJSON(data, &packet)
					suite.Require().NoError(err)
					suite.Require().Equal(signer.AccAddress().String(), packet.Sender)
					suite.Require().Equal(randMint.String(), packet.Amount)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestAccountFIP20CrossChain() {
	crossChainABI := crosschain.GetABI()
	otherABI := contract.MustABIJson(testJsonABI)

	testCases := []struct {
		name     string
		malleate func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "failed - call with address - method not found",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				data, err := otherABI.Pack(
					"fip20CrossChainV2",
					signer.Address(),
					signer.Address(),
					signer.Address().String(),
					big.NewInt(10),
					big.NewInt(0),
					fxtypes.MustStrToByte32("abc"),
					"",
				)
				suite.Require().NoError(err)

				return data, []string{}
			},
			error: func(args []string) string {
				return "unknown method"
			},
			result: false,
		},
		{
			name: "failed - call with address - pair not found",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				data, err := crossChainABI.Pack(
					crosschain.FIP20CrossChainMethod.Name,
					signer.Address(),
					signer.Address().String(),
					big.NewInt(10),
					big.NewInt(0),
					fxtypes.MustStrToByte32("abc"),
					"",
				)
				suite.Require().NoError(err)

				return data, []string{signer.Address().String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("token pair not found: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed - unpack error",
			malleate: func(pair *types.TokenPair, md Metadata, signer *helpers.Signer, randMint *big.Int) ([]byte, []string) {
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, types.NewTokenPair(signer.Address(), "abc", true, types.OWNER_MODULE))

				method := otherABI.Methods[crosschain.FIP20CrossChainMethod.Name]
				data, err := otherABI.Pack(
					crosschain.FIP20CrossChainMethod.Name,
					signer.Address(),
					signer.Address(),
					signer.Address().String(),
					big.NewInt(10),
					big.NewInt(0),
					"12",
				)
				suite.Require().NoError(err)

				dateTrimPrefix := bytes.TrimPrefix(data, method.ID)
				return append(crosschain.FIP20CrossChainMethod.ID, dateTrimPrefix...), []string{hex.EncodeToString(data)}
			},
			error: func(args []string) string {
				return "abi: cannot marshal in to go slice: offset"
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
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))

			packData, errArgs := tc.malleate(pair, md, signer, randMint)

			tx, err := suite.PackEthereumTx(signer, crosschain.GetAddress(), big.NewInt(0), packData)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			// check result
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.error(errArgs))
			}
		})
	}
}
