package contracts

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
)

type TargetType int

const (
	TargetUnknown TargetType = iota
	TargetChain
	TargetIBC
)

const (
	TransferChainPrefix = "chain/"
	TransferIBCPrefix   = "ibc/"
)

const (
	FIP20EventTransfer      = "Transfer"
	FIP20EventTransferCross = "TransferCross"
	FIP20EventTransferChain = "TransferChain"
	FIP20EventTransferIBC   = "TransferIBC"
)

type TransferCrossEvent struct {
	To     string
	Value  *big.Int
	Fee    *big.Int
	Target string
}

func VerifyTransferCrossEvent(log *ethtypes.Log) bool {
	if len(log.Topics) < 2 {
		return false
	}
	eventID := log.Topics[0]
	event, err := FIP20Contract.ABI.EventByID(eventID)
	if err != nil {
		return false
	}
	if !(event.Name == FIP20EventTransferCross) {
		return false
	}
	return true
}

func ParseTransferCrossData(data []byte) (*TransferCrossEvent, error) {
	event := new(TransferCrossEvent)
	err := FIP20Contract.ABI.UnpackIntoInterface(event, FIP20EventTransferCross, data)
	return event, err
}

func ParseTransferCrossTarget(t string) (TargetType, string) {
	if strings.HasPrefix(t, TransferChainPrefix) {
		return TargetChain, strings.TrimPrefix(t, TransferChainPrefix)
	}
	if strings.HasPrefix(t, TransferIBCPrefix) {
		return TargetIBC, strings.TrimPrefix(t, TransferIBCPrefix)
	}
	return TargetUnknown, t
}

type TransferCross struct {
	Type   TargetType
	From   common.Address
	To     string
	Amount sdk.Coin
	Fee    sdk.Coin
	Target string
}

func (tc *TransferCross) TotalAmount() sdk.Coins {
	return sdk.NewCoins(tc.Amount).Add(tc.Fee)
}

func LogToTransferCross(log *ethtypes.Log, denom string) (*TransferCross, error) {
	event, err := ParseTransferCrossData(log.Data)
	if err != nil {
		return nil, fmt.Errorf("parse transfer cross event error %v", err)
	}
	from := common.BytesToAddress(log.Topics[1].Bytes())
	amount := sdk.NewCoin(denom, sdk.NewIntFromBigInt(event.Value))
	fee := sdk.NewCoin(denom, sdk.NewIntFromBigInt(event.Fee))
	targetType, targetProcessed := ParseTransferCrossTarget(event.Target)
	return &TransferCross{
		Type:   targetType,
		From:   from,
		To:     event.To,
		Amount: amount,
		Fee:    fee,
		Target: targetProcessed,
	}, nil
}
