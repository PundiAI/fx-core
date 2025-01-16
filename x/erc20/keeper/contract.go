package keeper

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
)

func (k Keeper) ERC20BaseInfo(ctx context.Context, contractAddr common.Address) (name, symbol string, decimals uint8, err error) {
	evmErc20Keeper := contract.NewERC20TokenKeeper(k.evmKeeper)
	name, err = evmErc20Keeper.Name(ctx, contractAddr)
	if err != nil {
		return name, symbol, decimals, err
	}
	symbol, err = evmErc20Keeper.Symbol(ctx, contractAddr)
	if err != nil {
		return name, symbol, decimals, err
	}
	decimals, err = evmErc20Keeper.Decimals(ctx, contractAddr)
	if err != nil {
		return name, symbol, decimals, err
	}
	return name, symbol, decimals, err
}
