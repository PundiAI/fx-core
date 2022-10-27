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
	FIP20EventTransferCrossChain = "TransferCrossChain"
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

type TransferCrossChainEvent struct {
	From      common.Address
	Recipient string
	Amount    *big.Int
	Fee       *big.Int
	Target    [32]byte
}

func ParseTransferCrossChainEvent(fip20ABI abi.ABI, log *ethtypes.Log) (*TransferCrossChainEvent, bool, error) {
	if len(log.Topics) != 2 {
		return nil, false, nil
	}
	tc := new(TransferCrossChainEvent)
	if log.Topics[0] != fip20ABI.Events[FIP20EventTransferCrossChain].ID {
		return nil, false, nil
	}
	if len(log.Data) > 0 {
		if err := fip20ABI.UnpackIntoInterface(tc, FIP20EventTransferCrossChain, log.Data); err != nil {
			return nil, false, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fip20ABI.Events[FIP20EventTransferCrossChain].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(tc, indexed, log.Topics[1:]); err != nil {
		return nil, false, err
	}
	return tc, true, nil
}

func (event *TransferCrossChainEvent) GetFrom() sdk.AccAddress {
	return event.From.Bytes()
}

func (event *TransferCrossChainEvent) GetAmount(denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewIntFromBigInt(event.Amount))
}

func (event *TransferCrossChainEvent) GetFee(denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewIntFromBigInt(event.Fee))
}

func (event *TransferCrossChainEvent) GetTarget() (Fip20TargetType, string) {
	target := Byte32ToString(event.Target)
	if strings.HasPrefix(target, FIP20TransferToChainPrefix) {
		return FIP20TargetChain, strings.TrimPrefix(target, FIP20TransferToChainPrefix)
	}
	if strings.HasPrefix(target, FIP20TransferToIBCPrefix) {
		return FIP20TargetIBC, strings.TrimPrefix(target, FIP20TransferToIBCPrefix)
	}
	return FIP20TargetUnknown, target
}

func (event *TransferCrossChainEvent) TotalAmount(denom string) sdk.Coins {
	return sdk.NewCoins(event.GetAmount(denom)).Add(event.GetFee(denom))
}

func StringToByte32(str string) [32]byte {
	var bz [32]byte
	copy(bz[:], str)
	return bz
}

func Byte32ToString(bytes [32]byte) string {
	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] != 0 {
			return string(bytes[:i+1])
		}
	}
	return ""
}
