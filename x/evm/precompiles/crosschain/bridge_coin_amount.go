package crosschain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v4/types"
	"github.com/functionx/fx-core/v4/x/evm/types"
)

func (c *Contract) BridgeCoinAmount(ctx sdk.Context, _ *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	cacheCtx, _ := ctx.CacheContext()
	// parse args
	var args BridgeCoinAmountArgs
	if err := types.ParseMethodArgs(BridgeCoinAmountMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	pair, has := c.erc20Keeper.GetTokenPair(cacheCtx, args.Token.Hex())
	if !has {
		return nil, fmt.Errorf("token not support: %s", args.Token.Hex())
	}
	md, has := c.bankKeeper.GetDenomMetaData(cacheCtx, pair.GetDenom())
	if !has {
		return nil, fmt.Errorf("denom not support: %s", pair.GetDenom())
	}
	denom := c.erc20Keeper.ToTargetDenom(
		cacheCtx,
		pair.GetDenom(),
		md.GetBase(),
		md.GetDenomUnits()[0].GetAliases(),
		fxtypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target)),
	)

	balance := c.bankKeeper.GetBalance(cacheCtx, c.erc20Keeper.ModuleAddress().Bytes(), pair.GetDenom())
	supply := c.bankKeeper.GetSupply(cacheCtx, denom)
	if balance.Amount.LT(supply.Amount) {
		return BridgeCoinAmountMethod.Outputs.Pack(balance.Amount.BigInt())
	}
	return BridgeCoinAmountMethod.Outputs.Pack(supply.Amount.BigInt())
}
