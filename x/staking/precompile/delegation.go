package precompile

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) Delegation(evm *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	var args DelegationArgs
	if err := types.ParseMethodArgs(DelegationMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	ctx := stateDB.CacheContext()

	valAddr := args.GetValidator()
	validator, found := c.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}

	delegation, found := c.stakingKeeper.GetDelegation(ctx, args.Delegator.Bytes(), valAddr)
	if !found {
		return DelegationMethod.Outputs.Pack(big.NewInt(0), big.NewInt(0))
	}

	delegationAmt := delegation.GetShares().MulInt(validator.GetTokens()).Quo(validator.GetDelegatorShares())
	return DelegationMethod.Outputs.Pack(delegation.GetShares().TruncateInt().BigInt(), delegationAmt.TruncateInt().BigInt())
}
