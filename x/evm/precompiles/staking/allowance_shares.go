package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var AllowanceSharesMethod = abi.NewMethod(
	AllowanceSharesMethodName,
	AllowanceSharesMethodName,
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

type AllowanceSharesArgs struct {
	Validator string         `abi:"_val"`
	Owner     common.Address `abi:"_owner"`
	Spender   common.Address `abi:"_spender"`
}

func (c *Contract) AllowanceShares(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	cacheCtx, _ := ctx.CacheContext()
	// parse args
	var args AllowanceSharesArgs
	if err := ParseMethodParams(AllowanceSharesMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(args.Validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	allowance := c.stakingKeeper.GetAllowance(cacheCtx, valAddr, args.Owner.Bytes(), args.Spender.Bytes())
	return AllowanceSharesMethod.Outputs.Pack(allowance)
}
