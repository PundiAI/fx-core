package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// FX defines the default coin denomination used in Ethermint in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM state transitions in Ethermint.
	FX string = "FX"

	// BaseDenomUnit defines the base denomination unit for Photons.
	// 1 FX = 1x10^{BaseDenomUnit} fx
	BaseDenomUnit = 18

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)

// NewFXCoinInt64 is a utility function that returns an "FX" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewFXCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(FX, amount)
}
