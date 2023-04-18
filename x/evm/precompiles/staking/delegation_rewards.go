package staking

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v4/x/evm/types"
)

var DelegationRewardsMethod = abi.NewMethod(
	DelegationRewardsMethodName,
	DelegationRewardsMethodName,
	abi.Function, "view", false, false,
	abi.Arguments{
		abi.Argument{Name: "_val", Type: types.TypeString},
		abi.Argument{Name: "_del", Type: types.TypeAddress},
	},
	abi.Arguments{
		abi.Argument{Name: "_reward", Type: types.TypeUint256},
	},
)

type DelegationRewardsArgs struct {
	Validator string         `abi:"_val"`
	Delegator common.Address `abi:"_del"`
}

func (c *Contract) DelegationRewards(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	// NOTE: function modify state, so cache context and not commit
	cacheCtx, _ := ctx.CacheContext()
	// parse args
	var args DelegationRewardsArgs
	if err := ParseMethodParams(DelegationRewardsMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(args.Validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	validator, found := c.stakingKeeper.GetValidator(cacheCtx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}

	delegation, found := c.stakingKeeper.GetDelegation(cacheCtx, sdk.AccAddress(args.Delegator.Bytes()), valAddr)
	if !found {
		return DelegationRewardsMethod.Outputs.Pack(big.NewInt(0))
	}

	evmDenom := c.evmKeeper.GetParams(cacheCtx).EvmDenom
	endingPeriod := c.distrKeeper.IncrementValidatorPeriod(cacheCtx, validator)
	rewards := c.distrKeeper.CalculateDelegationRewards(cacheCtx, validator, delegation, endingPeriod)

	return DelegationRewardsMethod.Outputs.Pack(rewards.AmountOf(evmDenom).TruncateInt().BigInt())
}
