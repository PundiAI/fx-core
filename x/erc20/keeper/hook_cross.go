package keeper

import (
	"fmt"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/contracts"
	"github.com/functionx/fx-core/x/erc20/types"
)

func (k Keeper) RelayTransferCrossProcessing(ctx sdk.Context, from common.Address, to *common.Address, receipt *ethtypes.Receipt) (err error) {
	fip20ABI := contracts.GetERC20(ctx.BlockHeight()).ABI
	for _, log := range receipt.Logs {
		if !contracts.VerifyTransferCrossEvent(fip20ABI, log) {
			continue
		}
		pair, found := k.GetTokenPairByAddress(ctx, log.Address)
		if !found {
			continue
		}
		tc, err := contracts.LogToTransferCross(fip20ABI, log, pair.Denom)
		if err != nil {
			return err
		}

		k.Logger(ctx).Info("transfer cross", "tx-hash", receipt.TxHash.Hex(), "from", from.Hex(), "to", to.Hex(), "token", pair.Erc20Address, "denom", pair.Denom)

		balances := k.bankKeeper.GetAllBalances(ctx, tc.From.Bytes())
		if !balances.IsAllGTE(tc.TotalAmount()) {
			return fmt.Errorf("insufficient balance, have %s expected %s", balances.String(), tc.TotalAmount().String())
		}

		switch tc.Type {
		case contracts.TargetChain:
			err = k.TransferChainHandler(ctx, tc, receipt)
		case contracts.TargetIBC:
			err = k.TransferIBCHandler(ctx, tc, receipt)
		default:
			err = fmt.Errorf("traget unknown %d", tc.Type)
		}

		if err != nil {
			k.Logger(ctx).Error("failed to transfer cross", "tx-hash", receipt.TxHash.Hex(), "error", err.Error())
			return err
		}
		k.Logger(ctx).Info("transfer cross success", "tx-hash", receipt.TxHash.Hex())

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "relay_transfer_cross"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", pair.Erc20Address),
				telemetry.NewLabel("denom", pair.Denom),
				telemetry.NewLabel("type", tc.Type.String()),
				telemetry.NewLabel("target", tc.Target),
				telemetry.NewLabel("amount", tc.Amount.Amount.String()),
			},
		)
	}
	return nil
}

func (k Keeper) TransferChainHandler(ctx sdk.Context, tc *contracts.TransferCross, _ *ethtypes.Receipt) error {
	k.Logger(ctx).Info("transfer chain handler", "from", tc.From.Hex(), "to", tc.To, "amount", tc.Amount.String(), "fee", tc.Fee.String(), "target", tc.Target)
	router := k.GetRouter()
	if router == nil || !router.HasRoute(tc.Target) {
		return fmt.Errorf("target %s not support", tc.Target)
	}
	route, _ := router.GetRoute(tc.Target)
	return route.TransferAfter(ctx, sdk.AccAddress(tc.From.Bytes()).String(), tc.To, tc.Amount, tc.Fee)
}

func (k Keeper) TransferIBCHandler(ctx sdk.Context, tc *contracts.TransferCross, receipt *ethtypes.Receipt) error {
	k.Logger(ctx).Info("transfer ibc handler", "from", tc.From.Hex(), "to", tc.To, "amount", tc.Amount.String(), "fee", tc.Fee.String(), "target", tc.Target)
	ibcPrefix, sourcePort, sourceChannel, ok := covertIBCData(tc.Target)
	if !ok {
		return fmt.Errorf("invalid target ibc %s", tc.Target)
	}
	if _, err := sdk.GetFromBech32(tc.To, ibcPrefix); err != nil {
		return fmt.Errorf("invalid to address %s", tc.To)
	}
	_, _, err := k.ibcChannelKeeper.GetChannelClientState(ctx, sourcePort, sourceChannel)
	if err != nil {
		return err
	}
	params := k.GetParams(ctx)
	ibcTimeoutHeight := ibcclienttypes.ZeroHeight()
	ibcTimeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + uint64(params.IbcTimeout)
	nextSequenceSend, found := k.ibcChannelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return fmt.Errorf("ibc channel next sequence send not found, port %s, channel %s", sourcePort, sourceChannel)
	}
	ctx.Logger().Info("ibc transfer", "port", sourcePort, "channel", sourceChannel, "sequence", nextSequenceSend, "timeout-height", ibcTimeoutHeight)
	if err := k.ibcTransferKeeper.SendTransfer(
		ctx, sourcePort, sourceChannel, tc.Amount, tc.From.Bytes(),
		tc.To, ibcTimeoutHeight, ibcTimeoutTimestamp, "", tc.Fee); err != nil {
		return err
	}
	k.setIBCTransferHash(ctx, sourcePort, sourceChannel, nextSequenceSend, receipt.TxHash)
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
