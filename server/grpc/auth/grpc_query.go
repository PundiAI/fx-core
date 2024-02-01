package auth

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/gogo/protobuf/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fxtypes "github.com/functionx/fx-core/v6/types"
)

var _ QueryServer = Querier{}

type Querier struct{}

func (Querier) ConvertAddress(_ context.Context, req *ConvertAddressRequest) (*ConvertAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if len(req.Address) == 0 {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}
	if len(req.Prefix) == 0 {
		req.Prefix = fxtypes.AddressPrefix
	}
	address, err := ConvertBech32Prefix(req.Address, req.Prefix)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &ConvertAddressResponse{Address: address}, nil
}

// ConvertBech32Prefix convert bech32 address to specified prefix.
func ConvertBech32Prefix(address, prefix string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", fmt.Errorf("cannot decode %s address: %w", address, err)
	}

	convertedAddress, err := bech32.ConvertAndEncode(prefix, bz)
	if err != nil {
		return "", fmt.Errorf("cannot convert %s address: %w", address, err)
	}
	return convertedAddress, nil
}

func RegisterGRPCGatewayRoutes(clientConn grpc.ClientConn, mux *runtime.ServeMux) {
	if err := RegisterQueryHandlerClient(context.Background(), mux, NewQueryClient(clientConn)); err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", "err", err.Error()))
	}
}
