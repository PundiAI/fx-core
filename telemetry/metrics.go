package telemetry

import (
	"fmt"
	"math"
	"math/big"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var maxFloat32, _ = big.NewInt(0).SetString(fmt.Sprintf("%.0f", math.MaxFloat32), 10)

func SetGaugeLabelsWithToken(keys []string, token string, amount *big.Int, labels ...metrics.Label) {
	if amount.Cmp(maxFloat32) == 1 {
		return
	}
	amountFloat32, _ := new(big.Float).SetInt(amount).Float32()
	telemetry.SetGaugeWithLabels(append(keys, token), amountFloat32,
		append(labels, telemetry.NewLabel(LabelToken, token)))
}

func SetGaugeLabelsWithCoin(keys []string, coin sdk.Coin, labels ...metrics.Label) {
	SetGaugeLabelsWithToken(keys, coin.Denom, coin.Amount.BigInt(), labels...)
}

func SetGaugeLabelsWithCoins(keys []string, coins sdk.Coins, labels ...metrics.Label) {
	for _, coin := range coins {
		SetGaugeLabelsWithCoin(keys, coin, labels...)
	}
}
