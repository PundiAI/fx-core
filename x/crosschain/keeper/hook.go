package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

// TransferAfter
// 1. Hook operation after transfer transaction triggered by IBC module
// 2. Hook operation after transferCrossChain triggered by ERC20 module
func (k Keeper) TransferAfter(ctx sdk.Context, sender, receive string, amount, fee sdk.Coin) error {
	sendAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	if err = fxtypes.ValidateEthereumAddress(receive); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receive address: %s", err)
	}

	_, err = k.AddToOutgoingPool(ctx, sendAddr, receive, amount, fee)
	return err
}
