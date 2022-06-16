package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/gravity/types"
)

// TransferAfter Hook operation after transfer transaction triggered by IBC module
func (k Keeper) TransferAfter(ctx sdk.Context, sender, receive string, amount, fee sdk.Coin) error {
	// Claim channel capability passed back by IBC module
	sendAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	if err = types.ValidateEthAddressAndValidateChecksum(receive); err != nil {
		return err
	}
	_, err = k.AddToOutgoingPool(ctx, sendAddr, receive, amount, fee)
	return err
}
