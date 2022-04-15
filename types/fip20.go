package types

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	FIP20EventTransferCrosschain = "TransferCrosschain"
)

const (
	FIP20TransferToChainPrefix = "chain/"
	FIP20TransferToIBCPrefix   = "ibc/"
)

const (
	FIP20TargetUnknown Fip20TargetType = iota
	FIP20TargetChain
	FIP20TargetIBC
)

type Fip20TargetType int

func (tt Fip20TargetType) String() string {
	switch tt {
	case FIP20TargetChain:
		return "chain"
	case FIP20TargetIBC:
		return "ibc"
	default:
		return "unknown"
	}
}

type TransferCrosschainEvent struct {
	From   common.Address
	To     string
	Amount *big.Int
	Fee    *big.Int
	Target [32]byte
}

func ParseTransferCrosschainEvent(fip20ABI abi.ABI, log *ethtypes.Log) (*TransferCrosschainEvent, bool, error) {
	if len(log.Topics) < 2 {
		return nil, false, nil
	}
	tc := new(TransferCrosschainEvent)
	if log.Topics[0] != fip20ABI.Events[FIP20EventTransferCrosschain].ID {
		return nil, false, nil
	}
	if len(log.Data) > 0 {
		if err := fip20ABI.UnpackIntoInterface(tc, FIP20EventTransferCrosschain, log.Data); err != nil {
			return nil, false, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fip20ABI.Events[FIP20EventTransferCrosschain].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(tc, indexed, log.Topics[1:]); err != nil {
		return nil, false, err
	}
	return tc, true, nil
}

func (event *TransferCrosschainEvent) GetFrom() sdk.AccAddress {
	return event.From.Bytes()
}

func (event *TransferCrosschainEvent) GetAmount(denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewIntFromBigInt(event.Amount))
}

func (event *TransferCrosschainEvent) GetFee(denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewIntFromBigInt(event.Fee))
}

func (event *TransferCrosschainEvent) GetTarget() (Fip20TargetType, string) {
	if strings.HasPrefix(string(event.Target[:]), FIP20TransferToChainPrefix) {
		return FIP20TargetChain, strings.TrimPrefix(string(event.Target[:]), FIP20TransferToChainPrefix)
	}
	if strings.HasPrefix(string(event.Target[:]), FIP20TransferToIBCPrefix) {
		return FIP20TargetIBC, strings.TrimPrefix(string(event.Target[:]), FIP20TransferToIBCPrefix)
	}
	return FIP20TargetUnknown, string(event.Target[:])
}

func (event *TransferCrosschainEvent) TotalAmount(denom string) sdk.Coins {
	return sdk.NewCoins(event.GetAmount(denom)).Add(event.GetFee(denom))
}
