package contract

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	AssetERC20 = "ERC20"
)

var (
	AssetTypeDecode = abi.Arguments{
		abi.Argument{Name: "_assetType", Type: TypeString},
		abi.Argument{Name: "_assetData", Type: TypeBytes},
	}

	Erc20AssetDecode = abi.Arguments{
		abi.Argument{Name: "_tokens", Type: TypeBytes},
		abi.Argument{Name: "_amounts", Type: TypeUint256Array},
	}
)

func UnpackAssetType(asset []byte) (string, []byte, error) {
	values, err := AssetTypeDecode.UnpackValues(asset)
	if err != nil {
		return "", nil, err
	}
	if len(values) != 2 {
		return "", nil, fmt.Errorf("invalid asset length")
	}
	assetType, ok := values[0].(string)
	if !ok {
		return "", nil, fmt.Errorf("invalid asset type type")
	}
	assetData, ok := values[1].([]byte)
	if !ok {
		return "", nil, fmt.Errorf("invalid asset data type")
	}
	return assetType, assetData, nil
}

func UnpackERC20Asset(asset []byte) ([]common.Address, []*big.Int, error) {
	values, err := Erc20AssetDecode.UnpackValues(asset)
	if err != nil {
		return nil, nil, err
	}
	if len(values) != 2 {
		return nil, nil, fmt.Errorf("invalid asset length")
	}
	addrBytes, ok := values[0].([]byte)
	if !ok {
		return nil, nil, fmt.Errorf("invalid token address type")
	}
	amounts, ok := values[1].([]*big.Int)
	if !ok {
		return nil, nil, fmt.Errorf("invalid amount type")
	}

	addrLength := len(addrBytes) / len(amounts)
	tokenAddrs := make([]common.Address, 0)
	for i := 0; i*addrLength < len(addrBytes); i++ {
		contract := common.BytesToAddress(addrBytes[i*addrLength : (i+1)*addrLength])
		tokenAddrs = append(tokenAddrs, contract)
	}
	return tokenAddrs, amounts, nil
}

func PackERC20AssetWithType(tokenAddrs []common.Address, amounts []*big.Int) (string, error) {
	addrBytes := make([]byte, 0)
	for i := 0; i < len(tokenAddrs); i++ {
		addrBytes = append(addrBytes, tokenAddrs[i].Bytes()...)
	}
	assetData, err := Erc20AssetDecode.Pack(addrBytes, amounts)
	if err != nil {
		return "", err
	}
	pack, err := AssetTypeDecode.Pack(AssetERC20, assetData)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(pack), nil
}
