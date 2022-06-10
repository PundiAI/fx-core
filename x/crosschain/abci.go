package crosschain

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/functionx/fx-core/x/crosschain/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	slashing(ctx, k, params)
	attestationTally(ctx, k)
	cleanupTimedOutBatches(ctx, k)
	createOracleSetRequest(ctx, k, params)
	pruneOracleSet(ctx, k, params.SignedWindow)
	pruneAttestations(ctx, k)
}

func createOracleSetRequest(ctx sdk.Context, k keeper.Keeper, params types.Params) {
	// Auto OracleSetRequest Creation.
	// WARNING: do not use k.GetLastObservedOracleSet in this function, it *will* result in losing control of the bridge
	if currentOracleSet, isNeed := isNeedOracleSetRequest(ctx, k, params.OracleSetUpdatePowerChangePercent); isNeed {
		k.AddOracleSetRequest(ctx, currentOracleSet, params.GravityId)
	}
}

func isNeedOracleSetRequest(ctx sdk.Context, k keeper.Keeper, oracleSetUpdatePowerChangePercent sdk.Dec) (*types.OracleSet, bool) {
	currentOracleSet := k.GetCurrentOracleSet(ctx)
	// 1. get latest OracleSet
	latestOracleSet := k.GetLatestOracleSet(ctx)
	if latestOracleSet == nil {
		return currentOracleSet, true
	}
	// 2. Oracle slash
	if k.GetLastOracleSlashBlockHeight(ctx) == uint64(ctx.BlockHeight()) {
		ctx.Logger().Info("oracle set change", "has oracle slash in block", ctx.BlockHeight())
		return currentOracleSet, true
	}
	// 3. Power diff
	powerDiff := types.BridgeValidators(currentOracleSet.Members).PowerDiff(latestOracleSet.Members)
	powerDiffStr := fmt.Sprintf("%.8f", powerDiff)
	powerDiffDec, err := sdk.NewDecFromStr(powerDiffStr)
	if err != nil {
		panic(fmt.Errorf("covert power diff to dec err!!!powerDiff: %v, err: %v", powerDiffStr, err))
	}

	if oracleSetUpdatePowerChangePercent.GT(sdk.OneDec()) {
		oracleSetUpdatePowerChangePercent = sdk.OneDec()
	}
	if powerDiffDec.GTE(oracleSetUpdatePowerChangePercent) {
		ctx.Logger().Info("oracle set change", "change threshold", oracleSetUpdatePowerChangePercent.String(), "powerDiff", powerDiff)
		return currentOracleSet, true
	}
	return currentOracleSet, false
}

func slashing(ctx sdk.Context, k keeper.Keeper, params types.Params) {
	if uint64(ctx.BlockHeight()) <= params.SignedWindow {
		return
	}
	// Slash oracle for not confirming oracle set requests, batch requests
	oracles := k.GetAllOracles(ctx, true)
	oracleSetHasSlash := oracleSetSlashing(ctx, k, oracles, params)
	batchHasSlash := batchSlashing(ctx, k, oracles, params)
	if oracleSetHasSlash || batchHasSlash {
		k.CommonSetOracleTotalPower(ctx)
	}
}

func oracleSetSlashing(ctx sdk.Context, k keeper.Keeper, oracles types.Oracles, params types.Params) (hasSlash bool) {
	maxHeight := uint64(ctx.BlockHeight()) - params.SignedWindow
	unSlashedOracleSets := k.GetUnSlashedOracleSets(ctx, maxHeight)
	logger := k.Logger(ctx)
	// Find all verifiers that meet the penalty to change the signature consensus
	for _, oracleSet := range unSlashedOracleSets {
		confirms := k.GetOracleSetConfirms(ctx, oracleSet.Nonce)
		confirmOracleMap := make(map[string]bool, len(confirms))
		for _, confirm := range confirms {
			confirmOracleMap[confirm.ExternalAddress] = true
		}
		for _, oracle := range oracles {
			if uint64(oracle.StartHeight) > oracleSet.Height {
				continue
			}
			if _, ok := confirmOracleMap[oracle.ExternalAddress]; !ok {
				logger.Info("slash oracle by oracle set", "oracleAddress", oracle.OracleAddress,
					"oracleSetNonce", oracleSet.Nonce, "oracleSetHeight", oracleSet.Height, "blockHeight", ctx.BlockHeight(), "slashFraction", params.SlashFraction.String())
				k.SlashOracle(ctx, oracle, params.SlashFraction)
				hasSlash = true
			}
		}
		// then we set the latest slashed oracleSet  nonce
		k.SetLastSlashedOracleSetNonce(ctx, oracleSet.Nonce)
	}
	return
}

func batchSlashing(ctx sdk.Context, k keeper.Keeper, oracles types.Oracles, params types.Params) (hasSlash bool) {
	maxHeight := uint64(ctx.BlockHeight()) - params.SignedWindow
	unSlashedBatches := k.GetUnSlashedBatches(ctx, maxHeight)
	logger := k.Logger(ctx)
	for _, batch := range unSlashedBatches {
		confirms := k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)
		confirmOracleMap := make(map[string]bool, len(confirms))
		for _, confirm := range confirms {
			confirmOracleMap[confirm.ExternalAddress] = true
		}
		for _, oracle := range oracles {
			if uint64(oracle.StartHeight) > batch.Block {
				continue
			}
			if _, ok := confirmOracleMap[oracle.ExternalAddress]; !ok {
				logger.Info("slash oracle by batch", "oracleAddress", oracle.OracleAddress,
					"batchNonce", batch.BatchNonce, "batchHeight", batch.Block, "blockHeight", ctx.BlockHeight(), "slashFraction", params.SlashFraction.String())
				k.SlashOracle(ctx, oracle, params.SlashFraction)
				hasSlash = true
			}
		}
		// then we set the latest slashed batch block
		k.SetLastSlashedBatchBlock(ctx, batch.Block)
	}
	return
}

// Iterate over all attestations currently being voted on in order of nonce and
// "Observe" those who have passed the threshold. Break the loop once we see
// an attestation that has not passed the threshold
func attestationTally(ctx sdk.Context, k keeper.Keeper) {
	attMap := k.GetAttestationMapping(ctx)
	// We make a slice with all the event nonces that are in the attestation mapping
	nonces := make([]uint64, 0, len(attMap))
	for k := range attMap {
		nonces = append(nonces, k)
	}
	// Then we sort it
	sort.Slice(nonces, func(i, j int) bool {
		return nonces[i] < nonces[j]
	})

	// This iterates over all nonces (event nonces) in the attestation mapping. Each value contains
	// a slice with one or more attestations at that event nonce. There can be multiple attestations
	// at one event nonce when validators disagree about what event happened at that nonce.
	for _, nonce := range nonces {
		// This iterates over all attestations at a particular event nonce.
		// They are ordered by when the first attestation at the event nonce was received.
		// This order is not important.
		for _, att := range attMap[nonce] {
			// We check if the event nonce is exactly 1 higher than the last attestation that was
			// observed. If it is not, we just move on to the next nonce. This will skip over all
			// attestations that have already been observed.
			//
			// Once we hit an event nonce that is one higher than the last observed event, we stop
			// skipping over this conditional and start calling tryAttestation (counting votes)
			// Once an attestation at a given event nonce has enough votes and becomes observed,
			// every other attestation at that nonce will be skipped, since the lastObservedEventNonce
			// will be incremented.
			//
			// Then we go to the next event nonce in the attestation mapping, if there is one. This
			// nonce will once again be one higher than the lastObservedEventNonce.
			// If there is an attestation at this event nonce which has enough votes to be observed,
			// we skip the other attestations and move on to the next nonce again.
			// If no attestation becomes observed, when we get to the next nonce, every attestation in
			// it will be skipped. The same will happen for every nonce after that.
			if nonce == k.GetLastObservedEventNonce(ctx)+1 {
				k.TryAttestation(ctx, &att)
			}
		}
	}
}

// cleanupTimedOutBatches deletes batches that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning batch 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single batch at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last SendToExternal or SendToFx to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutBatches(ctx sdk.Context, k keeper.Keeper) {
	externalBlockHeight := k.GetLastObservedBlockHeight(ctx).ExternalBlockHeight
	batches := k.GetOutgoingTxBatches(ctx)
	for _, batch := range batches {
		if batch.BatchTimeout < externalBlockHeight {
			if err := k.CancelOutgoingTXBatch(ctx, batch.TokenContract, batch.BatchNonce); err != nil {
				panic(fmt.Sprintf("Failed cancel out batch %s %d while trying to execute failed: %s", batch.TokenContract, batch.BatchNonce, err))
			}
		}
	}
}

// pruneOracleSet
func pruneOracleSet(ctx sdk.Context, k keeper.Keeper, signedOracleSetsWindow uint64) {
	// Validator set pruning
	// prune all validator sets with a nonce less than the
	// last observed nonce, they can't be submitted any longer
	//
	// Only prune oracleSets after the signed oracleSets window has passed
	// so that slashing can occur the block before we remove them
	lastObserved := k.GetLastObservedOracleSet(ctx)
	currentBlock := uint64(ctx.BlockHeight())
	tooEarly := currentBlock < signedOracleSetsWindow
	if lastObserved != nil && !tooEarly {
		earliestToPrune := currentBlock - signedOracleSetsWindow
		sets := k.GetOracleSets(ctx)
		for _, set := range sets {
			if set.Nonce < lastObserved.Nonce && set.Nonce < earliestToPrune {
				k.DeleteOracleSet(ctx, set.Nonce)
			}
		}
	}
}

// Iterate over all attestations currently being voted on in order of nonce
// and prune those that are older than the current nonce and no longer have any
// use. This could be combined with create attestation and save some computation
// but (A) pruning keeps the iteration small in the first place and (B) there is
// already enough nuance in the other handler that it's best not to complicate it further
func pruneAttestations(ctx sdk.Context, k keeper.Keeper) {
	attMap := k.GetAttestationMapping(ctx)
	// We make a slice with all the event nonces that are in the attestation mapping
	keys := make([]uint64, 0, len(attMap))
	for k := range attMap {
		keys = append(keys, k)
	}
	// Then we sort it
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// we delete all attestations earlier than the current event nonce
	// minus some buffer value. This buffer value is purely to allow
	// frontends and other UI components to view recent oracle history
	const eventsToKeep = 100
	lastNonce := k.GetLastObservedEventNonce(ctx)
	var cutoff uint64
	if lastNonce <= eventsToKeep {
		return
	} else {
		cutoff = lastNonce - eventsToKeep
	}

	// This iterates over all keys (event nonces) in the attestation mapping. Each value contains
	// a slice with one or more attestations at that event nonce. There can be multiple attestations
	// at one event nonce when validators disagree about what event happened at that nonce.
	for _, nonce := range keys {
		// This iterates over all attestations at a particular event nonce.
		// They are ordered by when the first attestation at the event nonce was received.
		// This order is not important.
		for _, att := range attMap[nonce] {
			// delete all before the cutoff
			if nonce < cutoff {
				k.DeleteAttestation(ctx, att)
			}
		}
	}
}
