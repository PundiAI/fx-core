package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (k Keeper) BridgeCallHandler(ctx sdk.Context, msg *types.MsgBridgeCallClaim) error {
	k.CreateBridgeAccount(ctx, msg.TxOrigin)
	if senderAccount := k.ak.GetAccount(ctx, msg.GetSenderAddr().Bytes()); senderAccount != nil {
		if _, ok := senderAccount.(sdk.ModuleAccountI); ok {
			return types.ErrInvalid.Wrap("sender is module account")
		}
	}
	isMemoSendCallTo := msg.IsMemoSendCallTo()
	receiverAddr := msg.GetToAddr()
	if isMemoSendCallTo {
		receiverAddr = msg.GetSenderAddr()
	}

	baseCoins := sdk.NewCoins()
	for i, tokenAddr := range msg.TokenContracts {
		bridgeToken, err := k.DepositBridgeToken(ctx, receiverAddr.Bytes(), msg.Amounts[i], tokenAddr)
		if err != nil {
			return err
		}
		baseCoin, err := k.BridgeTokenToBaseCoin(ctx, receiverAddr.Bytes(), msg.Amounts[i], bridgeToken)
		if err != nil {
			return err
		}
		baseCoins = baseCoins.Add(baseCoin)
	}

	if err := k.handleBridgeCallInQuote(ctx, msg.GetSenderAddr(), msg.QuoteId.BigInt(), msg.GasLimit.BigInt()); err != nil {
		return err
	}

	cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()
	err := k.BridgeCallEvm(cacheCtx, msg.GetSenderAddr(), msg.GetRefundAddr(), msg.GetToAddr(),
		receiverAddr, baseCoins, msg.MustData(), msg.MustMemo(), isMemoSendCallTo, msg.GetGasLimit())
	if !ctx.IsCheckTx() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "bridge_call_in"},
			float32(1),
			[]metrics.Label{
				telemetry.NewLabel("module", k.moduleName),
				telemetry.NewLabel("success", strconv.FormatBool(err == nil)),
			},
		)
	}
	if err == nil {
		commit()
		return nil
	}
	revertMsg := err.Error()

	// refund bridge-call case of error
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeBridgeCallEvent, sdk.NewAttribute(types.AttributeKeyErrCause, err.Error())))

	if !baseCoins.Empty() {
		if err = k.bankKeeper.SendCoins(ctx, receiverAddr.Bytes(), msg.GetRefundAddr().Bytes(), baseCoins); err != nil {
			return err
		}

		for _, coin := range baseCoins {
			if fxtypes.IsOriginDenom(coin.Denom) {
				continue
			}
			if _, err = k.erc20Keeper.BaseCoinToEvm(ctx, msg.GetRefundAddr(), coin); err != nil {
				return err
			}
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallFailed,
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", msg.EventNonce)),
		sdk.NewAttribute(types.AttributeKeyBridgeCallFailedRefundAddr, msg.GetRefundAddr().Hex()),
	))

	// onRevert bridgecall
	_, err = k.AddOutgoingBridgeCall(ctx, msg.GetToAddr(), common.Address{}, sdk.NewCoins(),
		msg.GetSenderAddr(), []byte(revertMsg), []byte{}, msg.EventNonce)
	return err
}

func (k Keeper) BridgeCallEvm(ctx sdk.Context, sender, refundAddr, to, receiverAddr common.Address, baseCoins sdk.Coins, data, memo []byte, isMemoSendCallTo bool, gasLimit uint64) error {
	tokens := make([]common.Address, 0, baseCoins.Len())
	amounts := make([]*big.Int, 0, baseCoins.Len())
	for _, coin := range baseCoins {
		tokenContract, err := k.erc20Keeper.BaseCoinToEvm(ctx, receiverAddr, coin)
		if err != nil {
			return err
		}
		tokens = append(tokens, common.HexToAddress(tokenContract))
		amounts = append(amounts, coin.Amount.BigInt())
	}

	if !k.evmKeeper.IsContract(ctx, to) {
		return nil
	}
	var callEvmSender common.Address
	var args []byte

	if isMemoSendCallTo {
		args = data
		callEvmSender = sender
	} else {
		var err error
		args, err = contract.PackOnBridgeCall(sender, refundAddr, tokens, amounts, data, memo)
		if err != nil {
			return err
		}
		callEvmSender = k.GetCallbackFrom()
	}

	if gasLimit == 0 {
		gasLimit = k.GetBridgeCallMaxGasLimit(ctx)
	}
	txResp, err := k.evmKeeper.ExecuteEVM(ctx, callEvmSender, &to, nil, gasLimit, args)
	if err != nil {
		return err
	}
	if txResp.Failed() {
		return types.ErrInvalid.Wrap(txResp.VmError)
	}
	return nil
}

func (k Keeper) handleBridgeCallInQuote(ctx sdk.Context, from common.Address, quoteId, gasLimit *big.Int) error {
	if quoteId == nil || quoteId.Sign() <= 0 || gasLimit == nil || gasLimit.Sign() <= 0 {
		return nil
	}

	contractQuote, err := k.validatorQuoteGasLimit(ctx, quoteId, gasLimit)
	if err != nil {
		return err
	}

	// transfer fee to quote oracle
	bridgeToken, err := k.erc20Keeper.GetBridgeToken(ctx, k.moduleName, contractQuote.TokenName)
	if err != nil {
		return err
	}

	if bridgeToken.IsOrigin() {
		fees := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(contractQuote.Fee)))
		return k.bankKeeper.SendCoins(ctx, from.Bytes(), contractQuote.Oracle.Bytes(), fees)
	}

	_, err = k.erc20TokenKeeper.Transfer(ctx, bridgeToken.GetContractAddress(), from, contractQuote.Oracle, contractQuote.Fee)
	return err
}

func (k Keeper) validatorQuoteGasLimit(ctx sdk.Context, quoteId, gasLimit *big.Int) (contract.IBridgeFeeQuoteQuoteInfo, error) {
	contractQuote, err := k.bridgeFeeQuoteKeeper.GetQuoteById(ctx, quoteId)
	if err != nil {
		return contract.IBridgeFeeQuoteQuoteInfo{}, err
	}
	if contractQuote.IsTimeout(ctx.BlockTime()) {
		return contract.IBridgeFeeQuoteQuoteInfo{}, types.ErrInvalid.Wrapf("quote has timed out")
	}
	if contractQuote.GasLimit.Cmp(gasLimit) < 0 {
		return contract.IBridgeFeeQuoteQuoteInfo{}, types.ErrInvalid.Wrapf("quote gas limit is less than gas limit")
	}
	return contractQuote, nil
}
