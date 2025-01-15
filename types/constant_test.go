package types

import (
	"math"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Address(t *testing.T) {
	accAddress := sdk.AccAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	addr, err := sdk.AccAddressFromHexUnsafe("0102030405060708090a0b0c0d0e0f1011121314")
	require.NoError(t, err)
	assert.Equal(t, addr, accAddress)

	bech32, err := sdk.AccAddressFromBech32(accAddress.String())
	require.NoError(t, err)
	assert.Equal(t, bech32, accAddress)
}

func Test_NormalizeCoin(t *testing.T) {
	SetConfig(false)
	denom, err := sdk.GetBaseDenom()
	require.NoError(t, err)
	assert.Equal(t, DefaultDenom, denom)
	coin := sdk.NewCoin(DefaultDenom, sdkmath.NewInt(1e18))
	assert.Equal(t, coin, sdk.NormalizeCoin(coin))

	myCoin, err := sdk.ParseCoinNormalized("1" + strings.Repeat("0", 18) + DefaultDenom)
	require.NoError(t, err)
	assert.Equal(t, myCoin, sdk.NewCoin(DefaultDenom, sdkmath.NewInt(1e18)))

	myCoin, err = sdk.ParseCoinNormalized("1" + DefaultDenom)
	require.NoError(t, err)
	assert.Equal(t, myCoin, sdk.NewCoin(DefaultDenom, sdkmath.NewInt(1)))

	myCoin, err = sdk.ParseCoinNormalized("1" + strings.ToLower(DefaultSymbol))
	require.NoError(t, err)
	assert.Equal(t, myCoin, sdk.NewCoin(DefaultDenom, sdkmath.NewInt(1e18)))
}

func Test_PowerReduction(t *testing.T) {
	token := sdkmath.NewInt(math.MaxInt64).MulRaw(1e18)
	assert.Equal(t, int64(math.MaxInt64), sdk.TokensToConsensusPower(token, sdk.DefaultPowerReduction))
	assert.Panics(t, func() {
		sdk.TokensToConsensusPower(token.MulRaw(1e18).MulRaw(2), sdk.DefaultPowerReduction)
	})
}
