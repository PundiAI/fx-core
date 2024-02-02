package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) AllowanceShares(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	cacheCtx, _ := ctx.CacheContext()
	// parse args
	var args AllowanceSharesArgs
	if err := types.ParseMethodArgs(AllowanceSharesMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	allowance := c.stakingKeeper.GetAllowance(cacheCtx, args.GetValidator(), args.Owner.Bytes(), args.Spender.Bytes())
	return AllowanceSharesMethod.Outputs.Pack(allowance)
}
