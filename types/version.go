package types

import (
	"math"
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
	mainnetCrossChainSupportBscBlock            = 1354000
	mainnetCrossChainSupportTronAndPolygonBlock = 2062000

	// gravity prune validator set
	mainnetGravityPruneValsetAndAttestationBlock = 610000
	// gravity not slash no set eth address validator
	mainnetGravityValsetSlashBlock = 1685000

	mainnetEvmChainID        = 1
	mainnetSupportEvmV0Block = math.MaxInt64
	mainnetSupportEvmV1Block = math.MaxInt64
)

// testnet constant
const (
	testnetCrossChainSupportBscBlock            = 1
	testnetCrossChainSupportTronAndPolygonBlock = 1

	testnetGravityPruneValsetAndAttestationBlock = 1
	testnetGravityValsetSlashBlock               = 1

	testnetEvmChainID        = 90001
	testnetSupportEvmV0Block = 408000
	testnetSupportEvmV1Block = math.MaxInt64
)

// devnet constant
const (
	devnetCrossChainSupportBscBlock            = 1
	devnetCrossChainSupportTronAndPolygonBlock = 1

	devnetGravityPruneValsetAndAttestationBlock = 1
	devnetGravityValsetSlashBlock               = 1

	devnetEvmChainID        = 221
	devnetSupportEvmV0Block = math.MaxInt64
	devnetSupportEvmV1Block = 10
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

func GravityPruneValsetAndAttestationBlock() int64 {
	if networkDevnet == network {
		return devnetGravityPruneValsetAndAttestationBlock
	} else if networkTestnet == network {
		return testnetGravityPruneValsetAndAttestationBlock
	}
	return mainnetGravityPruneValsetAndAttestationBlock
}

func GravityValsetSlashBlock() int64 {
	if networkDevnet == network {
		return devnetGravityValsetSlashBlock
	} else if networkTestnet == network {
		return testnetGravityValsetSlashBlock
	}
	return mainnetGravityValsetSlashBlock
}

func CrossChainSupportBscBlock() int64 {
	if networkDevnet == network {
		return devnetCrossChainSupportBscBlock
	} else if networkTestnet == network {
		return testnetCrossChainSupportBscBlock
	}
	return mainnetCrossChainSupportBscBlock
}

func CrossChainSupportPolygonAndTronBlock() int64 {
	if networkDevnet == network {
		return devnetCrossChainSupportTronAndPolygonBlock
	} else if networkTestnet == network {
		return testnetCrossChainSupportTronAndPolygonBlock
	}
	return mainnetCrossChainSupportTronAndPolygonBlock
}

func EIP155ChainID() *big.Int {
	if networkDevnet == network {
		return big.NewInt(devnetEvmChainID)
	} else if networkTestnet == network {
		return big.NewInt(testnetEvmChainID)
	}
	return big.NewInt(mainnetEvmChainID)
}

func EvmV0SupportBlock() int64 {
	if networkDevnet == network {
		return devnetSupportEvmV0Block
	} else if networkTestnet == network {
		return testnetSupportEvmV0Block
	}
	return mainnetSupportEvmV0Block
}

func EvmV1SupportBlock() int64 {
	if networkDevnet == network {
		return devnetSupportEvmV1Block
	} else if networkTestnet == network {
		return testnetSupportEvmV1Block
	}
	return mainnetSupportEvmV1Block
}

func RequestBatchBaseFeeBlock() int64 {
	return EvmV1SupportBlock()
}

func IsRequestBatchBaseFee(height int64) bool {
	return height >= RequestBatchBaseFeeBlock()
}

// ChangeNetworkForTest change network for test
func ChangeNetworkForTest(newNetwork string) {
	if network != networkDevnet && network != networkTestnet && network != networkMainnet {
		panic("Unsupported network:" + newNetwork)
	}
	network = newNetwork
}
