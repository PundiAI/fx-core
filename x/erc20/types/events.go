package types

const (
	EventTypeConvertCoin      = "convert_coin"
	EventTypeRegisterCoin     = "register_coin"
	EventTypeRegisterERC20    = "register_erc20"
	EventTypeToggleTokenRelay = "toggle_token_relay" // #nosec G101

	AttributeKeyDenom        = "coin"
	AttributeKeyTokenAddress = "token_address"
	AttributeKeyReceiver     = "receiver"
)
