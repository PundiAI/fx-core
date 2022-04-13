package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// erc20 events
const (
	EventTypeTokenLock             = "token_lock"
	EventTypeTokenUnlock           = "token_unlock"
	EventTypeMint                  = "mint"
	EventTypeConvertCoin           = "convert_coin"
	EventTypeConvertERC20          = "convert_erc20"
	EventTypeBurn                  = "burn"
	EventTypeRegisterCoin          = "register_coin"
	EventTypeRegisterERC20         = "register_erc20"
	EventTypeToggleTokenRelay      = "toggle_token_relay" // #nosec
	EventTypeUpgradeSystemContract = "upgrade_system_contract"

	AttributeKeyCosmosCoin      = "cosmos_coin"
	AttributeKeyERC20Token      = "erc20_token" // #nosec
	AttributeKeyReceiver        = "receiver"
	AttributeKeyContractAddress = "contract_address"

	ERC20EventTransfer = "Transfer"

	EventTypeRelayToken = "relay_token"
	EventEthereumTxHash = "ethereum_tx_hash"
)

// Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
