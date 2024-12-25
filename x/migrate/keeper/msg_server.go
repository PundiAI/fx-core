package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/x/migrate/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) MigrateAccount(goCtx context.Context, msg *types.MsgMigrateAccount) (*types.MsgMigrateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	toAddress := common.HexToAddress(msg.To)

	// check migrated
	if k.HasMigrateRecord(ctx, fromAddress) {
		return nil, sdkerrors.ErrInvalidSequence.Wrapf("address %s has been migrated", msg.From)
	}
	if k.HasMigrateRecord(ctx, toAddress.Bytes()) {
		return nil, sdkerrors.ErrInvalidSequence.Wrapf("address %s has been migrated", msg.To)
	}

	// check from address
	_, err = k.checkMigrateFrom(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	// migrate Validate
	for _, m := range k.GetMigrateI() {
		if err = m.Validate(ctx, k.cdc, fromAddress, toAddress); err != nil {
			return nil, err
		}
	}
	// migrate Execute
	for _, m := range k.GetMigrateI() {
		if err = m.Execute(ctx, k.cdc, fromAddress, toAddress); err != nil {
			return nil, err
		}
	}

	// set record
	k.SetMigrateRecord(ctx, fromAddress, toAddress)

	defer func() {
		telemetry.IncrCounter(1,
			types.ModuleName, sdk.MsgTypeURL(msg),
		)
	}()

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeMigrate,
		sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
		sdk.NewAttribute(types.AttributeKeyTo, msg.To),
	))
	return &types.MsgMigrateAccountResponse{}, nil
}
