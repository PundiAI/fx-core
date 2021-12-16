package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// intrarelayer events
const (
	EventTypeTokenLock            = "token_lock"
	EventTypeTokenUnlock          = "token_unlock"
	EventTypeMint                 = "mint"
	EventTypeRelay                = "relay"
	EventTypeConvertCoin          = "convert_coin"
	EventTypeConvertERC20         = "convert_erc20"
	EventTypeBurn                 = "burn"
	EventTypeRegisterCoin         = "register_coin"
	EventTypeRegisterERC20        = "register_erc20"
	EventTypeToggleTokenRelay     = "toggle_token_relay" // #nosec
	EventTypeUpdateTokenPairERC20 = "update_token_pair_erc20"

	AttributeKeyCosmosCoin = "cosmos_coin"
	AttributeKeyERC20Token = "erc20_token" // #nosec
	AttributeKeyReceiver   = "receiver"

	ERC20EventTransfer   = "Transfer"
	ERC20RelayEventRelay = "Relay"

	EventTypeERC20Relay = "erc20_relay"
	EventERC20RelayHash = "hash"
)

// Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
