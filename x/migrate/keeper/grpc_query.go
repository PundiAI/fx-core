package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	typescommon "github.com/functionx/fx-core/x/migrate/types/common"
	typesv1 "github.com/functionx/fx-core/x/migrate/types/v1"
)

var _ typesv1.QueryServer = Keeper{}

func (k Keeper) QueryMigrateRecord(ctx context.Context, req *typesv1.QueryMigrateRecordRequest) (*typesv1.MigrateRecord, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(typescommon.ErrInvalidAddress, err.Error())
	}
	record, _ := k.GetMigrateRecord(sdk.UnwrapSDKContext(ctx), addr)
	return &record, nil
}
