package helpers

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v8/types"
)

func IsLocalTest() bool {
	return os.Getenv("LOCAL_TEST") == "true"
}

func SkipTest(t *testing.T, msg ...any) {
	t.Helper()
	if !IsLocalTest() {
		t.Skip(append(msg, "#Please set env LOCAL_TEST=true#")...)
	}
}

func NewRandSymbol() string {
	return strings.ToUpper(fmt.Sprintf("a%sb", tmrand.Str(5)))
}

func NewRandDenom() string {
	return strings.ToLower(fmt.Sprintf("a%sb", tmrand.Str(5)))
}

func NewRandAmount() sdkmath.Int {
	return sdkmath.NewIntFromUint64(tmrand.Uint64() + 1)
}

func AssertJsonFile(t *testing.T, filePath string, result interface{}) {
	t.Helper()
	expected, err := os.ReadFile(filePath)
	require.NoError(t, err, filePath)

	actual, err := json.MarshalIndent(result, "", "  ")
	require.NoError(t, err)

	if !assert.JSONEqf(t, string(expected), string(actual), filePath) {
		require.NoError(t, os.WriteFile(filePath, actual, 0o600))
	}
}

func NewStakingCoin(amount, power int64) sdk.Coin {
	powerBig := new(big.Int).Exp(big.NewInt(10), big.NewInt(power), nil)
	return sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(amount).Mul(sdkmath.NewIntFromBigInt(powerBig)))
}

func NewStakingCoins(amount, power int64) sdk.Coins {
	return sdk.NewCoins(NewStakingCoin(amount, power))
}
