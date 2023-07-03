package keeper

import (
	"time"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/functionx/fx-core/v5/x/staking/types"
)

type Keeper struct {
	stakingkeeper.Keeper
	storeKey storetypes.StoreKey

	authzKeeper types.AuthzKeeper
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(k stakingkeeper.Keeper, storeKey storetypes.StoreKey) Keeper {
	return Keeper{
		Keeper:   k,
		storeKey: storeKey,
	}
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
