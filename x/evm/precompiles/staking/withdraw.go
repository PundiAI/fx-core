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

var (
	WithdrawMethod = abi.NewMethod(
		WithdrawMethodName,
		WithdrawMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: types.TypeString},
		},
		abi.Arguments{
			abi.Argument{Name: "_reward", Type: types.TypeUint256},
		},
	)

	WithdrawEvent = abi.NewEvent(
		WithdrawEventName,
		WithdrawEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "reward", Type: types.TypeUint256, Indexed: false},
		},
	)
)

type WithdrawArgs struct {
	Validator string `abi:"_val"`
}

func (c *Contract) Withdraw(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("withdraw method not readonly")
	}
	// parse args
	var args WithdrawArgs
	if err := ParseMethodParams(WithdrawMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(args.Validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", args.Validator)
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

	// add withdraw log
	if err := c.AddLog(WithdrawEvent, []common.Hash{sender.Hash()}, valAddr.String(), rewardAmount); err != nil {
		return nil, err
	}

	return rewardAmount, nil
}
