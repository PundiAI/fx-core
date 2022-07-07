package crosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/functionx/fx-core/x/crosschain/types"
)

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
