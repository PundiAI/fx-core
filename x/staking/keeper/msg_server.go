package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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

func (k Keeper) EditConsensusPubKey(goCtx context.Context, msg *types.MsgEditConsensusPubKey) (*types.MsgEditConsensusPubKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	fromAddr := sdk.MustAccAddressFromBech32(msg.From)

	if _, found := k.GetValidator(ctx, valAddr); !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "validator %s not found", msg.ValidatorAddress)
	}
	if !k.HasValidatorGrant(ctx, fromAddr, valAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "from address not authorized")
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", msg.Pubkey.GetCachedValue())
	}
	newConsAddr := sdk.GetConsAddress(pk)
	if _, found := k.GetValidatorByConsAddr(ctx, newConsAddr); found {
		return nil, stakingtypes.ErrValidatorPubKeyExists
	}

	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		pkType := pk.Type()
		hasKeyType := false
		for _, keyType := range cp.Validator.PubKeyTypes {
			if pkType == keyType {
				hasKeyType = true
				break
			}
		}
		if !hasKeyType {
			return nil, errorsmod.Wrapf(
				stakingtypes.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	// todo validator jailed/inactive, can update consensus pubkey?
	// todo one block can not update more than 1/3

	// set validator new consensus pubkey
	if err = k.SetValidatorNewConsensusPubKey(ctx, valAddr, pk); err != nil {
		return nil, err
	}

	// todo the update process(2 block) can be delegate(undelegate,redelegate)?

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditConsensusPubKey,
		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
		sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
		sdk.NewAttribute(types.AttributeKeyPubKey, newConsAddr.String()),
	))

	return &types.MsgEditConsensusPubKeyResponse{}, err
}
