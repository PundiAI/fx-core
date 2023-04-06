package staking

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var (
	TransferMethod = abi.NewMethod(
		TransferMethodName,
		TransferMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "validator", Type: types.TypeString},
			abi.Argument{Name: "to", Type: types.TypeAddress},
			abi.Argument{Name: "shares", Type: types.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "token", Type: types.TypeUint256},
			abi.Argument{Name: "reward", Type: types.TypeUint256},
		},
	)
	TransferFromMethod = abi.NewMethod(
		TransferFromMethodName,
		TransferFromMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "validator", Type: types.TypeString},
			abi.Argument{Name: "from", Type: types.TypeAddress},
			abi.Argument{Name: "to", Type: types.TypeAddress},
			abi.Argument{Name: "shares", Type: types.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "token", Type: types.TypeUint256},
			abi.Argument{Name: "reward", Type: types.TypeUint256},
		},
	)

	TransferEvent = abi.NewEvent(
		TransferEventName,
		TransferEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "from", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "to", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "token", Type: types.TypeUint256, Indexed: false},
		},
	)
)

func (c *Contract) Transfer(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("transfer method not readonly")
	}
	args, err := TransferMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}

	valAddrStr, ok0 := args[0].(string)
	toAddr, ok1 := args[1].(common.Address)
	shares, ok2 := args[2].(*big.Int)
	if !ok0 || !ok1 || !ok2 {
		return nil, errors.New("unexpected arg type")
	}
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	if shares.Cmp(big.NewInt(0)) < 0 {
		return nil, errors.New("shares cannot be negative")
	}

	token, reward, err := c.handlerTransfer(ctx, evm, valAddr, contract.Caller(), toAddr, shares)
	if err != nil {
		return nil, err
	}
	return TransferMethod.Outputs.Pack(token, reward)
}

func (c *Contract) TransferFrom(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("transferFrom method not readonly")
	}
	args, err := TransferFromMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}

	valAddrStr, ok0 := args[0].(string)
	fromAddr, ok1 := args[1].(common.Address)
	toAddr, ok2 := args[2].(common.Address)
	shares, ok3 := args[3].(*big.Int)
	if !ok0 || !ok1 || !ok2 || !ok3 {
		return nil, errors.New("unexpected arg type")
	}
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	if shares.Cmp(big.NewInt(0)) < 0 {
		return nil, errors.New("shares cannot be negative")
	}

	spender := contract.Caller()
	if err = c.decrementAllowance(ctx, valAddr, fromAddr.Bytes(), spender.Bytes(), shares); err != nil {
		return nil, err
	}

	token, reward, err := c.handlerTransfer(ctx, evm, valAddr, fromAddr, toAddr, shares)
	if err != nil {
		return nil, err
	}
	return TransferFromMethod.Outputs.Pack(token, reward)
}

func (c *Contract) decrementAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress, decrease *big.Int) error {
	allowance := c.stakingKeeper.GetAllowance(ctx, valAddr, owner, spender)
	if allowance.Cmp(decrease) < 0 {
		return fmt.Errorf("transfer shares exceeds allowance(%s < %s)", allowance.String(), decrease.String())
	}
	newAllowance := big.NewInt(0).Sub(allowance, decrease)
	c.stakingKeeper.SetAllowance(ctx, valAddr, owner, spender, newAllowance)
	return nil
}

func (c *Contract) handlerTransfer(
	ctx sdk.Context,
	evm *vm.EVM,
	valAddr sdk.ValAddress,
	from, to common.Address,
	sharesInt *big.Int,
) (*big.Int, *big.Int, error) {
	validator, found := c.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}
	fromDel, found := c.stakingKeeper.GetDelegation(ctx, from.Bytes(), valAddr)
	if !found {
		return nil, nil, errors.New("from delegation not found")
	}

	shares := sdk.NewDecFromBigInt(sharesInt) // TODO share with sdk.Precision
	if fromDel.GetShares().LT(shares) {
		return nil, nil, fmt.Errorf("insufficient shares(%s < %s)", shares.TruncateInt().String(), fromDel.GetShares().TruncateInt().String())
	}

	// withdraw reward
	evmDenom := c.evmKeeper.GetEVMDenom(ctx)
	_, err := c.withdraw(ctx, evm, from, valAddr, evmDenom)
	if err != nil {
		return nil, nil, fmt.Errorf("withdraw failed: %s", err.Error())
	}

	// get to delegation
	toReward := big.NewInt(0)
	toDel, toDelFound := c.stakingKeeper.GetDelegation(ctx, to.Bytes(), valAddr)
	if !toDelFound {
		toDel = stakingtypes.NewDelegation(to.Bytes(), valAddr, sdk.ZeroDec())
		// if address to not delegate before, increase validator period
		_ = c.distrKeeper.IncrementValidatorPeriod(ctx, validator)
	} else {
		// withdraw to address reward
		toReward, err = c.withdraw(ctx, evm, to, valAddr, evmDenom)
		if err != nil {
			return nil, nil, fmt.Errorf("withdraw failed: %s", err.Error())
		}
	}

	// update from delegate, delete it if shares zero
	fromDel.Shares = fromDel.Shares.Sub(shares)
	if fromDel.GetShares().IsZero() {
		// if shares zero, remove delegation and delete starting info and reference count
		if err := c.stakingKeeper.RemoveDelegation(ctx, fromDel); err != nil {
			return nil, nil, err
		}
		// decrement previous period reference count
		startingPeriod := c.distrKeeper.GetDelegatorStartingInfo(ctx, valAddr, from.Bytes()).PreviousPeriod
		decrementReferenceCount(c.distrKeeper, ctx, valAddr, startingPeriod)
		c.distrKeeper.DeleteDelegatorStartingInfo(ctx, valAddr, from.Bytes())
	} else {
		c.stakingKeeper.SetDelegation(ctx, fromDel)
	}

	// update to delegate, set starting info if to not delegate before
	toDel.Shares = toDel.Shares.Add(shares)
	c.stakingKeeper.SetDelegation(ctx, toDel)
	if !toDelFound {
		// if to not delegate before, last period reference count - 1 and set starting info
		previousPeriod := c.distrKeeper.GetValidatorCurrentRewards(ctx, valAddr).Period - 1
		incrementReferenceCount(c.distrKeeper, ctx, valAddr, previousPeriod)

		stakeToken := validator.TokensFromSharesTruncated(shares)
		toDelStartingInfo := distrtypes.NewDelegatorStartingInfo(previousPeriod, stakeToken, uint64(ctx.BlockHeight()))
		c.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, to.Bytes(), toDelStartingInfo)
	}

	// calculate token from shares
	token := validator.TokensFromShares(shares).TruncateInt()

	// add Transfer event
	if err := c.AddLog(TransferEvent, []common.Hash{from.Hash(), to.Hash()},
		valAddr.String(), shares.TruncateInt().BigInt(), token.BigInt()); err != nil {
		return nil, nil, err
	}

	return token.BigInt(), toReward, nil
}

// increment the reference count for a historical rewards value
func incrementReferenceCount(k DistrKeeper, ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func decrementReferenceCount(k DistrKeeper, ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	} else {
		k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
	}
}
