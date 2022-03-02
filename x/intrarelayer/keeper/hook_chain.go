package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
	trontypes "github.com/functionx/fx-core/x/tron/types"
	"math/big"
)

const (
	TransferChainETH = "eth"
)

func (k Keeper) RelayTransferChainProcessing(ctx sdk.Context, txHash common.Hash, logs []*ethtypes.Log) (err error) {
	//TODO check height support relay
	for _, log := range logs {
		if !isTransferChainEvent(log) {
			continue
		}
		pair, found := k.GetTokenPairByAddress(ctx, log.Address)
		if !found {
			continue
		}
		event, err := parseTransferChainData(log.Data)
		if err != nil {
			return fmt.Errorf("parse transfer chain event error %v", err)
		}
		from := common.BytesToAddress(log.Topics[1].Bytes())
		k.Logger(ctx).Info("relay transfer chain", "hash", txHash.Hex(), "from", from.Hex(), "to", event.To, "target",
			event.Target, "amount", event.Value.String(), "fee", event.Fee, "denom", pair.Denom, "token", pair.Fip20Address)
		//check balance
		balances := k.bankKeeper.GetAllBalances(ctx, from.Bytes())
		totalAmount := big.NewInt(0).Add(event.Value, event.Fee)
		if balances.AmountOf(pair.Denom).BigInt().Cmp(totalAmount) < 0 {
			return errors.New("insufficient balance")
		}
		//transfer chain
		err = k.transferChainHandler(ctx, event.Target, from, event.To,
			sdk.NewCoin(pair.Denom, sdk.NewIntFromBigInt(event.Value)),
			sdk.NewCoin(pair.Denom, sdk.NewIntFromBigInt(event.Fee)))
		if err != nil {
			k.Logger(ctx).Error("failed relay transfer chain", "hash", txHash.Hex(), "error", err.Error())
			return err
		}
		k.Logger(ctx).Info("relay transfer chain success", "hash", txHash.Hex())
	}
	return nil
}

func isTransferChainEvent(log *ethtypes.Log) bool {
	if len(log.Topics) < 2 {
		return false
	}
	eventID := log.Topics[0] // event ID
	event, err := contracts.FIP20Contract.ABI.EventByID(eventID)
	if err != nil {
		return false
	}
	if !(event.Name == types.FIP20EventTransferChain) {
		return false
	}
	return true
}

type TransferChainData struct {
	To     string
	Value  *big.Int
	Fee    *big.Int
	Target string
}

func parseTransferChainData(data []byte) (*TransferChainData, error) {
	event := new(TransferChainData)
	err := contracts.FIP20Contract.ABI.UnpackIntoInterface(event, types.FIP20EventTransferChain, data)
	return event, err
}

func (k Keeper) transferChainHandler(ctx sdk.Context, chain string, from common.Address, to string, amount, fee sdk.Coin) error {
	if !k.ChainSupport(chain) {
		return fmt.Errorf("chain %s not support", chain)
	}
	if !k.ValidateToAddress(chain, to) {
		return fmt.Errorf("invalid address %s", to)
	}
	return k.AddToOutgoingPool(ctx, chain, from, to, amount, fee)
}

func (k Keeper) ChainSupport(chain string) bool {
	if chain == TransferChainETH {
		return true
	}
	return k.crossChainKeepers[chain] != nil
}

func (k Keeper) ValidateToAddress(chain, to string) bool {
	if chain == trontypes.ModuleName {
		return trontypes.ValidateExternalAddress(to) == nil
	}
	return crosschaintypes.ValidateExternalAddress(to) == nil
}

func (k Keeper) AddToOutgoingPool(ctx sdk.Context, chain string, from common.Address, to string, amount, fee sdk.Coin) (err error) {
	var poolId uint64
	//eth
	if chain == TransferChainETH {
		poolId, err = k.gravityKeeper.AddToOutgoingPool(ctx, from.Bytes(), to, amount, fee)
	} else {
		//cross chain
		crossChainKeeper := k.crossChainKeepers[chain]
		poolId, err = crossChainKeeper.AddToOutgoingPool(ctx, from.Bytes(), to, amount, fee)
	}
	if err != nil {
		return err
	}
	//TODO save pool id
	_ = poolId
	return nil
}
