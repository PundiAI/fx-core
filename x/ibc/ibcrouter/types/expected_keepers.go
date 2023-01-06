package types

import (
	"context"

	"github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

// TransferKeeper defines the expected transfer keeper
type TransferKeeper interface {
	Transfer(ctx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error)
}
