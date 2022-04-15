package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/migrate/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryMigrateRecord(ctx context.Context, req *types.QueryMigrateRecordRequest) (*types.MigrateRecord, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	record, _ := k.GetMigrateRecord(sdk.UnwrapSDKContext(ctx), addr)
	return &record, nil
}
