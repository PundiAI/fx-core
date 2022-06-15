package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// x/gov module sentinel errors
var (
	ErrInitialAmountTooLow = sdkerrors.Register(govtypes.ModuleName, 10, "initial amount too low")
)
