package precompile

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) BridgeCoinAmount(ctx sdk.Context, evm *vm.EVM, contractAddr *vm.Contract, _ bool) ([]byte, error) {
	cacheCtx, _ := ctx.CacheContext()

	var args crosschaintypes.BridgeCoinAmountArgs
	if err := evmtypes.ParseMethodArgs(crosschaintypes.BridgeCoinAmountMethod, &args, contractAddr.Input[4:]); err != nil {
		return nil, err
	}
	pair, has := c.erc20Keeper.GetTokenPair(cacheCtx, args.Token.Hex())
	if !has {
		return nil, fmt.Errorf("token not support: %s", args.Token.Hex())
	}
	// FX
	if contract.IsZeroEthAddress(args.Token) {
		supply := c.bankKeeper.GetSupply(cacheCtx, fxtypes.DefaultDenom)
		balance := c.bankKeeper.GetBalance(cacheCtx, c.accountKeeper.GetModuleAddress(ethtypes.ModuleName), fxtypes.DefaultDenom)
		return crosschaintypes.BridgeCoinAmountMethod.Outputs.Pack(supply.Amount.Sub(balance.Amount).BigInt())
	}
	// OriginDenom
	if c.erc20Keeper.IsOriginDenom(cacheCtx, pair.GetDenom()) {
		erc20Call := contract.NewERC20Call(evm, c.Address(), args.Token, c.GetBlockGasLimit())
		supply, err := erc20Call.TotalSupply()
		if err != nil {
			return nil, err
		}
		return crosschaintypes.BridgeCoinAmountMethod.Outputs.Pack(supply)
	}
	// one to one
	_, has = c.erc20Keeper.HasDenomAlias(cacheCtx, pair.GetDenom())
	if !has && pair.GetDenom() != fxtypes.DefaultDenom {
		return crosschaintypes.BridgeCoinAmountMethod.Outputs.Pack(
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
		return crosschaintypes.BridgeCoinAmountMethod.Outputs.Pack(balance.Amount.BigInt())
	}
	return crosschaintypes.BridgeCoinAmountMethod.Outputs.Pack(supply.Amount.BigInt())
}
