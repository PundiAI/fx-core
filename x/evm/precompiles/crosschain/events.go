package crosschain

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

var (
	CancelSendToExternalEvent = abi.NewEvent(
		CancelSendToExternalEventName,
		CancelSendToExternalEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "chain", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "txID", Type: types.TypeUint256, Indexed: false},
		})

	CrossChainEvent = abi.NewEvent(
		CrossChainEventName,
		CrossChainEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "token", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "denom", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "receipt", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "amount", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "fee", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "target", Type: types.TypeBytes32, Indexed: false},
			abi.Argument{Name: "memo", Type: types.TypeString, Indexed: false},
		})

	IncreaseBridgeFeeEvent = abi.NewEvent(
		IncreaseBridgeFeeEventName,
		IncreaseBridgeFeeEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "token", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "chain", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "txID", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "fee", Type: types.TypeUint256, Indexed: false},
		})
)
