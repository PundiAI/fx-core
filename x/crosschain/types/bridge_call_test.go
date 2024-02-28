package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
)

func TestBridgeCall_UnpackAssetType(t *testing.T) {
	asset := "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000005455243323000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"
	assetType, _, err := UnpackAssetType(asset)
	require.NoError(t, err)
	require.Equal(t, assetType, AssetERC20)
}

func testPackAsset(token [][]byte, amount []*big.Int) []byte {
	var tokenBz []byte
	for _, t := range token {
		tokenBz = append(tokenBz, t...)
	}
	pack, err := erc20AssetDecode.Pack(tokenBz, amount)
	if err != nil {
		panic(err)
	}
	return pack
}

func TestBridgeCall_UnpackERC20Asset(t *testing.T) {
	// ERC20
	assetBytes := testPackAsset([][]byte{common.BigToAddress(big.NewInt(1)).Bytes()}, []*big.Int{big.NewInt(1)})
	tokens, amounts, err := UnpackERC20Asset(assetBytes)
	require.NoError(t, err)
	require.Equal(t, common.BytesToAddress(tokens[0]).String(), common.BigToAddress(big.NewInt(1)).String())
	require.Equal(t, amounts[0].String(), "1")

	// Tron
	assetBytes = testPackAsset([][]byte{tronaddress.BigToAddress(big.NewInt(1))}, []*big.Int{big.NewInt(1)})
	tokens, amounts, err = UnpackERC20Asset(assetBytes)
	require.NoError(t, err)
	require.Equal(t, tokens[0], tronaddress.BigToAddress(big.NewInt(1)).Bytes())
	require.Equal(t, amounts[0].String(), "1")
}

func TestBridgeCall_MergeDuplicationERC20(t *testing.T) {
	tokens := [][]byte{
		common.BigToAddress(big.NewInt(1)).Bytes(),
		common.BigToAddress(big.NewInt(2)).Bytes(),
		common.BigToAddress(big.NewInt(3)).Bytes(),
		common.BigToAddress(big.NewInt(1)).Bytes(),
		common.BigToAddress(big.NewInt(3)).Bytes(),
	}
	expectedTokens := [][]byte{
		common.BigToAddress(big.NewInt(1)).Bytes(),
		common.BigToAddress(big.NewInt(2)).Bytes(),
		common.BigToAddress(big.NewInt(3)).Bytes(),
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
