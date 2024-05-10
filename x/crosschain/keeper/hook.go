package keeper

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// TransferAfter
// 1. Hook operation after transfer transaction triggered by IBC module
// 2. Hook operation after transferCrossChain triggered by ERC20 module
func (k Keeper) TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, amount, fee sdk.Coin, originToken, insufficientLiquidity bool) error {
	if err := types.ValidateExternalAddr(k.moduleName, receive); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receive address: %s", err)
	}

	var err error
	var txID uint64
	if insufficientLiquidity {
		txID, err = k.AddToOutgoingPendingPool(ctx, sender, receive, amount, fee)
	} else {
		txID, err = k.AddToOutgoingPool(ctx, sender, receive, amount, fee)
	}

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

func (k Keeper) PrecompileBridgeCall(
	ctx sdk.Context,
	sender common.Address,
	receiver common.Address,
	coins sdk.Coins,
	to common.Address,
	data []byte,
	memo string,
) (eventNonce uint64, err error) {
	tokens, err := k.bridgeCallCoinsToERC20Token(ctx, sender.Bytes(), coins)
	if err != nil {
		return 0, err
	}

	outCall, err := k.AddOutgoingBridgeCall(
		ctx,
		sender.Bytes(),
		types.ExternalAddrToStr(k.moduleName, receiver.Bytes()),
		tokens,
		types.ExternalAddrToStr(k.moduleName, to.Bytes()),
		hex.EncodeToString(data),
		memo,
	)
	if err != nil {
		return 0, err
	}

	return outCall.Nonce, nil
}
