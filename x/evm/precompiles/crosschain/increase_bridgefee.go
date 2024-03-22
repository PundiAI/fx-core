package crosschain

import (
	"errors"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/types"
)

// IncreaseBridgeFee add bridge fee to unbatched tx
//
//gocyclo:ignore
func (c *Contract) IncreaseBridgeFee(ctx sdk.Context, evm *vm.EVM, contractAddr *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("increase bridge fee method not readonly")
	}

	// args
	var args IncreaseBridgeFeeArgs
	err := types.ParseMethodArgs(IncreaseBridgeFeeMethod, &args, contractAddr.Input[4:])
	if err != nil {
		return nil, err
	}

	if c.router == nil {
		return nil, errors.New("cross chain router empty")
	}

	fxTarget := fxtypes.ParseFxTarget(args.Chain)
	route, has := c.router.GetRoute(fxTarget.GetTarget())
	if !has {
		return nil, fmt.Errorf("chain not support: %s", args.Chain)
	}

	value := contractAddr.Value()
	sender := contractAddr.Caller()
	crossChainDenom := ""
	if value.Cmp(big.NewInt(0)) == 1 && args.Token.String() == contract.EmptyEvmAddress {
		if args.Fee.Cmp(value) != 0 {
			return nil, errors.New("add bridge fee not equal msg.value")
		}
		crossChainDenom, err = c.handlerOriginToken(ctx, evm, sender, args.Fee)
		if err != nil {
			return nil, err
		}
	} else {
		crossChainDenom, err = c.handlerERC20Token(ctx, evm, args.Token, sender, args.Fee)
		if err != nil {
			return nil, err
		}
	}

	// convert token to bridge fee token
	feeCoin := sdk.NewCoin(crossChainDenom, sdkmath.NewIntFromBigInt(args.Fee))
	addBridgeFee, err := c.erc20Keeper.ConvertDenomToTarget(ctx, sender.Bytes(), feeCoin, fxTarget)
	if err != nil {
		return nil, err
	}

	if err := route.PrecompileIncreaseBridgeFee(ctx, args.TxID.Uint64(), sender.Bytes(), addBridgeFee); err != nil {
		return nil, err
	}

	// add event log
	if err := c.AddLog(evm, IncreaseBridgeFeeEvent, []common.Hash{sender.Hash(), args.Token.Hash()},
		args.Chain, args.TxID, args.Fee); err != nil {
		return nil, err
	}

	return IncreaseBridgeFeeMethod.Outputs.Pack(true)
}
