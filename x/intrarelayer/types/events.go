package types

import (
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// intrarelayer events
const (
	EventTypeTokenLock        = "token_lock"
	EventTypeTokenUnlock      = "token_unlock"
	EventTypeMint             = "mint"
	EventTypeRelay            = "relay"
	EventTypeConvertCoin      = "convert_coin"
	EventTypeConvertFIP20     = "convert_fip20"
	EventTypeBurn             = "burn"
	EventTypeRegisterCoin     = "register_coin"
	EventTypeRegisterFIP20    = "register_fip20"
	EventTypeToggleTokenRelay = "toggle_token_relay" // #nosec

	AttributeKeyCosmosCoin  = "cosmos_coin"
	AttributeKeyFIP20Token  = "fip20_token" // #nosec
	AttributeKeyFIP20Symbol = "fip20_symbol"
	AttributeKeyReceiver    = "receiver"

	FIP20EventTransfer      = contracts.FIP20EventTransfer
	FIP20EventTransferCross = contracts.FIP20EventTransferCross
	FIP20EventTransferChain = contracts.FIP20EventTransferChain
	FIP20EventTransferIBC   = contracts.FIP20EventTransferIBC

	EventTypeRelayToken = "relay_token"
	EventEthereumTxHash = "ethereum_tx_hash"
)

// Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
