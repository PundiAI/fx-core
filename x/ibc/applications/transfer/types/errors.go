package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC channel sentinel errors
var (
	ErrInvalidPacketTimeout       = sdkerrors.Register(ModuleName, 102, "invalid packet timeout")
	ErrInvalidDenomForTransfer    = sdkerrors.Register(ModuleName, 103, "invalid denomination for cross-chain transfer")
	ErrInvalidVersion             = sdkerrors.Register(ModuleName, 104, "invalid ICS20 version")
	ErrInvalidAmount              = sdkerrors.Register(ModuleName, 105, "invalid token amount")
	ErrTraceNotFound              = sdkerrors.Register(ModuleName, 106, "denomination trace not found")
	ErrSendDisabled               = sdkerrors.Register(ModuleName, 107, "fungible token transfers from this chain are disabled")
	ErrReceiveDisabled            = sdkerrors.Register(ModuleName, 108, "fungible token transfers to this chain are disabled")
	ErrMaxTransferChannels        = sdkerrors.Register(ModuleName, 109, "max transfer channels")
	ErrFeeDenomNotMatchTokenDenom = sdkerrors.Register(ModuleName, 110, "invalid fee denom, must match token denom")
)
