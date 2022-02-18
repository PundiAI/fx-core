package types

import "math/big"

// FIP20Data represents the ERC20 token details used to map
// the token to a Cosmos Coin
type FIP20Data struct {
	Name     string
	Symbol   string
	Decimals uint8
}

// FIP20StringResponse defines the string value from the call response
type FIP20StringResponse struct {
	Value string
}

// FIP20Uint8Response defines the uint8 value from the call response
type FIP20Uint8Response struct {
	Value uint8
}

// FIP20BoolResponse defines the bool value from the call response
type FIP20BoolResponse struct {
	Value bool
}

// NewFIP20Data creates a new FIP20Data instance
func NewFIP20Data(name, symbol string, decimals uint8) FIP20Data {
	return FIP20Data{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}
}

type FIP20Uint256Response struct {
	Value *big.Int
}
