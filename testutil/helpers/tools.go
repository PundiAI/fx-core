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

	fxtypes "github.com/functionx/fx-core/v7/types"
)

func IsLocalTest() bool {
	return os.Getenv("LOCAL_TEST") == "true"
}

func SkipTest(t *testing.T, msg ...any) {
	if !IsLocalTest() {
		t.Skip(msg...)
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
	expected, err := os.ReadFile(filePath)
	assert.NoError(t, err, filePath)

	var actual []byte
	switch res := result.(type) {
	case []byte:
		actual = res
	default:
		actual, err = json.MarshalIndent(result, "", " ")
		assert.NoError(t, err)
	}

	if !assert.JSONEqf(t, string(expected), string(actual), filePath) {
		assert.NoError(t, os.WriteFile(filePath, actual, 0o600))
	}
}

func NewStakingCoin(amount int64, power int64) sdk.Coin {
	powerBig := new(big.Int).Exp(big.NewInt(10), big.NewInt(power), nil)
	return sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(amount).Mul(sdk.NewIntFromBigInt(powerBig)))
}

func NewStakingCoins(amount int64, power int64) sdk.Coins {
	return sdk.NewCoins(NewStakingCoin(amount, power))
}
