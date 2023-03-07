package types

import (
	stdlog "log"
	"math/big"
	"os"
	"path/filepath"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Name          = "fxcore"
	AddressPrefix = "fx"
	EnvPrefix     = "FX"

	DefaultDenom = "FX"
	DenomUnit    = 18
)

// defaultNodeHome default home directories for the application daemon
var defaultNodeHome string

func init() {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	// votingPower = delegateToken / sdk.PowerReduction  --  sdk.TokensToConsensusPower(tokens Int)
	sdk.DefaultPowerReduction = sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil))

	fxHome := os.ExpandEnv("$FX_HOME")
	if len(fxHome) > 0 {
		defaultNodeHome = fxHome
		return
	}
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		stdlog.Println("Failed to get home dir %2", err)
	}

	defaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

func GetDefaultNodeHome() string {
	return defaultNodeHome
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
