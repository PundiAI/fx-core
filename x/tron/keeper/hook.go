package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/functionx/fx-core/v7/types"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

var _ fxtypes.TransactionHook = &Keeper{}

// TransferAfter
// 1. Hook operation after transfer transaction triggered by IBC module
// 2. Hook operation after transferCrossChain triggered by ERC20 module
func (k Keeper) TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, amount, fee sdk.Coin, originToken bool) error {
	if err := trontypes.ValidateTronAddress(receive); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receive address: %s", err)
	}

	txID, err := k.Keeper.AddToOutgoingPool(ctx, sender, receive, amount, fee)
	if err != nil {
		return err
	}
	if !originToken {
		k.erc20Keeper.SetOutgoingTransferRelation(ctx, k.ModuleName(), txID)
	}
	return nil
}
