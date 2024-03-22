package contract

import "github.com/ethereum/go-ethereum/accounts/abi"

var (
	TypeString, _       = abi.NewType("string", "", nil)
	TypeBytes, _        = abi.NewType("bytes", "", nil)
	TypeUint256Array, _ = abi.NewType("uint256[]", "", nil)
)
