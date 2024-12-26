package contract

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	DefaultMaxQuoteIndex = 3

	// TransferModuleRole is keccak256("TRANSFER_MODULE_ROLE")
	TransferModuleRole = "0x4845f2571489e4ee59e15b11b74598e4330ef896ebb57513ebdbdb3b260a4671"
)

type BridgeDenoms struct {
	ChainName string
	Denoms    []string
}

func GetTransferModuleRole() common.Hash {
	return common.HexToHash(TransferModuleRole)
}
