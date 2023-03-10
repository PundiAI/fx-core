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

var DelegationMethod = abi.NewMethod(DelegationMethodName, DelegationMethodName, abi.Function, "", false, false,
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
			Name: "delegate",
			Type: types.TypeUint256,
		},
	},
)

func (c *Contract) Delegation(_ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	args, err := DelegationMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr := args[0].(string)
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	validator, found := c.stakingKeeper.GetValidator(c.ctx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}

	delAddr := args[1].(common.Address)
	delegation, found := c.stakingKeeper.GetDelegation(c.ctx, sdk.AccAddress(delAddr.Bytes()), valAddr)
	if !found {
		return DelegationMethod.Outputs.Pack(big.NewInt(0))
	}

	delegationAmt := delegation.GetShares().MulInt(validator.GetTokens()).Quo(validator.GetDelegatorShares())
	return DelegationMethod.Outputs.Pack(delegationAmt.TruncateInt().BigInt())
}
