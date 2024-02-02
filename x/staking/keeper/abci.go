package keeper

import (
	"errors"
	"fmt"
	"strconv"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"

	fxstakingtypes "github.com/functionx/fx-core/v7/x/staking/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	// staking EndBlocker
	valUpdates := staking.EndBlocker(ctx, k.Keeper)

	// process pk update previous block
	k.ConsensusProcess(ctx)

	// validators update current block, delayed next block
	if len(valUpdates) > 0 {
		return valUpdates
	}

	// update validator consensus pubkey
	return k.ConsensusPubKeyUpdate(ctx)
}

func (k Keeper) ConsensusProcess(ctx sdk.Context) {
	k.IteratorConsensusProcess(ctx, fxstakingtypes.ProcessEnd, func(valAddr sdk.ValAddress, pkBytes []byte) {
		// update signing info and remove old consensus pubkey
		if err := k.processEnd(ctx, valAddr, pkBytes); err != nil {
			panic(err)
		}
	})

	k.IteratorConsensusProcess(ctx, fxstakingtypes.ProcessStart, func(valAddr sdk.ValAddress, pkBytes []byte) {
		// update signing info
		if err := k.processStart(ctx, valAddr, pkBytes); err != nil {
			panic(err)
		}
	})
}

func (k Keeper) processEnd(ctx sdk.Context, valAddr sdk.ValAddress, pkBytes []byte) error {
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
		k.slashingKeeper.DeleteConsensusPubKey(ctx, oldConsAddr)
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

func (k Keeper) processStart(ctx sdk.Context, valAddr sdk.ValAddress, pkBytes []byte) error {
	k.DeleteConsensusProcess(ctx, valAddr, fxstakingtypes.ProcessStart)

	var oldPubKey cryptotypes.PubKey
	if err := k.cdc.UnmarshalInterfaceJSON(pkBytes, &oldPubKey); err != nil {
		return fmt.Errorf("invalid pubkey")
	}
	oldConsAddr := sdk.ConsAddress(oldPubKey.Address())

	// set process end and delete process start
	if err := k.SetConsensusProcess(ctx, valAddr, oldPubKey, fxstakingtypes.ProcessEnd); err != nil {
		return err
	}

	// update validator new pk signing info
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		// NOTE: validator not found, because it deleted in below case
		// 1. validator unbonded and undelegate all (tx)
		// 2. unbonding to unbonded and validator share is zero (end block)
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
	if newSigningInfo.Tombstoned {
		newSigningInfo.IndexOffset = signingInfo.IndexOffset
		newSigningInfo.MissedBlocksCounter = signingInfo.MissedBlocksCounter
	} else {
		newSigningInfo = signingInfo
		newSigningInfo.Address = consAddr.String()
	}
	k.slashingKeeper.SetValidatorSigningInfo(ctx, consAddr, newSigningInfo)
	return nil
}

func (k Keeper) ConsensusPubKeyUpdate(ctx sdk.Context) []abci.ValidatorUpdate {
	valUpdates := make([]abci.ValidatorUpdate, 0, 50)
	k.IteratorConsensusPubKey(ctx, func(valAddr sdk.ValAddress, pkBytes []byte) bool {
		// no matter what happens, clear new consensus pubkey
		k.RemoveConsensusPubKey(ctx, valAddr)

		cacheCtx, commit := ctx.CacheContext()
		pkUpdates, err := k.updateConsensusPubKey(cacheCtx, valAddr, pkBytes)
		if err == nil {
			valUpdates = append(valUpdates, pkUpdates...)
			commit()
		}

		// event
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			fxstakingtypes.EventTypeEditingConsensusPubKey,
			sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeResult, strconv.FormatBool(err == nil)),
		))
		return false
	})
	return valUpdates
}

func (k Keeper) updateConsensusPubKey(ctx sdk.Context, valAddr sdk.ValAddress, pkBytes []byte) ([]abci.ValidatorUpdate, error) {
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		// NOTE: validator not found, because it deleted in below case
		// 1. validator unbonded and undelegate all (tx)
		// 2. unbonding to unbonded and validator share is zero (end block)
		return nil, errors.New("validator not found")
	}
	oldPubKey, err := validator.ConsPubKey()
	if err != nil {
		return nil, fmt.Errorf("invalid old consensus pubkey")
	}
	oldConsAddr := sdk.ConsAddress(oldPubKey.Address())

	var newPubKey cryptotypes.PubKey
	if err = k.cdc.UnmarshalInterfaceJSON(pkBytes, &newPubKey); err != nil {
		return nil, err
	}

	// validator pubkey and consaddr
	if err = k.updateValidator(ctx, validator, newPubKey); err != nil {
		return nil, err
	}
	// pubkey and signing info
	if err = k.updateSlashing(ctx, newPubKey, oldConsAddr); err != nil {
		return nil, err
	}

	// validator unjailed and bonded
	var valUpdates []abci.ValidatorUpdate
	if validator.IsBonded() {
		// new validator updates
		valUpdates, err = k.updateABCIValidator(ctx, validator, newPubKey, oldPubKey)
		if err != nil {
			return nil, err
		}
	}

	// set consensus process, next 2 block will delete old consensus pubkey
	if err = k.SetConsensusProcess(ctx, valAddr, oldPubKey, fxstakingtypes.ProcessStart); err != nil {
		return nil, err
	}
	return valUpdates, nil
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
		return fmt.Errorf("consensus address %s not found signing info", oldConsAddr.String())
	}
	signingInfo.Address = newConsAddr.String()
	k.slashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, signingInfo)
	return nil
}

func (k Keeper) updateABCIValidator(ctx sdk.Context, val types.Validator, newPk, oldPk cryptotypes.PubKey) ([]abci.ValidatorUpdate, error) {
	oldTmProtoPk, err := cryptocodec.ToTmProtoPublicKey(oldPk)
	if err != nil {
		return nil, err
	}
	newTmProtoPk, err := cryptocodec.ToTmProtoPublicKey(newPk)
	if err != nil {
		return nil, err
	}
	valPower := val.ConsensusPower(k.PowerReduction(ctx))
	// set old pk power to 0
	oldPkUpdate := abci.ValidatorUpdate{PubKey: oldTmProtoPk, Power: 0}
	// add new pk with power
	newPkUpdate := abci.ValidatorUpdate{PubKey: newTmProtoPk, Power: valPower}

	return []abci.ValidatorUpdate{oldPkUpdate, newPkUpdate}, nil
}
