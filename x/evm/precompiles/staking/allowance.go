package staking

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var AllowanceMethod = abi.NewMethod(
	AllowanceMethodName,
	AllowanceMethodName,
	abi.Function, "view", false, false,
	abi.Arguments{
		abi.Argument{Name: "_val", Type: types.TypeString},
		abi.Argument{Name: "_owner", Type: types.TypeAddress},
		abi.Argument{Name: "_spender", Type: types.TypeAddress},
	},
	abi.Arguments{
		abi.Argument{Name: "_shares", Type: types.TypeUint256},
	},
)

func (c *Contract) Allowance(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	cacheCtx, _ := ctx.CacheContext()
	args, err := AllowanceMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr, ok0 := args[0].(string)
	owner, ok1 := args[1].(common.Address)
	spender, ok2 := args[2].(common.Address)
	if !ok0 || !ok1 || !ok2 {
		return nil, errors.New("unexpected arg type")
	}
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	allowance := c.stakingKeeper.GetAllowance(cacheCtx, valAddr, owner.Bytes(), spender.Bytes())
	return AllowanceMethod.Outputs.Pack(allowance)
}
