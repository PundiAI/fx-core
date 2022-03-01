package types

import "fmt"

func NewGenesisState(params Params, pairs []TokenPair) GenesisState {
	return GenesisState{
		Params:     params,
		TokenPairs: pairs,
	}
}

func (gs GenesisState) Validate() error {
	seenFip20 := make(map[string]bool)
	seenDenom := make(map[string]bool)

	for _, b := range gs.TokenPairs {
		if seenFip20[b.Fip20Address] {
			return fmt.Errorf("token FIP20 contract duplicated on genesis '%s'", b.Fip20Address)
		}
		if seenDenom[b.Denom] {
			return fmt.Errorf("coin denomination duplicated on genesis: '%s'", b.Denom)
		}

		if err := b.Validate(); err != nil {
			return err
		}

		seenFip20[b.Fip20Address] = true
		seenDenom[b.Denom] = true
	}

	return gs.Params.Validate()
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}
