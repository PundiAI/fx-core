package staking

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var WithdrawMethod = abi.NewMethod(WithdrawMethodName, WithdrawMethodName, abi.Function, "", false, false,
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

	sender := sdk.AccAddress(contract.CallerAddress.Bytes())
	snapshot := evm.StateDB.Snapshot()
	cacheCtx, commit := c.ctx.CacheContext()
	bondDenom := c.stakingKeeper.BondDenom(cacheCtx)

	rewards, err := c.distrKeeper.WithdrawDelegationRewards(cacheCtx, sender, valAddr)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		return nil, err
	}
	evm.StateDB.AddBalance(contract.CallerAddress, rewards.AmountOf(bondDenom).BigInt())

	commit()
	return WithdrawMethod.Outputs.Pack(rewards.AmountOf(bondDenom).BigInt())
}
