package crosschain

import (
	"errors"
	"math/big"

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
	if c.router == nil {
		return nil, errors.New("bridge call router is empty")
	}

	var args BridgeCallArgs
	if err := types.ParseMethodArgs(BridgeCallMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	route, has := c.router.GetRoute(args.DstChain)
	if !has {
		return nil, errors.New("invalid dstChain")
	}
	sender := contract.Caller()

	coins := sdk.NewCoins()
	for i, token := range args.Tokens {
		coin, err := c.handlerERC20Token(ctx, evm, sender, token, args.Amounts[i])
		if err != nil {
			return nil, err
		}
		coins = coins.Add(coin)
	}
	value := contract.Value()
	if value.Cmp(big.NewInt(0)) == 1 {
		totalCoin, err := c.handlerOriginToken(ctx, evm, sender, value)
		if err != nil {
			return nil, err
		}
		coins = coins.Add(totalCoin)
	}

	eventNonce, err := route.PrecompileBridgeCall(
		ctx,
		sender,
		args.Receiver,
		coins,
		args.To,
		args.Data,
		args.Memo,
	)
	if err != nil {
		return nil, err
	}

	if err = c.AddLog(evm, BridgeCallEvent,
		[]common.Hash{sender.Hash(), args.Receiver.Hash(), args.To.Hash()},
		evm.Origin,
		args.Value,
		sdkmath.NewIntFromUint64(eventNonce).BigInt(),
		args.DstChain,
		args.Tokens,
		args.Amounts,
		args.Data,
		args.Memo,
	); err != nil {
		return nil, err
	}
	return BridgeCallMethod.Outputs.Pack(eventNonce)
}
