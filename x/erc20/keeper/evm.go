package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

// QueryERC20 returns the data of a deployed ERC20 contract
func (k Keeper) QueryERC20(ctx sdk.Context, contractAddr common.Address) (types.ERC20Data, error) {
	name, err := k.evmErc20Keeper.ERE20Name(ctx, contractAddr)
	if err != nil {
		return types.ERC20Data{}, err
	}

	symbol, err := k.evmErc20Keeper.ERE20Symbol(ctx, contractAddr)
	if err != nil {
		return types.ERC20Data{}, err
	}

	decimals, err := k.evmErc20Keeper.ERE20Decimals(ctx, contractAddr)
	if err != nil {
		return types.ERC20Data{}, err
	}

	return types.NewERC20Data(name, symbol, decimals), nil
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

	return k.evmKeeper.DeployUpgradableContract(ctx, from, tokenContract.Address, nil, &tokenContract.ABI, name, symbol, decimals, k.moduleAddress)
}
