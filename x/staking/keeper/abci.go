package keeper

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
)

func (k Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	valUpdates := staking.EndBlocker(ctx, k.Keeper)
	// remove consensus record after consensus update
	k.RemoveValidatorConsensusKey(ctx)
	// update validator consensus after abci update
	k.UpdateValidatorConsensusKey(ctx)
	// abci validator update
	return k.ABCIValidatorUpdate(ctx, valUpdates)
}

// RemoveValidatorConsensusKey delete validator record after UpdateValidatorConsensusKey
func (k Keeper) RemoveValidatorConsensusKey(ctx sdk.Context) {
	k.IteratorValidatorDelConsensusAddr(ctx, func(valAddr sdk.ValAddress, delConsAddr sdk.ConsAddress) bool {
		// remove del consensus address record
		k.RemoveValidatorDelConsensusAddr(ctx, valAddr)

		// delete validator by old consensus address
		k.RemoveValidatorOperatorByConsAddr(ctx, delConsAddr)

		// todo delete old consensus pubkey

		// todo delete old validator signing info

		return false
	})
}

// UpdateValidatorConsensusKey update validator consensus pubkey after ABCIValidatorUpdate
func (k Keeper) UpdateValidatorConsensusKey(ctx sdk.Context) {
	k.IteratorValidatorOldConsensusAddr(ctx, func(valAddr sdk.ValAddress, oldConsAddr sdk.ConsAddress) bool {
		// delete old consensus address record
		k.RemoveValidatorOldConsensusAddr(ctx, valAddr)

		// update validator consensus pubkey
		newPubkey, found := k.GetValidatorNewConsensusPubKey(ctx, valAddr)
		if !found {
			return false
		}

		// delete new consensus pubkey record
		k.RemoveValidatorNewConsensusPubKey(ctx, valAddr)

		// update validator consensus pubkey
		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			return false
		}
		newConsAddr := sdk.ConsAddress(newPubkey.Address())
		pkAny, err := codectypes.NewAnyWithValue(newPubkey)
		if err != nil {
			panic(err)
		}
		validator.ConsensusPubkey = pkAny
		k.SetValidator(ctx, validator)

		// add consensus pubkey record
		if err = k.slashingKeeper.AddPubkey(ctx, newPubkey); err != nil {
			panic(err)
		}

		// Update the signing info start height or create a new signing info
		signingInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, oldConsAddr)
		if !found {
			panic("expected signing info to be found")
		}
		signingInfo.Address = newConsAddr.String()
		k.slashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, signingInfo)

		// set del consensus address
		k.SetValidatorDelConsensusAddr(ctx, valAddr, oldConsAddr)
		return false
	})
}

// ABCIValidatorUpdate update old consensus pubkey to zero power, add new consensus pubkey to validator power
func (k Keeper) ABCIValidatorUpdate(ctx sdk.Context, valUpdates []abci.ValidatorUpdate) []abci.ValidatorUpdate {
	proposer := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)

	pkPower := make(map[crypto.PublicKey]int64, len(valUpdates))
	for _, valUpdate := range valUpdates {
		pkPower[valUpdate.PubKey] = valUpdate.Power
	}

	pkUpdate := make([]abci.ValidatorUpdate, 0, 50)
	k.IteratorValidatorNewConsensusPubKey(ctx, func(valAddr sdk.ValAddress, newPubKey cryptotypes.PubKey) bool {
		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			// validator not found, remove new pubkey update
			k.RemoveValidatorNewConsensusPubKey(ctx, valAddr)
			return false
		}

		oldPubKey, err := validator.ConsPubKey()
		if err != nil {
			panic(err)
		}
		oldConsAddr := sdk.ConsAddress(oldPubKey.Address())
		// if validator is proposer, skip this block
		if oldConsAddr.Equals(proposer) {
			return false
		}
		oldTmPk, err := cryptocodec.ToTmProtoPublicKey(oldPubKey)
		if err != nil {
			panic(err)
		}
		newTmPk, err := cryptocodec.ToTmProtoPublicKey(newPubKey)
		if err != nil {
			panic(err)
		}

		power, ok := pkPower[oldTmPk]
		// if power not found, use validator power
		if !ok {
			power = validator.ConsensusPower(k.PowerReduction(ctx))
		}
		// set old tmPk power to 0, remove
		pkUpdate = append(pkUpdate, abci.ValidatorUpdate{PubKey: oldTmPk, Power: 0})
		// add new tmPk with power
		pkUpdate = append(pkUpdate, abci.ValidatorUpdate{PubKey: newTmPk, Power: power})

		newConsAddr := sdk.ConsAddress(newPubKey.Address())
		// set new consensus address with validator
		k.SetValidatorOperatorByConsAddr(ctx, newConsAddr, validator.GetOperator())
		// record validator new consensus address
		k.SetValidatorOldConsensusAddr(ctx, valAddr, oldConsAddr)
		return false
	})

	// joint pkPower and pkUpdate
	newValUpdates := make([]abci.ValidatorUpdate, 0, len(pkPower)+len(pkUpdate))
	for tmPk, power := range pkPower {
		newValUpdates = append(newValUpdates, abci.ValidatorUpdate{PubKey: tmPk, Power: power})
	}
	newValUpdates = append(newValUpdates, pkUpdate...)
	return newValUpdates
}
