package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// erc20 events
const (
	EventTypeConvertCoin             = "convert_coin"
	EventTypeConvertERC20            = "convert_erc20"
	EventTypeConvertDenom            = "convert_denom"
	EventTypeConvertDenomToOne       = "convert_denom_to_one"
	EventTypeConvertDenomToMany      = "convert_denom_to_many"
	EventTypeRegisterCoin            = "register_coin"
	EventTypeRegisterERC20           = "register_erc20"
	EventTypeToggleTokenRelay        = "toggle_token_relay"
	EventTypeRelayToken              = "relay_token"
	EventTypeRelayTransferCrossChain = "relay_transfer_cross_chain"
	EventUpdateContractCode          = "update_contract_code"

	AttributeKeyDenom        = "coin"
	AttributeKeyTokenAddress = "token_address"
	AttributeKeyReceiver     = "receiver"
	AttributeKeyTo           = "to"
	AttributeKeyFrom         = "from"
	AttributeKeyRecipient    = "recipient"
	AttributeKeyTarget       = "target"
	AttributeKeyEvmTxHash    = "evm_tx_hash"
	AttributeKeyTargetDenom  = "target_coin"
	AttributeKeyAlias        = "alias"
	AttributeKeyUpdateFlag   = "update_flag"
	AttributeKeyContract     = "contract"
	AttributeKeyVersion      = "version"

	ERC20EventTransfer = "Transfer"
)

// LogTransfer Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
