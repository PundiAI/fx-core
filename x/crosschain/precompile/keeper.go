package precompile

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

type Keeper struct {
	router      *Router
	bankKeeper  BankKeeper
	erc20Keeper Erc20Keeper
}

func (c *Keeper) EvmTokenToBase(ctx sdk.Context, evm *vm.EVM, crossChainKeeper CrosschainKeeper, holder, token common.Address, amount *big.Int) (sdk.Coin, error) {
	baseDenom, isNativeCoin, err := crossChainKeeper.GetBaseDenomByErc20(ctx, token)
	if err != nil {
		return sdk.Coin{}, err
	}
	baseCoin := sdk.NewCoin(baseDenom, sdkmath.NewIntFromBigInt(amount))

	if isNativeCoin {
		erc20Call := contract.NewERC20Call(evm, c.erc20Keeper.ModuleAddress(), token, 0)
		if err = erc20Call.Burn(holder, amount); err != nil {
			return sdk.Coin{}, err
		}
		if baseDenom == fxtypes.DefaultDenom {
			err = c.bankKeeper.SendCoinsFromAccountToModule(ctx, token.Bytes(), erc20types.ModuleName, sdk.NewCoins(baseCoin))
			if err != nil {
				return sdk.Coin{}, err
			}
		}
	} else {
		// transferFrom to erc20 module
		erc20Call := contract.NewERC20Call(evm, crosschaintypes.GetAddress(), token, 0)
		if err := erc20Call.TransferFrom(holder, c.erc20Keeper.ModuleAddress(), amount); err != nil {
			return sdk.Coin{}, err
		}
		if err = c.bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(baseCoin)); err != nil {
			return sdk.Coin{}, err
		}
	}
	if err = c.bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, holder.Bytes(), sdk.NewCoins(baseCoin)); err != nil {
		return sdk.Coin{}, err
	}
	return baseCoin, nil
}
