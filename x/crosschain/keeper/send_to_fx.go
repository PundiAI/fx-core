package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/pundiai/fx-core/v8/contract"
	fxtelemetry "github.com/pundiai/fx-core/v8/telemetry"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (k Keeper) SendToFxExecuted(ctx sdk.Context, caller contract.Caller, claim *types.MsgSendToFxClaim) error {
	if !ctx.IsCheckTx() {
		defer func() {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "send_to_fx"},
				float32(1),
				[]metrics.Label{
					telemetry.NewLabel("module", k.moduleName),
				},
			)
			fxtelemetry.SetGaugeLabelsWithDenom(
				[]string{types.ModuleName, "send_to_fx_amount"},
				claim.TokenContract, claim.Amount.BigInt(),
				telemetry.NewLabel("module", k.moduleName),
			)
		}()
	}

	receiveAddr, err := sdk.AccAddressFromBech32(claim.Receiver)
	if err != nil {
		return types.ErrInvalid.Wrapf("receiver address")
	}

	amount := claim.Amount
	bridgeToken, err := k.DepositBridgeToken(ctx, receiveAddr, amount, claim.TokenContract)
	if err != nil {
		return err
	}
	bridgeToken, amount, err = k.SwapBridgeToken(ctx, receiveAddr, bridgeToken, amount)
	if err != nil {
		return err
	}
	baseCoin, err := k.BridgeTokenToBaseCoin(ctx, receiveAddr, amount, bridgeToken)
	if err != nil {
		return err
	}

	event := sdk.NewEvent(
		types.EventTypeEvmTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
	)

	if target, err := types.ParseFxTarget(claim.GetTargetIbc(), true); err == nil && target.IsIBC() {
		if ibcReceiveAddress, err := bech32.ConvertAndEncode(target.Bech32Prefix, receiveAddr); err == nil {
			cacheCtx, commit := ctx.CacheContext()
			ibcSeq, err := k.IBCTransfer(cacheCtx, receiveAddr, ibcReceiveAddress, baseCoin, target.IBCChannel, "")
			if err == nil {
				commit()
				event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyIBCSequence, strconv.FormatUint(ibcSeq, 10)))
				ctx.EventManager().EmitEvent(event)
				return nil
			}
		}
	}

	if !bridgeToken.IsOrigin() && baseCoin.IsPositive() {
		_, err = k.erc20Keeper.BaseCoinToEvm(ctx, caller, common.BytesToAddress(receiveAddr.Bytes()), baseCoin)
		if err != nil {
			return err
		}
	}
	ctx.EventManager().EmitEvent(event)
	return nil
}
