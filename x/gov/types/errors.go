package types

import (
	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var ErrInitialAmountTooLow = errorsmod.Register(govtypes.ModuleName, 116, "initial amount too low")
