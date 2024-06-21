package contract

import "github.com/ethereum/go-ethereum/accounts/abi"

const (
	DefaultGasCap uint64 = 30000000
)

var TypeString, _ = abi.NewType("string", "", nil)
