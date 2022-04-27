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
	prefixMigratedRecord = iota + 1

	MigrateAccountSignaturePrefix = "MigrateAccount:"

	EventTypeMigrate = "migrate"
	AttributeKeyFrom = "from"
	AttributeKeyTo   = "to"
)

var (
	PrefixMigrateFromFlag = []byte{0x1}
	PrefixMigrateToFlag   = []byte{0x2}
)

// GetMigratedRecordKey returns the following key format
func GetMigratedRecordKey(addr sdk.AccAddress) []byte {
	return append([]byte{prefixMigratedRecord}, addr.Bytes()...)
}
