package precompile

import (
	"errors"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type TransferShare struct {
	*Keeper
	abi.Method
	abi.Event
}

type TransferShares struct {
	*TransferShare
}

func NewTransferSharesMethod(keeper *Keeper) *TransferShares {
	return &TransferShares{
		TransferShare: &TransferShare{
			Keeper: keeper,
			Method: fxstakingtypes.GetABI().Methods["transferShares"],
			Event:  fxstakingtypes.GetABI().Events["TransferShares"],
		},
	}
}

func (m *TransferShares) IsReadonly() bool {
	return false
}

func (m *TransferShares) GetMethodId() []byte {
	return m.Method.ID
}

func (m *TransferShares) RequiredGas() uint64 {
	return 50_000
}

func (m *TransferShares) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	valAddr := args.GetValidator()
	stateDB := evm.StateDB.(types.ExtStateDB)
	var result []byte
	if err := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		token, reward, err := m.handlerTransferShares(ctx, evm, valAddr, contract.Caller(), args.To, args.Shares)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
				sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(contract.Caller().Bytes()).String()),
			),
		)
		result, err = m.PackOutput(token, reward)
		return err
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *TransferShares) PackInput(args fxstakingtypes.TransferSharesArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.To, args.Shares)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *TransferShares) UnpackInput(data []byte) (*fxstakingtypes.TransferSharesArgs, error) {
	args := new(fxstakingtypes.TransferSharesArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

type TransferFromShares struct {
	*TransferShare
}

func NewTransferFromSharesMethod(keeper *Keeper) *TransferFromShares {
	return &TransferFromShares{
		TransferShare: &TransferShare{
			Keeper: keeper,
			Method: fxstakingtypes.GetABI().Methods["transferFromShares"],
			Event:  fxstakingtypes.GetABI().Events["TransferShares"],
		},
	}
}

func (m *TransferFromShares) IsReadonly() bool {
	return false
}

func (m *TransferFromShares) GetMethodId() []byte {
	return m.Method.ID
}

func (m *TransferFromShares) RequiredGas() uint64 {
	return 60_000
}

func (m *TransferFromShares) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		valAddr := args.GetValidator()
		spender := contract.Caller()
		if err := m.decrementAllowance(ctx, valAddr, args.From.Bytes(), spender.Bytes(), args.Shares); err != nil {
			return err
		}
		token, reward, err := m.handlerTransferShares(ctx, evm, valAddr, args.From, args.To, args.Shares)
		if err != nil {
			return err
		}

		result, err = m.PackOutput(token, reward)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
				sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(contract.Caller().Bytes()).String()),
			),
		)
		return nil
	})
	return result, err
}

func (m *TransferFromShares) PackInput(args fxstakingtypes.TransferFromSharesArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.From, args.To, args.Shares)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *TransferFromShares) UnpackInput(data []byte) (*fxstakingtypes.TransferFromSharesArgs, error) {
	args := new(fxstakingtypes.TransferFromSharesArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *TransferShare) decrementAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress, decrease *big.Int) error {
	allowance := m.stakingKeeper.GetAllowance(ctx, valAddr, owner, spender)
	if allowance.Cmp(decrease) < 0 {
		return fmt.Errorf("transfer shares exceeds allowance(%s < %s)", allowance.String(), decrease.String())
	}
	newAllowance := big.NewInt(0).Sub(allowance, decrease)
	m.stakingKeeper.SetAllowance(ctx, valAddr, owner, spender, newAllowance)
	return nil
}

//nolint:gocyclo // need to refactor
func (m *TransferShare) handlerTransferShares(
	ctx sdk.Context,
	evm *vm.EVM,
	valAddr sdk.ValAddress,
	from, to common.Address,
	sharesInt *big.Int,
) (*big.Int, *big.Int, error) {
	validator, err := m.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, nil, err
	}
	fromDel, err := m.stakingKeeper.GetDelegation(ctx, from.Bytes(), valAddr)
	if err != nil {
		return nil, nil, err
	}
	// if from has receiving redelegation, can't transfer shares
	has, err := m.stakingKeeper.HasReceivingRedelegation(ctx, from.Bytes(), valAddr)
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
	withdrawAddr, err := m.distrKeeper.GetDelegatorWithdrawAddr(ctx, to.Bytes())
	if err != nil {
		return nil, nil, err
	}
	beforeDelBalance := m.bankKeeper.GetBalance(ctx, withdrawAddr, m.stakingDenom)

	// withdraw reward
	withdrawRewardRes, err := m.distrMsgServer.WithdrawDelegatorReward(ctx, &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: valAddr.String(),
	})
	if err != nil {
		return nil, nil, err
	}

	withdrawMethod := NewWithdrawMethod(nil)
	data, topic, err := withdrawMethod.NewWithdrawEvent(from, valAddr.String(), withdrawRewardRes.Amount.AmountOf(m.stakingDenom).BigInt())
	if err != nil {
		return nil, nil, err
	}
	EmitEvent(evm, data, topic)

	// get to delegation
	toDel, err := m.stakingKeeper.GetDelegation(ctx, to.Bytes(), valAddr)
	toDelFound := false
	if err != nil {
		if !errors.Is(err, stakingtypes.ErrNoDelegation) {
			return nil, nil, err
		}
		toDel = stakingtypes.NewDelegation(sdk.AccAddress(to.Bytes()).String(), valAddr.String(), sdkmath.LegacyZeroDec())
		// if address to not delegate before, increase validator period
		if _, err = m.distrKeeper.IncrementValidatorPeriod(ctx, validator); err != nil {
			return nil, nil, err
		}
	} else {
		toDelFound = true
		toWithdrawRewardsRes, err := m.distrMsgServer.WithdrawDelegatorReward(ctx, &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: sdk.AccAddress(to.Bytes()).String(),
			ValidatorAddress: valAddr.String(),
		})
		if err != nil {
			return nil, nil, err
		}
		data, topic, err = withdrawMethod.NewWithdrawEvent(to, valAddr.String(), toWithdrawRewardsRes.Amount.AmountOf(m.stakingDenom).BigInt())
		if err != nil {
			return nil, nil, err
		}
		EmitEvent(evm, data, topic)
	}

	// update from delegate, delete it if shares zero
	fromDelStartingInfo, err := m.distrKeeper.GetDelegatorStartingInfo(ctx, valAddr, from.Bytes())
	if err != nil {
		return nil, nil, err
	}
	fromDel.Shares = fromDel.Shares.Sub(shares)
	if fromDel.GetShares().IsZero() {
		// if shares zero, remove delegation and delete starting info and reference count
		if err = m.stakingKeeper.RemoveDelegation(ctx, fromDel); err != nil {
			return nil, nil, err
		}
		// decrement previous period reference count
		if err = decrementReferenceCount(m.distrKeeper, ctx, valAddr, fromDelStartingInfo.PreviousPeriod); err != nil {
			return nil, nil, err
		}
		if err = m.distrKeeper.DeleteDelegatorStartingInfo(ctx, valAddr, from.Bytes()); err != nil {
			return nil, nil, err
		}
	} else {
		if err = m.stakingKeeper.SetDelegation(ctx, fromDel); err != nil {
			return nil, nil, err
		}
		// update from starting info
		fromDelStartingInfo.Stake = validator.TokensFromSharesTruncated(fromDel.GetShares())
		if err = m.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, from.Bytes(), fromDelStartingInfo); err != nil {
			return nil, nil, err
		}
	}

	// update to delegate, set starting info if to not delegate before
	toDel.Shares = toDel.Shares.Add(shares)
	if err = m.stakingKeeper.SetDelegation(ctx, toDel); err != nil {
		return nil, nil, err
	}
	if !toDelFound {
		// if to not delegate before, last period reference count - 1 and set starting info
		validatorCurrentRewards, err := m.distrKeeper.GetValidatorCurrentRewards(ctx, valAddr)
		if err != nil {
			return nil, nil, err
		}
		previousPeriod := validatorCurrentRewards.Period - 1
		if err = incrementReferenceCount(m.distrKeeper, ctx, valAddr, previousPeriod); err != nil {
			return nil, nil, err
		}

		stakeToken := validator.TokensFromSharesTruncated(shares)
		toDelStartingInfo := distrtypes.NewDelegatorStartingInfo(previousPeriod, stakeToken, uint64(ctx.BlockHeight()))
		if err = m.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, to.Bytes(), toDelStartingInfo); err != nil {
			return nil, nil, err
		}
	} else {
		// update to starting info
		toDelStartingInfo, err := m.distrKeeper.GetDelegatorStartingInfo(ctx, valAddr, to.Bytes())
		if err != nil {
			return nil, nil, err
		}
		toDelStartingInfo.Stake = validator.TokensFromSharesTruncated(toDel.GetShares())
		if err = m.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, to.Bytes(), toDelStartingInfo); err != nil {
			return nil, nil, err
		}
	}

	// calculate token from shares
	token := validator.TokensFromShares(shares).TruncateInt()

	afterDelBalance := m.bankKeeper.GetBalance(ctx, withdrawAddr, m.stakingDenom)
	toRewardCoin := afterDelBalance.Sub(beforeDelBalance)

	// add log
	data, topic, err = m.NewTransferShareEvent(from, to, valAddr.String(), shares.TruncateInt().BigInt(), token.BigInt())
	if err != nil {
		return nil, nil, err
	}
	EmitEvent(evm, data, topic)

	return token.BigInt(), toRewardCoin.Amount.BigInt(), nil
}

func (m *TransferShare) PackOutput(amount, reward *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(amount, reward)
}

func (m *TransferShare) UnpackOutput(data []byte) (*big.Int, *big.Int, error) {
	unpacks, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, nil, err
	}
	return unpacks[0].(*big.Int), unpacks[1].(*big.Int), nil
}

func (m *TransferShare) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingTransferShares, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseTransferShares(*log)
}

func (m *TransferShare) NewTransferShareEvent(sender, to common.Address, validator string, shares, amount *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash(), to.Hash()}, validator, shares, amount)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

// increment the reference count for a historical rewards value
func incrementReferenceCount(k DistrKeeper, ctx sdk.Context, valAddr sdk.ValAddress, period uint64) error {
	historical, err := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if err != nil {
		return err
	}
	if historical.ReferenceCount > 2 {
		return errors.New("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	return k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func decrementReferenceCount(k DistrKeeper, ctx sdk.Context, valAddr sdk.ValAddress, period uint64) error {
	historical, err := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if err != nil {
		return err
	}
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		return k.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	} else {
		return k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
	}
}
