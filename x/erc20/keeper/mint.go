package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/v7/x/erc20/types"
)

// MintingEnabled checks that:
//   - the global parameter for intrarelaying is enabled
//   - minting is enabled for the given (erc20,coin) token pair
//   - recipient address is not on the blocked list
//   - bank module transfers are enabled for the Cosmos coin
func (k Keeper) MintingEnabled(ctx sdk.Context, receiver sdk.AccAddress, token string) (types.TokenPair, error) {
	if !k.GetEnableErc20(ctx) {
		return types.TokenPair{}, errorsmod.Wrap(types.ErrERC20Disabled, "module is currently disabled by governance")
	}

	pair, found := k.GetTokenPair(ctx, token)
	if !found {
		return types.TokenPair{}, errorsmod.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered", token)
	}

	if !pair.Enabled {
		return pair, errorsmod.Wrapf(types.ErrERC20TokenPairDisabled, "minting token '%s' is not enabled by governance", token)
	}

	if k.bankKeeper.BlockedAddr(receiver.Bytes()) {
		return pair, errorsmod.Wrapf(errortypes.ErrUnauthorized, "%s is not allowed to receive transactions", receiver)
	}

	// NOTE: ignore amount as only denom is checked on IsSendEnabledCoin
	coin := sdk.Coin{Denom: pair.Denom}

	// check if minting to a recipient address other than the sender is enabled for the given coin denom
	// if coin disable and sender not equal receiver, can not convert denom
	if !k.bankKeeper.IsSendEnabledCoin(ctx, coin) {
		return pair, errorsmod.Wrapf(banktypes.ErrSendDisabled, "minting '%s' coins to an external address is currently disabled", token)
	}

	return pair, nil
}
