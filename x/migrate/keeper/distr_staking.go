package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/x/migrate/types"
)

type DistrStakingMigrate struct {
	distrKey      storetypes.StoreKey
	stakingKey    storetypes.StoreKey
	stakingKeeper types.StakingKeeper
}

func NewDistrStakingMigrate(distrKey, stakingKey storetypes.StoreKey, stakingKeeper types.StakingKeeper) MigrateI {
	return &DistrStakingMigrate{
		distrKey:      distrKey,
		stakingKey:    stakingKey,
		stakingKeeper: stakingKeeper,
	}
}

func (m *DistrStakingMigrate) Validate(ctx sdk.Context, _ codec.BinaryCodec, from sdk.AccAddress, to common.Address) error {
	// check validator
	if _, err := m.stakingKeeper.GetValidator(ctx, sdk.ValAddress(from)); err == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("can not migrate, %s is the validator address", from.String())
	}
	if _, err := m.stakingKeeper.GetValidator(ctx, to.Bytes()); err == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("can not migrate, %s is the validator address", to.String())
	}

	// check delegation
	delegations, err := m.stakingKeeper.GetDelegatorDelegations(ctx, to.Bytes(), 1)
	if err != nil {
		return err
	}
	if len(delegations) > 0 {
		return sdkerrors.ErrInvalidAddress.Wrapf("can not migrate, address %s has delegation record", to.String())
	}

	// check undelegatetion
	undelegations, err := m.stakingKeeper.GetUnbondingDelegations(ctx, to.Bytes(), 1)
	if err != nil {
		return err
	}
	if len(undelegations) > 0 {
		return sdkerrors.ErrInvalidAddress.Wrapf("can not migrate, address %s has undelegate record", to.String())
	}

	// check redelegation
	redelegations, err := m.stakingKeeper.GetRedelegations(ctx, to.Bytes(), 1)
	if err != nil {
		return err
	}
	if len(redelegations) > 0 {
		return sdkerrors.ErrInvalidAddress.Wrapf("can not migrate, address %s has redelegation record", to.String())
	}

	return nil
}

//nolint:gocyclo // copy from cosmos-sdk
func (m *DistrStakingMigrate) Execute(ctx sdk.Context, cdc codec.BinaryCodec, from sdk.AccAddress, to common.Address) error {
	stakingStore := ctx.KVStore(m.stakingKey)
	distrStore := ctx.KVStore(m.distrKey)

	events := make([]sdk.Event, 0, 10)

	// migrate delegate info
	delegateIterator := storetypes.KVStorePrefixIterator(stakingStore, stakingtypes.GetDelegationsKey(from))
	defer delegateIterator.Close()
	for ; delegateIterator.Valid(); delegateIterator.Next() {
		info := stakingtypes.MustUnmarshalDelegation(cdc, delegateIterator.Value())

		// distribution starting info
		validatorAddrStr := info.GetValidatorAddr()
		validatorAddr, err := m.stakingKeeper.ValidatorAddressCodec().StringToBytes(validatorAddrStr)
		if err != nil {
			return err
		}
		key := distrtypes.GetDelegatorStartingInfoKey(validatorAddr, from)
		startingInfo := distrStore.Get(key)
		distrStore.Delete(key)
		distrStore.Set(distrtypes.GetDelegatorStartingInfoKey(validatorAddr, to.Bytes()), startingInfo)

		// staking delegate
		info.DelegatorAddress = sdk.AccAddress(to.Bytes()).String()
		stakingStore.Delete(delegateIterator.Key())
		stakingStore.Set(stakingtypes.GetDelegationKey(to.Bytes(), validatorAddr), stakingtypes.MustMarshalDelegation(cdc, info))

		events = append(events,
			sdk.NewEvent(
				types.EventTypeMigrateStakingDelegate,
				sdk.NewAttribute(types.AttributeKeyValidatorAddr, sdk.ValAddress(validatorAddr).String()),
			),
		)
	}

	// migrate unbonding delegation
	unbondingDelegationIterator := storetypes.KVStorePrefixIterator(stakingStore, stakingtypes.GetUBDsKey(from))
	defer unbondingDelegationIterator.Close()
	for ; unbondingDelegationIterator.Valid(); unbondingDelegationIterator.Next() {
		ubd := stakingtypes.MustUnmarshalUBD(cdc, unbondingDelegationIterator.Value())
		ubd.DelegatorAddress = sdk.AccAddress(to.Bytes()).String()

		valAddr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		stakingStore.Delete(unbondingDelegationIterator.Key())
		stakingStore.Set(stakingtypes.GetUBDKey(to.Bytes(), valAddr), stakingtypes.MustMarshalUBD(cdc, ubd))

		stakingStore.Delete(stakingtypes.GetUBDByValIndexKey(from, valAddr))
		stakingStore.Set(stakingtypes.GetUBDByValIndexKey(to.Bytes(), valAddr), []byte{})

		// migrate unbonding queue
		for _, entry := range ubd.Entries {
			var ubdFlag bool
			dvPairs, err := m.stakingKeeper.GetUBDQueueTimeSlice(ctx, entry.CompletionTime)
			if err != nil {
				panic(err)
			}
			for i := range dvPairs {
				if dvPairs[i].DelegatorAddress == from.String() {
					dvPairs[i].DelegatorAddress = sdk.AccAddress(to.Bytes()).String()
					ubdFlag = true
				}
			}
			if ubdFlag {
				key := stakingtypes.GetUnbondingDelegationTimeKey(entry.CompletionTime)
				value := cdc.MustMarshal(&stakingtypes.DVPairs{Pairs: dvPairs})
				stakingStore.Set(key, value)
			}
		}

		events = append(events,
			sdk.NewEvent(
				types.EventTypeMigrateStakingUndelegate,
				sdk.NewAttribute(types.AttributeKeyValidatorAddr, valAddr.String()),
			),
		)
	}

	// migrate redelegate
	redelegateIterator := storetypes.KVStorePrefixIterator(stakingStore, stakingtypes.GetREDsKey(from))
	defer redelegateIterator.Close()
	for ; redelegateIterator.Valid(); redelegateIterator.Next() {
		red := stakingtypes.MustUnmarshalRED(cdc, redelegateIterator.Value())
		red.DelegatorAddress = sdk.AccAddress(to.Bytes()).String()

		valSrcAddr, err := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}

		stakingStore.Delete(redelegateIterator.Key())
		stakingStore.Set(stakingtypes.GetREDKey(to.Bytes(), valSrcAddr, valDstAddr), stakingtypes.MustMarshalRED(cdc, red))

		stakingStore.Delete(stakingtypes.GetREDByValSrcIndexKey(from, valSrcAddr, valDstAddr))
		stakingStore.Set(stakingtypes.GetREDByValSrcIndexKey(to.Bytes(), valSrcAddr, valDstAddr), []byte{})

		stakingStore.Delete(stakingtypes.GetREDByValDstIndexKey(from, valSrcAddr, valDstAddr))
		stakingStore.Set(stakingtypes.GetREDByValDstIndexKey(to.Bytes(), valSrcAddr, valDstAddr), []byte{})

		// migrate redelegate queue
		for _, entry := range red.Entries {
			var redFlag bool
			redQueue, err := m.stakingKeeper.GetRedelegationQueueTimeSlice(ctx, entry.CompletionTime)
			if err != nil {
				panic(err)
			}
			for i := range redQueue {
				if redQueue[i].DelegatorAddress == from.String() {
					redQueue[i].DelegatorAddress = sdk.AccAddress(to.Bytes()).String()
					redFlag = true
				}
			}
			if redFlag {
				key := stakingtypes.GetRedelegationTimeKey(entry.CompletionTime)
				value := cdc.MustMarshal(&stakingtypes.DVVTriplets{Triplets: redQueue})
				stakingStore.Set(key, value)
			}
		}

		events = append(events,
			sdk.NewEvent(
				types.EventTypeMigrateStakingRedelegate,
				sdk.NewAttribute(types.AttributeKeyValidatorSrcAddr, valSrcAddr.String()),
				sdk.NewAttribute(types.AttributeKeyValidatorDstAddr, valDstAddr.String()),
			),
		)
	}

	if len(events) > 0 {
		ctx.EventManager().EmitEvents(events)
	}
	return nil
}
