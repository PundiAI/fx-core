package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
)

type EthereumAddress struct{}

func (b EthereumAddress) ValidateExternalAddress(addr string) error {
	return contract.ValidateEthereumAddress(addr)
}

func (b EthereumAddress) ExternalAddressToAccAddress(addr string) (sdk.AccAddress, error) {
	if err := contract.ValidateEthereumAddress(addr); err != nil {
		return nil, err
	}
	return common.HexToAddress(addr).Bytes(), nil
}
