package types

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi"

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

func UnpackERC20Asset(module string, asset []byte) ([]ERC20Token, error) {
	values, err := erc20AssetDecode.UnpackValues(asset)
	if err != nil {
		return nil, err
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("invalid length")
	}
	tokenBytes, ok := values[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid token type")
	}
	amounts, ok := values[1].([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid amount type")
	}
	tokens := make([]ERC20Token, 0)
	addrLength := len(tokenBytes) / len(amounts)
	for i := 0; i*addrLength < len(tokenBytes); i++ {
		var found bool
		contract := fxtypes.AddressToStr(tokenBytes[i*addrLength:(i+1)*addrLength], module)
		amount := sdkmath.NewIntFromBigInt(amounts[i])
		for j := 0; j < len(tokens); j++ {
			if contract == tokens[j].Contract {
				tokens[j].Amount = tokens[j].Amount.Add(amount)
				found = true
				break
			}
		}
		if !found {
			tokens = append(tokens, ERC20Token{
				Contract: contract,
				Amount:   amount,
			})
		}
	}
	return tokens, nil
}

func PackERC20Asset(tokens [][]byte, amounts []*big.Int) ([]byte, error) {
	if len(tokens) != len(amounts) {
		return nil, fmt.Errorf("token not match amount")
	}
	tokenBytes := make([]byte, 0)
	for _, token := range tokens {
		tokenBytes = append(tokenBytes, token...)
	}
	assetData, err := erc20AssetDecode.Pack(tokenBytes, amounts)
	if err != nil {
		return nil, err
	}
	return assetTypeDecode.Pack("ERC20", assetData)
}
