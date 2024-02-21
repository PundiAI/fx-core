package types

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

const (
	AssetERC20 = "ERC20"
)

var FxcoreChainID = fxtypes.EIP155ChainID().String()

var (
	TypeString, _       = abi.NewType("string", "", nil)
	TypeBytes, _        = abi.NewType("bytes", "", nil)
	TypeUint256Array, _ = abi.NewType("uint256[]", "", nil)
)

var (
	assetTypeDecode = abi.Arguments{
		abi.Argument{Name: "_assetType", Type: TypeString},
		abi.Argument{Name: "_assetData", Type: TypeBytes},
	}

	erc20AssetDecode = abi.Arguments{
		abi.Argument{Name: "_tokens", Type: TypeBytes},
		abi.Argument{Name: "_amounts", Type: TypeUint256Array},
	}
)

func UnpackAssetType(asset string) (string, []byte, error) {
	assetBytes, err := hex.DecodeString(asset)
	if err != nil {
		return "", nil, err
	}
	values, err := assetTypeDecode.UnpackValues(assetBytes)
	if err != nil {
		return "", nil, err
	}
	if len(values) != 2 {
		return "", nil, fmt.Errorf("invalid length")
	}
	assetType, ok := values[0].(string)
	if !ok {
		return "", nil, fmt.Errorf("invalid type type")
	}
	assetData, ok := values[1].([]byte)
	if !ok {
		return "", nil, fmt.Errorf("invalid data type")
	}
	return assetType, assetData, nil
}

func UnpackERC20Asset(asset []byte) ([]common.Address, []*big.Int, error) {
	values, err := erc20AssetDecode.UnpackValues(asset)
	if err != nil {
		return nil, nil, err
	}
	if len(values) != 2 {
		return nil, nil, fmt.Errorf("invalid length")
	}
	tokenBytes, ok := values[0].([]byte)
	if !ok {
		return nil, nil, fmt.Errorf("invalid token type")
	}
	amounts, ok := values[1].([]*big.Int)
	if !ok {
		return nil, nil, fmt.Errorf("invalid amount type")
	}
	tokens := make([]common.Address, 0, len(amounts))
	for i := 0; i*common.AddressLength < len(tokenBytes); i++ {
		token := tokenBytes[i*common.AddressLength : (i+1)*common.AddressLength]
		tokens = append(tokens, common.BytesToAddress(token))
	}
	if len(tokens) != len(amounts) {
		return nil, nil, fmt.Errorf("token not match amount")
	}
	return tokens, amounts, nil
}

func MergeDuplicationERC20(tokens []common.Address, amounts []*big.Int) ([]common.Address, []*big.Int) {
	tokenArray := make([]common.Address, 0, len(tokens))
	amountArray := make([]*big.Int, 0, len(amounts))
	for i := 0; i < len(tokens); i++ {
		found := false
		j := 0
		for ; j < len(tokenArray); j++ {
			if tokens[i] == tokenArray[j] {
				found = true
				break
			}
		}
		if found {
			amountArray[j] = new(big.Int).Add(amountArray[j], amounts[i])
		} else {
			tokenArray = append(tokenArray, tokens[i])
			amountArray = append(amountArray, amounts[i])
		}
	}
	return tokenArray, amountArray
}

func PackERC20Asset(tokens [][]byte, amounts []*big.Int) ([]byte, error) {
	if len(tokens) != len(amounts) {
		return nil, fmt.Errorf("token not match amount")
	}
	tokenBytes := make([]byte, 0, common.AddressLength*len(tokens))
	for _, token := range tokens {
		tokenBytes = append(tokenBytes, token...)
	}
	assetData, err := erc20AssetDecode.Pack(tokenBytes, amounts)
	if err != nil {
		return nil, err
	}
	return assetTypeDecode.Pack("ERC20", assetData)
}

func MustDecodeMessage(message string) []byte {
	if len(message) == 0 {
		return []byte{}
	}
	bz, err := hex.DecodeString(message)
	if err != nil {
		panic(err)
	}
	return bz
}
