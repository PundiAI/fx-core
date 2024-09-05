package precompile_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func TestAddPendingPoolRewardsABI(t *testing.T) {
	addPendingPoolRewards := precompile.NewAddPendingPoolRewardsMethod(nil)

	require.Equal(t, 4, len(addPendingPoolRewards.Method.Inputs))
	require.Equal(t, 1, len(addPendingPoolRewards.Method.Outputs))

	require.Equal(t, 5, len(addPendingPoolRewards.Event.Inputs))
}

func (suite *PrecompileTestSuite) TestAddPendingPoolRewards() {
	addRewardFee := big.NewInt(int64(tmrand.Uint32() + 10))
	testCases := []struct {
		name     string
		malleate func(moduleName string, pair *types.TokenPair, signer *helpers.Signer) (crosschaintypes.AddPendingPoolRewardArgs, error)
		result   bool
	}{
		{
			name: "success",
			malleate: func(moduleName string, _ *types.TokenPair, signer *helpers.Signer) (crosschaintypes.AddPendingPoolRewardArgs, error) {
				pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)

				coin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(addRewardFee))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)
				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), addRewardFee)

				return crosschaintypes.AddPendingPoolRewardArgs{
					Chain:  moduleName,
					TxID:   big.NewInt(1),
					Token:  pair.GetERC20Contract(),
					Reward: addRewardFee,
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid chain name",
			malleate: func(moduleName string, pair *types.TokenPair, _ *helpers.Signer) (crosschaintypes.AddPendingPoolRewardArgs, error) {
				chain := "123"
				return crosschaintypes.AddPendingPoolRewardArgs{
					Chain:  chain,
					TxID:   big.NewInt(1),
					Token:  pair.GetERC20Contract(),
					Reward: addRewardFee,
				}, fmt.Errorf("invalid module name: %s", chain)
			},
			result: false,
		},
		{
			name: "failed - invalid tx id",
			malleate: func(moduleName string, pair *types.TokenPair, _ *helpers.Signer) (crosschaintypes.AddPendingPoolRewardArgs, error) {
				return crosschaintypes.AddPendingPoolRewardArgs{
					Chain:  moduleName,
					TxID:   big.NewInt(0),
					Token:  pair.GetERC20Contract(),
					Reward: addRewardFee,
				}, errors.New("invalid tx id")
			},
			result: false,
		},
		{
			name: "failed - invalid bridge fee",
			malleate: func(moduleName string, pair *types.TokenPair, _ *helpers.Signer) (crosschaintypes.AddPendingPoolRewardArgs, error) {
				return crosschaintypes.AddPendingPoolRewardArgs{
					Chain:  moduleName,
					TxID:   big.NewInt(1),
					Token:  pair.GetERC20Contract(),
					Reward: big.NewInt(0),
				}, errors.New("invalid add reward")
			},
			result: false,
		},
		{
			name: "failed - tx id not found",
			malleate: func(moduleName string, _ *types.TokenPair, signer *helpers.Signer) (crosschaintypes.AddPendingPoolRewardArgs, error) {
				pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)

				coin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(addRewardFee))
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
				_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
					Coin:     coin,
					Receiver: signer.Address().Hex(),
					Sender:   signer.AccAddress().String(),
				})
				suite.Require().NoError(err)
				suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), addRewardFee)

				return crosschaintypes.AddPendingPoolRewardArgs{
					Chain:  moduleName,
					TxID:   big.NewInt(10),
					Token:  pair.GetERC20Contract(),
					Reward: addRewardFee,
				}, errors.New("not found pending record: invalid request")
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
			suite.SetCorsschainEnablePending(moduleName, true)

			amount := big.NewInt(0).Add(randMint, big.NewInt(1))
			coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(amount))
			helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
			_, err = suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertCoin{
				Coin:     coin,
				Receiver: signer.Address().Hex(),
				Sender:   signer.AccAddress().String(),
			})
			suite.Require().NoError(err)
			suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), amount)

			crossChainPack, err := precompile.NewCrossChainMethod(nil).PackInput(crosschaintypes.CrossChainArgs{
				Token:   pair.GetERC20Contract(),
				Receipt: helpers.GenExternalAddr(moduleName),
				Amount:  randMint,
				Fee:     big.NewInt(1),
				Target:  fxtypes.MustStrToByte32(moduleName),
				Memo:    "",
			})
			suite.Require().NoError(err)

			res := suite.EthereumTx(signer, crosschaintypes.GetAddress(), big.NewInt(0), crossChainPack)
			suite.Require().False(res.Failed(), res.VmError)

			pendingTxBefore, found := suite.CrossChainKeepers()[moduleName].GetPendingPoolTxById(suite.ctx, 1)
			suite.Require().True(found)

			addPendingPoolRewards, errResult := tc.malleate(moduleName, pair, signer)
			packData, err := precompile.NewAddPendingPoolRewardsMethod(nil).PackInput(addPendingPoolRewards)
			suite.Require().NoError(err)

			res = suite.EthereumTx(signer, crosschaintypes.GetAddress(), big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				pendingTxAfter, found := suite.CrossChainKeepers()[moduleName].GetPendingPoolTxById(suite.ctx, 1)
				suite.Require().True(found)

				suite.Require().Equal(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(addRewardFee)).String(),
					sdk.NewCoins(pendingTxAfter.Rewards...).Sub(sdk.NewCoins(pendingTxBefore.Rewards...)...).String())
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}

func TestNewAddPendingPoolRewardsEvent(t *testing.T) {
	method := precompile.NewAddPendingPoolRewardsMethod(nil)
	args := &crosschaintypes.AddPendingPoolRewardArgs{
		Chain:  "eth",
		TxID:   big.NewInt(1000),
		Token:  common.BytesToAddress([]byte{0x11}),
		Reward: big.NewInt(2000),
	}
	sender := common.BytesToAddress([]byte{0x1})

	data, topic, err := method.NewAddPendingPoolRewardsEvent(args, sender)
	require.NoError(t, err)

	expectedData := "000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000003e800000000000000000000000000000000000000000000000000000000000007d000000000000000000000000000000000000000000000000000000000000000036574680000000000000000000000000000000000000000000000000000000000"
	expectedTopic := []common.Hash{
		common.HexToHash("3afbebaebe58f01b574a31dcb1a2186714107461ff1efebbf3eef3aa79ced285"),
		common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
		common.HexToHash("0000000000000000000000000000000000000000000000000000000000000011"),
	}

	require.Equal(t, expectedData, hex.EncodeToString(data))
	require.Equal(t, expectedTopic, topic)
}
