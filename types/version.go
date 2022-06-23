package types

import (
	"math/big"
)

// network constant
const (
	networkMainnet = "mainnet"
	networkTestnet = "testnet"
	networkDevnet  = "devnet"
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
	testnetChainId    = "dhobyghaut"
	testnetEvmChainID = 90001
)

// devnet constant
const (
	devnetChainId    = "boonlay"
	devnetEvmChainID = 221
)

var (
	// network config network, default mainnet
	network = networkMainnet
)

func init() {
	if network != networkTestnet && network != networkMainnet && network != networkDevnet {
		network = networkMainnet
	}
}

func Network() string {
	return network
}

func NetworkMainnet() string {
	return networkMainnet
}
func NetworkTestnet() string {
	return networkTestnet
}
func NetworkDevnet() string {
	return networkDevnet
}

func ChainId() string {
	if networkDevnet == network {
		return devnetChainId
	} else if networkTestnet == network {
		return testnetChainId
	}
	return mainnetChainId
}
func EIP155ChainID() *big.Int {
	if networkDevnet == network {
		return big.NewInt(devnetEvmChainID)
	} else if networkTestnet == network {
		return big.NewInt(testnetEvmChainID)
	}
	return big.NewInt(mainnetEvmChainID)
}
