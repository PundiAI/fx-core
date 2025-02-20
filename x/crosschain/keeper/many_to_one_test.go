package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func (suite *KeeperTestSuite) TestKeeper_SwapBridgeToken() {
	type args struct {
		name            string
		initBridgeToken func(*erc20types.BridgeToken)
		amount          sdkmath.Int
		err             error
		swapAmount      sdkmath.Int
	}
	tests := []args{
		{
			name: "swap fx token",
			initBridgeToken: func(token *erc20types.BridgeToken) {
				token.Denom = fxtypes.FXDenom
			},
			amount:     sdkmath.NewInt(100),
			err:        nil,
			swapAmount: sdkmath.NewInt(1),
		},
		{
			name: "swap origin token",
			initBridgeToken: func(token *erc20types.BridgeToken) {
				token.Denom = fxtypes.DefaultDenom
			},
			amount:     sdkmath.NewInt(100),
			err:        nil,
			swapAmount: sdkmath.NewInt(100),
		},
		{
			name: "swap bridge token",
			initBridgeToken: func(token *erc20types.BridgeToken) {
				token.Denom = "usdt"
			},
			amount:     sdkmath.NewInt(100),
			err:        nil,
			swapAmount: sdkmath.NewInt(100),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			bridgeToken := erc20types.BridgeToken{
				ChainName: suite.chainName,
				Contract:  helpers.GenExternalAddr(suite.chainName),
				Denom:     helpers.NewRandDenom(),
				IsNative:  false,
			}
			tt.initBridgeToken(&bridgeToken)
			err := suite.App.Erc20Keeper.AddBridgeToken(suite.Ctx, bridgeToken.Denom, bridgeToken.ChainName, bridgeToken.Contract, bridgeToken.IsNative)
			suite.Require().NoError(err)

			expBridgeToken := bridgeToken
			if bridgeToken.Denom == fxtypes.FXDenom {
				expBridgeToken = erc20types.BridgeToken{
					ChainName: suite.chainName,
					Contract:  helpers.GenExternalAddr(suite.chainName),
					Denom:     fxtypes.DefaultDenom,
					IsNative:  false,
				}
				err = suite.App.Erc20Keeper.AddBridgeToken(suite.Ctx, expBridgeToken.Denom, expBridgeToken.ChainName, expBridgeToken.Contract, expBridgeToken.IsNative)
				suite.Require().NoError(err)
				suite.MintTokenToModule(ethtypes.ModuleName, sdk.NewCoin(expBridgeToken.Denom, tt.amount))
			}

			from := helpers.GenAccAddress()
			suite.MintToken(from, sdk.NewCoin(bridgeToken.BridgeDenom(), tt.amount))
			bridgeToken, swapAmount, err := suite.Keeper().SwapBridgeToken(suite.Ctx, from, bridgeToken, tt.amount)
			if tt.err != nil {
				suite.Require().Error(err)
				suite.Require().Equal(tt.err.Error(), err.Error())
				return
			}
			suite.Require().NoError(err)
			suite.Equal(tt.swapAmount, swapAmount)
			suite.Equal(expBridgeToken, bridgeToken)
			suite.AssertBalance(from, sdk.NewCoin(expBridgeToken.BridgeDenom(), tt.swapAmount))
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_IBCCoinToEvm() {
	type args struct {
		name string
		init func() (holderAddr string, ibcCoin sdk.Coin)
		err  error
	}
	tests := []args{
		{
			name: "success origin coin",
			init: func() (string, sdk.Coin) {
				holder := helpers.GenHexAddress()
				return holder.String(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100))
			},
		},
		{
			name: "success not register ibc coin",
			init: func() (string, sdk.Coin) {
				return helpers.GenHexAddress().String(), sdk.NewCoin("ibc/coin", sdkmath.NewInt(100))
			},
		},
		{
			name: "success native coin",
			init: func() (string, sdk.Coin) {
				holder := helpers.GenHexAddress()
				symbol := helpers.NewRandSymbol()
				erc20Token, err := suite.App.Erc20Keeper.RegisterNativeCoin(suite.Ctx, symbol, symbol, 18)
				suite.Require().NoError(err)
				coin := sdk.NewCoin(erc20Token.Denom, sdkmath.NewInt(100))
				suite.MintToken(holder.Bytes(), coin)
				return holder.String(), coin
			},
		},
		{
			name: "success erc20 token",
			init: func() (string, sdk.Coin) {
				holder := helpers.GenHexAddress()
				symbol := helpers.NewRandSymbol()
				erc20TokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, suite.signer.Address(), symbol)
				erc20Token, err := suite.App.Erc20Keeper.RegisterNativeERC20(suite.Ctx, erc20TokenAddr)
				suite.Require().NoError(err)
				coin := sdk.NewCoin(erc20Token.Denom, sdkmath.NewInt(100))
				suite.MintToken(holder.Bytes(), coin)
				erc20ModuleAddr := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes())
				suite.erc20TokenSuite.WithContract(erc20TokenAddr).
					Mint(suite.Ctx, suite.signer.Address(), erc20ModuleAddr, big.NewInt(100))
				return holder.String(), coin
			},
		},
	}
	tests = append(tests, []args{
		{
			name: "success erc20 token with fx address and sequence",
			init: func() (string, sdk.Coin) {
				test := tests[3]
				suite.Require().Equal("success erc20 token", test.name)
				holderAddr, coin := test.init()
				hexAddr := common.HexToAddress(holderAddr)
				account := suite.App.AccountKeeper.GetAccount(suite.Ctx, hexAddr.Bytes())
				suite.Require().NoError(account.SetSequence(1))
				suite.App.AccountKeeper.SetAccount(suite.Ctx, account)
				return sdk.AccAddress(hexAddr.Bytes()).String(), coin
			},
		},
		{
			name: "success erc20 token with fx address and pubkey",
			init: func() (string, sdk.Coin) {
				test := tests[3]
				suite.Require().Equal("success erc20 token", test.name)
				holderAddr, coin := test.init()
				hexAddr := common.HexToAddress(holderAddr)
				account := suite.App.AccountKeeper.GetAccount(suite.Ctx, hexAddr.Bytes())
				suite.Require().NoError(account.SetSequence(1))
				suite.Require().NoError(account.SetPubKey(helpers.NewEthPrivKey().PubKey()))
				suite.App.AccountKeeper.SetAccount(suite.Ctx, account)
				return sdk.AccAddress(hexAddr.Bytes()).String(), coin
			},
		},
	}...)

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			holderAddr, ibcCoin := tt.init()
			err := suite.Keeper().IBCCoinToEvm(suite.Ctx, holderAddr, ibcCoin)
			if tt.err != nil {
				suite.Require().Error(err)
				suite.Require().Equal(tt.err.Error(), err.Error())
				return
			}
			suite.Require().NoError(err)
			erc20Token, err := suite.App.Erc20Keeper.GetERC20Token(suite.Ctx, ibcCoin.Denom)
			if err != nil {
				suite.Require().ErrorIs(err, collections.ErrNotFound)
			} else {
				addr, _, err := fxtypes.ParseAddress(holderAddr)
				suite.Require().NoError(err)
				balance := suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract()).
					BalanceOf(suite.Ctx, common.BytesToAddress(addr))
				suite.Equal(balance.String(), ibcCoin.Amount.String())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_DepositBridgeTokenToBaseCoin() {
	type args struct {
		name            string
		initBridgeToken func(*erc20types.BridgeToken)
		amount          sdkmath.Int
		err             error
	}
	tests := []args{
		{
			name:            "success bridge token",
			initBridgeToken: func(token *erc20types.BridgeToken) {},
			amount:          sdkmath.NewInt(100),
		},
		{
			name: "success origin coin",
			initBridgeToken: func(token *erc20types.BridgeToken) {
				token.Denom = fxtypes.DefaultDenom
			},
			amount: sdkmath.NewInt(100),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			bridgeToken := erc20types.BridgeToken{
				ChainName: suite.chainName,
				Contract:  helpers.GenExternalAddr(suite.chainName),
				Denom:     helpers.NewRandDenom(),
				IsNative:  false,
			}
			tt.initBridgeToken(&bridgeToken)
			err := suite.App.Erc20Keeper.AddBridgeToken(suite.Ctx, bridgeToken.Denom, bridgeToken.ChainName, bridgeToken.Contract, bridgeToken.IsNative)
			suite.Require().NoError(err)

			from := helpers.GenAccAddress()
			if bridgeToken.IsOrigin() {
				suite.MintTokenToModule(ethtypes.ModuleName, sdk.NewCoin(bridgeToken.BridgeDenom(), tt.amount))
			}
			baseCoin, err := suite.Keeper().DepositBridgeTokenToBaseCoin(suite.Ctx, from, tt.amount, bridgeToken.Contract)
			if tt.err != nil {
				suite.Require().Error(err)
				suite.Require().Equal(tt.err.Error(), err.Error())
				return
			}
			suite.Require().NoError(err)
			suite.Equal(sdk.NewCoin(bridgeToken.Denom, tt.amount), baseCoin)

			suite.AssertBalance(from, sdk.NewCoin(bridgeToken.Denom, tt.amount))
		})
	}
}

func TestIsEthSecp256k1(t *testing.T) {
	tests := []struct {
		name       string
		getAccount func() sdk.AccountI
		want       bool
	}{
		{
			name: "account is nil",
			getAccount: func() sdk.AccountI {
				return nil
			},
			want: false,
		},
		{
			name: "account pubkey is nil",
			getAccount: func() sdk.AccountI {
				return authtypes.NewBaseAccountWithAddress(helpers.GenAccAddress())
			},
			want: false,
		},
		{
			name: "account pubkey is nil and sequence > 0",
			getAccount: func() sdk.AccountI {
				account := authtypes.NewBaseAccountWithAddress(helpers.GenAccAddress())
				require.NoError(t, account.SetSequence(1))
				return account
			},
			want: true,
		},
		{
			name: "account pubkey is eth_secp256k1",
			getAccount: func() sdk.AccountI {
				account, err := authtypes.NewBaseAccountWithPubKey(helpers.NewEthPrivKey().PubKey())
				require.NoError(t, err)
				require.NoError(t, account.SetSequence(1))
				return account
			},
			want: true,
		},
		{
			name: "accunt pubkey is secp256k1",
			getAccount: func() sdk.AccountI {
				account, err := authtypes.NewBaseAccountWithPubKey(helpers.NewPriKey().PubKey())
				require.NoError(t, err)
				require.NoError(t, account.SetSequence(1))
				return account
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keeper.IsEthSecp256k1(tt.getAccount()); got != tt.want {
				t.Errorf("IsEthSecp256k1() = %v, want %v", got, tt.want)
			}
		})
	}
}
