package legacy

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Deprecated:
type Querier struct{}

var _ QueryServer = Querier{}

func (q Querier) FxGasPrice(ctx context.Context, request *GasPriceRequest) (*GasPriceResponse, error) {
	return q.GasPrice(ctx, request)
}

func (q Querier) GasPrice(c context.Context, _ *GasPriceRequest) (*GasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var gasPrices sdk.Coins
	for _, coin := range ctx.MinGasPrices() {
		gasPrices = append(gasPrices, sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
	}
	return &GasPriceResponse{GasPrices: gasPrices}, nil
}
