package common

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidSignature = sdkerrors.Register(ModuleName, 2, "invalid signature")
	ErrInvalidAddress   = sdkerrors.Register(ModuleName, 3, "invalid address")
	ErrAlreadyMigrate   = sdkerrors.Register(ModuleName, 4, "already migrate")
	InvalidRequest      = sdkerrors.Register(ModuleName, 5, "invalid request")
	ErrMigrateValidate  = sdkerrors.Register(ModuleName, 6, "migrate validate error")
	ErrMigrateExecute   = sdkerrors.Register(ModuleName, 7, "migrate execute error")
)
