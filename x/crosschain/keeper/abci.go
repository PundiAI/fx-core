package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
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
	powerDiffDec := sdkmath.LegacyMustNewDecFromStr(powerDiff)

	oracleSetUpdatePowerChangePercent := k.GetOracleSetUpdatePowerChangePercent(ctx)
	if oracleSetUpdatePowerChangePercent.GT(sdkmath.LegacyOneDec()) {
		oracleSetUpdatePowerChangePercent = sdkmath.LegacyOneDec()
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
	bridgeCallHasSlash := k.bridgeCallSlashing(ctx, oracles, signedWindow)
	if oracleSetHasSlash || batchHasSlash || bridgeCallHasSlash {
		k.SetLastTotalPower(ctx)
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

func (k Keeper) bridgeCallSlashing(ctx sdk.Context, oracles types.Oracles, signedWindow uint64) (hasSlash bool) {
	maxHeight := uint64(ctx.BlockHeight()) - signedWindow
	unSlashOutgoingBridgeCalls := k.GetUnSlashedBridgeCalls(ctx, maxHeight)

	for _, record := range unSlashOutgoingBridgeCalls {
		confirmOracleMap := make(map[string]struct{})
		k.IterBridgeCallConfirmByNonce(ctx, record.Nonce, func(confirm *types.MsgBridgeCallConfirm) bool {
			confirmOracleMap[confirm.ExternalAddress] = struct{}{}
			return false
		})

		for i := 0; i < len(oracles); i++ {
			if uint64(oracles[i].StartHeight) > record.BlockHeight {
				continue
			}
			if _, ok := confirmOracleMap[oracles[i].ExternalAddress]; !ok {
				k.SlashOracle(ctx, oracles[i].OracleAddress)
				k.Logger(ctx).Info("slash oracle by outgoing bridge call", "oracleAddress", oracles[i].OracleAddress,
					"nonce", record.Nonce, "bridgeCallHeight", record.BlockHeight, "blockHeight", ctx.BlockHeight())
				hasSlash = true
			}
		}

		k.SetLastSlashedBridgeCallNonce(ctx, record.Nonce)
	}
	return hasSlash
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
