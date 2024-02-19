package types

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBridgeCall_UnpackAssetType(t *testing.T) {
	asset := "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000005455243323000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"
	assetType, _, err := UnpackAssetType(asset)
	require.NoError(t, err)
	require.Equal(t, assetType, AssetERC20)
}

func TestBridgeCall_UnpackERC20Asset(t *testing.T) {
	asset := "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"
	assetBytes, err := hex.DecodeString(asset)
	require.NoError(t, err)
	tokens, amounts, err := UnpackERC20Asset(assetBytes)
	require.NoError(t, err)
	require.Equal(t, tokens[0].String(), common.BigToAddress(big.NewInt(1)).String())
	require.Equal(t, amounts[0].String(), "1")
}

func TestBridgeCall_MergeDuplicationERC20(t *testing.T) {
	tokens := []common.Address{
		common.BigToAddress(big.NewInt(1)),
		common.BigToAddress(big.NewInt(2)),
		common.BigToAddress(big.NewInt(3)),
		common.BigToAddress(big.NewInt(1)),
		common.BigToAddress(big.NewInt(3)),
	}
	expectedTokens := []common.Address{
		common.BigToAddress(big.NewInt(1)),
		common.BigToAddress(big.NewInt(2)),
		common.BigToAddress(big.NewInt(3)),
	}
	amounts := []*big.Int{
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(1),
	}
	expectedAmounts := []*big.Int{
		big.NewInt(2),
		big.NewInt(1),
		big.NewInt(2),
	}

	newTokens, newAmounts := MergeDuplicationERC20(tokens, amounts)

	require.Equal(t, expectedTokens, newTokens)
	require.Equal(t, expectedAmounts, newAmounts)
}
