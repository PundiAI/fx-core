package grpc

import (
	"context"

	"github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc/connectivity"
)

type ClientConn interface {
	grpc.ClientConn
	// Connect begins connecting the StateChanger.
	Connect()
	// GetState returns the current state of the StateChanger.
	GetState() connectivity.State
	// WaitForStateChange returns true when the state becomes s, or returns
	// false if ctx is canceled first.
	WaitForStateChange(ctx context.Context, s connectivity.State) bool
}
