package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

// TransferAfter
// 1. Hook operation after transfer transaction triggered by IBC module
// 2. Hook operation after transferCrossChain triggered by ERC20 module
func (k Keeper) TransferAfter(ctx sdk.Context, sender, receive string, amount, fee sdk.Coin) error {
	sendAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = trontypes.ValidateTronAddress(receive); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid receive address: %s", err)
	}

	_, err = k.Keeper.AddToOutgoingPool(ctx, sendAddr, receive, amount, fee)
	return err
}
