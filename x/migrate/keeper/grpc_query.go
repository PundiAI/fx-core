package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/migrate/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) MigrateRecord(ctx context.Context, req *types.QueryMigrateRecordRequest) (*types.QueryMigrateRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	record, found := k.GetMigrateRecord(sdk.UnwrapSDKContext(ctx), addr)
	return &types.QueryMigrateRecordResponse{MigrateRecord: record, Found: found}, nil
}
