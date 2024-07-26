package v1

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Querier
// Deprecated
type Querier struct{}

var _ QueryServer = Querier{}

// FxGasPrice
// Deprecated
func (q Querier) FxGasPrice(ctx context.Context, request *GasPriceRequest) (*GasPriceResponse, error) {
	return q.GasPrice(ctx, request)
}

// GasPrice
// Deprecated
func (q Querier) GasPrice(c context.Context, _ *GasPriceRequest) (*GasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var gasPrices sdk.Coins
	for _, coin := range ctx.MinGasPrices() {
		gasPrices = append(gasPrices, sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
	}
	return &GasPriceResponse{GasPrices: gasPrices}, nil
}

// RegisterGRPCGatewayRoutes
// Deprecated
func RegisterGRPCGatewayRoutes(clientConn grpc.ClientConn, mux *runtime.ServeMux) {
	if err := RegisterQueryHandlerClient(context.Background(), mux, NewQueryClient(clientConn)); err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", "gas price", err.Error()))
	}
}
