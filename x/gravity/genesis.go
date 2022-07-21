package gravity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v2/x/gravity/keeper"

	"github.com/functionx/fx-core/v2/x/gravity/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)

	// reset delegate keys in state
	for _, keys := range data.DelegateKeys {
		err := keys.ValidateBasic()
		if err != nil {
			panic("Invalid delegate key in Genesis!")
		}
		if err := types.ValidateEthAddressAndValidateChecksum(keys.EthAddress); err != nil {
			panic(fmt.Errorf("invalid ethereum address: %s", keys.EthAddress))
		}
		val, err := sdk.ValAddressFromBech32(keys.Validator)
		if err != nil {
			panic(err)
		}
		orch, err := sdk.AccAddressFromBech32(keys.Orchestrator)
		if err != nil {
			panic(err)
		}
		// set the orchestrator address
		k.SetOrchestratorValidator(ctx, val, orch)
		// set the ethereum address
		k.SetEthAddressForValidator(ctx, val, keys.EthAddress)
	}

	// populate state with cosmos originated denom-erc20 mapping
	for _, item := range data.Erc20ToDenoms {
		if err := types.ValidateEthAddressAndValidateChecksum(item.Erc20); err != nil {
			panic(fmt.Errorf("invalid erc20 address in Erc20ToDenoms for item %s: %s", item.Denom, item.Erc20))
		}
		k.SetFxOriginatedDenomToERC20(ctx, item.Denom, item.Erc20)
	}

	// reset valsets in state
	for _, vs := range data.Valsets {
		k.StoreValset(ctx, &vs)
	}

	// reset valset confirmations in state
	for _, conf := range data.ValsetConfirms {
		k.SetValsetConfirm(ctx, conf)
	}

	// reset batches in state
	for _, batch := range data.Batches {
		k.StoreBatchUnsafe(ctx, &batch)
	}

	// reset batch confirmations in state
	for _, conf := range data.BatchConfirms {
		k.SetBatchConfirm(ctx, &conf)
	}

	// reset pool transactions in state
	for _, tx := range data.UnbatchedTransfers {
		k.SetPoolEntry(ctx, &tx)
	}

	// reset attestations in state
	for _, att := range data.Attestations {
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
	}
	k.SetLastObservedEventNonce(ctx, data.LastObservedNonce)

	// reset attestation state of specific validators
	// this must be done after the above to be correct
	for _, att := range data.Attestations {
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
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.GetLastEventNonceByValidator(ctx, val)
			if claim.GetEventNonce() > last {
				k.SetLastEventNonceByValidator(ctx, val, claim.GetEventNonce())
			}
		}
	}
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesisState := &types.GenesisState{
		Params:            k.GetParams(ctx),
		LastObservedNonce: k.GetLastObservedEventNonce(ctx),
	}
	for _, address := range k.GetDelegateKeys(ctx) {
		genesisState.DelegateKeys = append(genesisState.DelegateKeys, *address)
	}
	for _, valset := range k.GetValsets(ctx) {
		genesisState.Valsets = append(genesisState.Valsets, *valset)
	}
	k.IterateERC20ToDenom(ctx, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
		genesisState.Erc20ToDenoms = append(genesisState.Erc20ToDenoms, *erc20ToDenom)
		return false
	})
	for _, batch := range k.GetOutgoingTxBatches(ctx) {
		genesisState.Batches = append(genesisState.Batches, *batch)
	}
	for _, tx := range k.GetPoolTransactions(ctx) {
		genesisState.UnbatchedTransfers = append(genesisState.UnbatchedTransfers, *tx)
	}
	for _, vs := range genesisState.Valsets {
		for _, cfg := range k.GetValsetConfirms(ctx, vs.Nonce) {
			genesisState.ValsetConfirms = append(genesisState.ValsetConfirms, *cfg)
		}
	}
	for _, batch := range genesisState.Batches {
		genesisState.BatchConfirms = append(genesisState.BatchConfirms, k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)...)
	}
	for _, attestations := range k.GetAttestationMapping(ctx) {
		genesisState.Attestations = append(genesisState.Attestations, attestations...)
	}
	return genesisState
}
