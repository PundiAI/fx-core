package keeper

import (
	"errors"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"

	fxstakingtypes "github.com/functionx/fx-core/v5/x/staking/types"
)

func (k *Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	// staking EndBlocker
	valUpdates := staking.EndBlocker(ctx, k.Keeper)
	// convert valUpdate and lastCommit to map
	pkPowerUpdate := make(map[string]int64, len(valUpdates))
	for _, valUpdate := range valUpdates {
		pkPowerUpdate[valUpdate.PubKey.String()] = valUpdate.Power
	}

	// consensus process start and end
	k.ConsensusProcess(ctx)

	// update validator consensus pubkey
	return k.ValidatorUpdate(ctx, valUpdates, pkPowerUpdate)
}

func (k Keeper) ConsensusProcess(ctx sdk.Context) {
	k.IteratorConsensusProcess(ctx, fxstakingtypes.ProcessEnd, func(valAddr sdk.ValAddress, pkBytes []byte) {
		if err := k.consesnsusProcessEnd(ctx, valAddr, pkBytes); err != nil {
			panic(err)
		}
	})

	k.IteratorConsensusProcess(ctx, fxstakingtypes.ProcessStart, func(valAddr sdk.ValAddress, pkBytes []byte) {
		if err := k.consensusProcessStart(ctx, valAddr, pkBytes); err != nil {
			panic(err)
		}
	})
}

func (k Keeper) consesnsusProcessEnd(ctx sdk.Context, valAddr sdk.ValAddress, pkBytes []byte) error {
	k.DeleteConsensusProcess(ctx, valAddr, fxstakingtypes.ProcessEnd)

	var oldPubKey cryptotypes.PubKey
	if err := k.cdc.UnmarshalInterfaceJSON(pkBytes, &oldPubKey); err != nil {
		return fmt.Errorf("invalid pubkey")
	}
	oldConsAddr := sdk.ConsAddress(oldPubKey.Address())

	// old pk signing info
	oldSigningInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, oldConsAddr)
	if !found {
		return fmt.Errorf("validator %s not found signing info", oldConsAddr.String())
	}

	// remove old consensus pubkey
	k.RemoveValidatorConsAddr(ctx, oldConsAddr)
	k.slashingKeeper.DeleteConsensusPubKey(ctx, oldConsAddr)
	k.slashingKeeper.DeleteValidatorSigningInfo(ctx, oldConsAddr)

	// update validator new pk signing info
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		// NOTE: validator not found, because it deleted in below case
		// 1. validator unbonded and undelegate all (tx)
		// 2. unbonding to unbonded and validator share is zero (end block)
		// remove old consensus pubkey and return
		k.Logger(ctx).Info("validator not found", "address", valAddr.String(), "process", "end")
		return nil
	}
	newConsAddr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	if err = k.updateSigningInfo(ctx, newConsAddr, oldSigningInfo); err != nil {
		return err
	}

	// success event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		fxstakingtypes.EventTypeEditedConsensusPubKey,
		sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
	))

	return nil
}

func (k Keeper) consensusProcessStart(ctx sdk.Context, valAddr sdk.ValAddress, pkBytes []byte) error {
	var oldPubKey cryptotypes.PubKey
	if err := k.cdc.UnmarshalInterfaceJSON(pkBytes, &oldPubKey); err != nil {
		return fmt.Errorf("invalid pubkey")
	}
	oldConsAddr := sdk.ConsAddress(oldPubKey.Address())

	// set process end and delete process start
	if err := k.SetConsensusProcess(ctx, valAddr, oldPubKey, fxstakingtypes.ProcessEnd); err != nil {
		return err
	}
	k.DeleteConsensusProcess(ctx, valAddr, fxstakingtypes.ProcessStart)

	// update validator new pk signing info
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		// NOTE: validator not found, because it deleted in below case
		// 1. validator unbonded and undelegate all (tx)
		// 2. unbonding to unbonded and validator share is zero (end block)
		// return nil and delete old consensus pubkey in end process
		k.Logger(ctx).Info("validator not found", "address", valAddr.String(), "process", "start")
		return nil
	}

	newConsAddr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	oldSigningInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, oldConsAddr)
	if !found {
		return fmt.Errorf("validator %s not found signing info", oldConsAddr.String())
	}
	return k.updateSigningInfo(ctx, newConsAddr, oldSigningInfo)
}

func (k Keeper) updateSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress, signingInfo slashingtypes.ValidatorSigningInfo) error {
	newSigningInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		return fmt.Errorf("validator %s not found signing info", consAddr.String())
	}
	// double sign
	if newSigningInfo.JailedUntil == evidencetypes.DoubleSignJailEndTime {
		newSigningInfo.IndexOffset = signingInfo.IndexOffset
		newSigningInfo.MissedBlocksCounter = signingInfo.MissedBlocksCounter
	} else {
		newSigningInfo = signingInfo
		newSigningInfo.Address = consAddr.String()
	}
	k.slashingKeeper.SetValidatorSigningInfo(ctx, consAddr, newSigningInfo)
	return nil
}

func (k Keeper) ValidatorUpdate(ctx sdk.Context, valUpdates []abci.ValidatorUpdate, pkPowerUpdate map[string]int64) []abci.ValidatorUpdate {
	pkUpdate := make([]abci.ValidatorUpdate, 0, 50)

	k.IteratorConsensusPubKey(ctx, func(valAddr sdk.ValAddress, pkBytes []byte) bool {
		// no matter what happens, clear new consensus pubkey
		k.RemoveConsensusPubKey(ctx, valAddr)

		// check validator exist
		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			// NOTE: validator not found, because it deleted in below case
			// 1. validator unbonded and undelegate all (tx)
			// 2. unbonding to unbonded and validator share is zero (end block)
			k.Logger(ctx).Error("validator not found", "address", valAddr.String(), "process", "update")
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
		newValUpdates, err := k.updateABICValidator(cacheCtx, pkPowerUpdate, validator, newPubKey, oldPubKey)
		if err != nil {
			k.Logger(ctx).Error("update abci validator", "address", valAddr.String(), "error", err.Error())
			return false
		}
		// set consensus process start
		if err := k.SetConsensusProcess(ctx, valAddr, oldPubKey, fxstakingtypes.ProcessStart); err != nil {
			return false
		}

		k.Logger(ctx).Info("update consensus pubkey", "address", valAddr.String(),
			"oldConsAddr", oldConsAddr.String(), "newConsAddr", sdk.ConsAddress(newPubKey.Address()).String())

		// event
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			fxstakingtypes.EventTypeEditingConsensusPubKey,
			sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeOldConsAddr, oldConsAddr.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeNewConsAddr, sdk.ConsAddress(newPubKey.Address()).String()),
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

func (k Keeper) updateValidator(ctx sdk.Context, validator types.Validator, newPubKey cryptotypes.PubKey) error {
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
		// NOTE: validator create but not bonded enough token
		return fmt.Errorf("validator %s not found signing info", oldConsAddr.String())
	}
	signingInfo.Address = newConsAddr.String()
	k.slashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, signingInfo)
	return nil
}

//gocyclo:ignore
func (k Keeper) updateABICValidator(ctx sdk.Context, pkUpdate map[string]int64, val types.Validator, newPk, oldPk cryptotypes.PubKey) ([]abci.ValidatorUpdate, error) {
	oldTmProtoPk, err := cryptocodec.ToTmProtoPublicKey(oldPk)
	if err != nil {
		return nil, err
	}
	newTmProtoPk, err := cryptocodec.ToTmProtoPublicKey(newPk)
	if err != nil {
		return nil, err
	}
	power, ok := pkUpdate[oldTmProtoPk.String()]
	// if power not found, validator not update current block, cal validator power
	if !ok {
		power = val.ConsensusPower(k.PowerReduction(ctx))
	} else {
		// remove old pk power
		delete(pkUpdate, oldTmProtoPk.String())
	}
	// set old pk power to 0
	oldPkUpdate := abci.ValidatorUpdate{PubKey: oldTmProtoPk, Power: 0}
	// add new pk with power
	newPkUpdate := abci.ValidatorUpdate{PubKey: newTmProtoPk, Power: power}

	/*
		a1: 1-block 2-block-edit 3-block					------ oldpk-0,newpk-power
		a2: 1-block 2-block-(jailed|edit) 3-block 4-unblock	------ oldpk-0

		b1: 1-unblock(jailed) 2-unblock-edit 3-unblock						------ nil
		b2: 1-unblock(inactive) 2-unblock-edit 3-unblock					------ nil
		b3: 1-unblock 2-unblock-(edit|(unjailed/active)) 3-unblock 4-block	------ newpk-power

		c1: 1-block-jailed 2-block-edit 3-unblock						------ nil
		c2: 1-block-jailed 2-block-(edit|unjailed) 3-unblock 4-block	------ newpk-power

		d1: 1-unblock-(unjailed/active) 2-unblock-edit 3-block						------ oldpk-0,newpk-power
		d2: 1-unblock-(unjailed/active) 2-unblock-(edit|jailed) 3-block 4-unblock	------ oldpk-0

		e1: 1-block-jailed 2-block 3-unblock-edit 4-unblock						------ nil
		e2: 1-block-jailed 2-block 3-unblock-(edit|unjailed) 4-unblock 5-block	------ newpk-power

		f1: 1-unblock-(unjailed/active) 2-unblock 3-block-edit 4-block						------ oldpk-0,newpk-power
		f2: 1-unblock-(unjailed/active) 2-unblock 3-block-(edit|jailed) 4-block 5-unblock	------ oldpk-0


		// validator status
		ok=true,power==0,jailed=true	// jailed current block
		ok=true,power==0,jailed=false	// impossible
		ok=true,power!=0,jailed=true	// impossible
		ok=true,power!=0,jailed=false	// unjailed/active current block
		ok=false,power==0,jailed=true	// jailed previous block
		ok=false,power==0,jailed=false	// inactive validator
		ok=false,power!=0,jailed=true	// impossible
		ok=false,power!=0,jailed=false	// online validator
	*/

	// validator jailed current block // a2,d2,f2
	if ok && power == 0 && val.Jailed {
		return []abci.ValidatorUpdate{oldPkUpdate}, nil
	}
	// validator unjailed/active current block // b3,c2,e2
	if ok && power != 0 && !val.Jailed {
		return []abci.ValidatorUpdate{newPkUpdate}, nil
	}
	// validator jailed previous block // b1,c1,e1
	if !ok && power == 0 && val.Jailed {
		return []abci.ValidatorUpdate{}, nil
	}
	// validator inactive // b2
	if !ok && power == 0 && !val.Jailed {
		return []abci.ValidatorUpdate{}, nil
	}
	// validator online // a1,d1,f1
	if !ok && power != 0 && !val.Jailed {
		return []abci.ValidatorUpdate{oldPkUpdate, newPkUpdate}, nil
	}
	return nil, errors.New("impossible case")
}
