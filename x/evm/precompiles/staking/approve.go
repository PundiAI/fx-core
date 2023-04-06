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
	ApproveMethod = abi.NewMethod(
		ApproveMethodName,
		ApproveMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "validator", Type: types.TypeString},
			abi.Argument{Name: "spender", Type: types.TypeAddress},
			abi.Argument{Name: "shares", Type: types.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "result", Type: types.TypeBool},
		},
	)

	ApproveEvent = abi.NewEvent(
		ApproveEventName,
		ApproveEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "owner", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "spender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
		},
	)
)

func (c *Contract) Approve(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("approve method not readonly")
	}
	// parse args
	args, err := ApproveMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr, ok0 := args[0].(string)
	spender, ok1 := args[1].(common.Address)
	shares, ok2 := args[2].(*big.Int)
	if !ok0 || !ok1 || !ok2 {
		return nil, errors.New("unexpected arg type")
	}
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	if shares.Cmp(big.NewInt(0)) < 0 {
		return nil, errors.New("allowance cannot be negative")
	}
	// owner
	owner := contract.Caller()

	// set allowance
	c.stakingKeeper.SetAllowance(ctx, valAddr, owner.Bytes(), spender.Bytes(), shares)

	// emit event
	if err := c.AddLog(ApproveEvent, []common.Hash{owner.Hash(), spender.Hash()}, valAddrStr, shares); err != nil {
		return nil, err
	}
	return ApproveMethod.Outputs.Pack(true)
}
