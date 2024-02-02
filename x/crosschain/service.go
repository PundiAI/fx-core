package crosschain

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
)

func RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := crosschaintypes.RegisterQueryHandlerClient(context.Background(), mux, crosschaintypes.NewQueryClient(clientCtx))
	if err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", crosschaintypes.ModuleName, err.Error()))
	}
}
