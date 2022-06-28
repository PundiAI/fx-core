package types

import (
	"math/big"
	"sync"
)

// mainnet constant
const (

	/*
		mainnetCrossChainSupportBscBlock            = 1354000
		mainnetCrossChainSupportTronAndPolygonBlock = 2062000
		// gravity prune validator set
		mainnetGravityPruneValsetAndAttestationBlock = 610000
		// gravity not slash no set eth address validator
		mainnetGravityValsetSlashBlock = 1685000
	*/

	mainnetChainId    = "fxcore"
	mainnetEvmChainID = 1
)

// testnet constant
const (
	testnetChainId        = "dhobyghaut"
	testnetEvmChainID     = 90001
	testnetIBCRouterBlock = 3433511
)

var (
	chainId = mainnetChainId
	once    sync.Once
)

func SetChainId(id string) {
	if id != mainnetChainId && id != testnetChainId {
		panic("invalid chainId: " + id)
	}
	once.Do(func() {
		chainId = id
	})
}

func ChainId() string {
	return chainId
}

func NetworkMainnet() string {
	return mainnetChainId
}
func TestnetChainId() string {
	return testnetChainId
}

func EIP155ChainID() *big.Int {
	if testnetChainId == chainId {
		return big.NewInt(testnetEvmChainID)
	}
	return big.NewInt(mainnetEvmChainID)
}

func IBCRouteBlock() int64 {
	if testnetChainId == chainId {
		return testnetIBCRouterBlock
	}
	return 0
}
