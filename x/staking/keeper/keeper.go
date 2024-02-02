package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	"github.com/functionx/fx-core/v7/x/staking/types"
)

type Keeper struct {
	stakingkeeper.Keeper
	storeKey storetypes.StoreKey
	cdc      codec.Codec

	accountKeeper  types.AccountKeeper
	authzKeeper    types.AuthzKeeper
	slashingKeeper types.SlashingKeeper
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(k stakingkeeper.Keeper, storeKey storetypes.StoreKey, cdc codec.Codec, ak types.AccountKeeper) Keeper {
	return Keeper{
		Keeper:        k,
		storeKey:      storeKey,
		cdc:           cdc,
		accountKeeper: ak,
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

// ValidatorGrant related functions

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

// ValidatorOperator related functions

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

func (k Keeper) DisableValidatorAddress(ctx sdk.Context, val sdk.ValAddress) error {
	account := k.accountKeeper.GetAccount(ctx, sdk.AccAddress(val))
	if account == nil || account.GetPubKey() == nil {
		return sdkerrors.ErrInvalidAddress.Wrap("invalid account")
	}
	pk := account.GetPubKey()
	var disablePK cryptotypes.PubKey
	disablePK = &ethsecp256k1.PubKey{Key: types.DisablePKBytes[:]}
	if pk.Type() == (&secp256k1.PubKey{}).Type() {
		disablePK = &secp256k1.PubKey{Key: types.DisablePKBytes[:]}
	}
	if err := account.SetPubKey(disablePK); err != nil {
		return err
	}
	k.accountKeeper.SetAccount(ctx, account)
	return nil
}

// ValidatorConsAddr related functions

func (k Keeper) SetValidatorConsAddr(ctx sdk.Context, newConsAddr sdk.ConsAddress, valOperator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(stakingtypes.GetValidatorByConsAddrKey(newConsAddr), valOperator)
}

func (k Keeper) RemoveValidatorConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(stakingtypes.GetValidatorByConsAddrKey(consAddr))
}

// ConsensusPubKey related functions

func (k Keeper) GetConsensusPubKey(ctx sdk.Context, valAddr sdk.ValAddress) (cryptotypes.PubKey, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetConsensusPubKey(valAddr))
	if bz == nil {
		return nil, false
	}
	var pubKey cryptotypes.PubKey
	if err := k.cdc.UnmarshalInterfaceJSON(bz, &pubKey); err != nil {
		return nil, false
	}
	return pubKey, true
}

func (k Keeper) HasConsensusPubKey(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetConsensusPubKey(valAddr))
}

func (k Keeper) SetConsensusPubKey(ctx sdk.Context, valAddr sdk.ValAddress, pubKey cryptotypes.PubKey) error {
	bz, err := k.cdc.MarshalInterfaceJSON(pubKey)
	if err != nil {
		return sdkerrors.ErrJSONMarshal.Wrapf("failed to marshal pubkey: %s", err.Error())
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetConsensusPubKey(valAddr), bz)
	return nil
}

func (k Keeper) RemoveConsensusPubKey(ctx sdk.Context, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetConsensusPubKey(valAddr))
}

func (k Keeper) IteratorConsensusPubKey(ctx sdk.Context, h func(valAddr sdk.ValAddress, pkBytes []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ConsensusPubKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		valAddr := sdk.ValAddress(types.AddressFromConsensusPubKey(iter.Key()))
		if h(valAddr, iter.Value()) {
			break
		}
	}
}

// ConsensusProcess related functions

func (k Keeper) GetConsensusProcess(ctx sdk.Context, valAddr sdk.ValAddress, process types.CProcess) (cryptotypes.PubKey, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetConsensusProcessKey(process, valAddr))
	if bz == nil {
		return nil, nil
	}
	var pubKey cryptotypes.PubKey
	if err := k.cdc.UnmarshalInterfaceJSON(bz, &pubKey); err != nil {
		return nil, err
	}
	return pubKey, nil
}

func (k Keeper) HasConsensusProcess(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetConsensusProcessKey(types.ProcessStart, valAddr)) ||
		store.Has(types.GetConsensusProcessKey(types.ProcessEnd, valAddr))
}

func (k Keeper) SetConsensusProcess(ctx sdk.Context, valAddr sdk.ValAddress, pubKey cryptotypes.PubKey, process types.CProcess) error {
	bz, err := k.cdc.MarshalInterfaceJSON(pubKey)
	if err != nil {
		return sdkerrors.ErrJSONMarshal.Wrapf("failed to marshal pubkey: %s", err.Error())
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetConsensusProcessKey(process, valAddr), bz)
	return nil
}

func (k Keeper) DeleteConsensusProcess(ctx sdk.Context, valAddr sdk.ValAddress, process types.CProcess) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetConsensusProcessKey(process, valAddr))
}

func (k Keeper) IteratorConsensusProcess(ctx sdk.Context, process types.CProcess, h func(valAddr sdk.ValAddress, pkBytes []byte)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, append(types.ConsensusProcessKey, process...))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		valAddr := sdk.ValAddress(types.AddressFromConsensusProcessKey(iter.Key()))
		h(valAddr, iter.Value())
	}
}
