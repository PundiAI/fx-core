package crosschain

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/types"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

type Keeper struct {
	router     *Router
	bankKeeper types.BankKeeper
}

func (c *Keeper) EvmTokenToBaseCoin(ctx sdk.Context, evm *vm.EVM, crosschainKeeper types.CrosschainKeeper, holder, tokenAddr common.Address, amount *big.Int) (sdk.Coin, error) {
	erc20Token, err := crosschainKeeper.GetERC20TokenByAddr(ctx, tokenAddr)
	if err != nil {
		return sdk.Coin{}, err
	}
	baseCoin := sdk.NewCoin(erc20Token.Denom, sdkmath.NewIntFromBigInt(amount))
	erc20ModuleAddress := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName))

	erc20Call := contract.NewERC20Call(evm, erc20ModuleAddress, tokenAddr, 0)
	if erc20Token.IsNativeCoin() {
		if err = erc20Call.Burn(holder, amount); err != nil {
			return sdk.Coin{}, err
		}
		if erc20Token.Denom == fxtypes.DefaultDenom {
			err = c.bankKeeper.SendCoins(ctx, tokenAddr.Bytes(), holder.Bytes(), sdk.NewCoins(baseCoin))
			return baseCoin, err
		}
	} else if erc20Token.IsNativeERC20() {
		if err = erc20Call.TransferFrom(holder, erc20ModuleAddress, amount); err != nil {
			return sdk.Coin{}, err
		}
		if err = c.bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(baseCoin)); err != nil {
			return sdk.Coin{}, err
		}
	} else {
		return sdk.Coin{}, fmt.Errorf("invalid erc20 token owner: %s", tokenAddr)
	}
	err = c.bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, holder.Bytes(), sdk.NewCoins(baseCoin))
	return baseCoin, err
}
