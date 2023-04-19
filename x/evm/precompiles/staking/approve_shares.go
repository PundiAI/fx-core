package staking

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v4/x/evm/types"
)

func (c *Contract) ApproveShares(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("approve method not readonly")
	}
	// parse args
	var args ApproveSharesArgs
	if err := types.ParseMethodArgs(ApproveSharesMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	// owner
	owner := contract.Caller()
	// set allowance
	c.stakingKeeper.SetAllowance(ctx, args.GetValidator(), owner.Bytes(), args.Spender.Bytes(), args.Shares)

	// emit event
	if err := c.AddLog(ApproveSharesEvent, []common.Hash{owner.Hash(), args.Spender.Hash()}, args.Validator, args.Shares); err != nil {
		return nil, err
	}
	return ApproveSharesMethod.Outputs.Pack(true)
}
