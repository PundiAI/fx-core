package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	fxtype "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/migrate/types"
)

var (
	_ types.MsgServer = &Keeper{}
)

func (k Keeper) MigrateAccount(ctx context.Context, msg *types.MsgMigrateAccount) (*types.MsgMigrateAccountResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	//check module enable
	if sdkCtx.BlockHeight() < fxtype.MigrateSupportBlock() {
		return nil, sdkerrors.Wrap(types.InvalidRequest, "migrate module not enable")
	}

	fromAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	toAddress, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}
	//migrated
	if k.HasMigrateRecord(sdkCtx, fromAddress) {
		return nil, sdkerrors.Wrapf(types.ErrAlreadyMigrate, "address %s has been migrated", msg.From)
	}
	if k.HasMigrateRecord(sdkCtx, toAddress) {
		return nil, sdkerrors.Wrapf(types.ErrAlreadyMigrate, "address %s has been migrated", msg.To)
	}

	//migrate Validate
	for _, m := range k.GetMigrateI() {
		if err := m.Validate(sdkCtx, k, fromAddress, toAddress); err != nil {
			return nil, sdkerrors.Wrap(types.ErrMigrateValidate, err.Error())
		}
	}
	//migrate Execute
	for _, m := range k.GetMigrateI() {
		err := m.Execute(sdkCtx, k, fromAddress, toAddress)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrMigrateExecute, err.Error())
		}
	}
	//set record
	k.SetMigrateRecord(sdkCtx, fromAddress, toAddress)

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMigrate,
			sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
			sdk.NewAttribute(types.AttributeKeyTo, msg.To),
		),
	})
	return &types.MsgMigrateAccountResponse{}, nil
}
