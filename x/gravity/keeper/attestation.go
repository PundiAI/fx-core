package keeper

import (
	"encoding/hex"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v2/x/gravity/types"
)

func (k Keeper) Attest(ctx sdk.Context, claim types.EthereumClaim, anyClaim *codectypes.Any) (*types.Attestation, error) {
	valAddr, found := k.GetOrchestratorValidator(ctx, claim.GetClaimer())
	if !found {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce.
	// This prevents there being two attestations with the same nonce that get 2/3s of the votes
	// in the endBlocker.
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, valAddr)
	if claim.GetEventNonce() != lastEventNonce+1 {
		return nil, types.ErrNonContiguousEventNonce
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

	// Add the validator's vote to this attestation
	att.Votes = append(att.Votes, valAddr.String())

	k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)
	k.SetLastEventNonceByValidator(ctx, valAddr, claim.GetEventNonce())
	k.setLastEventBlockHeightByValidator(ctx, valAddr, claim.GetBlockHeight())

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
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.
	if !att.Observed {
		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
		var totalPower uint64
		var valPowers = map[string]uint64{}
		for _, validator := range validators {
			val := validator.GetOperator()
			_, found := k.GetEthAddressByValidator(ctx, val)
			if !found {
				continue
			}
			power := uint64(k.StakingKeeper.GetLastValidatorPower(ctx, val))
			totalPower += power
			valPowers[val.String()] = power
		}

		requiredPower := types.AttestationVotesPowerThreshold.Mul(sdk.NewIntFromUint64(totalPower)).Quo(sdk.NewInt(100))
		attestationPower := sdk.NewInt(0)
		for _, validator := range att.Votes {
			// Add it to the attestation power's sum
			attestationPower = attestationPower.Add(sdk.NewIntFromUint64(valPowers[validator]))
			// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
			// process the attestation, set Observed to true, and break
			if attestationPower.GTE(requiredPower) {
				lastEventNonce := k.GetLastObservedEventNonce(ctx)
				// this check is performed at the next level up so this should never panic
				// outside of programmer error.
				if claim.GetEventNonce() != lastEventNonce+1 {
					panic("attempting to apply events to state out of order")
				}
				k.SetLastObservedEventNonce(ctx, claim.GetEventNonce())
				k.SetLastObservedEthBlockHeight(ctx, claim.GetBlockHeight())

				att.Observed = true
				k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)

				err = k.processAttestation(ctx, att, claim)
				ctx.EventManager().EmitEvent(sdk.NewEvent(
					types.EventTypeContractEvnet,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
					sdk.NewAttribute(types.AttributeKeyClaimType, claim.GetType().String()),
					sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.GetEventNonce())),
					sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprint(claim.GetBlockHeight())),
					sdk.NewAttribute(types.AttributeKeyStateSuccess, fmt.Sprint(err == nil)),
				))
				break
			}
		}
	} else {
		// We panic here because this should never happen
		panic("attempting to process observed attestation")
	}
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation, claim types.EthereumClaim) error {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.attestationHandler.Handle(xCtx, *att, claim); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.Logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", hex.EncodeToString(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash())),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return err
	}
	commit() // persist transient storage
	ctx.EventManager().EmitEvents(xCtx.EventManager().Events())
	return nil
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

// GetAttestationMapping returns a mapping of eventnonce -> attestations at that nonce
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

// GetLastObservedEthBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedEthBlockHeight(ctx sdk.Context) types.LastObservedEthereumBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEthereumBlockHeightKey)

	if len(bytes) == 0 {
		return types.LastObservedEthereumBlockHeight{
			FxBlockHeight:  0,
			EthBlockHeight: 0,
		}
	}
	height := types.LastObservedEthereumBlockHeight{}
	k.cdc.MustUnmarshal(bytes, &height)
	return height
}

// SetLastObservedEthBlockHeight sets the block height in the store.
func (k Keeper) SetLastObservedEthBlockHeight(ctx sdk.Context, ethereumHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	height := types.LastObservedEthereumBlockHeight{
		EthBlockHeight: ethereumHeight,
		FxBlockHeight:  uint64(ctx.BlockHeight()),
	}
	store.Set(types.LastObservedEthereumBlockHeightKey, k.cdc.MustMarshal(&height))
}

// GetLastObservedValset retrieves the last observed validator set from the store
// WARNING: This value is not an up to date validator set on Ethereum, it is a validator set
// that AT ONE POINT was the one in the Gravity bridge on Ethereum. If you assume that it's up
// to date you may break the bridge
func (k Keeper) GetLastObservedValset(ctx sdk.Context) *types.Valset {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedValsetKey)

	if len(bytes) == 0 {
		return nil
	}
	valset := types.Valset{}
	k.cdc.MustUnmarshal(bytes, &valset)
	return &valset
}

// SetLastObservedValset updates the last observed validator set in the store
func (k Keeper) SetLastObservedValset(ctx sdk.Context, valset types.Valset) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedValsetKey, k.cdc.MustMarshal(&valset))
}

// SetLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) SetLastObservedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEventNonceKey, types.UInt64Bytes(nonce))
}

// GetLastEventNonceByValidator returns the latest event nonce for a given validator
func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByValidatorKey(validator))

	if len(bytes) == 0 {
		// in the case that we have no existing value this is the first
		// time a validator is submitting a claim. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		//
		// We could start at the LastObservedEventNonce but if we do that this
		// validator will be slashed, because they are responsible for making a claim
		// on any attestation that has not yet passed the slashing window.
		//
		// Therefore we need to return to them the lowest attestation that is still within
		// the slashing window. Since we delete attestations after the slashing window that's
		// just the lowest observed event in the store. If no claims have been submitted in for
		// params.SignedClaimsWindow we may have no attestations in our nonce. At which point
		// the last observed which is a persistent and never cleaned counter will suffice.
		lowestObserved := k.GetLastObservedEventNonce(ctx)
		attMap := k.GetAttestationMapping(ctx)
		// no new claims in params.SignedClaimsWindow, we can return the current value
		// because the validator can't be slashed for an event that has already passed.
		// so they only have to worry about the *next* event to occur
		if len(attMap) == 0 {
			return lowestObserved
		}
		for nonce, atts := range attMap {
			for att := range atts {
				if atts[att].Observed && nonce < lowestObserved {
					lowestObserved = nonce
				}
			}
		}
		// return the latest event minus one so that the validator
		// can submit that event and avoid slashing. special case
		// for zero
		if lowestObserved > 0 {
			return lowestObserved - 1
		} else {
			return 0
		}
	}
	return types.UInt64FromBytes(bytes)
}

// SetLastEventNonceByValidator sets the latest event nonce for a give validator
func (k Keeper) SetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByValidatorKey(validator), types.UInt64Bytes(nonce))
}
