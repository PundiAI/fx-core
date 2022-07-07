package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

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
	prefixMigratedDirectionFrom
	prefixMigratedDirectionTo
)

const (
	MigrateAccountSignaturePrefix = "MigrateAccount:"

	EventTypeMigrate = "migrate"
	AttributeKeyFrom = "from"
	AttributeKeyTo   = "to"
)

var (
	KeyPrefixMigratedRecord        = []byte{prefixMigratedRecord}
	KeyPrefixMigratedDirectionFrom = []byte{prefixMigratedDirectionFrom}
	KeyPrefixMigratedDirectionTo   = []byte{prefixMigratedDirectionTo}

	ValuePrefixMigrateFromFlag = []byte{0x1}
	ValuePrefixMigrateToFlag   = []byte{0x2}
)

// GetMigratedRecordKey returns the following key format
func GetMigratedRecordKey(addr []byte) []byte {
	return append(KeyPrefixMigratedRecord, addr...)
}

func GetMigratedDirectionFrom(addr sdk.AccAddress) []byte {
	return append(KeyPrefixMigratedDirectionFrom, addr.Bytes()...)
}

func GetMigratedDirectionTo(addr common.Address) []byte {
	return append(KeyPrefixMigratedDirectionTo, addr.Bytes()...)
}
