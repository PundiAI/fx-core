package gasprice

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v3/server/grpc/base/gasprice/legacy"
)

type Querier struct{}

var _ QueryServer = Querier{}

func (q Querier) GetGasPrice(c context.Context, _ *GetGasPriceRequest) (*GetGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var gasPrices sdk.Coins
	for _, coin := range ctx.MinGasPrices() {
		gasPrices = append(gasPrices, sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
	}
	return &GetGasPriceResponse{GasPrices: gasPrices}, nil
}

func RegisterGRPCGatewayRoutes(clientConn grpc.ClientConn, mux *runtime.ServeMux) {
	if err := RegisterQueryHandlerClient(context.Background(), mux, NewQueryClient(clientConn)); err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", "gas price", err.Error()))
	}
	if err := legacy.RegisterQueryHandlerClient(context.Background(), mux, legacy.NewQueryClient(clientConn)); err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", "legacy gas price", err.Error()))
	}
}

func QueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-prices",
		Short: "query node gas prices",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := NewQueryClient(clientCtx)
			res, err := queryClient.GetGasPrice(context.Background(), &GetGasPriceRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
