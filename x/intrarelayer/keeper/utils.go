package keeper

import (
	"fmt"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func equalMetadata(a, b banktypes.Metadata) error {
	if a.Base == b.Base && a.Description == b.Description && a.Display == b.Display {
		if len(a.DenomUnits) != len(b.DenomUnits) {
			return fmt.Errorf("metadata provided has different denom units from stored, %d ≠ %d", len(a.DenomUnits), len(b.DenomUnits))
		}
		for i, v := range a.DenomUnits {
			if v.Denom != b.DenomUnits[i].Denom {
				return fmt.Errorf("metadata provided has different denom from stored, %s ≠ %s", a.DenomUnits[i].Denom, b.DenomUnits[i].Denom)
			}
			if v.Exponent != b.DenomUnits[i].Exponent {
				return fmt.Errorf("metadata provided has different denom exponent from stored, %s(%d) ≠ %s(%d)",
					a.DenomUnits[i].Denom, a.DenomUnits[i].Exponent, b.DenomUnits[i].Denom, b.DenomUnits[i].Exponent)
			}
			if !equalStringArray(v.Aliases, b.DenomUnits[i].Aliases) {
				return fmt.Errorf("metadata provided has different denom aliases from stored, %s(%s) ≠ %s(%s)",
					a.DenomUnits[i].Denom, a.DenomUnits[i].Aliases, b.DenomUnits[i].Denom, b.DenomUnits[i].Aliases)
			}
		}
		return nil
	}
	return fmt.Errorf("metadata provided is different from stored")
}

func equalStringArray(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for k, v := range s1 {
		if v != s2[k] {
			return false
		}
	}
	return true
}
