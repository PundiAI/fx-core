package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	ContractAddr common.Address
	ABI          abi.ABI
	Bin          []byte
	Code         []byte
	Description  string
}

func GetERC20Config(height int64) Config {
	return Config{}
}
