package client

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"google.golang.org/grpc/connectivity"
)

type SDKContext struct {
	client.Context
}

func (d SDKContext) Connect() {}

func (d SDKContext) GetState() connectivity.State {
	return connectivity.Connecting
}

func (d SDKContext) WaitForStateChange(ctx context.Context, s connectivity.State) bool {
	return true
}
