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

func (k Keeper) ValidateQuote(ctx sdk.Context, caller contract.Caller, quoteId *big.Int, gasLimit uint64) (contract.IBridgeFeeQuoteQuoteInfo, error) {
	bridgeFeeQuoteKeeper := contract.NewBridgeFeeQuoteKeeper(caller, contract.BridgeFeeAddress)
	quote, err := bridgeFeeQuoteKeeper.GetQuoteById(ctx, quoteId)
	if err != nil {
		return contract.IBridgeFeeQuoteQuoteInfo{}, err
	}
	if quote.IsTimeout(ctx.BlockTime()) {
		return contract.IBridgeFeeQuoteQuoteInfo{}, types.ErrInvalid.Wrapf("quote has timed out")
	}
	if quote.GasLimit < gasLimit {
		return contract.IBridgeFeeQuoteQuoteInfo{}, types.ErrInvalid.Wrapf("quote gas limit is less than gas limit")
	}
	return quote, nil
}

func (k Keeper) TransferBridgeFee(ctx sdk.Context, caller contract.Caller, from, to common.Address, bridgeFee *big.Int, bridgeTokenName string) error {
	if strings.ToUpper(bridgeTokenName) == fxtypes.DefaultDenom {
		fees := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(bridgeFee)))
		return k.bankKeeper.SendCoins(ctx, from.Bytes(), to.Bytes(), fees)
	}
	erc20Token, err := k.erc20Keeper.GetERC20Token(ctx, bridgeTokenName)
	if err != nil {
		return err
	}
	erc20TokenKeeper := contract.NewERC20TokenKeeper(caller)
	_, err = erc20TokenKeeper.Transfer(ctx, erc20Token.GetERC20Contract(), from, to, bridgeFee)
	return err
}

func (k Keeper) HandlerBridgeCallInFee(ctx sdk.Context, caller contract.Caller, from common.Address, quoteId *big.Int, gasLimit uint64) error {
	if quoteId == nil || quoteId.Sign() <= 0 {
		// Allow free bridgeCall
		return nil
	}

	quote, err := k.ValidateQuote(ctx, caller, quoteId, gasLimit)
	if err != nil {
		return err
	}

	return k.TransferBridgeFee(ctx, caller, from, quote.Oracle, quote.Amount, quote.GetTokenName())
}

func (k Keeper) HandlerBridgeCallOutFee(ctx sdk.Context, caller contract.Caller, from common.Address, bridgeCallNonce uint64, quoteId *big.Int, gasLimit uint64) error {
	if quoteId == nil || quoteId.Sign() <= 0 {
		// Users can send submitBridgeCall by themselves without paying
		return nil
	}

	quote, err := k.ValidateQuote(ctx, caller, quoteId, gasLimit)
	if err != nil {
		return err
	}

	bridgeFeeAddr := common.BytesToAddress(k.bridgeFeeCollector)
	if err = k.TransferBridgeFee(ctx, caller, from, bridgeFeeAddr, quote.Amount, quote.GetTokenName()); err != nil {
		return err
	}

	k.SetOutgoingBridgeCallQuoteInfo(ctx, bridgeCallNonce, types.NewQuoteInfo(quote))
	return nil
}

func (k Keeper) TransferBridgeFeeToRelayer(ctx sdk.Context, caller contract.Caller, bridgeCallNonce uint64) error {
	quote, found := k.GetOutgoingBridgeCallQuoteInfo(ctx, bridgeCallNonce)
	if !found {
		return nil
	}

	k.DeleteOutgoingBridgeCallQuoteInfo(ctx, bridgeCallNonce)

	bridgeFeeAddr := common.BytesToAddress(k.bridgeFeeCollector)
	return k.TransferBridgeFee(ctx, caller, bridgeFeeAddr, quote.OracleAddress(), quote.Fee.BigInt(), quote.Token)
}

func (k Keeper) SetOutgoingBridgeCallQuoteInfo(ctx sdk.Context, nonce uint64, quoteInfo types.QuoteInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeCallQuoteKey(nonce), k.cdc.MustMarshal(&quoteInfo))
}

func (k Keeper) GetOutgoingBridgeCallQuoteInfo(ctx sdk.Context, nonce uint64) (types.QuoteInfo, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetBridgeCallQuoteKey(nonce))
	if bz == nil {
		return types.QuoteInfo{}, false
	}

	quoteInfo := types.QuoteInfo{}
	k.cdc.MustUnmarshal(bz, &quoteInfo)
	return quoteInfo, true
}

func (k Keeper) DeleteOutgoingBridgeCallQuoteInfo(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBridgeCallQuoteKey(nonce))
}
