package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// CompatibleModuleName is the query and tx module name
	CompatibleModuleName = "fxtransfer"
)

func IntermediateSender(sourcePort, sourceChannel, sender string) common.Address {
	prefix := fmt.Sprintf("%s/%s", sourcePort, sourceChannel)
	senderHash32 := address.Hash(prefix, []byte(sender))
	return common.BytesToAddress(senderHash32)
}
