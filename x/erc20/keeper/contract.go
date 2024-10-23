package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func (k Keeper) ERC20BaseInfo(ctx context.Context, contractAddr common.Address) (name, symbol string, decimals uint8, err error) {
	name, err = k.evmErc20Keeper.Name(ctx, contractAddr)
	if err != nil {
		return name, symbol, decimals, err
	}
	symbol, err = k.evmErc20Keeper.Symbol(ctx, contractAddr)
	if err != nil {
		return name, symbol, decimals, err
	}
	decimals, err = k.evmErc20Keeper.Decimals(ctx, contractAddr)
	if err != nil {
		return name, symbol, decimals, err
	}
	return name, symbol, decimals, err
}

func (k Keeper) DeployUpgradableToken(ctx sdk.Context, from common.Address, name, symbol string, decimals uint8) (common.Address, error) {
	var tokenContract contract.Contract
	if symbol == fxtypes.DefaultDenom {
		tokenContract = contract.GetWFX()
		name = fmt.Sprintf("Wrapped %s", name)
		symbol = fmt.Sprintf("W%s", symbol)
	} else {
		tokenContract = contract.GetFIP20()
	}
	k.Logger(ctx).Info("deploy token contract", "name", name, "symbol", symbol, "decimals", decimals)

	return k.evmKeeper.DeployUpgradableContract(ctx, from, tokenContract.Address, nil, &tokenContract.ABI, name, symbol, decimals, k.contractOwner)
}
