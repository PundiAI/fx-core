package crosschain

import (
	"errors"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v4/types"
	crosschaintypes "github.com/functionx/fx-core/v4/x/crosschain/types"
	"github.com/functionx/fx-core/v4/x/evm/types"
)

var (
	IncreaseBridgeFeeMethod = abi.NewMethod(
		IncreaseBridgeFeeMethodName,
		IncreaseBridgeFeeMethodName,
		abi.Function, "payable", false, false,
		abi.Arguments{
			abi.Argument{Name: "_chain", Type: types.TypeString},
			abi.Argument{Name: "_txID", Type: types.TypeUint256},
			abi.Argument{Name: "_token", Type: types.TypeAddress},
			abi.Argument{Name: "_fee", Type: types.TypeUint256},
		},
		abi.Arguments{
			abi.Argument{Name: "_result", Type: types.TypeBool},
		},
	)

	IncreaseBridgeFeeEvent = abi.NewEvent(
		IncreaseBridgeFeeEventName,
		IncreaseBridgeFeeEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "token", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "chain", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "txID", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "fee", Type: types.TypeUint256, Indexed: false},
		})
)

type IncreaseBridgeFeeArgs struct {
	Chain string         `abi:"_chain"`
	TxID  *big.Int       `abi:"_txID"`
	Token common.Address `abi:"_token"`
	Fee   *big.Int       `abi:"_fee"`
}

// IncreaseBridgeFee add bridge fee to unbatched tx
//
//gocyclo:ignore
func (c *Contract) IncreaseBridgeFee(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("increase bridge fee method not readonly")
	}

	// args
	var args IncreaseBridgeFeeArgs
	err := ParseMethodParams(IncreaseBridgeFeeMethod, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}

	if err = crosschaintypes.ValidateModuleName(args.Chain); err != nil {
		return nil, err
	}
	if args.TxID.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid tx id: %s", args.TxID.String())
	}
	if args.Fee.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid add bridge fee: %s", args.Fee.String())
	}

	if c.router == nil {
		return nil, errors.New("cross chain router empty")
	}

	fxTarget := fxtypes.ParseFxTarget(args.Chain)
	route, has := c.router.GetRoute(fxTarget.GetTarget())
	if !has {
		return nil, fmt.Errorf("chain not support: %s", args.Chain)
	}

	value := contract.Value()
	sender := contract.Caller()
	crossChainDenom := ""
	if value.Cmp(big.NewInt(0)) == 1 && args.Token.String() == fxtypes.EmptyEvmAddress {
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
	if err := c.AddLog(IncreaseBridgeFeeEvent, []common.Hash{sender.Hash(), args.Token.Hash()},
		args.Chain, args.TxID, args.Fee); err != nil {
		return nil, err
	}

	return IncreaseBridgeFeeMethod.Outputs.Pack(true)
}
