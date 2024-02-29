package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the IBC transfer name
	ModuleName = "transfer"

	// CompatibleModuleName is the query and tx module name
	CompatibleModuleName = "fxtransfer"

	// RouterKey is the message route for IBC transfer
	RouterKey = CompatibleModuleName
)

func IntermediateSender(sourcePort, sourceChannel, sender string) common.Address {
	prefix := fmt.Sprintf("%s/%s", sourcePort, sourceChannel)
	senderHash32 := address.Hash(prefix, []byte(sender))
	return common.BytesToAddress(senderHash32)
}
