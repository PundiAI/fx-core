package crosschain

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v4/types"
)

const (
	// FIP20CrossChainGas default send to external fee 0.3FX
	FIP20CrossChainGas = 200000 // if set gas price 500Gwei, about use token 0.1FX
	CrossChainGas      = 200000
	// CancelSendToExternalGas default cancel send to external fee 0.7FX
	CancelSendToExternalGas = 400000 // if set gas price 500Gwei, about use token 0.2FX
	IncreaseBridgeFeeGas    = 400000

	FIP20CrossChainMethodName      = "fip20CrossChain"
	CrossChainMethodName           = "crossChain"
	CancelSendToExternalMethodName = "cancelSendToExternal"
	IncreaseBridgeFeeMethodName    = "increaseBridgeFee"

	CrossChainEventName           = "CrossChain"
	CancelSendToExternalEventName = "CancelSendToExternal"
	IncreaseBridgeFeeEventName    = "IncreaseBridgeFee"
)

const (
	// EventTypeRelayTransferCrossChain
	// Deprecated
	EventTypeRelayTransferCrossChain = "relay_transfer_cross_chain"
	// EventTypeCrossChain new cross chain event type
	EventTypeCrossChain = "cross_chain"

	AttributeKeyDenom        = "coin"
	AttributeKeyTokenAddress = "token_address"
	AttributeKeyFrom         = "from"
	AttributeKeyRecipient    = "recipient"
	AttributeKeyTarget       = "target"
	AttributeKeyMemo         = "memo"
)

var (
	crossChainAddress = common.HexToAddress(fxtypes.CrossChainAddress)
	crossChainABI     = fxtypes.MustABIJson(CrosschainMetaData.ABI)
)

func GetAddress() common.Address {
	return crossChainAddress
}

func GetABI() abi.ABI {
	return crossChainABI
}
