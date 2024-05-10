package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
)

var _ ExternalAddress = EthereumAddress{}

type EthereumAddress struct{}

func (b EthereumAddress) ValidateExternalAddr(addr string) error {
	return contract.ValidateEthereumAddress(addr)
}

func (b EthereumAddress) ExternalAddrToAccAddr(addr string) sdk.AccAddress {
	return common.HexToAddress(addr).Bytes()
}

func (b EthereumAddress) ExternalAddrToHexAddr(addr string) common.Address {
	return common.HexToAddress(addr)
}

func (b EthereumAddress) ExternalAddrToStr(bz []byte) string {
	return common.BytesToAddress(bz).String()
}
