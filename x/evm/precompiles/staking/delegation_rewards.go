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

var DelegationRewardsMethod = abi.NewMethod(DelegationRewardsMethodName, DelegationRewardsMethodName, abi.Function, "nonpayable", false, false,
	abi.Arguments{
		abi.Argument{
			Name: "validator",
			Type: types.TypeString,
		},
		abi.Argument{
			Name: "delegator",
			Type: types.TypeAddress,
		},
	},
	abi.Arguments{
		abi.Argument{
			Name: "rewards",
			Type: types.TypeUint256,
		},
	},
)

func (c *Contract) DelegationRewards(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	args, err := DelegationRewardsMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr, ok := args[0].(string)
	if !ok {
		return nil, errors.New("unexpected arg type")
	}
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	validator, found := c.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}

	delAddr, ok := args[1].(common.Address)
	if !ok {
		return nil, errors.New("unexpected arg type")
	}
	delegation, found := c.stakingKeeper.GetDelegation(ctx, sdk.AccAddress(delAddr.Bytes()), valAddr)
	if !found {
		return DelegationRewardsMethod.Outputs.Pack(big.NewInt(0))
	}

	evmDenom := c.evmKeeper.GetEVMDenom(ctx)
	endingPeriod := c.distrKeeper.IncrementValidatorPeriod(ctx, validator)
	rewards := c.distrKeeper.CalculateDelegationRewards(ctx, validator, delegation, endingPeriod)

	return DelegationRewardsMethod.Outputs.Pack(rewards.AmountOf(evmDenom).TruncateInt().BigInt())
}
