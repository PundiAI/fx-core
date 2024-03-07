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
	assetBytes := testPackAsset(
		[][]byte{
			common.BigToAddress(big.NewInt(1)).Bytes(),
			common.BigToAddress(big.NewInt(1)).Bytes(),
			common.BigToAddress(big.NewInt(1)).Bytes(),
		},
		[]*big.Int{
			big.NewInt(1),
			big.NewInt(0),
			big.NewInt(1),
		})
	tokens, err := UnpackERC20Asset("eth", assetBytes)
	require.NoError(t, err)
	require.Equal(t, tokens[0].Contract, common.BigToAddress(big.NewInt(1)).String())
	require.Equal(t, tokens[0].Amount.String(), "2")

	// Tron
	assetBytes = testPackAsset(
		[][]byte{
			tronaddress.BigToAddress(big.NewInt(1)),
			tronaddress.BigToAddress(big.NewInt(2)),
			tronaddress.BigToAddress(big.NewInt(1)),
		},
		[]*big.Int{
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
		},
	)
	tokens, err = UnpackERC20Asset("tron", assetBytes)
	require.NoError(t, err)
	require.Equal(t, tokens[0].Contract, tronaddress.BigToAddress(big.NewInt(1)).String())
	require.Equal(t, tokens[0].Amount.String(), "2")
}
