package staking

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v4/x/evm/types"
)

var (
	ApproveSharesMethod = abi.NewMethod(
		ApproveSharesMethodName,
		ApproveSharesMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_val", Type: types.TypeString},
			abi.Argument{Name: "_spender", Type: types.TypeAddress},
			abi.Argument{Name: "_shares", Type: types.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: types.TypeBool},
		},
	)

	ApproveSharesEvent = abi.NewEvent(
		ApproveSharesEventName,
		ApproveSharesEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "owner", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "spender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
		},
	)
)

type ApproveSharesArgs struct {
	Validator string         `abi:"_val"`
	Spender   common.Address `abi:"_spender"`
	Shares    *big.Int       `abi:"_shares"`
}

func (c *Contract) ApproveShares(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("approve method not readonly")
	}
	// parse args
	var args ApproveSharesArgs
	if err := ParseMethodParams(ApproveSharesMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(args.Validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", args.Validator)
	}
	if args.Shares.Cmp(big.NewInt(0)) < 0 {
		return nil, errors.New("allowance cannot be negative")
	}
	// owner
	owner := contract.Caller()

	// set allowance
	c.stakingKeeper.SetAllowance(ctx, valAddr, owner.Bytes(), args.Spender.Bytes(), args.Shares)

	// emit event
	if err := c.AddLog(ApproveSharesEvent, []common.Hash{owner.Hash(), args.Spender.Hash()}, args.Validator, args.Shares); err != nil {
		return nil, err
	}
	return ApproveSharesMethod.Outputs.Pack(true)
}
