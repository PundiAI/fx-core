package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidSignature = sdkerrors.Register(ModuleName, 1, "invalid signature")
	ErrInvalidAddress   = sdkerrors.Register(ModuleName, 2, "invalid address")
	ErrAlreadyMigrate   = sdkerrors.Register(ModuleName, 3, "already migrate")
	InvalidRequest      = sdkerrors.Register(ModuleName, 4, "invalid request")
	ErrMigrateValidate  = sdkerrors.Register(ModuleName, 5, "migrate validate error")
	ErrMigrateExecute   = sdkerrors.Register(ModuleName, 6, "migrate execute error")
)
