package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	fxtypes "github.com/functionx/fx-core/types"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"

	"github.com/functionx/fx-core/x/crosschain/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

func covertIbcData(targetIbc string) (prefix, sourcePort, sourceChannel string, isOk bool) {
	// fx/transfer/channel-0
	targetIbcBytes, err := hex.DecodeString(targetIbc)
	if err != nil {
		return
	}
	ibcData := strings.Split(string(targetIbcBytes), "/")
	if len(ibcData) < 3 {
		isOk = false
		return
	}
	prefix = ibcData[0]
	sourcePort = ibcData[1]
	sourceChannel = ibcData[2]
	isOk = true
	return
}

type Target int

const (
	TargetUnknown Target = iota
	TargetEvm
	TargetIBC
)

func (k Keeper) handleIbcTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiveAddr sdk.AccAddress, coin sdk.Coin) ([]sdk.Attribute, bool) {
	ibcPrefix, sourcePort, sourceChannel, ok := covertIbcData(claim.TargetIbc)
	if !ok {
		return nil, false
	}
	logger := k.Logger(ctx)
	ibcReceiveAddress, err := bech32.ConvertAndEncode(ibcPrefix, receiveAddr)
	if err != nil {
		logger.Error("convert ibc transfer receive address error!!!", "fxReceive:", claim.Receiver,
			"ibcPrefix:", ibcPrefix, "sourcePort:", sourcePort, "sourceChannel:", sourceChannel, "error:", err)
		return nil, false
	}
	wrapSdkContext := sdk.WrapSDKContext(ctx)

	_, clientState, err := k.ibcChannelKeeper.GetChannelClientState(ctx, sourcePort, sourceChannel)
	if err != nil {
		logger.Error("get channel client state error!!!", "sourcePort", sourcePort, "sourceChannel", sourceChannel)
		return nil, false
	}

	params := k.GetParams(ctx)
	clientStateHeight := clientState.GetLatestHeight()
	ibcTimeoutHeight := ibcclienttypes.Height{
		RevisionNumber: clientStateHeight.GetRevisionNumber(),
		RevisionHeight: clientStateHeight.GetRevisionHeight() + params.IbcTransferTimeoutHeight,
	}

	nextSequenceSend, found := k.ibcChannelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		logger.Error("ibc channel next sequence send not found!!!", "source port:", sourcePort, "source channel:", sourceChannel)
		return nil, false
	}
	logger.Info("gravity start ibc transfer", "sender:", receiveAddr, "receive:", ibcReceiveAddress, "coin:", coin, "timeout:", params.IbcTransferTimeoutHeight, "nextSequenceSend:", nextSequenceSend)

	ibcTransferMsg := ibctransfertypes.NewMsgTransfer(sourcePort, sourceChannel, coin, receiveAddr, ibcReceiveAddress, ibcTimeoutHeight, 0, "", sdk.NewCoin(coin.Denom, sdk.ZeroInt()))
	if _, err = k.ibcTransferKeeper.Transfer(wrapSdkContext, ibcTransferMsg); err != nil {
		logger.Error("gravity ibc transfer fail. ", "sender:", receiveAddr, "receive:", ibcReceiveAddress, "coin:", coin, "err:", err)
		return nil, false
	}
	k.SetIbcSequenceHeight(ctx, sourcePort, sourceChannel, nextSequenceSend, uint64(ctx.BlockHeight()))

	attributes := make([]sdk.Attribute, 0, 3)
	attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyAttestationHandlerIbcChannelSendSequence, fmt.Sprintf("%d", nextSequenceSend)))
	attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyAttestationHandlerIbcChannelSourcePort, sourcePort))
	attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyAttestationHandlerIbcChannelSourceChannel, sourceChannel))
	return attributes, true
}

func (k Keeper) handlerEvmTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) ([]sdk.Attribute, bool) {
	logger := k.Logger(ctx)
	if !k.IntrarelayerKeeper.HasInit(ctx) {
		logger.Error("emv transfer, module not init", "module", "intrarelayer")
		return nil, false
	}
	if !k.IntrarelayerKeeper.IsDenomRegistered(ctx, coin.Denom) {
		logger.Error("evm transfer, denom not registered", "denom", coin.Denom)
		return nil, false
	}
	receiverEthType := common.BytesToAddress(receiver.Bytes())
	logger.Info("convert denom to fip20", "sender", claim.Sender, "receiver", claim.Receiver,
		"receiver-eth-type", receiverEthType.String(), "amount", coin.String(), "target", claim.TargetIbc)
	err := k.IntrarelayerKeeper.ConvertDenomToFIP20(ctx, receiver, receiverEthType, coin)
	if err != nil {
		logger.Error("evm transfer, convert denom to fip20 failed", "error", err.Error())
		return nil, false
	}
	attributes := make([]sdk.Attribute, 0, 1)
	attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyAttestationHandlerEvmTransfer, claim.TargetIbc))
	return nil, true
}

func verifyTarget(targetHex string) Target {
	targetBZ, err := hex.DecodeString(targetHex)
	if err != nil {
		return TargetUnknown
	}

	target := string(targetBZ)
	if strings.HasPrefix(target, "module/evm") {
		return TargetEvm
	}
	if _, _, _, ok := covertIbcData(target); ok {
		return TargetIBC
	}
	return TargetUnknown
}

func (k Keeper) handlerRelayTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) ([]sdk.Attribute, bool) {
	if ctx.BlockHeight() < fxtypes.IntrarelayerSupportBlock() {
		return k.handleIbcTransfer(ctx, claim, receiver, coin)
	} else {
		switch verifyTarget(claim.TargetIbc) {
		case TargetEvm:
			return k.handlerEvmTransfer(ctx, claim, receiver, coin)
		case TargetIBC:
			return k.handleIbcTransfer(ctx, claim, receiver, coin)
		default:
			return nil, false
		}
	}
}
