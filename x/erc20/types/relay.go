package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type RelayTransfer struct {
	From          common.Address
	Amount        *big.Int
	TokenContract common.Address
	Denom         string
	ContractOwner Owner
}

type RelayTransferCrossChain struct {
	*TransferCrossChainEvent
	TokenContract common.Address
	Denom         string
	ContractOwner Owner
}
