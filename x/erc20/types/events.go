package types

const (
	EventTypeConvertCoin             = "convert_coin"
	EventTypeConvertERC20            = "convert_erc20"
	EventTypeConvertDenom            = "convert_denom"
	EventTypeRegisterCoin            = "register_coin"
	EventTypeRegisterERC20           = "register_erc20"
	EventTypeToggleTokenRelay        = "toggle_token_relay" //nolint:gosec #nosec G101
	EventTypeERC20Processing         = "erc20_processing"
	EventTypeRelayTransfer           = "relay_transfer"
	EventTypeRelayTransferCrossChain = "relay_transfer_cross_chain"

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
)
