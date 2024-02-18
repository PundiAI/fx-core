package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v7/x/migrate/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) MigrateRecord(ctx context.Context, req *types.QueryMigrateRecordRequest) (*types.QueryMigrateRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}
	var addr []byte
	if common.IsHexAddress(req.Address) {
		addr = common.HexToAddress(req.Address).Bytes()
	} else {
		if acc, err := sdk.AccAddressFromBech32(req.Address); err == nil {
			addr = acc.Bytes()
		}
	}
	if len(addr) == 0 {
		return nil, status.Error(codes.InvalidArgument, "must be bech32 or hex address")
	}
	record, found := k.GetMigrateRecord(sdk.UnwrapSDKContext(ctx), addr)
	return &types.QueryMigrateRecordResponse{MigrateRecord: record, Found: found}, nil
}

func (k Keeper) MigrateCheckAccount(goCtx context.Context, req *types.QueryMigrateCheckAccountRequest) (*types.QueryMigrateCheckAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.From == "" || req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	from, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "form")
	}

	if !common.IsHexAddress(req.To) {
		return nil, status.Error(codes.InvalidArgument, "not hex address")
	}
	to := common.HexToAddress(req.To)

	ctx := sdk.UnwrapSDKContext(goCtx)
	// check migrated
	if k.HasMigrateRecord(ctx, from) {
		return nil, status.Errorf(codes.AlreadyExists, "address %s has been migrated", req.From)
	}
	if k.HasMigrateRecord(ctx, to.Bytes()) {
		return nil, status.Errorf(codes.AlreadyExists, "address %s has been migrated", req.To)
	}
	for _, m := range k.GetMigrateI() {
		if err := m.Validate(ctx, k.cdc, from, to); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return &types.QueryMigrateCheckAccountResponse{}, nil
}
