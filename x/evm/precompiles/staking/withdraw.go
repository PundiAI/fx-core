package staking

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var WithdrawMethod = abi.NewMethod(WithdrawMethodName, WithdrawMethodName, abi.Function, "nonpayable", false, false,
	abi.Arguments{
		abi.Argument{
			Name: "validator",
			Type: types.TypeString,
		},
	},
	abi.Arguments{
		abi.Argument{
			Name: "reward",
			Type: types.TypeUint256,
		},
	},
)

func (c *Contract) Withdraw(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("withdraw method not readonly")
	}

	args, err := WithdrawMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr := args[0].(string)
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}

	snapshot := evm.StateDB.Snapshot()
	cacheCtx, commit := c.ctx.CacheContext()
	bondDenom := c.stakingKeeper.BondDenom(cacheCtx)

	rewardAmount, err := c.withdraw(cacheCtx, evm, contract.Caller(), valAddr, bondDenom)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		return nil, err
	}

	commit()
	c.ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

	return WithdrawMethod.Outputs.Pack(rewardAmount)
}

func (c *Contract) withdraw(ctx sdk.Context, evm *vm.EVM, sender common.Address, valAddr sdk.ValAddress, withDrawDenom string) (*big.Int, error) {
	delAddr := sdk.AccAddress(sender.Bytes())
	withdrawAddr := c.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
	if !withdrawAddr.Equals(delAddr) {
		// cache withdraw address state, before withdraw rewards
		evm.StateDB.GetBalance(common.BytesToAddress(withdrawAddr.Bytes()))
	}

	rewards, err := c.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}

	rewardAmount := rewards.AmountOf(withDrawDenom).BigInt()
	if rewardAmount.Cmp(big.NewInt(0)) == 0 {
		return rewardAmount, nil
	}

	if withdrawAddr.Equals(delAddr) {
		evm.StateDB.AddBalance(sender, rewardAmount)
	} else {
		evm.StateDB.AddBalance(common.BytesToAddress(withdrawAddr.Bytes()), rewardAmount)
	}
	return rewardAmount, nil
}
