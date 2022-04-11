package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewFXCoinInt64 is a utility function that returns an "FX" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewFXCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(MintDenom, amount)
}
