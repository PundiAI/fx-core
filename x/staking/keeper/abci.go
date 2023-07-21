package keeper

import (
	"errors"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	cryptoenc "github.com/tendermint/tendermint/crypto/encoding"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	"github.com/functionx/fx-core/v5/x/staking/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	staking.BeginBlocker(ctx, k.Keeper)
	k.lastCommit = req.LastCommitInfo.GetVotes()
}

func (k *Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	// staking EndBlocker
	valUpdates := staking.EndBlocker(ctx, k.Keeper)
	if len(k.lastCommit) == 0 {
		return valUpdates
	}

	// clear lastCommit after EndBlocker
	defer func() { k.lastCommit = make([]abci.VoteInfo, 0) }()

	// convert valUpdate and lastCommit to map
	pkPowerUpdate := make(map[string]int64, len(valUpdates))
	for _, valUpdate := range valUpdates {
		pkPowerUpdate[valUpdate.PubKey.String()] = valUpdate.Power
	}
	lastVote := make(map[string]bool, len(k.lastCommit))
	for _, voteInfo := range k.lastCommit {
		lastVote[string(voteInfo.Validator.Address)] = true
	}

	// consensus process start and end
	k.ConsensusProcess(ctx, pkPowerUpdate, lastVote)

	// update validator consensus pubkey
	return k.ValidatorUpdate(ctx, valUpdates, pkPowerUpdate, lastVote)
}

func (k Keeper) ConsensusProcess(ctx sdk.Context, pkPowerUpdate map[string]int64, lastVote map[string]bool) {
	k.IteratorConsensusProcess(ctx, types.ProcessEnd, func(valAddr sdk.ValAddress, pkBytes []byte) {
		err := k.consesnsusProcessEnd(ctx, lastVote, valAddr, pkBytes)
		if err != nil {
			panic(err)
		}
	})

	k.IteratorConsensusProcess(ctx, types.ProcessStart, func(valAddr sdk.ValAddress, pkBytes []byte) {
		err := k.consensusProcessStart(ctx, pkPowerUpdate, lastVote, valAddr, pkBytes)
		if err != nil {
			panic(err)
		}
	})
}

func (k Keeper) consesnsusProcessEnd(ctx sdk.Context, lastVote map[string]bool, valAddr sdk.ValAddress, pkBytes []byte) error {
	k.DeleteConsensusProcess(ctx, valAddr, types.ProcessEnd)

	_, _, newConsAddr, err := k.getValidatorKey(ctx, valAddr)
	if err != nil {
		return err
	}
	_, oldTmPk, oldConsAddr, err := k.unmarshalPubKey(pkBytes)
	if err != nil {
		return err
	}

	// case1: validator jailed current block
	// case2: validator jailed previous block
	// impossible case: power == 0 && !validator.Jailed, power != 0 && validator.Jailed
	// case3: validator unjailed current block
	// case4: validator always online

	lastVoted := lastVote[string(oldTmPk.Address())]
	if lastVoted {
		if err := k.updateSigningInfo(ctx, oldConsAddr, newConsAddr); err != nil {
			panic(err)
		}
	}

	// remove validator by old consensus address
	k.RemoveValidatorConsAddr(ctx, oldConsAddr)

	// remove old consensus pubkey
	k.slashingKeeper.DeleteConsensusPubKey(ctx, oldConsAddr)

	// remove old validator signing info
	k.slashingKeeper.DeleteValidatorSigningInfo(ctx, oldConsAddr)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditedConsensusPubKey,
		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, valAddr.String()),
	))

	return nil
}

func (k Keeper) consensusProcessStart(
	ctx sdk.Context,
	pkPowerUpdate map[string]int64,
	lastVote map[string]bool,
	valAddr sdk.ValAddress,
	pkBytes []byte,
) error {
	k.DeleteConsensusProcess(ctx, valAddr, types.ProcessStart)

	validator, newTmProtoPk, newConsAddr, err := k.getValidatorKey(ctx, valAddr)
	if err != nil {
		return err
	}
	oldPubKey, oldTmPk, oldConsAddr, err := k.unmarshalPubKey(pkBytes)
	if err != nil {
		return err
	}

	// set process end first
	if err = k.SetConsensusProcess(ctx, valAddr, oldPubKey, types.ProcessEnd); err != nil {
		return err
	}

	power, ok := pkPowerUpdate[newTmProtoPk.String()] // only new pk, validator consensus pk update previous block
	if !ok {
		power = validator.ConsensusPower(k.PowerReduction(ctx))
	}

	// case1: validator jailed current block(missed blocks/double sign/low power)
	if ok && power == 0 && validator.Jailed {
		return k.updateSigningInfo(ctx, oldConsAddr, newConsAddr)
	}

	lastVoted := lastVote[string(oldTmPk.Address())] // only old pk, validator new pk vote in next block

	// case2: validator jailed previous block
	if power == 0 && validator.Jailed {
		// jailed last-1 or last block
		if lastVoted {
			return k.updateSigningInfo(ctx, oldConsAddr, newConsAddr)
		}
		return nil
	}

	// impossible case: power != 0 && validator jailed, power == 0 && validator unjailed
	// power !=0 && validator unjailed

	// case3: validator jailed previous block, unjailed current block
	if !lastVoted && !validator.Jailed {
		return nil
	}

	// case4: validator always online or jailed(low power) last block, unjailed current block
	// maybe miss block current height, update signing info
	return k.updateSigningInfo(ctx, oldConsAddr, newConsAddr)
}

func (k Keeper) getValidatorKey(ctx sdk.Context, valAddr sdk.ValAddress) (stakingtypes.Validator, crypto.PublicKey, sdk.ConsAddress, error) {
	// todo validator not found ??
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.Validator{}, crypto.PublicKey{}, nil,
			fmt.Errorf("validator %s not found", valAddr.String())
	}
	pubKey, err := validator.ConsPubKey()
	if err != nil {
		return stakingtypes.Validator{}, crypto.PublicKey{}, nil,
			fmt.Errorf("invalid validator %s pubkey", valAddr.String())
	}
	tmProtoPk, err := cryptocodec.ToTmProtoPublicKey(pubKey)
	if err != nil {
		return stakingtypes.Validator{}, crypto.PublicKey{}, nil,
			fmt.Errorf("invalid validator %s pubkey", valAddr.String())
	}
	consAddr := sdk.ConsAddress(pubKey.Address())
	return validator, tmProtoPk, consAddr, nil
}

func (k Keeper) unmarshalPubKey(pkBytes []byte) (cryptotypes.PubKey, tmcrypto.PubKey, sdk.ConsAddress, error) {
	var pubKey cryptotypes.PubKey
	if err := k.cdc.UnmarshalInterfaceJSON(pkBytes, &pubKey); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid pubkey")
	}
	tmProtoPk, err := cryptocodec.ToTmProtoPublicKey(pubKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid pubkey")
	}
	tmPk, err := cryptoenc.PubKeyFromProto(tmProtoPk)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid pubkey")
	}
	consAddr := sdk.ConsAddress(pubKey.Address())
	return pubKey, tmPk, consAddr, nil
}

func (k Keeper) updateSigningInfo(ctx sdk.Context, oldConsAddr, newConsAddr sdk.ConsAddress) error {
	oldSigningInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, oldConsAddr)
	if !found {
		return fmt.Errorf("validator %s not found signing info", oldConsAddr.String())
	}
	newSigningInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, newConsAddr)
	if !found {
		return fmt.Errorf("validator %s not found signing info", newConsAddr.String())
	}
	// double sign
	if newSigningInfo.JailedUntil == evidencetypes.DoubleSignJailEndTime {
		newSigningInfo.IndexOffset = oldSigningInfo.IndexOffset
		newSigningInfo.MissedBlocksCounter = oldSigningInfo.MissedBlocksCounter
	} else {
		newSigningInfo = oldSigningInfo
		newSigningInfo.Address = newConsAddr.String()
	}
	k.slashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newSigningInfo)
	return nil
}

func (k Keeper) ValidatorUpdate(ctx sdk.Context, valUpdates []abci.ValidatorUpdate, pkPowerUpdate map[string]int64, lastVote map[string]bool) []abci.ValidatorUpdate {
	pkUpdate := make([]abci.ValidatorUpdate, 0, 50)

	k.IteratorConsensusPubKey(ctx, func(valAddr sdk.ValAddress, pkBytes []byte) bool {
		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			k.Logger(ctx).Error("validator not found", "address", valAddr.String())
			k.RemoveConsensusPubKey(ctx, valAddr)
			return false
		}
		oldPubKey, err := validator.ConsPubKey()
		if err != nil {
			k.Logger(ctx).Error("invalid consensus pubkey", "address", valAddr.String(), "error", err.Error())
			k.RemoveConsensusPubKey(ctx, valAddr)
			return false
		}
		oldConsAddr := sdk.ConsAddress(oldPubKey.Address())

		// no matter what happens next, clear new consensus pubkey
		k.RemoveConsensusPubKey(ctx, valAddr)

		// unmarshal failed, remove new consensus pubkey
		var newPubKey cryptotypes.PubKey
		if err = k.cdc.UnmarshalInterfaceJSON(pkBytes, &newPubKey); err != nil {
			k.Logger(ctx).Error("unmarshal new consensus pubkey", "validator", valAddr.String(), "err", err.Error())
			return false
		}

		cacheCtx, commit := ctx.CacheContext()
		// update validator pubkey
		if err = k.updateValidator(cacheCtx, validator, newPubKey); err != nil {
			k.Logger(ctx).Error("update validator", "address", valAddr.String(), "error", err.Error())
			return false
		}
		// slash update
		if err = k.updateSlashing(cacheCtx, newPubKey, oldConsAddr); err != nil {
			k.Logger(ctx).Error("update slashing", "address", valAddr.String(), "error", err.Error())
			return false
		}
		// new validator updates
		newValUpdates, err := k.updateABICValidator(cacheCtx, pkPowerUpdate, lastVote, validator, newPubKey, oldPubKey)
		if err != nil {
			k.Logger(ctx).Error("update abci validator", "address", valAddr.String(), "error", err.Error())
			return false
		}
		// set consensus process start
		if err := k.SetConsensusProcess(ctx, valAddr, oldPubKey, types.ProcessStart); err != nil {
			return false
		}

		k.Logger(ctx).Info("update consensus pubkey", "address", valAddr.String(),
			"oldConsAddr", oldConsAddr.String(), "newConsAddr", sdk.ConsAddress(newPubKey.Address()).String())
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeEditingConsensusPubKey,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, valAddr.String()),
			sdk.NewAttribute(types.AttributeOldConsAddr, oldConsAddr.String()),
			sdk.NewAttribute(types.AttributeNewConsAddr, sdk.ConsAddress(newPubKey.Address()).String()),
		))
		// commit cache context
		commit()

		pkUpdate = append(pkUpdate, newValUpdates...)
		return false
	})
	// joint pkPowerUpdate and pkUpdate
	newValUpdates := make([]abci.ValidatorUpdate, 0, len(pkPowerUpdate)+len(pkUpdate))
	for _, vu := range valUpdates {
		if power, ok := pkPowerUpdate[vu.PubKey.String()]; ok {
			newValUpdates = append(newValUpdates, abci.ValidatorUpdate{PubKey: vu.PubKey, Power: power})
		}
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
		return errors.New("expected signing info not found")
	}
	signingInfo.Address = newConsAddr.String()
	k.slashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, signingInfo)
	return nil
}

func (k Keeper) updateABICValidator(
	ctx sdk.Context,
	pkPower map[string]int64,
	lastVote map[string]bool,
	validator stakingtypes.Validator,
	newPubKey, oldPubKey cryptotypes.PubKey,
) ([]abci.ValidatorUpdate, error) {
	oldTmProtoPk, err := cryptocodec.ToTmProtoPublicKey(oldPubKey)
	if err != nil {
		return nil, err
	}
	oldTmPk, err := cryptoenc.PubKeyFromProto(oldTmProtoPk)
	if err != nil {
		return nil, err
	}
	newTmProtoPk, err := cryptocodec.ToTmProtoPublicKey(newPubKey)
	if err != nil {
		return nil, err
	}
	power, ok := pkPower[oldTmProtoPk.String()]
	// if power not found, validator not update current block, cal validator power
	if !ok {
		power = validator.ConsensusPower(k.PowerReduction(ctx))
	} else {
		// remove old pk power
		delete(pkPower, oldTmProtoPk.String())
	}
	// set old pk power to 0
	oldPkUpdate := abci.ValidatorUpdate{PubKey: oldTmProtoPk, Power: 0}
	// add new pk with power
	newPkUpdate := abci.ValidatorUpdate{PubKey: newTmProtoPk, Power: power}

	// case1: validator jailed current block
	if ok && power == 0 && validator.Jailed {
		return []abci.ValidatorUpdate{oldPkUpdate}, nil
	}

	// case2: validator jailed previous block
	if power == 0 && validator.Jailed {
		return []abci.ValidatorUpdate{}, nil
	}

	// impossible case: power == 0 && validator unjailed, power != 0 && validator jailed
	// power !=0 && validator unjailed

	// case3: validator jailed previous block, unjailed current block
	lastVoted := lastVote[string(oldTmPk.Address())]
	if !lastVoted && !validator.Jailed {
		return []abci.ValidatorUpdate{newPkUpdate}, nil
	}

	// case4: validator already online
	return []abci.ValidatorUpdate{oldPkUpdate, newPkUpdate}, nil
}
