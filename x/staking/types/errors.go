package types

import (
	errorsmod "cosmossdk.io/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	ErrUnexpectedEvent    = errorsmod.Register(stakingtypes.ModuleName, 10001, "unexpected event")
	ErrTinyTransferAmount = errorsmod.Register(stakingtypes.ModuleName, 10002, "too few tokens to transfer (truncates to zero tokens)")
	ErrLPTokenNotFound    = errorsmod.Register(stakingtypes.ModuleName, 10003, "lp token not found")
)
