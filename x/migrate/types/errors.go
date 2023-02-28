package types

import errorsmod "cosmossdk.io/errors"

// todo: remove unused code

var (
	ErrInvalidSignature = errorsmod.Register(ModuleName, 2, "invalid signature")
	ErrInvalidAddress   = errorsmod.Register(ModuleName, 3, "invalid address")
	ErrAlreadyMigrate   = errorsmod.Register(ModuleName, 4, "already migrate")
	ErrMigrateValidate  = errorsmod.Register(ModuleName, 6, "migrate validate error")
	ErrMigrateExecute   = errorsmod.Register(ModuleName, 7, "migrate execute error")
	ErrSameAccount      = errorsmod.Register(ModuleName, 8, "same account")
	ErrInvalidPublicKey = errorsmod.Register(ModuleName, 9, "invalid public key")

	// InvalidRequest      = errorsmod.Register(ModuleName, 5, "invalid request")
)
