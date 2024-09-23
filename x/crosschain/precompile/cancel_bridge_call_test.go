package precompile_test

import (
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
	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func TestCancelPendingBridgeCallABI(t *testing.T) {
	cancelPendingBridgeCall := precompile.NewCancelPendingBridgeCallMethod(nil)

	require.Equal(t, 2, len(cancelPendingBridgeCall.Method.Inputs))
	require.Equal(t, 1, len(cancelPendingBridgeCall.Method.Outputs))

	require.Equal(t, 3, len(cancelPendingBridgeCall.Event.Inputs))
}

func (suite *PrecompileTestSuite) TestCancelPendingBridgeCall() {
	testCases := []struct {
		name     string
		malleate func(moduleName string) (crosschaintypes.CancelPendingBridgeCallArgs, error)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "success",
			malleate: func(moduleName string) (crosschaintypes.CancelPendingBridgeCallArgs, error) {
				return crosschaintypes.CancelPendingBridgeCallArgs{
					Chain: moduleName,
					TxID:  big.NewInt(1),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid chain name",
			malleate: func(moduleName string) (crosschaintypes.CancelPendingBridgeCallArgs, error) {
				chain := "123"
				return crosschaintypes.CancelPendingBridgeCallArgs{
					Chain: chain,
					TxID:  big.NewInt(1),
				}, fmt.Errorf("invalid module name: %s", chain)
			},
			result: false,
		},
		{
			name: "failed - invalid tx id",
			malleate: func(moduleName string) (crosschaintypes.CancelPendingBridgeCallArgs, error) {
				return crosschaintypes.CancelPendingBridgeCallArgs{
					Chain: moduleName,
					TxID:  big.NewInt(0),
				}, errors.New("invalid tx id")
			},
			result: false,
		},
		{
			name: "failed - tx id not found",
			malleate: func(moduleName string) (crosschaintypes.CancelPendingBridgeCallArgs, error) {
				txID := big.NewInt(10)
				return crosschaintypes.CancelPendingBridgeCallArgs{
					Chain: moduleName,
					TxID:  txID,
				}, fmt.Errorf("not found, nonce: %s: invalid", txID.String())
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.InitObservedBlockHeight()
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
			_, err = suite.app.Erc20Keeper.ConvertCoin(suite.ctx, &types.MsgConvertCoin{
				Coin:     coin,
				Receiver: signer.Address().Hex(),
				Sender:   signer.AccAddress().String(),
			})
			suite.Require().NoError(err)
			suite.ERC20Approve(signer, pair.GetERC20Contract(), crosschaintypes.GetAddress(), amount)

			args := crosschaintypes.BridgeCallArgs{
				DstChain: moduleName,
				Refund:   suite.signer.Address(),
				Tokens:   []common.Address{pair.GetERC20Contract()},
				Amounts:  []*big.Int{amount},
				To:       common.Address{},
				Data:     []byte{},
				Value:    big.NewInt(0),
				Memo:     []byte{},
			}
			bridgeCallPack, err := precompile.NewBridgeCallMethod(nil).PackInput(args)
			suite.Require().NoError(err)
			res := suite.EthereumTx(signer, crosschaintypes.GetAddress(), big.NewInt(0), bridgeCallPack)
			suite.Require().False(res.Failed(), res.VmError)

			balanceBefore := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())

			cancelArgs, errResult := tc.malleate(moduleName)
			packData, err := precompile.NewCancelPendingBridgeCallMethod(nil).PackInput(cancelArgs)
			suite.Require().NoError(err)

			res = suite.EthereumTx(signer, crosschaintypes.GetAddress(), big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				balanceAfter := suite.BalanceOf(pair.GetERC20Contract(), signer.Address())
				suite.Equal(big.NewInt(0).Add(balanceBefore, amount).String(), balanceAfter.String())
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
