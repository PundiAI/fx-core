package precompile

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) DelegationRewards(evm *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	var args DelegationRewardsArgs
	if err := types.ParseMethodArgs(DelegationRewardsMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	cacheCtx := stateDB.CacheContext()

	valAddr := args.GetValidator()
	validator, found := c.stakingKeeper.GetValidator(cacheCtx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}
	delegation, found := c.stakingKeeper.GetDelegation(cacheCtx, args.Delegator.Bytes(), valAddr)
	if !found {
		return DelegationRewardsMethod.Outputs.Pack(big.NewInt(0))
	}

	evmDenom := c.evmKeeper.GetParams(cacheCtx).EvmDenom
	endingPeriod := c.distrKeeper.IncrementValidatorPeriod(cacheCtx, validator)
	rewards := c.distrKeeper.CalculateDelegationRewards(cacheCtx, validator, delegation, endingPeriod)

	return DelegationRewardsMethod.Outputs.Pack(rewards.AmountOf(evmDenom).TruncateInt().BigInt())
}
