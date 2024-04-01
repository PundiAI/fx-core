package crosschain_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/crosschain"
)

func TestBridgeCallABI(t *testing.T) {
	crosschainABI := crosschain.GetABI()

	method := crosschainABI.Methods[crosschain.BridgeCallMethodName]
	require.Equal(t, method, crosschain.BridgeCallMethod)
	require.Equal(t, 7, len(crosschain.BridgeCallMethod.Inputs))
	require.Equal(t, 1, len(crosschain.BridgeCallMethod.Outputs))
}

func (suite *PrecompileTestSuite) TestBridgeCall() {
	testCases := []struct {
		name     string
		malleate func(tokenPair *erc20types.TokenPair, randMint *big.Int, moduleName string) []byte
		error    func(args []string) string
		result   bool
	}{
		{
			name: "pass",
			malleate: func(tokenPair *erc20types.TokenPair, randMint *big.Int, moduleName string) []byte {
				asset, err := contract.PackERC20AssetWithType([]common.Address{tokenPair.GetERC20Contract()}, []*big.Int{randMint})
				suite.NoError(err)
				assetBytes, err := hex.DecodeString(asset)
				suite.NoError(err)

				data, err := crosschain.GetABI().Pack(
					"bridgeCall",
					moduleName,
					big.NewInt(1000),
					helpers.GenerateAddress(),
					helpers.GenerateAddress(),
					[]byte{},
					big.NewInt(0),
					assetBytes,
				)
				suite.Require().NoError(err)
				return data
			},
			error: func(args []string) string {
				return ""
			},
			result: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// token pair
			md := suite.GenerateCrossChainDenoms()
			moduleName := md.RandModule()

			// deploy fip20 external
			fip20External, err := suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, signer.Address(), "Test token", "TEST", 18)
			suite.Require().NoError(err)
			// token pair
			pair, err := suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, fip20External, md.GetMetadata().DenomUnits[0].Aliases...)
			suite.Require().NoError(err)

			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))

			coin := sdk.NewCoin(pair.GetDenom(), sdkmath.NewIntFromBigInt(randMint))
			helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(coin))
			suite.MintERC20Token(signer, pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress(), randMint)

			_, err = suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), &erc20types.MsgConvertCoin{
				Coin:     coin,
				Receiver: signer.Address().Hex(),
				Sender:   signer.AccAddress().String(),
			})
			suite.Require().NoError(err)
			suite.CrossChainKeepers()[moduleName].SetLastObservedBlockHeight(suite.ctx, 100, uint64(suite.ctx.BlockHeight()))

			packData := tc.malleate(pair, randMint, moduleName)

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
			}
		})
	}
}
