package telemetry

import (
	"fmt"
	"math"
	"math/big"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/hashicorp/go-metrics"
)

var maxFloat32, _ = big.NewInt(0).SetString(fmt.Sprintf("%.0f", math.MaxFloat32), 10)

func SetGaugeLabelsWithDenom(keys []string, denom string, amount *big.Int, labels ...metrics.Label) {
	if amount.Cmp(maxFloat32) == 1 {
		return
	}
	amountFloat32, _ := new(big.Float).SetInt(amount).Float32()
	telemetry.SetGaugeWithLabels(append(keys, denom), amountFloat32,
		append(labels, telemetry.NewLabel("denom", denom)))
}
