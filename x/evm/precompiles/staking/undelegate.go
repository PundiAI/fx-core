package staking

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var (
	UndelegateMethod = abi.NewMethod(UndelegateMethodName, UndelegateMethodName, abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{
				Name: "validator",
				Type: types.TypeString,
			},
			abi.Argument{
				Name: "shares",
				Type: types.TypeUint256,
			},
		},
		abi.Arguments{
			abi.Argument{
				Name: "amount",
				Type: types.TypeUint256,
			},
			abi.Argument{
				Name: "reward",
				Type: types.TypeUint256,
			},
			abi.Argument{
				Name: "completionTime",
				Type: types.TypeUint256,
			},
		},
	)

	UndelegateEvent = abi.NewEvent(
		UndelegateEventName,
		UndelegateEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "amount", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "completionTime", Type: types.TypeUint256, Indexed: false},
		},
	)
)

func (c *Contract) Undelegate(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("undelegate method not readonly")
	}

	args, err := UndelegateMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr, ok1 := args[0].(string)
	shareAmount, ok2 := args[1].(*big.Int)
	if !ok1 || !ok2 {
		return nil, errors.New("unexpected arg type")
	}

	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	if shareAmount.Sign() <= 0 {
		return nil, fmt.Errorf("invalid shares: %s", shareAmount.String())
	}

	_, found := c.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}

	sender := sdk.AccAddress(contract.Caller().Bytes())
	evmDenom := c.evmKeeper.GetEVMDenom(ctx)

	// withdraw rewards if delegation exist, add reward to evm state balance
	reward := big.NewInt(0)
	if _, found = c.stakingKeeper.GetDelegation(ctx, sender, valAddr); found {
		if reward, err = c.withdraw(ctx, evm, contract.Caller(), valAddr, evmDenom); err != nil {
			return nil, err
		}
	}

	unDelAmount, completionTime, err := Undelegate(ctx, c.stakingKeeper, c.bankKeeper, sender, valAddr, sdk.NewDecFromBigInt(shareAmount), evmDenom)
	if err != nil {
		return nil, fmt.Errorf("undelegate failed: %s", err.Error())
	}

	// add undelegate log
	if err := c.AddLog(UndelegateEvent, []common.Hash{contract.Caller().Hash()},
		valAddrStr, shareAmount, unDelAmount.BigInt(), big.NewInt(completionTime.Unix())); err != nil {
		return nil, err
	}

	return UndelegateMethod.Outputs.Pack(unDelAmount.BigInt(), reward, big.NewInt(completionTime.Unix()))
}

func Undelegate(ctx sdk.Context, sk StakingKeeper, bk BankKeeper, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress, shares sdk.Dec, bondDenom string,
) (sdkmath.Int, time.Time, error) {
	validator, found := sk.GetValidator(ctx, valAddr)
	if !found {
		return sdkmath.Int{}, time.Time{}, stakingtypes.ErrNoDelegatorForAddress
	}

	if sk.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return sdkmath.Int{}, time.Time{}, stakingtypes.ErrMaxUnbondingDelegationEntries
	}

	returnAmount, err := sk.Unbond(ctx, delAddr, valAddr, shares)
	if err != nil {
		return sdkmath.Int{}, time.Time{}, err
	}

	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() {
		coins := sdk.NewCoins(sdk.NewCoin(bondDenom, returnAmount))
		if err := bk.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, coins); err != nil {
			return sdkmath.Int{}, time.Time{}, err
		}
	}

	completionTime := ctx.BlockHeader().Time.Add(sk.UnbondingTime(ctx))
	ubd := sk.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	sk.InsertUBDQueue(ctx, ubd, completionTime)

	return returnAmount, completionTime, nil
}
