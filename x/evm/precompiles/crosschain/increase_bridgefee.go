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

	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/evm/types"
)

var IncreaseBridgeFeeMethod = abi.NewMethod(
	IncreaseBridgeFeeMethodName,
	IncreaseBridgeFeeMethodName,
	abi.Function, "payable", false, true,
	abi.Arguments{
		abi.Argument{Name: "chain", Type: types.TypeString},
		abi.Argument{Name: "txid", Type: types.TypeUint256},
		abi.Argument{Name: "token", Type: types.TypeAddress},
		abi.Argument{Name: "fee", Type: types.TypeUint256},
	},
	abi.Arguments{
		abi.Argument{Name: "result", Type: types.TypeBool},
	},
)

// IncreaseBridgeFee add bridge fee to unbatched tx
//
//gocyclo:ignore
func (c *Contract) IncreaseBridgeFee(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("increase bridge fee method not readonly")
	}

	// args
	args, err := IncreaseBridgeFeeMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	chain, ok0 := args[0].(string)
	txID, ok1 := args[1].(*big.Int)
	token, ok2 := args[2].(common.Address)
	feeAmount, ok3 := args[3].(*big.Int)
	if !ok0 || !ok1 || !ok2 || !ok3 {
		return nil, errors.New("unexpected arg type")
	}

	if err = crosschaintypes.ValidateModuleName(chain); err != nil {
		return nil, err
	}
	if txID.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid tx id: %s", txID.String())
	}
	if feeAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid add bridge fee: %s", feeAmount.String())
	}

	if c.router == nil {
		return nil, errors.New("cross chain router empty")
	}

	fxTarget := fxtypes.ParseFxTarget(chain)
	route, has := c.router.GetRoute(fxTarget.GetTarget())
	if !has {
		return nil, fmt.Errorf("chain not support: %s", chain)
	}

	value := contract.Value()
	sender := contract.Caller()
	crossChainDenom := ""

	if value.Cmp(big.NewInt(0)) == 1 && token.String() == fxtypes.EmptyEvmAddress {
		if feeAmount.Cmp(value) != 0 {
			return nil, errors.New("add bridge fee not equal msg.value")
		}
		crossChainDenom, err = c.handlerOriginToken(ctx, evm, sender, feeAmount)
		if err != nil {
			return nil, err
		}
	} else {
		crossChainDenom, err = c.handlerERC20Token(ctx, evm, token, sender, feeAmount)
		if err != nil {
			return nil, err
		}
	}

	// convert token to bridge fee token
	feeCoin := sdk.NewCoin(crossChainDenom, sdkmath.NewIntFromBigInt(feeAmount))
	addBridgeFee, err := c.erc20Keeper.ConvertDenomToTarget(ctx, sender.Bytes(), feeCoin, fxTarget)
	if err != nil {
		return nil, err
	}

	if err := route.PrecompileIncreaseBridgeFee(ctx, txID.Uint64(), sender.Bytes(), addBridgeFee); err != nil {
		return nil, err
	}

	return IncreaseBridgeFeeMethod.Outputs.Pack(true)
}
