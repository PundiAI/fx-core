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
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var UndelegateMethod = abi.NewMethod(UndelegateMethodName, UndelegateMethodName, abi.Function, "nonpayable", false, false,
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
			Name: "endTime",
			Type: types.TypeUint256,
		},
	},
)

func (c *Contract) Undelegate(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("undelegate method not readonly")
	}

	args, err := UndelegateMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr := args[0].(string)
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	_, found := c.stakingKeeper.GetValidator(c.ctx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}

	shareAmount := args[1].(*big.Int)
	if shareAmount.Sign() <= 0 {
		return nil, fmt.Errorf("invalid shares: %s", shareAmount.String())
	}
	sender := sdk.AccAddress(contract.CallerAddress.Bytes())

	snapshot := evm.StateDB.Snapshot()
	cacheCtx, commit := c.ctx.CacheContext()
	evmDenom := c.evmKeeper.GetEVMDenom(cacheCtx)

	// withdraw rewards if delegation exist, add reward to evm state balance
	reward := big.NewInt(0)
	if _, found = c.stakingKeeper.GetDelegation(cacheCtx, sender, valAddr); found {
		if reward, err = c.withdraw(cacheCtx, evm, contract.Caller(), valAddr, evmDenom); err != nil {
			evm.StateDB.RevertToSnapshot(snapshot)
			return nil, err
		}
	}

	unDelAmount, endTime, err := Undelegate(cacheCtx, c.stakingKeeper,
		c.bankKeeper, sender, valAddr, sdk.NewDecFromBigInt(shareAmount), evmDenom)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		return nil, fmt.Errorf("undelegate failed: %s", err.Error())
	}
	commit()
	c.ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

	return UndelegateMethod.Outputs.Pack(unDelAmount.BigInt(), reward, big.NewInt(endTime.Unix()))
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
		if err := bk.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName,
			stakingtypes.NotBondedPoolName, coins); err != nil {
			return sdkmath.Int{}, time.Time{}, err
		}
	}

	completionTime := ctx.BlockHeader().Time.Add(sk.UnbondingTime(ctx))
	ubd := sk.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	sk.InsertUBDQueue(ctx, ubd, completionTime)

	return returnAmount, completionTime, nil
}
