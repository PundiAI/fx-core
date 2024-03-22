package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/contract"
)

func TestBridgeCall_UnpackAssetType(t *testing.T) {
	asset := "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000005455243323000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"
	assetType, _, err := UnpackAssetType(asset)
	require.NoError(t, err)
	require.Equal(t, assetType, contract.AssetERC20)
}

func TestBridgeCall_UnpackERC20Asset(t *testing.T) {
	asset, err := contract.PackERC20AssetWithType(
		[]common.Address{
			common.BigToAddress(big.NewInt(1)),
			common.BigToAddress(big.NewInt(1)),
			common.BigToAddress(big.NewInt(1)),
		},
		[]*big.Int{
			big.NewInt(1),
			big.NewInt(0),
			big.NewInt(1),
		})
	require.NoError(t, err)
	_, assetBytes, err := UnpackAssetType(asset)
	require.NoError(t, err)
	tokenAddrs, amounts, err := contract.UnpackERC20Asset(assetBytes)
	require.NoError(t, err)
	require.Equal(t, tokenAddrs[0].String(), common.BigToAddress(big.NewInt(1)).String())
	require.Equal(t, amounts[0].String(), "1")
}
