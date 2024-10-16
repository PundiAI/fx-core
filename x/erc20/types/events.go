package types

const (
	EventTypeConvertCoin      = "convert_coin"
	EventTypeConvertERC20     = "convert_erc20"
	EventTypeRegisterCoin     = "register_coin"
	EventTypeRegisterERC20    = "register_erc20"
	EventTypeToggleTokenRelay = "toggle_token_relay" // #nosec G101

	AttributeKeyDenom        = "coin"
	AttributeKeyTokenAddress = "token_address"
	AttributeKeyReceiver     = "receiver"
)
