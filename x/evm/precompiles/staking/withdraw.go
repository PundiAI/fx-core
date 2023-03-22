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

func (c *Contract) Withdraw(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("withdraw method not readonly")
	}

	args, err := WithdrawMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr, ok := args[0].(string)
	if !ok {
		return nil, errors.New("unexpected arg type")
	}

	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	evmDenom := c.evmKeeper.GetEVMDenom(ctx)

	rewardAmount, err := c.withdraw(ctx, evm, contract.Caller(), valAddr, evmDenom)
	if err != nil {
		return nil, err
	}

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
