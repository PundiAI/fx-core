package staking

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/functionx/fx-core/v7/contract"
)

var (
	ApproveSharesEvent = abi.NewEvent(
		ApproveSharesEventName,
		ApproveSharesEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "owner", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "spender", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: contract.TypeUint256, Indexed: false},
		},
	)

	DelegateEvent = abi.NewEvent(
		DelegateEventName,
		DelegateEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "delegator", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "amount", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "shares", Type: contract.TypeUint256, Indexed: false},
		},
	)

	TransferSharesEvent = abi.NewEvent(
		TransferSharesEventName,
		TransferSharesEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "from", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "to", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "token", Type: contract.TypeUint256, Indexed: false},
		},
	)

	UndelegateEvent = abi.NewEvent(
		UndelegateEventName,
		UndelegateEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "amount", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "completionTime", Type: contract.TypeUint256, Indexed: false},
		},
	)

	WithdrawEvent = abi.NewEvent(
		WithdrawEventName,
		WithdrawEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "reward", Type: contract.TypeUint256, Indexed: false},
		},
	)

	RedelegateEvent = abi.NewEvent(
		RedelegateEventName,
		RedelegateEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: contract.TypeAddress, Indexed: true},
			abi.Argument{Name: "valSrc", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "valDst", Type: contract.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "amount", Type: contract.TypeUint256, Indexed: false},
			abi.Argument{Name: "completionTime", Type: contract.TypeUint256, Indexed: false},
		},
	)
)
