package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// todo: remove unused code

var (
	ErrInvalidSignature = sdkerrors.Register(ModuleName, 2, "invalid signature")
	ErrInvalidAddress   = sdkerrors.Register(ModuleName, 3, "invalid address")
	ErrAlreadyMigrate   = sdkerrors.Register(ModuleName, 4, "already migrate")
	ErrMigrateValidate  = sdkerrors.Register(ModuleName, 6, "migrate validate error")
	ErrMigrateExecute   = sdkerrors.Register(ModuleName, 7, "migrate execute error")
	ErrSameAccount      = sdkerrors.Register(ModuleName, 8, "same account")
	ErrInvalidPublicKey = sdkerrors.Register(ModuleName, 9, "invalid public key")

	// InvalidRequest      = sdkerrors.Register(ModuleName, 5, "invalid request")
)
