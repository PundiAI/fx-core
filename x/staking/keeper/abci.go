package keeper

import (
	"errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	"github.com/functionx/fx-core/v5/x/staking/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	valUpdates := staking.EndBlocker(ctx, k.Keeper)

	k.ConsensusProcess(ctx)
	return k.ValidatorUpdate(ctx, valUpdates)
}

func (k Keeper) ConsensusProcess(ctx sdk.Context) {
	// process end, remove old consensus key
	k.IteratorConsensusProcess(ctx, types.ProcessEnd, func(valAddr sdk.ValAddress, oldConsAddr sdk.ConsAddress) {
		k.DeleteConsensusProcess(ctx, valAddr, types.ProcessEnd)

		// remove validator by old consensus address
		k.RemoveValidatorConsAddr(ctx, oldConsAddr)

		// remove old consensus pubkey
		k.slashingKeeper.DeleteConsensusPubKey(ctx, oldConsAddr)

		// remove old validator signing info
		k.slashingKeeper.DeleteValidatorSigningInfo(ctx, oldConsAddr)

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeEndEditConsensusPubKey,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, valAddr.String()),
		))
	})

	// process start to end
	k.IteratorConsensusProcess(ctx, types.ProcessStart, func(valAddr sdk.ValAddress, oldConsAddr sdk.ConsAddress) {
		k.DeleteConsensusProcess(ctx, valAddr, types.ProcessStart)
		k.SetConsensusProcess(ctx, valAddr, oldConsAddr, types.ProcessEnd)
	})
}

func (k Keeper) ValidatorUpdate(ctx sdk.Context, valUpdates []abci.ValidatorUpdate) []abci.ValidatorUpdate {
	proposer := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	pkPower := make(map[crypto.PublicKey]int64, len(valUpdates))
	for _, valUpdate := range valUpdates {
		pkPower[valUpdate.PubKey] = valUpdate.Power
	}
	pkUpdate := make([]abci.ValidatorUpdate, 0, 50)

	k.IteratorConsensusPubKey(ctx, func(valAddr sdk.ValAddress, newPubKey cryptotypes.PubKey) {
		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			k.RemoveConsensusPubKey(ctx, valAddr)
			return
		}
		oldPubKey, err := validator.ConsPubKey()
		if err != nil {
			k.RemoveConsensusPubKey(ctx, valAddr)
			return
		}
		oldConsAddr := sdk.ConsAddress(oldPubKey.Address())
		// if validator is proposer, skip this block
		if oldConsAddr.Equals(proposer) {
			k.Logger(ctx).Info("validator is proposer, skip update", "address", valAddr.String())
			return
		}

		// no matter what happens next, clear new consensus pubkey
		k.RemoveConsensusPubKey(ctx, valAddr)

		cacheCtx, commit := ctx.CacheContext()
		// update validator pubkey
		if err = k.updateValidator(cacheCtx, validator, newPubKey); err != nil {
			k.Logger(ctx).Error("update validator", "address", valAddr.String(), "error", err.Error())
			return
		}
		// slash update
		if err = k.updateSlashing(cacheCtx, newPubKey, oldConsAddr); err != nil {
			k.Logger(ctx).Error("update slashing", "address", valAddr.String(), "error", err.Error())
			return
		}
		// new validator updates
		newValUpdates, err := k.updateABICValidator(cacheCtx, pkPower, validator, newPubKey, oldPubKey)
		if err != nil {
			k.Logger(ctx).Error("update abci validator", "address", valAddr.String(), "error", err.Error())
			return
		}
		pkUpdate = append(pkUpdate, newValUpdates...)

		// set consensus process start
		k.SetConsensusProcess(ctx, valAddr, oldConsAddr, types.ProcessStart)

		k.Logger(ctx).Info("update consensus pubkey", "address", valAddr.String(),
			"oldConsAddr", oldConsAddr.String(), "newConsAddr", sdk.ConsAddress(newPubKey.Address()).String())

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeStartEditConsensusPubKey,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, valAddr.String()),
			sdk.NewAttribute(types.AttributeOldConsAddr, oldConsAddr.String()),
			sdk.NewAttribute(types.AttributeNewConsAddr, sdk.ConsAddress(newPubKey.Address()).String()),
		))

		// commit cache context
		commit()
	})
	// joint pkPower and pkUpdate
	newValUpdates := make([]abci.ValidatorUpdate, 0, len(pkPower)+len(pkUpdate))
	for tmPk, power := range pkPower {
		newValUpdates = append(newValUpdates, abci.ValidatorUpdate{PubKey: tmPk, Power: power})
	}
	newValUpdates = append(newValUpdates, pkUpdate...)
	return newValUpdates
}

func (k Keeper) updateValidator(ctx sdk.Context, validator stakingtypes.Validator, newPubKey cryptotypes.PubKey) error {
	pkAny, err := codectypes.NewAnyWithValue(newPubKey)
	if err != nil {
		return err
	}
	validator.ConsensusPubkey = pkAny
	k.SetValidator(ctx, validator)

	// set new consensus address with validator
	k.SetValidatorConsAddr(ctx, sdk.ConsAddress(newPubKey.Address()), validator.GetOperator())
	return nil
}

func (k Keeper) updateSlashing(ctx sdk.Context, newPubKey cryptotypes.PubKey, oldConsAddr sdk.ConsAddress) error {
	// add new pubkey
	if err := k.slashingKeeper.AddPubkey(ctx, newPubKey); err != nil {
		return err
	}

	newConsAddr := sdk.ConsAddress(newPubKey.Address())
	// add signing info
	signingInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, oldConsAddr)
	if !found {
		return errors.New("expected signing info to be found")
	}
	signingInfo.Address = newConsAddr.String()
	k.slashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, signingInfo)

	return nil
}

func (k Keeper) updateABICValidator(ctx sdk.Context, pkPower map[crypto.PublicKey]int64, validator stakingtypes.Validator, newPubKey, oldPubKey cryptotypes.PubKey) ([]abci.ValidatorUpdate, error) {
	oldTmPk, err := cryptocodec.ToTmProtoPublicKey(oldPubKey)
	if err != nil {
		return nil, err
	}
	newTmPk, err := cryptocodec.ToTmProtoPublicKey(newPubKey)
	if err != nil {
		return nil, err
	}

	power, ok := pkPower[oldTmPk]
	// if power not found, use validator power
	if !ok {
		power = validator.ConsensusPower(k.PowerReduction(ctx))
	}
	// set old tmPk power to 0, remove
	oldTmPkUpdate := abci.ValidatorUpdate{PubKey: oldTmPk, Power: 0}
	// add new tmPk with power
	newTmPkUpdate := abci.ValidatorUpdate{PubKey: newTmPk, Power: power}
	return []abci.ValidatorUpdate{oldTmPkUpdate, newTmPkUpdate}, nil
}
