package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	fxtype "github.com/functionx/fx-core/types"
	typescommon "github.com/functionx/fx-core/x/migrate/types/common"
	typesv1 "github.com/functionx/fx-core/x/migrate/types/v1"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) *msgServer {
	return &msgServer{Keeper: keeper}
}

var (
	_ typesv1.MsgServer = msgServer{}
)

func (k msgServer) MigrateAccount(ctx context.Context, msg *typesv1.MsgMigrateAccount) (*typesv1.MsgMigrateAccountResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	//check module enable
	if sdkCtx.BlockHeight() < fxtype.MigrateSupportBlock() {
		return nil, sdkerrors.Wrap(typescommon.InvalidRequest, "migrate module not enable")
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
		return nil, sdkerrors.Wrapf(typescommon.ErrAlreadyMigrate, "address %s has been migrated", msg.From)
	}
	if k.HasMigrateRecord(sdkCtx, toAddress) {
		return nil, sdkerrors.Wrapf(typescommon.ErrAlreadyMigrate, "address %s has been migrated", msg.To)
	}

	//migrate Validate
	for _, m := range k.GetMigrateI() {
		if err := m.Validate(sdkCtx, k.Keeper, fromAddress, toAddress); err != nil {
			return nil, sdkerrors.Wrap(typescommon.ErrMigrateValidate, err.Error())
		}
	}
	//migrate Execute
	for _, m := range k.GetMigrateI() {
		err := m.Execute(sdkCtx, k.Keeper, fromAddress, toAddress)
		if err != nil {
			return nil, sdkerrors.Wrap(typescommon.ErrMigrateExecute, err.Error())
		}
	}
	//set record
	k.SetMigrateRecord(sdkCtx, fromAddress, toAddress)

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			typescommon.EventTypeMigrate,
			sdk.NewAttribute(typescommon.AttributeKeyFrom, msg.From),
			sdk.NewAttribute(typescommon.AttributeKeyTo, msg.To),
		),
	})
	return &typesv1.MsgMigrateAccountResponse{}, nil
}
