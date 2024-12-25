package staking

import (
	"errors"
	"fmt"
	"math/big"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/types"
)

type Validator struct {
	ValAddr      string
	MissedBlocks int64
}

type Keeper struct {
	bankKeeper       types.BankKeeper
	distrKeeper      types.DistrKeeper
	distrMsgServer   distrtypes.MsgServer
	stakingKeeper    types.StakingKeeper
	stakingMsgServer stakingtypes.MsgServer
	slashingKeeper   types.SlashingKeeper
	stakingDenom     string
}

//nolint:gocyclo // need to refactor
func (k Keeper) handlerTransferShares(
	ctx sdk.Context,
	evm *vm.EVM,
	valAddr sdk.ValAddress,
	from, to common.Address,
	sharesInt *big.Int,
) (*big.Int, *big.Int, error) {
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, nil, err
	}
	fromDel, err := k.stakingKeeper.GetDelegation(ctx, from.Bytes(), valAddr)
	if err != nil {
		return nil, nil, err
	}
	// if from has receiving redelegation, can't transfer shares
	has, err := k.stakingKeeper.HasReceivingRedelegation(ctx, from.Bytes(), valAddr)
	if err != nil {
		return nil, nil, err
	}
	if has {
		return nil, nil, errors.New("from has receiving redelegation")
	}

	shares := sdkmath.LegacyNewDecFromBigInt(sharesInt)
	if fromDel.GetShares().LT(shares) {
		return nil, nil, fmt.Errorf("insufficient shares(%s < %s)", fromDel.GetShares().TruncateInt().String(), shares.TruncateInt().String())
	}

	// withdraw reward
	withdrawAddr, err := k.distrKeeper.GetDelegatorWithdrawAddr(ctx, to.Bytes())
	if err != nil {
		return nil, nil, err
	}
	beforeDelBalance := k.bankKeeper.GetBalance(ctx, withdrawAddr, k.stakingDenom)

	// withdraw reward
	withdrawRewardRes, err := k.distrMsgServer.WithdrawDelegatorReward(ctx, &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: valAddr.String(),
	})
	if err != nil {
		return nil, nil, err
	}

	withdrawABI := NewWithdrawABI()
	data, topic, err := withdrawABI.NewWithdrawEvent(from, valAddr.String(), withdrawRewardRes.Amount.AmountOf(k.stakingDenom).BigInt())
	if err != nil {
		return nil, nil, err
	}
	fxcontract.EmitEvent(evm, stakingAddress, data, topic)

	// get to delegation
	toDel, err := k.stakingKeeper.GetDelegation(ctx, to.Bytes(), valAddr)
	toDelFound := false
	if err != nil {
		if !errors.Is(err, stakingtypes.ErrNoDelegation) {
			return nil, nil, err
		}
		toDel = stakingtypes.NewDelegation(sdk.AccAddress(to.Bytes()).String(), valAddr.String(), sdkmath.LegacyZeroDec())
		// if address to not delegate before, increase validator period
		if _, err = k.distrKeeper.IncrementValidatorPeriod(ctx, validator); err != nil {
			return nil, nil, err
		}
	} else {
		toDelFound = true
		toWithdrawRewardsRes, err := k.distrMsgServer.WithdrawDelegatorReward(ctx, &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: sdk.AccAddress(to.Bytes()).String(),
			ValidatorAddress: valAddr.String(),
		})
		if err != nil {
			return nil, nil, err
		}
		data, topic, err = withdrawABI.NewWithdrawEvent(to, valAddr.String(), toWithdrawRewardsRes.Amount.AmountOf(k.stakingDenom).BigInt())
		if err != nil {
			return nil, nil, err
		}
		fxcontract.EmitEvent(evm, stakingAddress, data, topic)
	}

	// update from delegate, delete it if shares zero
	fromDelStartingInfo, err := k.distrKeeper.GetDelegatorStartingInfo(ctx, valAddr, from.Bytes())
	if err != nil {
		return nil, nil, err
	}
	fromDel.Shares = fromDel.Shares.Sub(shares)
	if fromDel.GetShares().IsZero() {
		// if shares zero, remove delegation and delete starting info and reference count
		if err = k.stakingKeeper.RemoveDelegation(ctx, fromDel); err != nil {
			return nil, nil, err
		}
		// decrement previous period reference count
		if err = k.decrementReferenceCount(ctx, valAddr, fromDelStartingInfo.PreviousPeriod); err != nil {
			return nil, nil, err
		}
		if err = k.distrKeeper.DeleteDelegatorStartingInfo(ctx, valAddr, from.Bytes()); err != nil {
			return nil, nil, err
		}
	} else {
		if err = k.stakingKeeper.SetDelegation(ctx, fromDel); err != nil {
			return nil, nil, err
		}
		// update from starting info
		fromDelStartingInfo.Stake = validator.TokensFromSharesTruncated(fromDel.GetShares())
		if err = k.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, from.Bytes(), fromDelStartingInfo); err != nil {
			return nil, nil, err
		}
	}

	// update to delegate, set starting info if to not delegate before
	toDel.Shares = toDel.Shares.Add(shares)
	if err = k.stakingKeeper.SetDelegation(ctx, toDel); err != nil {
		return nil, nil, err
	}
	if !toDelFound {
		// if to not delegate before, last period reference count - 1 and set starting info
		validatorCurrentRewards, err := k.distrKeeper.GetValidatorCurrentRewards(ctx, valAddr)
		if err != nil {
			return nil, nil, err
		}
		previousPeriod := validatorCurrentRewards.Period - 1
		if err = k.incrementReferenceCount(ctx, valAddr, previousPeriod); err != nil {
			return nil, nil, err
		}

		stakeToken := validator.TokensFromSharesTruncated(shares)
		toDelStartingInfo := distrtypes.NewDelegatorStartingInfo(previousPeriod, stakeToken, uint64(ctx.BlockHeight()))
		if err = k.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, to.Bytes(), toDelStartingInfo); err != nil {
			return nil, nil, err
		}
	} else {
		// update to starting info
		toDelStartingInfo, err := k.distrKeeper.GetDelegatorStartingInfo(ctx, valAddr, to.Bytes())
		if err != nil {
			return nil, nil, err
		}
		toDelStartingInfo.Stake = validator.TokensFromSharesTruncated(toDel.GetShares())
		if err = k.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, to.Bytes(), toDelStartingInfo); err != nil {
			return nil, nil, err
		}
	}

	// calculate token from shares
	token := validator.TokensFromShares(shares).TruncateInt()

	afterDelBalance := k.bankKeeper.GetBalance(ctx, withdrawAddr, k.stakingDenom)
	toRewardCoin := afterDelBalance.Sub(beforeDelBalance)

	return token.BigInt(), toRewardCoin.Amount.BigInt(), nil
}

func (k Keeper) decrementAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress, decrease *big.Int) error {
	allowance := k.stakingKeeper.GetAllowance(ctx, valAddr, owner, spender)
	if allowance.Cmp(decrease) < 0 {
		return fmt.Errorf("transfer shares exceeds allowance(%s < %s)", allowance.String(), decrease.String())
	}
	newAllowance := big.NewInt(0).Sub(allowance, decrease)
	k.stakingKeeper.SetAllowance(ctx, valAddr, owner, spender, newAllowance)
	return nil
}

// increment the reference count for a historical rewards value
func (k Keeper) incrementReferenceCount(ctx sdk.Context, valAddr sdk.ValAddress, period uint64) error {
	historical, err := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if err != nil {
		return err
	}
	if historical.ReferenceCount > 2 {
		return errors.New("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	return k.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k Keeper) decrementReferenceCount(ctx sdk.Context, valAddr sdk.ValAddress, period uint64) error {
	historical, err := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if err != nil {
		return err
	}
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		return k.distrKeeper.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	} else {
		return k.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
	}
}

func (k Keeper) NewStakingCoin(amount *big.Int) sdk.Coin {
	return sdk.NewCoin(k.stakingDenom, sdkmath.NewIntFromBigInt(amount))
}

func (k Keeper) ValidatorListMissedBlock(ctx sdk.Context, bondedVals []stakingtypes.Validator) ([]string, error) {
	valList := make([]Validator, 0, len(bondedVals))
	for _, val := range bondedVals {
		consAddr, err := val.GetConsAddr()
		if err != nil {
			return nil, err
		}
		info, err := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
		if err != nil {
			return nil, err
		}
		valList = append(valList, Validator{
			ValAddr:      val.OperatorAddress,
			MissedBlocks: info.MissedBlocksCounter,
		})
	}
	sort.Slice(valList, func(i, j int) bool {
		return valList[i].MissedBlocks > valList[j].MissedBlocks
	})
	valAddrs := make([]string, 0, len(valList))
	for _, l := range valList {
		valAddrs = append(valAddrs, l.ValAddr)
	}
	return valAddrs, nil
}

func validatorListPower(bondedVals []stakingtypes.Validator) []string {
	valAddrs := make([]string, 0, len(bondedVals))
	for _, val := range bondedVals {
		valAddrs = append(valAddrs, val.OperatorAddress)
	}
	return valAddrs
}
