package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v2/x/other/types"
)

type Querier struct {
}

var _ types.QueryServer = Querier{}

func (q Querier) GasPrice(c context.Context, _ *types.GasPriceRequest) (*types.GasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var gasPrices sdk.Coins
	for _, coin := range ctx.MinGasPrices() {
		gasPrices = append(gasPrices, sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
	}
	return &types.GasPriceResponse{GasPrices: gasPrices}, nil
}
