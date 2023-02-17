package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ValidatorLPTokenKey = []byte{0xf0} // prefix for val to lp token
	LPTokenValidatorKey = []byte{0xf1} // prefix for lp token to val
)

func GetValidatorLPTokenKey(validator sdk.ValAddress) []byte {
	return append(ValidatorLPTokenKey, validator.Bytes()...)
}

func GetLPTokenValidatorKey(lpTokenContract common.Address) []byte {
	return append(LPTokenValidatorKey, lpTokenContract.Bytes()...)
}
