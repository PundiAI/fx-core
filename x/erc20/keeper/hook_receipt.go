package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

type RelayTokenEventLog struct {
	Event *RelayTokenEvent
	Log   *ethtypes.Log
	Pair  *types.TokenPair
}

type TransferCrossChainEventLog struct {
	Event *fxtypes.TransferCrossChainEvent
	Log   *ethtypes.Log
	Pair  *types.TokenPair
}

type EventLog struct {
	RelayToken         []*RelayTokenEventLog
	TransferCrossChain []*TransferCrossChainEventLog
}

func ParseEventLog(receipt *ethtypes.Receipt) (EventLog, bool) {
	fip20ABI := fxtypes.GetERC20().ABI

	relayTokenEvents := make([]*RelayTokenEventLog, 0, len(receipt.Logs))
	transferCrossChainEvents := make([]*TransferCrossChainEventLog, 0, len(receipt.Logs))

	for _, log := range receipt.Logs {
		rt, rtOk, rtErr := ParseRelayTokenEvent(fip20ABI, log)
		tc, tcOk, tcErr := fxtypes.ParseTransferCrossChainEvent(fip20ABI, log)

		if rtErr != nil || tcErr != nil {
			return EventLog{}, false
		}
		if !rtOk && !tcOk {
			continue
		}

		if rtOk {
			relayTokenEvents = append(relayTokenEvents, &RelayTokenEventLog{Event: rt, Log: log})
		}
		if tcOk {
			transferCrossChainEvents = append(transferCrossChainEvents, &TransferCrossChainEventLog{Event: tc, Log: log})
		}
	}

	el := EventLog{RelayToken: relayTokenEvents, TransferCrossChain: transferCrossChainEvents}
	return el, true
}

func (k Keeper) TokenPairEnable(ctx sdk.Context, eventLog EventLog) (EventLog, error) {
	rtels := eventLog.RelayToken
	tcels := eventLog.TransferCrossChain

	addressEnable := make(map[common.Address]*types.TokenPair, len(rtels)+len(tcels))
	addressNotFound := make(map[common.Address]bool, len(rtels)+len(tcels))

	rtelsNew := make([]*RelayTokenEventLog, 0, len(rtels))
	tcelsNew := make([]*TransferCrossChainEventLog, 0, len(tcels))

	for _, rtel := range rtels {
		// contract not found
		if addressNotFound[rtel.Log.Address] {
			continue
		}

		// contract enable
		if pair, ok := addressEnable[rtel.Log.Address]; ok {
			rtel.Pair = pair
			rtelsNew = append(rtelsNew, rtel)
			continue
		}

		// get contract token pair
		pair, found := k.GetTokenPairByAddress(ctx, rtel.Log.Address)
		if !found {
			addressNotFound[rtel.Log.Address] = true
			continue
		}
		if !pair.Enabled {
			return EventLog{}, fmt.Errorf("token pair not enable, contract %s, denom %s", pair.Erc20Address, pair.Denom)
		}
		//record contract token pair
		addressEnable[rtel.Log.Address] = &pair

		rtel.Pair = &pair
		rtelsNew = append(rtelsNew, rtel)
	}

	for _, tcel := range tcels {
		// contract not found
		if addressNotFound[tcel.Log.Address] {
			continue
		}
		// contract enable
		if pair, ok := addressEnable[tcel.Log.Address]; ok {
			tcel.Pair = pair
			tcelsNew = append(tcelsNew, tcel)
			continue
		}

		// get contract token pair
		pair, found := k.GetTokenPairByAddress(ctx, tcel.Log.Address)
		if !found {
			addressNotFound[tcel.Log.Address] = true
			continue
		}
		if !pair.Enabled {
			return EventLog{}, fmt.Errorf("token pair not enable, contract %s, denom %s", pair.Erc20Address, pair.Denom)
		}
		// record contract token pair
		addressEnable[tcel.Log.Address] = &pair

		tcel.Pair = &pair
		tcelsNew = append(tcelsNew, tcel)
	}

	events := EventLog{RelayToken: rtelsNew, TransferCrossChain: tcelsNew}
	return events, nil
}

func (k Keeper) GetTokenPairByAddress(ctx sdk.Context, address common.Address) (types.TokenPair, bool) {
	//check contract is registered
	pairID := k.GetERC20Map(ctx, address)
	if len(pairID) == 0 {
		// contract is not registered coin or fip20
		return types.TokenPair{}, false
	}
	return k.GetTokenPair(ctx, pairID)
}
