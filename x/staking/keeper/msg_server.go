package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "validator %s not found", msg.ValidatorAddress)
	}

	// pubkey and address
	newPubKey, err := k.validateAnyPubKey(ctx, msg.Pubkey)
	if err != nil {
		return nil, err
	}

	// validator jailed/inactive, update pubkey
	if validator.IsJailed() || validator.IsUnbonding() || validator.IsUnbonded() {
		if err := k.updateValidatorPubKey(ctx, validator, newPubKey); err != nil {
			return nil, err
		}
		emitEditConsensusPubKeyEvents(ctx, valAddr, fromAddr, newPubKey)
		return &types.MsgEditConsensusPubKeyResponse{}, nil
	}

	// update validator less than 1/3 total power
	updatePower := math.NewInt(validator.ConsensusPower(k.PowerReduction(ctx)))
	k.IteratorConsensusPubKey(ctx, func(addr sdk.ValAddress, _ cryptotypes.PubKey) {
		power := k.GetLastValidatorPower(ctx, addr)
		updatePower = updatePower.Add(math.NewInt(power))
	})
	totalPowerOneThird := k.GetLastTotalPower(ctx).QuoRaw(3)
	if updatePower.GTE(totalPowerOneThird) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"update power %s more than 1/3 total power %s", updatePower.String(), totalPowerOneThird.String())
	}

	// from authorized
	if !k.HasValidatorGrant(ctx, fromAddr, valAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "from address not authorized")
	}

	// set validator new consensus pubkey
	if err = k.SetConsensusPubKey(ctx, valAddr, newPubKey); err != nil {
		return nil, err
	}

	// todo can delegate/undelegate/redelegate when process update(complete in 3 block)?

	emitEditConsensusPubKeyEvents(ctx, valAddr, fromAddr, newPubKey)

	return &types.MsgEditConsensusPubKeyResponse{}, err
}

func (k Keeper) validateAnyPubKey(ctx sdk.Context, pubkey *codectypes.Any) (cryptotypes.PubKey, error) {
	// pubkey type
	pk, ok := pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pubkey.GetCachedValue())
	}
	// pubkey exist
	newConsAddr := sdk.GetConsAddress(pk)
	if _, found := k.GetValidatorByConsAddr(ctx, newConsAddr); found {
		return nil, stakingtypes.ErrValidatorPubKeyExists
	}
	// pubkey type supported
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
	return pk, nil
}

func (k Keeper) updateValidatorPubKey(ctx sdk.Context, validator stakingtypes.Validator, newPubKey cryptotypes.PubKey) error {
	newConsAddr := sdk.ConsAddress(newPubKey.Address())
	oldConsAddr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	//  add new pubkey
	if err := k.slashingKeeper.AddPubkey(ctx, newPubKey); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
	}
	// remove old pubkey
	k.slashingKeeper.DeleteConsensusPubKey(ctx, oldConsAddr)

	// add new sign info
	info, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, oldConsAddr)
	if !found {
		return errorsmod.Wrap(sdkerrors.ErrUnknownAddress, "validator signing info not found")
	}
	info.Address = newConsAddr.String()
	k.slashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, info)

	// remove old sign info
	k.slashingKeeper.DeleteValidatorSigningInfo(ctx, oldConsAddr)

	// remove old cons address
	k.RemoveValidatorConsAddr(ctx, oldConsAddr)

	// update new cons address
	pkAny, _ := codectypes.NewAnyWithValue(newPubKey)
	validator.ConsensusPubkey = pkAny
	k.SetValidator(ctx, validator)
	k.SetValidatorConsAddr(ctx, newConsAddr, validator.GetOperator())
	return nil
}

func emitEditConsensusPubKeyEvents(ctx sdk.Context, val sdk.ValAddress, from sdk.AccAddress, newPubKey cryptotypes.PubKey) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditConsensusPubKey,
		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, val.String()),
		sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		sdk.NewAttribute(types.AttributeKeyPubKey, newPubKey.String()),
	))
}
