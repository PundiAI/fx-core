package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "migrate"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

const (
	MigrateAccountSignaturePrefix = "MigrateAccount:"

	MigratedRecordKey = "migrated_record"

	EventTypeMigrate = "migrate"
	AttributeKeyFrom = "from"
	AttributeKeyTo   = "to"
)

// GetMigratedRecordKey returns the following key format
func GetMigratedRecordKey(addr sdk.AccAddress) []byte {
	return append([]byte(MigratedRecordKey), addr.Bytes()...)
}
