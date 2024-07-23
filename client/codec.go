package client

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ethcryptocodec "github.com/evmos/ethermint/crypto/codec"
	ethermint "github.com/evmos/ethermint/types"
)

func NewAccountCodec() types.InterfaceRegistry {
	interfaceRegistry := types.NewInterfaceRegistry()
	ethermint.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	ethcryptocodec.RegisterInterfaces(interfaceRegistry)
	return interfaceRegistry
}
