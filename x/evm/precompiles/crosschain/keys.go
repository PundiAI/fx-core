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

	JsonABI = `[{"type":"function","name":"fip20CrossChain","inputs":[{"name":"sender","type":"address"},{"name":"receipt","type":"string"},{"name":"amount","type":"uint256"},{"name":"fee","type":"uint256"},{"name":"target","type":"bytes32"},{"name":"memo","type":"string"}],"outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"crossChain","inputs":[{"name":"token","type":"address"},{"name":"receipt","type":"string"},{"name":"amount","type":"uint256"},{"name":"fee","type":"uint256"},{"name":"target","type":"bytes32"},{"name":"memo","type":"string"}],"outputs":[{"name":"result","type":"bool"}],"payable":true,"stateMutability":"payable"},{"type":"function","name":"cancelSendToExternal","inputs":[{"name":"chain","type":"string"},{"name":"txid","type":"uint256"}],"outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"increaseBridgeFee","inputs":[{"name":"chain","type":"string"},{"name":"txid","type":"uint256"},{"name":"token","type":"address"},{"name":"fee","type":"uint256"}],"outputs":[{"name":"result","type":"bool"}],"payable":true,"stateMutability":"payable"}]`
)

var precompileAddress = common.HexToAddress(fxtypes.CrossChainAddress)

func GetPrecompileAddress() common.Address {
	return precompileAddress
}
