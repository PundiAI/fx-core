package types

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

var FxcoreChainID = fxtypes.EIP155ChainID().String()

func UnpackAssetType(asset string) (string, []byte, error) {
	assetBytes, err := hex.DecodeString(asset)
	if err != nil {
		return "", nil, err
	}
	return contract.UnpackAssetType(assetBytes)
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
