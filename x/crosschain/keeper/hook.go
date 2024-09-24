package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
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

func (k Keeper) PrecompileBridgeCall(ctx sdk.Context, sender, refund common.Address, coins sdk.Coins, to common.Address, data, memo []byte) (nonce uint64, err error) {
	tokens, notLiquidCoins, err := k.BridgeCallCoinsToERC20Token(ctx, sender.Bytes(), coins)
	if err != nil {
		return 0, err
	}

	var outCallNonce uint64
	if len(notLiquidCoins) > 0 {
		outCallNonce, err = k.AddPendingOutgoingBridgeCall(ctx, sender, refund, tokens, to, data, memo, 0, notLiquidCoins)
	} else {
		outCallNonce, err = k.AddOutgoingBridgeCall(ctx, sender, refund, tokens, to, data, memo, 0)
	}

	if err != nil {
		return 0, err
	}

	return outCallNonce, nil
}
