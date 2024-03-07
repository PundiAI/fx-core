package types

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	assetData, err := erc20AssetDecode.Pack(addrBytes, amounts)
	if err != nil {
		return "", err
	}
	pack, err := assetTypeDecode.Pack(AssetERC20, assetData)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(pack), nil
}

func ExternalAddressToAccAddress(chainName, addr string) sdk.AccAddress {
	router, ok := msgValidateBasicRouter[chainName]
	if !ok {
		panic("unrecognized cross chain name")
	}
	accAddr, err := router.ExternalAddressToAccAddress(addr)
	if err != nil {
		panic(err)
	}
	return accAddr
}

func NewERC20Tokens(module string, tokenAddrs []common.Address, tokenAmounts []*big.Int) ([]ERC20Token, error) {
	if len(tokenAddrs) != len(tokenAmounts) {
		return nil, fmt.Errorf("invalid length")
	}
	tokens := make([]ERC20Token, 0)
	for i := 0; i < len(tokenAddrs); i++ {
		contract := fxtypes.AddressToStr(tokenAddrs[i].Bytes(), module)
		amount := sdkmath.NewIntFromBigInt(tokenAmounts[i])
		found := false
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
