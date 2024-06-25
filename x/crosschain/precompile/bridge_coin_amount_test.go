package precompile_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/precompile"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func TestBridgeCoinAmountMethod_ABI(t *testing.T) {
	bridgeCoinAmount := precompile.NewBridgeCoinAmountMethod(nil).Method
	assert.Equal(t, 2, len(bridgeCoinAmount.Inputs))
	assert.Equal(t, 1, len(bridgeCoinAmount.Outputs))
}

func (suite *PrecompileTestSuite) TestBridgeCoinAmount() {
	testCases := []struct {
		name     string
		malleate func(token common.Address, target string) (types.BridgeCoinAmountArgs, error)
		success  bool
	}{
		{
			name: "ok",
			malleate: func(token common.Address, target string) (types.BridgeCoinAmountArgs, error) {
				return types.BridgeCoinAmountArgs{
					Token:  token,
					Target: fxtypes.MustStrToByte32(target),
				}, nil
			},
			success: true,
		},
		{
			name: "failed - invalid target",
			malleate: func(token common.Address, target string) (types.BridgeCoinAmountArgs, error) {
				return types.BridgeCoinAmountArgs{
					Token:  token,
					Target: fxtypes.MustStrToByte32(""),
				}, errors.New("empty target: evm transaction execution failed")
			},
			success: false,
		},
		{
			name: "failed - invalid token",
			malleate: func(_ common.Address, target string) (types.BridgeCoinAmountArgs, error) {
				token := helpers.GenHexAddress()
				return types.BridgeCoinAmountArgs{
					Token:  token,
					Target: fxtypes.MustStrToByte32(target),
				}, fmt.Errorf("token not support: %s: evm transaction execution failed", token.String())
			},
			success: false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			signer := suite.RandSigner()
			bridgeCoinAmount := precompile.NewBridgeCoinAmountMethod(nil)

			md := suite.GenerateCrossChainDenoms()
			pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, md.GetMetadata())
			suite.Require().NoError(err)
			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))

			args, expectedErr := tc.malleate(pair.GetERC20Contract(), md.RandModule())
			packData, err := bridgeCoinAmount.PackInput(args)
			suite.Require().NoError(err)

			contractAddr := types.GetAddress()
			res, err := suite.app.EvmKeeper.CallEVMWithoutGas(suite.ctx, signer.Address(), &contractAddr, nil, packData, false)

			if tc.success {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				shares, err := bridgeCoinAmount.UnpackOutput(res.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal(shares.String(), randMint.String())
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().EqualError(err, expectedErr.Error())
				}
			}
		})
	}
}
