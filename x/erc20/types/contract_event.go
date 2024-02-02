package types

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

const (
	FIP20EventTransferCrossChain = "TransferCrossChain"
	ERC20EventTransfer           = "Transfer"
)

type TransferCrossChainEvent struct {
	From      common.Address
	Recipient string
	Amount    *big.Int
	Fee       *big.Int
	Target    [32]byte
}

func ParseTransferCrossChainEvent(log *ethtypes.Log) (*TransferCrossChainEvent, error) {
	if len(log.Topics) != 2 {
		return nil, nil
	}
	fip20ABI := fxtypes.GetFIP20().ABI
	tc := new(TransferCrossChainEvent)
	if log.Topics[0] != fip20ABI.Events[FIP20EventTransferCrossChain].ID {
		return nil, nil
	}
	if len(log.Data) > 0 {
		if err := fip20ABI.UnpackIntoInterface(tc, FIP20EventTransferCrossChain, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fip20ABI.Events[FIP20EventTransferCrossChain].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(tc, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	return tc, nil
}

func (event *TransferCrossChainEvent) GetFrom() sdk.AccAddress {
	return event.From.Bytes()
}

func (event *TransferCrossChainEvent) GetAmount(denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(event.Amount))
}

func (event *TransferCrossChainEvent) GetFee(denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(event.Fee))
}

func (event *TransferCrossChainEvent) TotalAmount(denom string) sdk.Coins {
	return sdk.NewCoins(event.GetAmount(denom)).Add(event.GetFee(denom))
}

type TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// ParseTransferEvent transfer event ---> event Transfer(address indexed from, address indexed to, uint256 value);
func ParseTransferEvent(log *ethtypes.Log) (*TransferEvent, error) {
	// Note: the `Transfer` event contains 3 topics (id, from, to)
	if len(log.Topics) != 3 {
		return nil, nil
	}
	fip20ABI := fxtypes.GetFIP20().ABI
	if log.Topics[0] != fip20ABI.Events[ERC20EventTransfer].ID {
		return nil, nil
	}
	transferEvent := new(TransferEvent)
	if len(log.Data) > 0 {
		if err := fip20ABI.UnpackIntoInterface(transferEvent, ERC20EventTransfer, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fip20ABI.Events[ERC20EventTransfer].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(transferEvent, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	return transferEvent, nil
}
