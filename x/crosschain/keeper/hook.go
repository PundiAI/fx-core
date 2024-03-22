package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/v7/contract"
)

// TransferAfter
// 1. Hook operation after transfer transaction triggered by IBC module
// 2. Hook operation after transferCrossChain triggered by ERC20 module
func (k Keeper) TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, amount, fee sdk.Coin, originToken bool) error {
	if err := contract.ValidateEthereumAddress(receive); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receive address: %s", err)
	}

	txID, err := k.AddToOutgoingPool(ctx, sender, receive, amount, fee)
	if err != nil {
		return err
	}

	if !originToken {
		k.erc20Keeper.SetOutgoingTransferRelation(ctx, k.moduleName, txID)
	}
	return nil
}

func (k Keeper) PrecompileCancelSendToExternal(ctx sdk.Context, txID uint64, sender sdk.AccAddress) (sdk.Coin, error) {
	return k.RemoveFromOutgoingPoolAndRefund(ctx, txID, sender)
}

func (k Keeper) PrecompileIncreaseBridgeFee(ctx sdk.Context, txID uint64, sender sdk.AccAddress, addBridgeFee sdk.Coin) error {
	return k.AddUnbatchedTxBridgeFee(ctx, txID, sender, addBridgeFee)
}
