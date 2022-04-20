package app

import (
	"math/big"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(fxtypes.AddressPrefix, fxtypes.AddressPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(fxtypes.AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator, fxtypes.AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(fxtypes.AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, fxtypes.AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	config.SetCoinType(60)
	config.Seal()

	// votingPower = delegateToken / sdk.PowerReduction  --  sdk.TokensToConsensusPower(tokens Int)
	sdk.PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil))

	if err := sdk.RegisterDenom(fxtypes.DefaultDenom, sdk.NewDec(18)); err != nil {
		panic(err)
	}
}
