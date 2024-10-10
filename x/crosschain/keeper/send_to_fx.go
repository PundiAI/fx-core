package keeper

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/functionx/fx-core/v8/contract"
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

	baseCoin, err := k.BridgeTokenToBaseCoin(ctx, claim.TokenContract, claim.Amount.BigInt(), receiveAddr)
	if err != nil {
		return err
	}

	return k.RelayTransferHandler(ctx, claim.EventNonce, claim.TargetIbc, receiveAddr, baseCoin)
}

func (k Keeper) RelayTransferHandler(ctx sdk.Context, eventNonce uint64, targetHex string, receiver sdk.AccAddress, coin sdk.Coin) error {
	// ignore hex decode error
	targetByte, _ := hex.DecodeString(targetHex)
	fxTarget := fxtypes.ParseFxTarget(string(targetByte))

	if fxTarget.IsIBC() {
		// transfer to ibc
		// todo convert to ibc token
		return k.transferIBCHandler(ctx, eventNonce, receiver, coin, fxTarget)
	}

	if fxTarget.GetTarget() == fxtypes.ERC20Target {
		// transfer to evm
		if err := k.erc20Keeper.TransferAfter(ctx, receiver, common.BytesToAddress(receiver.Bytes()).String(), coin, sdk.NewCoin(coin.Denom, sdkmath.ZeroInt()), false); err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeEvmTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
			sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
		))
	}
	return nil
}

func (k Keeper) transferIBCHandler(ctx sdk.Context, eventNonce uint64, receive sdk.AccAddress, coin sdk.Coin, target fxtypes.FxTarget) error {
	var ibcReceiveAddress string
	if strings.ToLower(target.Prefix) == contract.EthereumAddressPrefix {
		ibcReceiveAddress = common.BytesToAddress(receive.Bytes()).String()
	} else {
		var err error
		ibcReceiveAddress, err = bech32.ConvertAndEncode(target.Prefix, receive)
		if err != nil {
			return err
		}
	}

	ibcCoin, err := k.BaseCoinToIBCCoin(ctx, coin, receive, target.String())
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
