package types

import (
	"math/big"
	"sync"
)

// mainnet
const (
	MainnetChainId    = "fxcore"
	mainnetEvmChainID = 530
)

// testnet
const (
	TestnetChainId    = "dhobyghaut"
	testnetEvmChainID = 90001
)

var (
	chainId = MainnetChainId
	once    sync.Once
)

func SetChainId(id string) {
	if id != MainnetChainId && id != TestnetChainId {
		panic("invalid chainId: " + id)
	}
	once.Do(func() {
		chainId = id
	})
}

func ChainId() string {
	return chainId
}

func EIP155ChainID() *big.Int {
	if TestnetChainId == ChainId() {
		return big.NewInt(testnetEvmChainID)
	}
	return big.NewInt(mainnetEvmChainID)
}
