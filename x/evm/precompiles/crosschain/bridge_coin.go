package crosschain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v4/types"
	"github.com/functionx/fx-core/v4/x/evm/types"
)

var BridgeCoinMethod = abi.NewMethod(
	BridgeCoinMethodName,
	BridgeCoinMethodName,
	abi.Function, "view", false, false,
	abi.Arguments{
		abi.Argument{Name: "_token", Type: types.TypeAddress},
		abi.Argument{Name: "_target", Type: types.TypeBytes32},
	},
	abi.Arguments{
		abi.Argument{Name: "_amount", Type: types.TypeUint256},
	},
)

type BridgeCoinArgs struct {
	Token  common.Address `abi:"_token"`
	Target [32]byte       `abi:"_target"`
}

func (c *Contract) BridgeCoin(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	cacheCtx, _ := ctx.CacheContext()
	// parse args
	var args BridgeCoinArgs
	if err := ParseMethodParams(BridgeCoinMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	pair, has := c.erc20Keeper.GetTokenPair(sdk.UnwrapSDKContext(cacheCtx), args.Token.Hex())
	if !has {
		return nil, fmt.Errorf("token not support: %s", args.Token.Hex())
	}
	md, has := c.bankKeeper.GetDenomMetaData(sdk.UnwrapSDKContext(cacheCtx), pair.GetDenom())
	if !has {
		return nil, fmt.Errorf("denom not support: %s", pair.GetDenom())
	}
	denom := c.erc20Keeper.ToTargetDenom(
		sdk.UnwrapSDKContext(cacheCtx),
		pair.GetDenom(),
		pair.GetDenom(),
		md.GetDenomUnits()[0].GetAliases(),
		fxtypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target)),
	)

	balance := c.bankKeeper.GetBalance(sdk.UnwrapSDKContext(cacheCtx), c.erc20Keeper.ModuleAddress().Bytes(), pair.GetDenom())
	supply := c.bankKeeper.GetSupply(sdk.UnwrapSDKContext(cacheCtx), denom)
	if balance.Amount.LT(supply.Amount) {
		return BridgeCoinMethod.Outputs.Pack(balance.Amount.BigInt())
	}
	return BridgeCoinMethod.Outputs.Pack(supply.Amount.BigInt())
}
