package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// erc20 events
const (
	EventTypeConvertCoin      = "convert_coin"
	EventTypeConvertERC20     = "convert_erc20"
	EventTypeRegisterCoin     = "register_coin"
	EventTypeRegisterERC20    = "register_erc20"
	EventTypeToggleTokenRelay = "toggle_token_relay"
	EventTypeRelayToken       = "relay_token"

	AttributeKeyDenom        = "coin"
	AttributeKeyTokenAddress = "token_address"
	AttributeKeyReceiver     = "receiver"
	AttributeKeyEvmTxHash    = "evm_tx_hash"

	ERC20EventTransfer = "Transfer"
)

// LogTransfer Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
