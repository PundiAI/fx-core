package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
	"math/big"
	"strings"
)

func (k Keeper) RelayTransferIBCProcessing(ctx sdk.Context, from common.Address, to *common.Address, receipt *ethtypes.Receipt) error {
	//TODO check height support relay
	for _, log := range receipt.Logs {
		if !isTransferIBCEvent(log) {
			continue
		}
		pair, found := k.GetTokenPairByAddress(ctx, log.Address)
		if !found {
			continue
		}
		event, err := parseTransferIBCEvent(log.Data)
		if err != nil {
			return fmt.Errorf("parse transfer ibc event error %v", err)
		}
		from := common.BytesToAddress(log.Topics[1].Bytes())

		k.Logger(ctx).Info("relay transfer ibc", "hash", receipt.TxHash.Hex(), "from", from.Hex(), "to", event.To,
			"amount", event.Value.String(), "denom", pair.Denom, "token", pair.Fip20Address)

		//check balance
		balances := k.bankKeeper.GetAllBalances(ctx, from.Bytes())
		if balances.AmountOf(pair.Denom).BigInt().Cmp(event.Value) < 0 {
			return errors.New("insufficient balance")
		}
		err = k.transferIBCHandler(ctx, event.Target, from, event.To, sdk.NewCoin(pair.Denom, sdk.NewIntFromBigInt(event.Value)), receipt.TxHash)
		if err != nil {
			k.Logger(ctx).Error("failed to relay transfer ibc", "hash", receipt.TxHash.Hex(), "error", err.Error())
			return err
		}
		k.Logger(ctx).Info("relay transfer ibc success", "hash", receipt.TxHash.Hex())
	}
	return nil
}

func isTransferIBCEvent(log *ethtypes.Log) bool {
	if len(log.Topics) < 2 {
		return false
	}
	eventID := log.Topics[0] // event ID
	event, err := contracts.FIP20Contract.ABI.EventByID(eventID)
	if err != nil {
		return false
	}
	if !(event.Name == types.FIP20EventTransferIBC) {
		return false
	}
	return true
}

type TransferIBCEvent struct {
	To     string
	Value  *big.Int
	Target string
}

func parseTransferIBCEvent(data []byte) (*TransferIBCEvent, error) {
	event := new(TransferIBCEvent)
	err := contracts.FIP20Contract.ABI.UnpackIntoInterface(event, types.FIP20EventTransferIBC, data)
	return event, err
}

func (k Keeper) transferIBCHandler(ctx sdk.Context, targetIBC string, sender common.Address, to string, amount sdk.Coin, txHash common.Hash) error {
	ibcPrefix, sourcePort, sourceChannel, ok := covertIBCData(targetIBC)
	if !ok {
		return fmt.Errorf("invalid target ibc %s", targetIBC)
	}

	if _, err := sdk.GetFromBech32(to, ibcPrefix); err != nil {
		return fmt.Errorf("invalid to address %s", to)
	}

	_, clientState, err := k.ibcChannelKeeper.GetChannelClientState(ctx, sourcePort, sourceChannel)
	if err != nil {
		return err
	}
	params := k.GetParams(ctx)
	clientStateHeight := clientState.GetLatestHeight()
	ibcTimeoutHeight := ibcclienttypes.Height{
		RevisionNumber: clientStateHeight.GetRevisionNumber(),
		RevisionHeight: clientStateHeight.GetRevisionHeight() + params.IbcTransferTimeoutHeight,
	}
	nextSequenceSend, found := k.ibcChannelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return fmt.Errorf("ibc channel next sequence send not found, port %s, channel %s", sourcePort, sourceChannel)
	}
	ctx.Logger().Info("ibc transfer", "port", sourcePort, "channel", sourceChannel, "sender", sender.String(), "receiver", to, "amount", amount.String(), "timeout-height", ibcTimeoutHeight)
	goCtx := sdk.WrapSDKContext(ctx)
	ibcTransferMsg := ibctransfertypes.NewMsgTransfer(sourcePort, sourceChannel, amount, sender.Bytes(), to, ibcTimeoutHeight, 0, "", sdk.NewCoin(amount.Denom, sdk.ZeroInt()))
	if _, err = k.ibcTransferKeeper.Transfer(goCtx, ibcTransferMsg); err != nil {
		return err
	}
	k.setIBCTransferHash(ctx, sourcePort, sourceChannel, nextSequenceSend, txHash)
	return nil
}

func covertIBCData(targetIbc string) (prefix, sourcePort, sourceChannel string, isOk bool) {
	// pay/transfer/channel-0
	ibcData := strings.Split(targetIbc, "/")
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

func (k Keeper) setIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64, hash common.Hash) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIBCTransferKey(port, channel, sequence), hash.Bytes())
}

func (k Keeper) GetIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) (common.Hash, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIBCTransferKey(port, channel, sequence)
	if !store.Has(key) {
		return common.Hash{}, false
	}
	value := store.Get(key)
	return common.BytesToHash(value), true
}

func (k Keeper) HashIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetIBCTransferKey(port, channel, sequence))
}
