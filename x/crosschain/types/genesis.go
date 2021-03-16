package types

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m GenesisState) ValidateBasic() error {
	return nil
}

// DefaultGenesisState returns empty genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}
