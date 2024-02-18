package staking

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) DelegationRewards(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	// NOTE: function modify state, so cache context and not commit
	cacheCtx, _ := ctx.CacheContext()
	// parse args
	var args DelegationRewardsArgs
	if err := types.ParseMethodArgs(DelegationRewardsMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

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
