package keeper

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	signedWindow := k.GetSignedWindow(ctx)
	k.slashing(ctx, signedWindow)
	k.createOracleSetRequest(ctx)
	k.pruneOracleSet(ctx, signedWindow)
}

func (k Keeper) createOracleSetRequest(ctx sdk.Context) {
	if currentOracleSet, isNeed := k.isNeedOracleSetRequest(ctx); isNeed {
		k.AddOracleSetRequest(ctx, currentOracleSet)
	}
}

func (k Keeper) isNeedOracleSetRequest(ctx sdk.Context) (*types.OracleSet, bool) {
	currentOracleSet := k.GetCurrentOracleSet(ctx)
	// 1. get latest OracleSet
	latestOracleSet := k.GetLatestOracleSet(ctx)
	if latestOracleSet == nil {
		return currentOracleSet, true
	}
	// 2. Oracle slash
	if k.GetLastOracleSlashBlockHeight(ctx) == uint64(ctx.BlockHeight()) {
		k.Logger(ctx).Info("oracle set change", "has oracle slash in block", ctx.BlockHeight())
		return currentOracleSet, true
	}
	// 3. Power diff
	powerDiff := fmt.Sprintf("%.8f", types.BridgeValidators(currentOracleSet.Members).PowerDiff(latestOracleSet.Members))
	powerDiffDec, err := sdk.NewDecFromStr(powerDiff)
	if err != nil {
		panic(fmt.Errorf("covert power diff to dec err, powerDiff: %v, err: %w", powerDiff, err))
	}

	oracleSetUpdatePowerChangePercent := k.GetOracleSetUpdatePowerChangePercent(ctx)
	if oracleSetUpdatePowerChangePercent.GT(sdk.OneDec()) {
		oracleSetUpdatePowerChangePercent = sdk.OneDec()
	}
	if powerDiffDec.GTE(oracleSetUpdatePowerChangePercent) {
		k.Logger(ctx).Info("oracle set change", "change threshold", oracleSetUpdatePowerChangePercent.String(), "powerDiff", powerDiff)
		return currentOracleSet, true
	}
	return currentOracleSet, false
}

func (k Keeper) slashing(ctx sdk.Context, signedWindow uint64) {
	if uint64(ctx.BlockHeight()) <= signedWindow {
		return
	}
	// Slash oracle for not confirming oracle set requests, batch requests
	oracles := k.GetAllOracles(ctx, true)
	oracleSetHasSlash := k.oracleSetSlashing(ctx, oracles, signedWindow)
	batchHasSlash := k.batchSlashing(ctx, oracles, signedWindow)
	refundHasSlash := k.refundSlashing(ctx, signedWindow)
	if oracleSetHasSlash || batchHasSlash || refundHasSlash {
		k.CommonSetOracleTotalPower(ctx)
	}
}

func (k Keeper) oracleSetSlashing(ctx sdk.Context, oracles types.Oracles, signedWindow uint64) (hasSlash bool) {
	maxHeight := uint64(ctx.BlockHeight()) - signedWindow
	unSlashedOracleSets := k.GetUnSlashedOracleSets(ctx, maxHeight)

	// Find all verifiers that meet the penalty to change the signature consensus
	for _, oracleSet := range unSlashedOracleSets {
		confirmOracleMap := make(map[string]struct{})
		k.IterateOracleSetConfirmByNonce(ctx, oracleSet.Nonce, func(confirm *types.MsgOracleSetConfirm) bool {
			confirmOracleMap[confirm.ExternalAddress] = struct{}{}
			return false
		})
		for i := 0; i < len(oracles); i++ {
			if uint64(oracles[i].StartHeight) > oracleSet.Height {
				continue
			}
			if _, ok := confirmOracleMap[oracles[i].ExternalAddress]; !ok {
				k.Logger(ctx).Info("slash oracle by oracle set", "oracleAddress", oracles[i].OracleAddress,
					"oracleSetNonce", oracleSet.Nonce, "oracleSetHeight", oracleSet.Height, "blockHeight", ctx.BlockHeight())
				k.SlashOracle(ctx, oracles[i].OracleAddress)
				hasSlash = true
			}
		}
		// then we set the latest slashed oracleSet  nonce
		k.SetLastSlashedOracleSetNonce(ctx, oracleSet.Nonce)
	}
	return hasSlash
}

func (k Keeper) batchSlashing(ctx sdk.Context, oracles types.Oracles, signedWindow uint64) (hasSlash bool) {
	maxHeight := uint64(ctx.BlockHeight()) - signedWindow
	unSlashedBatches := k.GetUnSlashedBatches(ctx, maxHeight)

	for _, batch := range unSlashedBatches {
		confirmOracleMap := make(map[string]struct{})
		k.IterateBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract, func(confirm *types.MsgConfirmBatch) bool {
			confirmOracleMap[confirm.ExternalAddress] = struct{}{}
			return false
		})
		for i := 0; i < len(oracles); i++ {
			if uint64(oracles[i].StartHeight) > batch.Block {
				continue
			}
			if _, ok := confirmOracleMap[oracles[i].ExternalAddress]; !ok {
				k.Logger(ctx).Info("slash oracle by batch", "oracleAddress", oracles[i].OracleAddress,
					"batchNonce", batch.BatchNonce, "batchHeight", batch.Block, "blockHeight", ctx.BlockHeight())
				k.SlashOracle(ctx, oracles[i].OracleAddress)
				hasSlash = true
			}
		}
		// then we set the latest slashed batch block
		k.SetLastSlashedBatchBlock(ctx, batch.Block)
	}
	return hasSlash
}

func (k Keeper) refundSlashing(ctx sdk.Context, signedWindow uint64) (hasSlash bool) {
	maxHeight := uint64(ctx.BlockHeight()) - signedWindow
	unSlashRefunds := k.GetUnSlashedRefundRecords(ctx, maxHeight)

	snapshotOracleMap := make(map[uint64]*types.SnapshotOracle)
	for _, record := range unSlashRefunds {
		confirmOracleMap := make(map[string]bool)
		k.IterRefundConfirmByNonce(ctx, record.EventNonce, func(confirm *types.MsgConfirmRefund) bool {
			confirmOracleMap[confirm.ExternalAddress] = true
			return false
		})

		snapshotOracle, found := snapshotOracleMap[record.OracleSetNonce]
		if !found {
			snapshotOracle, found = k.GetSnapshotOracle(ctx, record.OracleSetNonce)
			if !found {
				k.Logger(ctx).Error("refund slashing", "oracle set not found", record.OracleSetNonce)
				continue
			}
			snapshotOracleMap[record.OracleSetNonce] = snapshotOracle
		}

		for _, members := range snapshotOracle.GetMembers() {
			if _, ok := confirmOracleMap[members.ExternalAddress]; !ok {
				oracle, found := k.GetOracleByExternalAddress(ctx, members.ExternalAddress)
				if !found {
					k.Logger(ctx).Error("refund slashing", "oracle not found", members.ExternalAddress)
					continue
				}
				k.Logger(ctx).Info("slash oracle by refund", "externalAddress", members.ExternalAddress,
					"oracleAddress", oracle.String(), "recordNonce", record.EventNonce, "blockHeight", ctx.BlockHeight())
				k.SlashOracle(ctx, oracle.String())
				hasSlash = true
			}
		}
		k.SetLastSlashedRefundNonce(ctx, record.EventNonce)
	}
	return hasSlash
}

// cleanupTimedOutBatches deletes batches that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning batch 5 can have a later timeout than batch 6
//
//	this means that we MUST only cleanup a single batch at a time
//
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//
//	here is the Ethereum block height at the time of the last SendToExternal or SendToFx to be observed. It's very important we do not
//	project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//	AND any deposit or withdraw has occurred to update the Ethereum block height.
func (k Keeper) cleanupTimedOutBatches(ctx sdk.Context) {
	externalBlockHeight := k.GetLastObservedBlockHeight(ctx).ExternalBlockHeight
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		if batch.BatchTimeout < externalBlockHeight {
			if err := k.CancelOutgoingTxBatch(ctx, batch.TokenContract, batch.BatchNonce); err != nil {
				panic(fmt.Sprintf("Failed cancel out batch %s %d while trying to execute failed: %s", batch.TokenContract, batch.BatchNonce, err))
			}
		}
		return false
	})
}

func (k Keeper) cleanupTimeOutRefund(ctx sdk.Context) {
	externalBlockHeight := k.GetLastObservedBlockHeight(ctx).ExternalBlockHeight
	k.IterRefundRecord(ctx, func(record *types.RefundRecord) bool {
		if record.Timeout > externalBlockHeight {
			return true
		}
		receiver, coins, err := k.refundTokenToReceiver(ctx, record)
		if err != nil {
			k.Logger(ctx).Error("clean up refund timeout", "event nonce", record.EventNonce, "error", err.Error())
			return false
		}
		k.DeleteRefundRecord(ctx, record)
		k.DeleteRefundConfirm(ctx, record.EventNonce)
		k.RemoveEventSnapshotOracle(ctx, record.OracleSetNonce, record.EventNonce)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(types.EventTypeRefundTimeout,
				sdk.NewAttribute(sdk.AttributeKeySender, record.Receiver),
				sdk.NewAttribute(types.AttributeKeyRefundAddress, receiver.String()),
				sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(record.EventNonce)),
				sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
			),
		})
		return false
	})
}

func (k Keeper) refundTokenToReceiver(ctx sdk.Context, record *types.RefundRecord) (sdk.AccAddress, sdk.Coins, error) {
	receiverAddr := types.ExternalAddressToAccAddress(k.moduleName, record.Receiver)
	cacheCtx, commit := ctx.CacheContext()
	coins, err := k.bridgeCallTransferToSender(cacheCtx, receiverAddr.Bytes(), record.Tokens)
	if err != nil {
		return nil, nil, err
	}
	if err = k.bridgeCallTransferToReceiver(cacheCtx, receiverAddr, receiverAddr.Bytes(), coins); err != nil {
		return nil, nil, err
	}
	commit()
	return receiverAddr, coins, nil
}

func (k Keeper) pruneOracleSet(ctx sdk.Context, signedOracleSetsWindow uint64) {
	// Oracle set pruning
	// prune all Oracle sets with a nonce less than the
	// last observed nonce, they can't be submitted any longer
	//
	// Only prune oracleSets after the signed oracleSets window has passed
	// so that slashing can occur the block before we remove them
	lastObserved := k.GetLastObservedOracleSet(ctx)
	currentBlock := uint64(ctx.BlockHeight())
	tooEarly := currentBlock < signedOracleSetsWindow
	if lastObserved != nil && !tooEarly {
		earliestToPrune := currentBlock - signedOracleSetsWindow
		k.IterateOracleSets(ctx, false, func(set *types.OracleSet) bool {
			if earliestToPrune > set.Height && lastObserved.Nonce > set.Nonce {
				k.DeleteOracleSet(ctx, set.Nonce)
				k.DeleteOracleSetConfirm(ctx, set.Nonce)
			}
			return false
		})
	}
}

// Iterate over all attestations currently being voted on in order of nonce
// and prune those that are older than the current nonce and no longer have any
// use. This could be combined with create attestation and save some computation
// but (A) pruning keeps the iteration small in the first place and (B) there is
// already enough nuance in the other handler that it's best not to complicate it further
func (k Keeper) pruneAttestations(ctx sdk.Context) {
	lastNonce := k.GetLastObservedEventNonce(ctx)
	if lastNonce <= types.MaxKeepEventSize {
		return
	}

	// we delete all attestations earlier than the current event nonce
	// minus some buffer value. This buffer value is purely to allow
	// frontends and other UI components to view recent oracle history
	cutoff := lastNonce - types.MaxKeepEventSize
	claimMap := make(map[uint64][]types.ExternalClaim)
	// We make a slice with all the event nonces that are in the attestation mapping
	var nonces []uint64
	k.IterateAttestationAndClaim(ctx, func(att *types.Attestation, claim types.ExternalClaim) bool {
		if claim.GetEventNonce() > cutoff {
			return true
		}
		if v, ok := claimMap[claim.GetEventNonce()]; !ok {
			claimMap[claim.GetEventNonce()] = []types.ExternalClaim{claim}
			nonces = append(nonces, claim.GetEventNonce())
		} else {
			claimMap[claim.GetEventNonce()] = append(v, claim)
		}
		return false
	})
	// Then we sort it
	sort.Slice(nonces, func(i, j int) bool {
		return nonces[i] < nonces[j]
	})

	// This iterates over all keys (event nonces) in the attestation mapping. Each value contains
	// a slice with one or more attestations at that event nonce. There can be multiple attestations
	// at one event nonce when Oracles disagree about what event happened at that nonce.
	for _, nonce := range nonces {
		// This iterates over all attestations at a particular event nonce.
		// They are ordered by when the first attestation at the event nonce was received.
		// This order is not important.
		for _, claim := range claimMap[nonce] {
			k.DeleteAttestation(ctx, claim)
		}
	}
}
