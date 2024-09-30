package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

// TransferAfter
// 1. Hook operation after transfer transaction triggered by IBC module
// 2. Hook operation after transferCrossChain triggered by ERC20 module
func (k Keeper) TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, amount, fee sdk.Coin, originToken bool) error {
	if err := types.ValidateExternalAddr(k.moduleName, receive); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid receive address: %s", err)
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

func (k Keeper) PrecompileBridgeCall(ctx sdk.Context, sender, refund common.Address, coins sdk.Coins, to common.Address, data, memo []byte) (nonce uint64, err error) {
	tokens, err := k.BridgeCallCoinsToERC20Token(ctx, sender.Bytes(), coins)
	if err != nil {
		return 0, err
	}

	outCallNonce, err := k.AddOutgoingBridgeCall(ctx, sender, refund, tokens, to, data, memo, 0)
	if err != nil {
		return 0, err
	}

	return outCallNonce, nil
}
