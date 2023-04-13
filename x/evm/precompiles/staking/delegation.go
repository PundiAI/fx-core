package staking

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var DelegationMethod = abi.NewMethod(
	DelegationMethodName,
	DelegationMethodName,
	abi.Function, "view", false, false,
	abi.Arguments{
		abi.Argument{Name: "_val", Type: types.TypeString},
		abi.Argument{Name: "_del", Type: types.TypeAddress},
	},
	abi.Arguments{
		abi.Argument{Name: "_shares", Type: types.TypeUint256},
		abi.Argument{Name: "_delegateAmount", Type: types.TypeUint256},
	},
)

type DelegationArgs struct {
	Validator string         `abi:"_val"`
	Delegator common.Address `abi:"_del"`
}

func (c *Contract) Delegation(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	cacheCtx, _ := ctx.CacheContext()
	// parse args
	var args DelegationArgs
	if err := ParseMethodParams(DelegationMethod, &args, contract.Input[4:]); err != nil {
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
		return DelegationMethod.Outputs.Pack(big.NewInt(0), big.NewInt(0))
	}

	delegationAmt := delegation.GetShares().MulInt(validator.GetTokens()).Quo(validator.GetDelegatorShares())
	// TODO truncate shares, decimal 18
	return DelegationMethod.Outputs.Pack(delegation.GetShares().TruncateInt().BigInt(), delegationAmt.TruncateInt().BigInt())
}
