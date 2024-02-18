package crosschain

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) CancelSendToExternal(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("cancel send to external method not readonly")
	}

	// args
	var args CancelSendToExternalArgs
	if err := types.ParseMethodArgs(CancelSendToExternalMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	sender := contract.Caller()
	route, has := c.router.GetRoute(args.Chain)
	if !has {
		return nil, fmt.Errorf("chain not support: %s", args.Chain)
	}

	originDenom := c.evmKeeper.GetParams(ctx).EvmDenom
	// NOTE: must be get relation before cancel, cancel will delete it if relation exist
	hasRelation := c.erc20Keeper.HasOutgoingTransferRelation(ctx, args.Chain, args.TxID.Uint64())

	refundCoin, err := route.PrecompileCancelSendToExternal(ctx, args.TxID.Uint64(), sender.Bytes())
	if err != nil {
		return nil, err
	}
	if !hasRelation && refundCoin.Denom == originDenom {
		// add refund to sender in evm state db, because bank keeper add refund to sender
		evm.StateDB.AddBalance(sender, refundCoin.Amount.BigInt())
	}

	// add event log
	if err := c.AddLog(evm, CancelSendToExternalEvent, []common.Hash{sender.Hash()}, args.Chain, args.TxID); err != nil {
		return nil, err
	}

	return CancelSendToExternalMethod.Outputs.Pack(true)
}
