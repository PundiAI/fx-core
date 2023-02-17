package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var ErrUnexpectedEvent = sdkerrors.Register(stakingtypes.ModuleName, 10001, "unexpected event")
