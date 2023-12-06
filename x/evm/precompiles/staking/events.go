package staking

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/functionx/fx-core/v6/x/evm/types"
)

var (
	ApproveSharesEvent = abi.NewEvent(
		ApproveSharesEventName,
		ApproveSharesEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "owner", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "spender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
		},
	)

	DelegateEvent = abi.NewEvent(
		DelegateEventName,
		DelegateEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "delegator", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "amount", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
		},
	)

	TransferSharesEvent = abi.NewEvent(
		TransferSharesEventName,
		TransferSharesEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "from", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "to", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "token", Type: types.TypeUint256, Indexed: false},
		},
	)

	UndelegateEvent = abi.NewEvent(
		UndelegateEventName,
		UndelegateEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "amount", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "completionTime", Type: types.TypeUint256, Indexed: false},
		},
	)

	WithdrawEvent = abi.NewEvent(
		WithdrawEventName,
		WithdrawEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "validator", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "reward", Type: types.TypeUint256, Indexed: false},
		},
	)

	RedelegateEvent = abi.NewEvent(
		RedelegateEventName,
		RedelegateEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "valSrc", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "valDst", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "shares", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "amount", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "completionTime", Type: types.TypeUint256, Indexed: false},
		},
	)
)
