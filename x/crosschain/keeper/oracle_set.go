package keeper

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// --- ORACLE SET REQUESTS --- //

// GetCurrentOracleSet gets powers from the store and normalizes them
// into an integer percentage with a resolution of uint32 Max meaning
// a given validators 'gravity power' is computed as
// Cosmos power / total cosmos power = x / uint32 Max
// where x is the voting power on the gravity contract. This allows us
// to only use integer division which produces a known rounding error
// from truncation equal to the ratio of the validators
// Cosmos power / total cosmos power ratio, leaving us at uint32 Max - 1
// total voting power. This is an acceptable rounding error since floating
// point may cause consensus problems if different floating point unit
// implementations are involved.
func (k Keeper) GetCurrentOracleSet(ctx sdk.Context) *types.OracleSet {
	allOracles := k.GetAllOracles(ctx, true)
	bridgeValidators := make([]types.BridgeValidator, 0, len(allOracles))
	var totalPower uint64

	for _, oracle := range allOracles {
		power := oracle.GetPower()
		if power.LTE(sdkmath.ZeroInt()) {
			continue
		}
		totalPower += power.Uint64()
		bridgeValidators = append(bridgeValidators, types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: oracle.ExternalAddress,
		})
	}
	// normalize power values
	for i := range bridgeValidators {
		bridgeValidators[i].Power = sdkmath.NewUint(bridgeValidators[i].Power).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}

	oracleSetNonce := k.GetLatestOracleSetNonce(ctx) + 1
	return types.NewOracleSet(oracleSetNonce, uint64(ctx.BlockHeight()), bridgeValidators)
}

// AddOracleSetRequest returns a new instance of the Gravity BridgeValidatorSet
func (k Keeper) AddOracleSetRequest(ctx sdk.Context, currentOracleSet *types.OracleSet) {
	// if currentOracleSet member is empty, not store OracleSet.
	if len(currentOracleSet.Members) == 0 {
		return
	}
	k.StoreOracleSet(ctx, currentOracleSet)
	k.SetLatestOracleSetNonce(ctx, currentOracleSet.Nonce)

	k.CommonSetOracleTotalPower(ctx)

	// checkpoint, err := currentOracleSet.GetCheckpoint(k.GetGravityID(ctx))
	// if err != nil {
	// 	panic(err)
	// }
	// k.SetPastExternalSignatureCheckpoint(ctx, checkpoint)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOracleSetUpdate,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOracleSetNonce, fmt.Sprint(currentOracleSet.Nonce)),
		sdk.NewAttribute(types.AttributeKeyOracleSetLen, fmt.Sprint(len(currentOracleSet.Members))),
	))
}

// StoreOracleSet is for storing a oracle set at a given height
func (k Keeper) StoreOracleSet(ctx sdk.Context, oracleSet *types.OracleSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleSetKey(oracleSet.Nonce), k.cdc.MustMarshal(oracleSet))
}

// HasOracleSetRequest returns true if a oracleSet defined by a nonce exists
func (k Keeper) HasOracleSetRequest(ctx sdk.Context, nonce uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOracleSetKey(nonce))
}

// DeleteOracleSet deletes the oracleSet at a given nonce from state
func (k Keeper) DeleteOracleSet(ctx sdk.Context, nonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetOracleSetKey(nonce))
}

// SetLatestOracleSetNonce sets the latest oracleSet nonce
func (k Keeper) SetLatestOracleSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestOracleSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLatestOracleSetNonce returns the latest oracleSet nonce
func (k Keeper) GetLatestOracleSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LatestOracleSetNonce)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// GetOracleSet returns a oracleSet by nonce
func (k Keeper) GetOracleSet(ctx sdk.Context, nonce uint64) *types.OracleSet {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOracleSetKey(nonce))
	if bz == nil {
		return nil
	}
	var oracleSet types.OracleSet
	k.cdc.MustUnmarshal(bz, &oracleSet)
	return &oracleSet
}

// IterateOracleSets returns all oracleSet
func (k Keeper) IterateOracleSets(ctx sdk.Context, reverse bool, cb func(*types.OracleSet) bool) {
	store := ctx.KVStore(k.storeKey)
	var iter sdk.Iterator
	if reverse {
		iter = sdk.KVStoreReversePrefixIterator(store, types.OracleSetRequestKey)
	} else {
		iter = sdk.KVStorePrefixIterator(store, types.OracleSetRequestKey)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var oracleSet types.OracleSet
		k.cdc.MustUnmarshal(iter.Value(), &oracleSet)
		// cb returns true to stop early
		if cb(&oracleSet) {
			break
		}
	}
}

// GetOracleSets used in testing
func (k Keeper) GetOracleSets(ctx sdk.Context) (oracleSets types.OracleSets) {
	k.IterateOracleSets(ctx, false, func(set *types.OracleSet) bool {
		oracleSets = append(oracleSets, set)
		return false
	})
	return
}

// GetLatestOracleSet returns the latest oracle set in state
func (k Keeper) GetLatestOracleSet(ctx sdk.Context) *types.OracleSet {
	latestOracleSetNonce := k.GetLatestOracleSetNonce(ctx)
	return k.GetOracleSet(ctx, latestOracleSetNonce)
}

// SetLastSlashedOracleSetNonce sets the latest slashed oracleSet nonce
func (k Keeper) SetLastSlashedOracleSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedOracleSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLastSlashedOracleSetNonce returns the latest slashed oracleSet nonce
func (k Keeper) GetLastSlashedOracleSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LastSlashedOracleSetNonce)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// GetUnSlashedOracleSets returns all the unSlashed oracle sets in state
func (k Keeper) GetUnSlashedOracleSets(ctx sdk.Context, maxHeight uint64) (oracleSets types.OracleSets) {
	lastSlashedOracleSetNonce := k.GetLastSlashedOracleSetNonce(ctx) + 1
	k.IterateOracleSetByNonce(ctx, lastSlashedOracleSetNonce, func(oracleSet *types.OracleSet) bool {
		if maxHeight > oracleSet.Height {
			oracleSets = append(oracleSets, oracleSet)
			return false
		}
		return true
	})
	return
}

// IterateOracleSetByNonce iterates through all oracleSet by nonce
func (k Keeper) IterateOracleSetByNonce(ctx sdk.Context, startNonce uint64, cb func(*types.OracleSet) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.OracleSetRequestKey, sdk.Uint64ToBigEndian(startNonce)...)
	endKey := append(types.OracleSetRequestKey, sdk.Uint64ToBigEndian(math.MaxUint64)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		oracleSet := new(types.OracleSet)
		k.cdc.MustUnmarshal(iter.Value(), oracleSet)
		// cb returns true to stop early
		if cb(oracleSet) {
			break
		}
	}
}

// --- ORACLE SET CONFIRMS --- //

// GetOracleSetConfirm returns a oracleSet confirmation by a nonce and external address
func (k Keeper) GetOracleSetConfirm(ctx sdk.Context, nonce uint64, oracleAddr sdk.AccAddress) *types.MsgOracleSetConfirm {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetOracleSetConfirmKey(nonce, oracleAddr))
	if entity == nil {
		return nil
	}
	confirm := types.MsgOracleSetConfirm{}
	k.cdc.MustUnmarshal(entity, &confirm)
	return &confirm
}

// SetOracleSetConfirm sets a oracleSet confirmation
func (k Keeper) SetOracleSetConfirm(ctx sdk.Context, oracleAddr sdk.AccAddress, oracleSetConfirm *types.MsgOracleSetConfirm) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleSetConfirmKey(oracleSetConfirm.Nonce, oracleAddr)
	store.Set(key, k.cdc.MustMarshal(oracleSetConfirm))
}

// IterateOracleSetConfirmByNonce iterates through all oracleSet confirms by nonce
func (k Keeper) IterateOracleSetConfirmByNonce(ctx sdk.Context, nonce uint64, cb func(*types.MsgOracleSetConfirm) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetOracleSetConfirmKey(nonce, sdk.AccAddress{}))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgOracleSetConfirm)
		k.cdc.MustUnmarshal(iter.Value(), confirm)
		// cb returns true to stop early
		if cb(confirm) {
			break
		}
	}
}

func (k Keeper) DeleteOracleSetConfirm(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetOracleSetConfirmKey(nonce, sdk.AccAddress{}))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// GetLastObservedOracleSet retrieves the last observed oracle set from the store
// WARNING: This value is not an up to date oracle set on Ethereum, it is a oracle set
// that AT ONE POINT was the one in the bridge on Ethereum. If you assume that it's up
// to date you may break the bridge
func (k Keeper) GetLastObservedOracleSet(ctx sdk.Context) *types.OracleSet {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedOracleSetKey)

	if len(bytes) == 0 {
		return nil
	}
	oracleSet := types.OracleSet{}
	k.cdc.MustUnmarshal(bytes, &oracleSet)
	return &oracleSet
}

// SetLastObservedOracleSet updates the last observed oracle set in the store
func (k Keeper) SetLastObservedOracleSet(ctx sdk.Context, oracleSet *types.OracleSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedOracleSetKey, k.cdc.MustMarshal(oracleSet))
}
