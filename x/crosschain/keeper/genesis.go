package keeper

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

// InitGenesis import module genesis
func InitGenesis(ctx sdk.Context, k Keeper, state *types.GenesisState) {
	k.SetParams(ctx, &state.Params)

	// 0x24
	k.SetLastObservedEventNonce(ctx, state.LastObservedEventNonce)
	// 0x32
	k.SetLastObservedBlockHeight(ctx, state.LastObservedBlockHeight.ExternalBlockHeight, state.LastObservedBlockHeight.BlockHeight)

	// 0x38
	k.SetProposalOracle(ctx, &state.ProposalOracle)

	// 0x33
	k.SetLastObservedOracleSet(ctx, state.LastObservedOracleSet)

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

	for _, set := range state.OracleSets {
		// 0x15 0x29
		k.StoreOracleSet(ctx, &set)
	}

	for _, bridgeToken := range state.BridgeTokens {
		// 0x26 0x27
		if _, err := k.AddBridgeToken(ctx, bridgeToken.Token, hex.EncodeToString([]byte(bridgeToken.ChannelIbc))); err != nil {
			panic(err)
		}
	}
	for _, confirm := range state.BatchConfirms {
		for _, oracle := range state.Oracles {
			if confirm.BridgerAddress == oracle.BridgerAddress {
				// 0x22
				k.SetBatchConfirm(ctx, oracle.GetOracle(), &confirm)
			}
		}
	}
	for _, confirm := range state.OracleSetConfirms {
		for _, oracle := range state.Oracles {
			if confirm.BridgerAddress == oracle.BridgerAddress {
				// 0x16
				k.SetOracleSetConfirm(ctx, oracle.GetOracle(), &confirm)
			}
		}
	}

	for _, transfer := range state.UnbatchedTransfers {
		// 0x18
		if err := k.AddUnbatchedTx(ctx, &transfer); err != nil {
			panic(err)
		}
	}

	for _, batch := range state.Batches {
		// 0x20 0x21
		if err := k.StoreBatch(ctx, &batch); err != nil {
			panic(err)
		}
	}

	// reset attestations in state
	for _, att := range state.Attestations {
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		// 0x17
		k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
	}

	// reset attestation state of specific validators
	// this must be done after the above to be correct
	for _, att := range state.Attestations {
		claim, err := k.UnpackAttestationClaim(&att)
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
			oracle, err := sdk.AccAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
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
	for _, oracleSet := range k.GetOracleSets(ctx) {
		state.OracleSets = append(state.OracleSets, *oracleSet)
	}
	for _, batch := range k.GetOutgoingTxBatches(ctx) {
		state.Batches = append(state.Batches, *batch)
	}
	for _, attestations := range k.GetAttestationMapping(ctx) {
		state.Attestations = append(state.Attestations, attestations...)
	}
	for _, tx := range k.GetUnbatchedTransactions(ctx) {
		state.UnbatchedTransfers = append(state.UnbatchedTransfers, *tx)
	}
	for _, vs := range state.OracleSets {
		for _, cfg := range k.GetOracleSetConfirms(ctx, vs.Nonce) {
			state.OracleSetConfirms = append(state.OracleSetConfirms, *cfg)
		}
	}
	for _, batch := range state.Batches {
		state.BatchConfirms = append(state.BatchConfirms, k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)...)
	}
	k.IterateBridgeTokenToDenom(ctx, func(key []byte, erc20ToDenom *types.BridgeToken) bool {
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
