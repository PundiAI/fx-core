package crosschain

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
)

const (
	FIP20CrossChainGas      = 40_000 // 80000 - 160000
	CrossChainGas           = 40_000 // 70000 - 155000
	CancelSendToExternalGas = 30_000 // 70000 - 126000
	IncreaseBridgeFeeGas    = 40_000 // 70000 - 140000
	BridgeCoinAmountFeeGas  = 10_000 // 9000
	BridgeCallFeeGas        = 50_000 // 50000

	FIP20CrossChainMethodName      = "fip20CrossChain"
	CrossChainMethodName           = "crossChain"
	CancelSendToExternalMethodName = "cancelSendToExternal"
	IncreaseBridgeFeeMethodName    = "increaseBridgeFee"
	BridgeCoinAmountMethodName     = "bridgeCoinAmount"
	BridgeCallMethodName           = "bridgeCall"

	CrossChainEventName           = "CrossChain"
	CancelSendToExternalEventName = "CancelSendToExternal"
	IncreaseBridgeFeeEventName    = "IncreaseBridgeFee"
	BridgeCallEventName           = "BridgeCallEvent"
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
	crossChainAddress = common.HexToAddress(contract.CrossChainAddress)
	crossChainABI     = contract.MustABIJson(contract.ICrossChainMetaData.ABI)
)

func GetAddress() common.Address {
	return crossChainAddress
}

func GetABI() abi.ABI {
	return crossChainABI
}
