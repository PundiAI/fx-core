package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v5/x/staking/types"
)

type Keeper struct {
	stakingkeeper.Keeper
	storeKey storetypes.StoreKey
	cdc      codec.Codec

	authzKeeper    types.AuthzKeeper
	slashingKeeper types.SlashingKeeper
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(k stakingkeeper.Keeper, storeKey storetypes.StoreKey, cdc codec.Codec) Keeper {
	return Keeper{
		Keeper:   k,
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (k *Keeper) SetSlashingKeeper(slashingKeeper types.SlashingKeeper) *Keeper {
	k.slashingKeeper = slashingKeeper
	return k
}

func (k *Keeper) SetAuthzKeeper(authzKeeper types.AuthzKeeper) *Keeper {
	k.authzKeeper = authzKeeper
	return k
}

func (k Keeper) HasValidatorGrant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.ValAddress) bool {
	operator, found := k.GetValidatorOperator(ctx, granter)
	if !found {
		return granter.Equals(grantee)
	}
	return operator.Equals(grantee)
}

func (k Keeper) RevokeAuthorization(ctx sdk.Context, grantee, granter sdk.AccAddress) error {
	authorizations, err := k.authzKeeper.GetAuthorizations(ctx, grantee, granter)
	if err != nil {
		return authz.ErrNoAuthorizationFound.Wrap(err.Error())
	}
	for _, a := range authorizations {
		if err = k.authzKeeper.DeleteGrant(ctx, grantee, granter, a.MsgTypeURL()); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) GrantAuthorization(ctx sdk.Context, grantee, granter sdk.AccAddress, auths []authz.Authorization, exp time.Time) error {
	for _, a := range auths {
		if err := k.authzKeeper.SaveGrant(ctx, grantee, granter, a, &exp); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) HasValidatorOperator(ctx sdk.Context, val sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetValidatorOperatorKey(val))
}

func (k Keeper) GetValidatorOperator(ctx sdk.Context, val sdk.ValAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorOperatorKey(val))
	if bz == nil {
		return nil, false
	}
	return bz, true
}

func (k Keeper) UpdateValidatorOperator(ctx sdk.Context, val sdk.ValAddress, from sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	// remove old operator
	store.Delete(types.GetValidatorOperatorKey(val))
	// add new operator
	store.Set(types.GetValidatorOperatorKey(val), from.Bytes())
}

func (k Keeper) SetValidatorOperatorByConsAddr(ctx sdk.Context, newConsAddr sdk.ConsAddress, valOperator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(stakingtypes.GetValidatorByConsAddrKey(newConsAddr), valOperator)
}

func (k Keeper) RemoveValidatorOperatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(stakingtypes.GetValidatorByConsAddrKey(consAddr))
}

func (k Keeper) SetValidatorOldConsensusAddr(ctx sdk.Context, valAddr sdk.ValAddress, newConsAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorOldConsensusAddrKey(valAddr), newConsAddr)
}

func (k Keeper) RemoveValidatorOldConsensusAddr(ctx sdk.Context, valAddrs sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorOldConsensusAddrKey(valAddrs))
}

func (k Keeper) IteratorValidatorOldConsensusAddr(ctx sdk.Context, handler func(valAddr sdk.ValAddress, consAddr sdk.ConsAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.ValidatorOldConsensusAddrKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		valAddr := sdk.ValAddress(types.AddressFromValidatorNewConsensusAddrKey(iter.Key()))
		consAddr := sdk.ConsAddress(iter.Value())

		if handler(valAddr, consAddr) {
			break
		}
	}
}

func (k Keeper) SetValidatorDelConsensusAddr(ctx sdk.Context, valAddr sdk.ValAddress, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorDelConsensusAddrKey(valAddr), consAddr)
}

func (k Keeper) RemoveValidatorDelConsensusAddr(ctx sdk.Context, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorDelConsensusAddrKey(valAddr))
}

func (k Keeper) IteratorValidatorDelConsensusAddr(ctx sdk.Context, handler func(valAddr sdk.ValAddress, address sdk.ConsAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.ValidatorDelConsensusAddrKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		valAddr := sdk.ValAddress(types.AddressFromValidatorDelConsensusAddrKey(iter.Key()))
		consAddr := sdk.ConsAddress(iter.Value())
		if handler(valAddr, consAddr) {
			break
		}
	}
}

func (k Keeper) GetValidatorNewConsensusPubKey(ctx sdk.Context, valAddr sdk.ValAddress) (cryptotypes.PubKey, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorNewConsensusPubKey(valAddr))
	if bz == nil {
		return nil, false
	}
	var pubKey cryptotypes.PubKey
	if err := k.cdc.UnmarshalInterfaceJSON(bz, &pubKey); err != nil {
		return nil, false
	}
	return pubKey, true
}

func (k Keeper) SetValidatorNewConsensusPubKey(ctx sdk.Context, valAddr sdk.ValAddress, pubKey cryptotypes.PubKey) error {
	bz, err := k.cdc.MarshalInterfaceJSON(pubKey)
	if err != nil {
		return sdkerrors.ErrJSONMarshal.Wrapf("failed to marshal pubkey: %s", err.Error())
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorNewConsensusPubKey(valAddr), bz)
	return nil
}

func (k Keeper) RemoveValidatorNewConsensusPubKey(ctx sdk.Context, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorNewConsensusPubKey(valAddr))
}

func (k Keeper) IteratorValidatorNewConsensusPubKey(ctx sdk.Context, handler func(valAddr sdk.ValAddress, pubKey cryptotypes.PubKey) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.ValidatorNewConsensusPubKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		valAddr := sdk.ValAddress(types.AddressFromValidatorNewConsensusPubKey(iter.Key()))

		var pk cryptotypes.PubKey
		if err := k.cdc.UnmarshalInterfaceJSON(iter.Value(), &pk); err != nil {
			k.Logger(ctx).Error("failed to unmarshal pubKey", "validator", valAddr.String(), "err", err.Error())
			continue
		}
		if handler(valAddr, pk) {
			break
		}
	}
}
