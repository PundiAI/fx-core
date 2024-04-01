package crosschain

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/functionx/fx-core/v7/contract"
)

var (
	CancelSendToExternalEvent = abi.NewEvent(
		CancelSendToExternalEventName,
		CancelSendToExternalEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "chain", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "txID", Type: contract.TypeUint256, Indexed: false},
		})

	CrossChainEvent = abi.NewEvent(
		CrossChainEventName,
		CrossChainEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "token", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "denom", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "receipt", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "amount", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "fee", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "target", Type: contract.TypeBytes32, Indexed: false},
			abi.Argument{Name: "memo", Type: contract.TypeString, Indexed: false},
		})

	IncreaseBridgeFeeEvent = abi.NewEvent(
		IncreaseBridgeFeeEventName,
		IncreaseBridgeFeeEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "token", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "chain", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "txID", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "fee", Type: contract.TypeUint256, Indexed: false},
		})

	BridgeCallEvent = abi.NewEvent(
		BridgeCallEventName,
		BridgeCallEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "receiver", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "to", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "eventNonce", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "dstChainId", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "gasLimit", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "value", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "message", Type: contract.TypeBytes, Indexed: false},
			abi.Argument{Name: "asset", Type: contract.TypeBytes, Indexed: false},
		},
	)
)
