package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/gravity/types"
)

// TransferAfter Hook operation after transfer transaction triggered by IBC module
func (k Keeper) TransferAfter(ctx sdk.Context, sender, receive string, amount, fee sdk.Coin) error {
	// Claim channel capability passed back by IBC module
	sendAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	// verify receive address 2022-04-12
	if ctx.BlockHeight() >= fxtypes.EvmV1SupportBlock() {
		if err = types.ValidateEthAddressAndValidateChecksum(receive); err != nil {
			return err
		}
	}
	_, err = k.AddToOutgoingPool(ctx, sendAddr, receive, amount, fee)
	return err
}
