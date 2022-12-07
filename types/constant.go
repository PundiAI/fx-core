package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Name          = "fxcore"
	AddressPrefix = "fx"

	DefaultDenom = "FX"
	DenomUnit    = 18
)

func init() {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	// votingPower = delegateToken / sdk.PowerReduction  --  sdk.TokensToConsensusPower(tokens Int)
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil))
}

func SetConfig(isCosmosCoinType bool) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AddressPrefix, AddressPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator, AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	if isCosmosCoinType {
		config.SetCoinType(sdk.CoinType)
	} else {
		config.SetCoinType(60)
	}
	config.Seal()

	if err := sdk.RegisterDenom(DefaultDenom, sdk.NewDecWithPrec(1, 18)); err != nil {
		panic(err)
	}
}
