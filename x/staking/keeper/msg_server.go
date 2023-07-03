package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v5/x/staking/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) GrantPrivilege(goCtx context.Context, msg *types.MsgGrantPrivilege) (*types.MsgGrantPrivilegeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	fromAddr := sdk.MustAccAddressFromBech32(msg.FromAddress)
	pk, err := types.ProtoAnyToAccountPubKey(msg.ToPubkey)
	if err != nil {
		return nil, err
	}
	toAddress := sdk.AccAddress(pk.Address())

	// 1. validator
	if _, found := k.GetValidator(ctx, valAddr); !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "validator %s not found", msg.ValidatorAddress)
	}
	// 2. from authorized
	if !k.HasValidatorGrant(ctx, fromAddr, valAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "from address not authorized")
	}
	// 3. revoke old privilege
	if err = k.RevokeAuthorization(ctx, fromAddr, sdk.AccAddress(valAddr)); err != nil {
		return nil, err
	}
	// 4. grant new privilege
	genericGrant := []authz.Authorization{authz.NewGenericAuthorization(sdk.MsgTypeURL(&authz.MsgGrant{}))}
	if err = k.GrantAuthorization(ctx, toAddress, sdk.AccAddress(valAddr), genericGrant, types.GrantExpirationTime); err != nil {
		return nil, err
	}
	// 5. update validator operator
	k.UpdateValidatorOperator(ctx, valAddr, toAddress)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGrantPrivilege,
		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
		sdk.NewAttribute(types.AttributeKeyFrom, msg.FromAddress),
		sdk.NewAttribute(types.AttributeKeyTo, toAddress.String()),
	))
	return &types.MsgGrantPrivilegeResponse{}, nil
}
