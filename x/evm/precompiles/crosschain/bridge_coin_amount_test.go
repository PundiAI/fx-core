package crosschain_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/crosschain"
)

func TestBridgeCoinAmountABI(t *testing.T) {
	crosschainABI := crosschain.GetABI()

	method := crosschainABI.Methods[crosschain.BridgeCoinAmountMethodName]
	require.Equal(t, method, crosschain.BridgeCoinAmountMethod)
	require.Equal(t, 2, len(crosschain.BridgeCoinAmountMethod.Inputs))
	require.Equal(t, 1, len(crosschain.BridgeCoinAmountMethod.Outputs))
}

func (suite *PrecompileTestSuite) TestBridgeCoinAmount() {
	prepareFunc := func() (Metadata, *erc20types.TokenPair, *big.Int) {
		// token pair
		md := suite.GenerateCrossChainDenoms()
		pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, md.GetMetadata())
		suite.Require().NoError(err)
		randMint := big.NewInt(int64(tmrand.Uint32() + 10))
		suite.MintLockNativeTokenToModule(md.GetMetadata(), sdkmath.NewIntFromBigInt(randMint))
		return md, pair, randMint
	}
	testCases := []struct {
		name     string
		prepare  func() (Metadata, *erc20types.TokenPair, *big.Int)
		malleate func(token common.Address, target string) ([]byte, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name:    "ok",
			prepare: prepareFunc,
			malleate: func(token common.Address, target string) ([]byte, []string) {
				pack, err := crosschain.GetABI().Pack(crosschain.BridgeCoinAmountMethodName, token, fxtypes.MustStrToByte32(target))
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name:    "failed - invalid target",
			prepare: prepareFunc,
			malleate: func(token common.Address, target string) ([]byte, []string) {
				pack, err := crosschain.GetABI().Pack(crosschain.BridgeCoinAmountMethodName, token, fxtypes.MustStrToByte32(""))
				suite.Require().NoError(err)
				return pack, nil
			},
			error: func(args []string) string {
				return "empty target"
			},
			result: false,
		},
		{
			name:    "failed - invalid token",
			prepare: prepareFunc,
			malleate: func(_ common.Address, target string) ([]byte, []string) {
				token := helpers.GenerateAddress()
				pack, err := crosschain.GetABI().Pack(crosschain.BridgeCoinAmountMethodName, token, fxtypes.MustStrToByte32(target))
				suite.Require().NoError(err)
				return pack, []string{token.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("token not support: %s", args[0])
			},
			result: false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// prepare
			md, pair, randMint := tc.prepare()
			// malleate
			packData, errArgs := tc.malleate(pair.GetERC20Contract(), md.RandModule())
			tx, err := suite.PackEthereumTx(signer, crosschain.GetAddress(), big.NewInt(0), packData)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}
			// check result
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				unpack, err := crosschain.BridgeCoinAmountMethod.Outputs.Unpack(res.Ret)
				suite.Require().NoError(err)
				shares := unpack[0].(*big.Int)
				suite.Require().Equal(shares.String(), randMint.String())
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().EqualError(err, tc.error(errArgs))
				}
			}
		})
	}
}
