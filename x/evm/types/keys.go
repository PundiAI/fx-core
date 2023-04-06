package types

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	TypeAddress, _ = abi.NewType("address", "", nil)
	TypeUint256, _ = abi.NewType("uint256", "", nil)
	TypeString, _  = abi.NewType("string", "", nil)
	TypeBool, _    = abi.NewType("bool", "", nil)
	TypeBytes32, _ = abi.NewType("bytes32", "", nil)
)
