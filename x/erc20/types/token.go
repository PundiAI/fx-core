package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

type TokenPair ERC20Token

func (m *ERC20Token) GetERC20Contract() common.Address {
	return common.HexToAddress(m.Erc20Address)
}

func (m *ERC20Token) IsNativeCoin() bool {
	return m.ContractOwner == OWNER_MODULE
}

func (m *ERC20Token) IsNativeERC20() bool {
	return m.ContractOwner == OWNER_EXTERNAL
}

func (m *BridgeToken) BridgeDenom() string {
	if m.IsOrigin() {
		return m.Denom
	}
	return NewBridgeDenom(m.ChainName, m.Contract)
}

func (m *BridgeToken) IsOrigin() bool {
	return m.Denom == fxtypes.DefaultDenom
}

func (m *BridgeToken) GetContractAddress() common.Address {
	return common.HexToAddress(m.Contract)
}

func NewBridgeDenom(moduleName, token string) string {
	return fmt.Sprintf("%s%s", moduleName, token)
}
