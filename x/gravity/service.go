package gravity

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/functionx/fx-core/v7/x/gravity/types"
)

// RegisterGRPCGatewayRoutes
// Deprecated
func RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", types.ModuleName, err.Error()))
	}
}
