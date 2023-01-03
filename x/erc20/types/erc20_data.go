package types

// ERC20Data represents the ERC20 token details used to map
// the token to a Cosmos Coin
type ERC20Data struct {
	Name     string
	Symbol   string
	Decimals uint8
}

// NewERC20Data creates a new ERC20Data instance
func NewERC20Data(name, symbol string, decimals uint8) ERC20Data {
	return ERC20Data{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}
}
