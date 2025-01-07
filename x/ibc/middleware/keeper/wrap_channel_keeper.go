package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channelkeeper "github.com/cosmos/ibc-go/v8/modules/core/04-channel/keeper"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

var _ porttypes.ICS4Wrapper = (*WarpChannelKeeper)(nil)

type WarpChannelKeeper struct {
	channelkeeper.Keeper
}

func NewWarpChannelKeeper(k channelkeeper.Keeper) WarpChannelKeeper {
	return WarpChannelKeeper{
		Keeper: k,
	}
}

func (k WarpChannelKeeper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	var packetData transfertypes.FungibleTokenPacketData
	if err = transfertypes.ModuleCdc.UnmarshalJSON(data, &packetData); err != nil {
		return k.Keeper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	if fxtypes.IsPundixChannel(sourcePort, sourceChannel) && packetData.Denom == fxtypes.PundixWrapDenom {
		packetData.Denom = fxtypes.GetPundixUnWrapDenom(ctx.ChainID())
	}

	return k.Keeper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetData.GetBytes())
}
