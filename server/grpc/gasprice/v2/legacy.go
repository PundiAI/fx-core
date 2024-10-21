package v2

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Deprecated: Querier
type Querier struct{}

var _ QueryServer = Querier{}

// Deprecated: GetGasPrice
func (q Querier) GetGasPrice(c context.Context, _ *GetGasPriceRequest) (*GetGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var gasPrices sdk.Coins
	for _, coin := range ctx.MinGasPrices() {
		gasPrices = append(gasPrices, sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
	}
	return &GetGasPriceResponse{GasPrices: gasPrices}, nil
}

// Deprecated: RegisterGRPCGatewayRoutes
func RegisterGRPCGatewayRoutes(clientConn grpc.ClientConn, mux *runtime.ServeMux) {
	if err := RegisterQueryHandlerClient(context.Background(), mux, NewQueryClient(clientConn)); err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", "legacy gas price", err.Error()))
	}
}
