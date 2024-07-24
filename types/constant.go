package types

import (
	stdlog "log"
	"math/big"
	"os"
	"path/filepath"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	Name          = "fxcore"
	AddressPrefix = "fx"
	EnvPrefix     = "FX"

	DefaultDenom = "FX"
	DenomUnit    = 18

	AddrLen = 20
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

func GetDefGasPrice() sdk.Coin {
	return sdk.NewCoin(DefaultDenom, sdk.NewInt(4_000).MulRaw(1e9))
}

func GetDefaultNodeHome() string {
	return defaultNodeHome
}

func SetConfig(isCosmosCoinType bool) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AddressPrefix, AddressPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator, AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	config.SetAddressVerifier(VerifyAddressFormat)
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

// VerifyAddressFormat verifies whether the address is compatible with Ethereum
func VerifyAddressFormat(bz []byte) error {
	if len(bz) == 0 {
		return sdkerrors.ErrUnknownAddress.Wrap("invalid address; cannot be empty")
	}
	if len(bz) != AddrLen {
		return sdkerrors.ErrUnknownAddress.Wrapf("invalid address length; got: %d, expect: %d", len(bz), AddrLen)
	}

	return nil
}
