package crosschain

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) BridgeCall(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("bridge call method not readonly")
	}

	// args
	var args BridgeCallArgs
	if err := types.ParseMethodArgs(BridgeCallMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	sender := contract.Caller()

	if c.router == nil {
		return nil, errors.New("cross chain router empty")
	}
	route, has := c.router.GetRoute(args.DstChainId)
	if !has {
		return nil, errors.New("invalid target")
	}
	eventNonce, err := route.PrecompileBridgeCall(ctx, args.DstChainId, args.GasLimit.Uint64(), sender,
		args.Receiver, args.To, args.Asset, args.Message, args.Value)
	if err != nil {
		return nil, err
	}

	// add event log
	if err = c.AddLog(evm, BridgeCallEvent, []common.Hash{sender.Hash(), args.Receiver.Hash(), args.To.Hash()},
		sdkmath.NewIntFromUint64(eventNonce).BigInt(), args.DstChainId, args.GasLimit, args.Value, args.Message, args.Asset); err != nil {
		return nil, err
	}
	return BridgeCallMethod.Outputs.Pack(true)
}
