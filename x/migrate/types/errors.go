package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidAddress   = errorsmod.Register(ModuleName, 3, "invalid address")
	ErrAlreadyMigrate   = errorsmod.Register(ModuleName, 4, "already migrate")
	ErrMigrateValidate  = errorsmod.Register(ModuleName, 6, "migrate validate error")
	ErrMigrateExecute   = errorsmod.Register(ModuleName, 7, "migrate execute error")
	ErrInvalidPublicKey = errorsmod.Register(ModuleName, 9, "invalid public key")

	// ErrInvalidSignature = errorsmod.Register(ModuleName, 2, "invalid signature")
	// InvalidRequest      = errorsmod.Register(ModuleName, 5, "invalid request")
	// ErrSameAccount      = errorsmod.Register(ModuleName, 8, "same account")
)
