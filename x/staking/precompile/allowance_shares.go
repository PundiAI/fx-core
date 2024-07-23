package precompile

import (
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) AllowanceShares(evm *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	var args AllowanceSharesArgs
	if err := types.ParseMethodArgs(AllowanceSharesMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	allowance := c.stakingKeeper.GetAllowance(stateDB.CacheContext(), args.GetValidator(), args.Owner.Bytes(), args.Spender.Bytes())
	return AllowanceSharesMethod.Outputs.Pack(allowance)
}
