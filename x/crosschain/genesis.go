package crosschain

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v2/x/crosschain/keeper"
	"github.com/functionx/fx-core/v2/x/crosschain/types"
)

// InitGenesis import module genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, &state.Params)

	k.SetLastObservedEventNonce(ctx, state.LastObservedEventNonce)
	k.SetLastObservedBlockHeight(ctx, state.LastObservedBlockHeight.ExternalBlockHeight, state.LastObservedBlockHeight.BlockHeight)

	for _, oracle := range state.Oracles {
		k.SetOracle(ctx, oracle)
		k.SetOracleByBridger(ctx, oracle.GetBridger(), oracle.GetOracle())
		k.SetOracleByExternalAddress(ctx, oracle.ExternalAddress, oracle.GetOracle())
		k.CommonSetOracleTotalPower(ctx)
	}

	for _, set := range state.OracleSets {
		k.StoreOracleSet(ctx, &set)
	}

	for _, bridgeToken := range state.BridgeTokens {
		if _, err := k.AddBridgeToken(ctx, bridgeToken.Token, hex.EncodeToString([]byte(bridgeToken.ChannelIbc))); err != nil {
			panic(err)
		}
	}

	for _, transfer := range state.UnbatchedTransfers {
		if err := k.AddUnbatchedTx(ctx, &transfer); err != nil {
			panic(err)
		}
	}

	for _, batch := range state.Batches {
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
				k.SetLastEventNonceByOracle(ctx, oracle, claim.GetEventNonce())
			}
		}
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
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
	return state
}
