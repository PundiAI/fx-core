package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid   = sdkerrors.Register(ModuleName, 2, "invalid")
	ErrEmpty     = sdkerrors.Register(ModuleName, 3, "empty")
	ErrUnknown   = sdkerrors.Register(ModuleName, 4, "unknown")
	ErrDuplicate = sdkerrors.Register(ModuleName, 5, "duplicate")

	ErrInvalidChainName = sdkerrors.Register(ModuleName, 6, "invalid chain name")
	ErrInvalidCoin      = sdkerrors.Register(ModuleName, 7, "invalid coin")

	ErrOracleAddress        = sdkerrors.Register(ModuleName, 10, "invalid oracles address ")
	ErrOrchestratorAddress  = sdkerrors.Register(ModuleName, 11, "invalid orchestrator address")
	ErrExternalAddress      = sdkerrors.Register(ModuleName, 12, "invalid external address")
	ErrTokenContractAddress = sdkerrors.Register(ModuleName, 13, "invalid token contract")

	ErrNonContiguousEventNonce = sdkerrors.Register(ModuleName, 20, "non contiguous event nonce")
	ErrNotOracle               = sdkerrors.Register(ModuleName, 21, "not oracle")
	ErrNoOracleFound           = sdkerrors.Register(ModuleName, 22, "oracle does not exist")
	ErrOracleJailed            = sdkerrors.Register(ModuleName, 23, "oracle for this address is currently jailed")

	ErrBadDepositDenom           = sdkerrors.Register(ModuleName, 24, "invalid coin denomination")
	ErrDepositAmountBelowMinimum = sdkerrors.Register(ModuleName, 25, "deposit amount must be greater than oracle deposit threshold")
	ErrDepositAmountBelowMaximum = sdkerrors.Register(ModuleName, 26, "deposit amount must be less than double oracle deposit threshold")
)
