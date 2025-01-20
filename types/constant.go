package types

import (
	stdlog "log"
	"os"
	"path/filepath"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	Name          = "fxcore"
	AddressPrefix = "fx"
	EnvPrefix     = "FX"

	DefaultDenom  = "apundiai"
	DefaultSymbol = "PUNDIAI"
	DenomUnit     = 18

	LegacyFXDenom = "FX"
	FXDenom       = "fx"

	AddrLen = 20

	DefaultGasPrice = 5 * 1e9
)

// defaultNodeHome default home directories for the application daemon
var defaultNodeHome string

func init() {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	// votingPower = delegateToken / sdk.PowerReduction  --  sdk.TokensToConsensusPower(tokens Int)
	sdk.DefaultPowerReduction = sdkmath.NewInt(1e18)

	nodeHome := os.ExpandEnv("$FX_HOME")
	if len(nodeHome) > 0 {
		defaultNodeHome = nodeHome
		return
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		stdlog.Println("Failed to get home dir %2", err)
	}

	defaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

func GetDefMinGasPrices() sdk.Coin {
	return sdk.NewCoin(DefaultDenom, sdkmath.NewInt(DefaultGasPrice))
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

	if err := sdk.RegisterDenom(DefaultDenom, sdkmath.LegacyNewDecWithPrec(1, 18)); err != nil {
		panic(err)
	}
	if err := sdk.RegisterDenom(strings.ToLower(DefaultSymbol), sdkmath.LegacyOneDec()); err != nil {
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

func SwapAmount(amount sdkmath.Int) sdkmath.Int {
	return amount.QuoRaw(100)
}

func SwapDecAmount(amount sdkmath.LegacyDec) sdkmath.LegacyDec {
	return amount.QuoInt64(100)
}

func SwapCoin(coin sdk.Coin) sdk.Coin {
	if coin.Denom != FXDenom {
		return coin
	}
	coin.Amount = SwapAmount(coin.Amount)
	return coin
}
