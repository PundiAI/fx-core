package crosschain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v7/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) BridgeCoinAmount(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
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
	// FX
	if fxtypes.IsZeroEthereumAddress(args.Token.Hex()) {
		supply := c.bankKeeper.GetSupply(cacheCtx, fxtypes.DefaultDenom)
		balance := c.bankKeeper.GetBalance(cacheCtx, c.accountKeeper.GetModuleAddress(ethtypes.ModuleName), fxtypes.DefaultDenom)
		return BridgeCoinAmountMethod.Outputs.Pack(supply.Amount.Sub(balance.Amount).BigInt())
	}
	// OriginDenom
	if c.erc20Keeper.IsOriginDenom(cacheCtx, pair.GetDenom()) {
		supply, err := NewContractCall(cacheCtx, evm, c.Address(), args.Token).ERC20TotalSupply()
		if err != nil {
			return nil, err
		}
		return BridgeCoinAmountMethod.Outputs.Pack(supply)
	}
	// one to one
	_, has = c.erc20Keeper.HasDenomAlias(cacheCtx, pair.GetDenom())
	if !has && pair.GetDenom() != fxtypes.DefaultDenom {
		return BridgeCoinAmountMethod.Outputs.Pack(
			c.bankKeeper.GetSupply(cacheCtx, pair.GetDenom()).Amount.BigInt(),
		)
	}
	// many to one
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
