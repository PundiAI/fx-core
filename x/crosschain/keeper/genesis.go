package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// InitGenesis import module genesis
//
//gocyclo:ignore
func InitGenesis(ctx sdk.Context, k Keeper, state *types.GenesisState) {
	if err := k.SetParams(ctx, &state.Params); err != nil {
		panic(err)
	}

	// 0x24
	k.SetLastObservedEventNonce(ctx, state.LastObservedEventNonce)
	// 0x32
	k.SetLastObservedBlockHeight(ctx, state.LastObservedBlockHeight.ExternalBlockHeight, state.LastObservedBlockHeight.BlockHeight)

	// 0x38
	k.SetProposalOracle(ctx, &state.ProposalOracle)

	// 0x33
	k.SetLastObservedOracleSet(ctx, &state.LastObservedOracleSet)

	// 0x28
	k.SetLastSlashedOracleSetNonce(ctx, state.LastSlashedOracleSetNonce)

	// 0x30
	k.SetLastSlashedBatchBlock(ctx, state.LastSlashedBatchBlock)

	for _, oracle := range state.Oracles {
		// 0x12
		k.SetOracle(ctx, oracle)
		// 0x14
		k.SetOracleByBridger(ctx, oracle.GetBridger(), oracle.GetOracle())
		// 0x13
		k.SetOracleByExternalAddress(ctx, oracle.ExternalAddress, oracle.GetOracle())
	}
	// 0x39
	k.CommonSetOracleTotalPower(ctx)

	latestOracleSetNonce := uint64(0)
	for i := 0; i < len(state.OracleSets); i++ {
		set := state.OracleSets[i]
		// 0x15 0x29
		if set.Nonce > latestOracleSetNonce {
			latestOracleSetNonce = set.Nonce
		}
		k.StoreOracleSet(ctx, &set)
	}
	k.SetLatestOracleSetNonce(ctx, latestOracleSetNonce)

	for _, bridgeToken := range state.BridgeTokens {
		// 0x26 0x27
		k.AddBridgeToken(ctx, bridgeToken.Token, bridgeToken.Denom)
	}
	for i := 0; i < len(state.BatchConfirms); i++ {
		confirm := state.BatchConfirms[i]
		for _, oracle := range state.Oracles {
			if confirm.BridgerAddress == oracle.BridgerAddress {
				// 0x22
				k.SetBatchConfirm(ctx, oracle.GetOracle(), &confirm)
			}
		}
	}
	for i := 0; i < len(state.OracleSetConfirms); i++ {
		confirm := state.OracleSetConfirms[i]
		for _, oracle := range state.Oracles {
			if confirm.BridgerAddress == oracle.BridgerAddress {
				// 0x16
				k.SetOracleSetConfirm(ctx, oracle.GetOracle(), &confirm)
			}
		}
	}

	for i := 0; i < len(state.UnbatchedTransfers); i++ {
		transfer := state.UnbatchedTransfers[i]
		// 0x18
		if err := k.AddUnbatchedTx(ctx, &transfer); err != nil {
			panic(err)
		}
	}

	for i := 0; i < len(state.Batches); i++ {
		batch := state.Batches[i]
		// 0x20 0x21
		if err := k.StoreBatch(ctx, &batch); err != nil {
			panic(err)
		}
	}

	// reset attestations in state
	for i := 0; i < len(state.Attestations); i++ {
		att := state.Attestations[i]
		claim, err := types.UnpackAttestationClaim(k.cdc, &att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		// 0x17
		k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
	}

	// reset attestation state of specific validators
	// this must be done after the above to be correct
	for i := 0; i < len(state.Attestations); i++ {

		att := state.Attestations[i]
		claim, err := types.UnpackAttestationClaim(k.cdc, &att)
		if err != nil {
			panic("couldn't cast to claim")
		}
		// reconstruct the latest event nonce for every validator
		// if somehow this genesis state is saved when all attestations
		// have been cleaned up GetLastEventNonceByOracle handles that case
		//
		// if we where to save and load the last event nonce for every validator
		// then we would need to carry that state forever across all chain restarts
		// but since we've already had to handle the edge case of new validators joining
		// while all attestations have already been cleaned up we can do this instead and
		// not carry around every validators event nonce counter forever.
		for _, vote := range att.Votes {
			oracle := sdk.MustAccAddressFromBech32(vote)
			last := k.GetLastEventNonceByOracle(ctx, oracle)
			if claim.GetEventNonce() > last {
				// 0x23
				k.SetLastEventNonceByOracle(ctx, oracle, claim.GetEventNonce())
				// 0x35
				k.SetLastEventBlockHeightByOracle(ctx, oracle, claim.GetBlockHeight())
			}
		}
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	state := &types.GenesisState{
		Params:                  k.GetParams(ctx),
		LastObservedEventNonce:  k.GetLastObservedEventNonce(ctx),
		LastObservedBlockHeight: k.GetLastObservedBlockHeight(ctx),
	}
	k.IterateOracle(ctx, func(oracle types.Oracle) bool {
		state.Oracles = append(state.Oracles, oracle)
		return false
	})
	k.IterateOracleSets(ctx, false, func(oracleSet *types.OracleSet) bool {
		state.OracleSets = append(state.OracleSets, *oracleSet)
		return false
	})
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		state.Batches = append(state.Batches, *batch)
		return false
	})
	k.IterateAttestations(ctx, func(attestation *types.Attestation) bool {
		state.Attestations = append(state.Attestations, *attestation)
		return false
	})
	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		state.UnbatchedTransfers = append(state.UnbatchedTransfers, *tx)
		return false
	})
	for _, vs := range state.OracleSets {
		k.IterateOracleSetConfirmByNonce(ctx, vs.Nonce, func(confirm *types.MsgOracleSetConfirm) bool {
			state.OracleSetConfirms = append(state.OracleSetConfirms, *confirm)
			return false
		})
	}
	for _, batch := range state.Batches {
		k.IterateBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract, func(confirm *types.MsgConfirmBatch) bool {
			state.BatchConfirms = append(state.BatchConfirms, *confirm)
			return false
		})
	}
	k.IterateBridgeTokenToDenom(ctx, func(erc20ToDenom *types.BridgeToken) bool {
		state.BridgeTokens = append(state.BridgeTokens, *erc20ToDenom)
		return false
	})
	state.ProposalOracle, _ = k.GetProposalOracle(ctx)
	if lastObserved := k.GetLastObservedOracleSet(ctx); lastObserved != nil {
		state.LastObservedOracleSet = *lastObserved
	}
	state.LastSlashedBatchBlock = k.GetLastSlashedBatchBlock(ctx)
	state.LastSlashedOracleSetNonce = k.GetLastSlashedOracleSetNonce(ctx)
	return state
}
