package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	fxtelemetry "github.com/functionx/fx-core/v8/telemetry"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) SendToFxExecuted(ctx sdk.Context, claim *types.MsgSendToFxClaim) error {
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

	baseCoin, err := k.BridgeTokenToBaseCoin(ctx, claim.TokenContract, claim.Amount, receiveAddr)
	if err != nil {
		return err
	}

	fxTarget := fxtypes.ParseFxTarget(claim.TargetIbc, true)
	if fxTarget.IsIBC() {
		return k.transferIBCHandler(ctx, claim.EventNonce, receiveAddr, baseCoin, fxTarget)
	}

	if fxTarget.GetTarget() == fxtypes.ERC20Target {
		_, err = k.BaseCoinToEvm(ctx, baseCoin, common.BytesToAddress(receiveAddr.Bytes()))
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeEvmTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
			sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
		))
	}
	return nil
}

func (k Keeper) transferIBCHandler(ctx sdk.Context, eventNonce uint64, receive sdk.AccAddress, coin sdk.Coin, target fxtypes.FxTarget) error {
	ibcCoin, err := k.BaseCoinToIBCCoin(ctx, coin, receive, target.String())
	if err != nil {
		return err
	}

	ibcReceiveAddress, err := target.ReceiveAddrToStr(receive)
	if err != nil {
		return err
	}

	// Note: Height is fixed for 5 seconds
	ibcTransferTimeoutHeight := k.GetIbcTransferTimeoutHeight(ctx) * 5
	ibcTimeoutTime := ctx.BlockTime().Add(time.Second * time.Duration(ibcTransferTimeoutHeight))

	response, err := k.ibcTransferKeeper.Transfer(ctx,
		transfertypes.NewMsgTransfer(
			target.SourcePort,
			target.SourceChannel,
			ibcCoin,
			receive.String(),
			ibcReceiveAddress,
			ibcclienttypes.ZeroHeight(),
			uint64(ibcTimeoutTime.UnixNano()),
			"",
		),
	)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIbcTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(response.Sequence)),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, target.SourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, target.SourceChannel),
	))
	return err
}
