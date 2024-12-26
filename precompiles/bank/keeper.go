package bank

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/types"
)

type Keeper struct {
	bankKeeper  types.BankKeeper
	erc20Keeper types.Erc20Keeper
}

func NewKeeper(bankKeeper types.BankKeeper, erc20Keeper types.Erc20Keeper) *Keeper {
	return &Keeper{
		bankKeeper:  bankKeeper,
		erc20Keeper: erc20Keeper,
	}
}

func (k *Keeper) TransferFromModuleToAccount(ctx sdk.Context, args *contract.TransferFromModuleToAccountArgs) error {
	denom, err := k.erc20Keeper.GetBaseDenom(ctx, args.Token.String())
	if err != nil {
		return err
	}
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(args.Amount)))
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, args.Module, args.Account.Bytes(), coins)
}

func (k *Keeper) TransferFromAccountToModule(ctx sdk.Context, args *contract.TransferFromAccountToModuleArgs) error {
	denom, err := k.erc20Keeper.GetBaseDenom(ctx, args.Token.String())
	if err != nil {
		return err
	}
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(args.Amount)))
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, args.Account.Bytes(), args.Module, coins)
}
