package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var LPTokenKey = []byte{0xf0} // prefix for each key to a lp token

func GetLPTokenKey(valAddr sdk.ValAddress) []byte {
	return append(LPTokenKey, valAddr.Bytes()...)
}
