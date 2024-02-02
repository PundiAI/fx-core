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
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v7/x/staking/types"
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

	// 6. disable validator address
	if err = k.DisableValidatorAddress(ctx, valAddr); err != nil {
		return nil, err
	}

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

	// jailed by double sign, can't update consensus pubkey
	if validator.IsJailed() {
		if err := k.validateDoubleSign(ctx, validator); err != nil {
			return nil, err
		}
	}

	// authorized from address
	if !k.HasValidatorGrant(ctx, fromAddr, valAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "from address not authorized")
	}
	// check validator is updating consensus pubkey
	if k.HasConsensusPubKey(ctx, valAddr) || k.HasConsensusProcess(ctx, valAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "validator %s is updating consensus pubkey", msg.ValidatorAddress)
	}

	// pubkey and address
	newPubKey, err := k.validateAnyPubKey(ctx, msg.Pubkey)
	if err != nil {
		return nil, err
	}

	// iterator edit validator with new pubkey
	if err := k.iteratorEditValidator(ctx, validator, newPubKey); err != nil {
		return nil, err
	}

	// set validator new consensus pubkey
	if err = k.SetConsensusPubKey(ctx, valAddr, newPubKey); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditConsensusPubKey,
		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
		sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
		sdk.NewAttribute(types.AttributeKeyPubKey, newPubKey.String()),
	))

	return &types.MsgEditConsensusPubKeyResponse{}, nil
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
	if cp == nil || cp.Validator == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid consensus params")
	}
	pkType := pk.Type()
	hasKeyType := false
	for _, keyType := range cp.Validator.PubKeyTypes {
		if pkType == keyType {
			hasKeyType = true
			break
		}
	}
	if !hasKeyType {
		return nil, stakingtypes.ErrValidatorPubKeyTypeNotSupported.Wrapf("got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes)
	}
	return pk, nil
}

func (k Keeper) validateDoubleSign(ctx sdk.Context, validator stakingtypes.Validator) error {
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	info, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		return sdkerrors.ErrUnknownAddress.Wrapf("consensus %s not found", consAddr.String())
	}
	if info.JailedUntil.Equal(evidencetypes.DoubleSignJailEndTime) {
		return sdkerrors.ErrInvalidRequest.Wrapf("validator %s is jailed for double sign", validator.OperatorAddress)
	}
	return nil
}

func (k Keeper) iteratorEditValidator(ctx sdk.Context, validator stakingtypes.Validator, newPk cryptotypes.PubKey) error {
	newConsAddr := sdk.ConsAddress(newPk.Address())

	newPkFound := false
	totalUpdatePower := math.NewInt(validator.ConsensusPower(k.PowerReduction(ctx)))
	k.IteratorConsensusPubKey(ctx, func(valAddr sdk.ValAddress, pkBytes []byte) bool {
		var pk cryptotypes.PubKey
		if err := k.cdc.UnmarshalInterfaceJSON(pkBytes, &pk); err != nil {
			k.Logger(ctx).Error("failed to unmarshal pubKey", "validator", valAddr.String(), "err", err.Error())
			return false
		}
		if newConsAddr.Equals(sdk.ConsAddress(pk.Address())) {
			newPkFound = true
			return true
		}
		power := k.GetLastValidatorPower(ctx, valAddr)
		totalUpdatePower = totalUpdatePower.Add(math.NewInt(power))
		return false
	})
	if newPkFound { // new pk already exists
		return stakingtypes.ErrValidatorPubKeyExists.Wrapf("new consensus pubkey %s already exists", newConsAddr.String())
	}

	// iterate validator consensus process start
	k.IteratorConsensusProcess(ctx, types.ProcessStart, func(valAddr sdk.ValAddress, _ []byte) {
		// NOTE: not need check pk, already update validator consensus pubkey
		power := k.GetLastValidatorPower(ctx, valAddr)
		totalUpdatePower = totalUpdatePower.Add(math.NewInt(power))
	})

	// iterate validator consensus process end
	k.IteratorConsensusProcess(ctx, types.ProcessEnd, func(valAddr sdk.ValAddress, _ []byte) {
		// NOTE: not need check pk, already update validator consensus pubkey
		power := k.GetLastValidatorPower(ctx, valAddr)
		totalUpdatePower = totalUpdatePower.Add(math.NewInt(power))
	})

	totalPowerOneThird := k.GetLastTotalPower(ctx).QuoRaw(3) // less than 1/3 total power
	if totalUpdatePower.GTE(totalPowerOneThird) {
		return sdkerrors.ErrInvalidRequest.Wrapf("total update power %s more than 1/3 total power %s",
			totalUpdatePower.String(), totalPowerOneThird.String())
	}
	return nil
}
