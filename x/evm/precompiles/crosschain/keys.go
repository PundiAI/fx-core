package crosschain

import (
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
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

	JsonABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"string","name":"chain","type":"string"},{"indexed":false,"internalType":"uint256","name":"txid","type":"uint256"}],"name":"CancelSendToExternal","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":true,"internalType":"address","name":"token","type":"address"},{"indexed":false,"internalType":"string","name":"denom","type":"string"},{"indexed":false,"internalType":"string","name":"receipt","type":"string"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"},{"indexed":false,"internalType":"bytes32","name":"target","type":"bytes32"},{"indexed":false,"internalType":"string","name":"memo","type":"string"}],"name":"CrossChain","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":true,"internalType":"address","name":"token","type":"address"},{"indexed":false,"internalType":"string","name":"chain","type":"string"},{"indexed":false,"internalType":"uint256","name":"txid","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"}],"name":"IncreaseBridgeFee","type":"event"},{"type":"function","name":"fip20CrossChain","inputs":[{"name":"sender","type":"address"},{"name":"receipt","type":"string"},{"name":"amount","type":"uint256"},{"name":"fee","type":"uint256"},{"name":"target","type":"bytes32"},{"name":"memo","type":"string"}],"outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"crossChain","inputs":[{"name":"token","type":"address"},{"name":"receipt","type":"string"},{"name":"amount","type":"uint256"},{"name":"fee","type":"uint256"},{"name":"target","type":"bytes32"},{"name":"memo","type":"string"}],"outputs":[{"name":"result","type":"bool"}],"payable":true,"stateMutability":"payable"},{"type":"function","name":"cancelSendToExternal","inputs":[{"name":"chain","type":"string"},{"name":"txid","type":"uint256"}],"outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"increaseBridgeFee","inputs":[{"name":"chain","type":"string"},{"name":"txid","type":"uint256"},{"name":"token","type":"address"},{"name":"fee","type":"uint256"}],"outputs":[{"name":"result","type":"bool"}],"payable":true,"stateMutability":"payable"}]`
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

var precompileAddress = common.HexToAddress(fxtypes.CrossChainAddress)

func GetPrecompileAddress() common.Address {
	return precompileAddress
}
