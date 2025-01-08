package keeper

import (
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (k Keeper) ValidateQuote(ctx sdk.Context, quoteId, gasLimit *big.Int) (contract.IBridgeFeeQuoteQuoteInfo, error) {
	quote, err := k.bridgeFeeQuoteKeeper.GetQuoteById(ctx, quoteId)
	if err != nil {
		return contract.IBridgeFeeQuoteQuoteInfo{}, err
	}
	if quote.IsTimeout(ctx.BlockTime()) {
		return contract.IBridgeFeeQuoteQuoteInfo{}, types.ErrInvalid.Wrapf("quote has timed out")
	}
	if quote.GasLimit.Cmp(gasLimit) < 0 {
		return contract.IBridgeFeeQuoteQuoteInfo{}, types.ErrInvalid.Wrapf("quote gas limit is less than gas limit")
	}
	return quote, nil
}

func (k Keeper) TransferBridgeFee(ctx sdk.Context, from, to common.Address, bridgeFee *big.Int, bridgeTokenName string) error {
	if strings.ToUpper(bridgeTokenName) == fxtypes.DefaultDenom {
		fees := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(bridgeFee)))
		return k.bankKeeper.SendCoins(ctx, from.Bytes(), to.Bytes(), fees)
	}
	erc20Token, err := k.erc20Keeper.GetERC20Token(ctx, bridgeTokenName)
	if err != nil {
		return err
	}
	_, err = k.erc20TokenKeeper.Transfer(ctx, erc20Token.GetERC20Contract(), from, to, bridgeFee)
	return err
}
