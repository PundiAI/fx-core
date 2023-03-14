package types

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	GenesisCoinbase = common.HexToAddress("0x0000000000000000000000000000000000000000")

	TypeAddress, _ = abi.NewType("address", "", nil)
	TypeUint256, _ = abi.NewType("uint256", "", nil)
	TypeString, _  = abi.NewType("string", "", nil)
)
