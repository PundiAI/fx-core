package keeper

import (
	"encoding/hex"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

func (k Keeper) Attest(ctx sdk.Context, claim types.ExternalClaim, anyClaim *codectypes.Any) (*types.Attestation, error) {
	oracle, found := k.GetOracleAddressByOrchestratorKey(ctx, claim.GetClaimer())
	if !found {
		panic("Could not find Oracle for delegate key, should be checked by now")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this oracle.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce.
	// This prevents there being two attestations with the same nonce that get 2/3s of the votes
	// in the endBlocker.
	lastEventNonce := k.GetLastEventNonceByOracle(ctx, oracle)
	if claim.GetEventNonce() != lastEventNonce+1 {
		return nil, sdkerrors.Wrapf(types.ErrNonContiguousEventNonce, "got %v, expected %v", claim.GetEventNonce(), lastEventNonce+1)
	}

	// Tries to get an attestation with the same eventNonce and claim as the claim that was submitted.
	att := k.GetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash())

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			Observed: false,
			Height:   uint64(ctx.BlockHeight()),
			Claim:    anyClaim,
		}
	}

	// Add the oracle's vote to this attestation
	att.Votes = append(att.Votes, oracle.String())

	k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)
	k.SetLastEventNonceByOracle(ctx, oracle, claim.GetEventNonce())
	k.setLastEventBlockHeightByOracle(ctx, oracle, claim.GetBlockHeight())

	return att, nil
}

// TryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryAttestation(ctx sdk.Context, att *types.Attestation) {
	claim, err := k.UnpackAttestationClaim(att)
	if err != nil {
		panic("could not cast to claim")
	}
	if att.Observed {
		// We panic here because this should never happen
		panic("attempting to process observed attestation")
	}
	logger := k.Logger(ctx)
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.
	// Sum the current powers of all validators who have voted and see if it passes the current threshold
	totalPower := k.GetLastTotalPower(ctx)
	requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
	attestationPower := sdk.NewInt(0)

	for _, oracleStr := range att.Votes {
		oracleAddr, err := sdk.AccAddressFromBech32(oracleStr)
		if err != nil {
			panic(fmt.Errorf("invalid oracle address %s", err.Error()))
		}
		oracle, found := k.GetOracle(ctx, oracleAddr)
		if !found {
			//panic(fmt.Sprintf("not found oracle:%s", oracleAddr.String()))
			logger.Error("TryAttestation", "not found oracle", oracleAddr.String(), "claimEventNonce",
				claim.GetEventNonce(), "claimType", claim.GetEventNonce(), "claimHeight", claim.GetBlockHeight())
			continue
		}
		oraclePower := oracle.GetPower()
		// Add it to the attestation power's sum
		attestationPower = attestationPower.Add(oraclePower)
		if attestationPower.LT(requiredPower) {
			continue
		}
		// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
		// process the attestation, set Observed to true, and break
		lastEventNonce := k.GetLastObservedEventNonce(ctx)
		// this check is performed at the next level up so this should never panic
		// outside of programmer error.
		if claim.GetEventNonce() != lastEventNonce+1 {
			panic("attempting to apply events to state out of order")
		}
		k.SetLastObservedEventNonce(ctx, claim.GetEventNonce())
		k.SetLastObservedBlockHeight(ctx, claim.GetBlockHeight())

		att.Observed = true
		k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)

		k.processAttestation(ctx, att, claim)
		k.emitObservedEvent(ctx, att, claim)
		break
	}
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation, claim types.ExternalClaim) {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler(xCtx, *att, claim); err != nil {
		// execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.Logger(ctx).Error("attestation failed", "cause", err.Error(), "claim type", claim.GetType(),
			"id", hex.EncodeToString(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash())),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	} else {
		commit() // persist transient storage
		ctx.EventManager().EmitEvents(xCtx.EventManager().Events())
	}
}

// emitObservedEvent emits an event with information about an attestation that has been applied to
// consensus state.
func (k Keeper) emitObservedEvent(ctx sdk.Context, _ *types.Attestation, claim types.ExternalClaim) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyAttestationType, claim.GetType().String()),
		sdk.NewAttribute(types.AttributeKeyAttestationID, hex.EncodeToString(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()))),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.GetEventNonce())),
	))
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte, att *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	store.Set(aKey, k.cdc.MustMarshal(att))
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte) *types.Attestation {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.Attestation
	k.cdc.MustUnmarshal(bz, &att)
	return &att
}

// DeleteAttestation deletes an attestation given an event nonce and claim
func (k Keeper) DeleteAttestation(ctx sdk.Context, att types.Attestation) {
	claim, err := k.UnpackAttestationClaim(&att)
	if err != nil {
		panic("Bad Attestation in DeleteAttestation")
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAttestationKeyWithHash(claim.GetEventNonce(), claim.ClaimHash()))
}

// GetAttestationMapping returns a mapping of eventNonce -> attestations at that nonce
func (k Keeper) GetAttestationMapping(ctx sdk.Context) (out map[uint64][]types.Attestation) {
	out = make(map[uint64][]types.Attestation)
	k.IterateAttestations(ctx, func(_ []byte, att types.Attestation) bool {
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		if val, ok := out[claim.GetEventNonce()]; !ok {
			out[claim.GetEventNonce()] = []types.Attestation{att}
		} else {
			out[claim.GetEventNonce()] = append(val, att)
		}
		return false
	})
	return
}

// IterateAttestations iterates through all attestations
func (k Keeper) IterateAttestations(ctx sdk.Context, cb func([]byte, types.Attestation) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.OracleAttestationKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := types.Attestation{}
		k.cdc.MustUnmarshal(iter.Value(), &att)
		// cb returns true to stop early
		if cb(iter.Key(), att) {
			return
		}
	}
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEventNonceKey)

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// GetLastObservedBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedBlockHeight(ctx sdk.Context) types.LastObservedBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedBlockHeightKey)

	if len(bytes) == 0 {
		return types.LastObservedBlockHeight{
			ExternalBlockHeight: 0,
			BlockHeight:         0,
		}
	}
	height := types.LastObservedBlockHeight{}
	k.cdc.MustUnmarshal(bytes, &height)
	return height
}

// SetLastObservedBlockHeight sets the block height in the store.
func (k Keeper) SetLastObservedBlockHeight(ctx sdk.Context, externalBlockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	height := types.LastObservedBlockHeight{
		ExternalBlockHeight: externalBlockHeight,
		BlockHeight:         uint64(ctx.BlockHeight()),
	}
	store.Set(types.LastObservedBlockHeightKey, k.cdc.MustMarshal(&height))
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
	valset := types.OracleSet{}
	k.cdc.MustUnmarshal(bytes, &valset)
	return &valset
}

// SetLastObservedOracleSet updates the last observed oracle set in the store
func (k Keeper) SetLastObservedOracleSet(ctx sdk.Context, valset types.OracleSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedOracleSetKey, k.cdc.MustMarshal(&valset))
}

// SetLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) SetLastObservedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEventNonceKey, types.UInt64Bytes(nonce))
}

// GetLastEventNonceByOracle returns the latest event nonce for a given oracle
func (k Keeper) GetLastEventNonceByOracle(ctx sdk.Context, oracle sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByOracleKey(oracle))

	if len(bytes) == 0 {
		// in the case that we have no existing value this is the first
		// time a oracle is submitting a claim. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		lastEventNonce := k.GetLastObservedEventNonce(ctx)
		if lastEventNonce >= 1 {
			return lastEventNonce - 1
		} else {
			return 0
		}
	}
	return types.UInt64FromBytes(bytes)
}

// DelLastEventNonceByOracle delete the latest event nonce for a given oracle
func (k Keeper) DelLastEventNonceByOracle(ctx sdk.Context, oracle sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLastEventNonceByOracleKey(oracle)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

// SetLastEventNonceByOracle sets the latest event nonce for a give oracle
func (k Keeper) SetLastEventNonceByOracle(ctx sdk.Context, oracle sdk.AccAddress, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByOracleKey(oracle), types.UInt64Bytes(nonce))
}
