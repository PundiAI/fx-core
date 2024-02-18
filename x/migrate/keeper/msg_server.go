package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/x/migrate/types"
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
		return nil, errorsmod.Wrapf(types.ErrAlreadyMigrate, "address %s has been migrated", msg.From)
	}
	if k.HasMigrateRecord(ctx, toAddress.Bytes()) {
		return nil, errorsmod.Wrapf(types.ErrAlreadyMigrate, "address %s has been migrated", msg.To)
	}

	// check from address
	_, err = k.checkMigrateFrom(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	// migrate Validate
	for _, m := range k.GetMigrateI() {
		if err := m.Validate(ctx, k.cdc, fromAddress, toAddress); err != nil {
			return nil, errorsmod.Wrap(types.ErrMigrateValidate, err.Error())
		}
	}
	// migrate Execute
	for _, m := range k.GetMigrateI() {
		err := m.Execute(ctx, k.cdc, fromAddress, toAddress)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrMigrateExecute, err.Error())
		}
	}

	// set record
	k.SetMigrateRecord(ctx, fromAddress, toAddress)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, msg.Type()},
			1,
			[]metrics.Label{
				telemetry.NewLabel("from", msg.From),
				telemetry.NewLabel("to", msg.To),
			},
		)
	}()

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeMigrate,
		sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
		sdk.NewAttribute(types.AttributeKeyTo, msg.To),
	))
	return &types.MsgMigrateAccountResponse{}, nil
}
