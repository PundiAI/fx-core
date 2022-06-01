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

	mainnetEvmChainID      = 1
	mainnetSupportEvmBlock = math.MaxInt64
)

// testnet constant
const (
	testnetEvmChainID      = 90001
	testnetSupportEvmBlock = 2940000
)

// devnet constant
const (
	devnetEvmChainID      = 221
	devnetSupportEvmBlock = 10
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

//func GravityPruneValsetAndAttestationBlock() int64 {
//	if networkDevnet == network {
//		return devnetGravityPruneValsetAndAttestationBlock
//	} else if networkTestnet == network {
//		return testnetGravityPruneValsetAndAttestationBlock
//	}
//	return mainnetGravityPruneValsetAndAttestationBlock
//}
//
//func GravityValsetSlashBlock() int64 {
//	if networkDevnet == network {
//		return devnetGravityValsetSlashBlock
//	} else if networkTestnet == network {
//		return testnetGravityValsetSlashBlock
//	}
//	return mainnetGravityValsetSlashBlock
//}
//
//func CrossChainSupportBscBlock() int64 {
//	if networkDevnet == network {
//		return devnetCrossChainSupportBscBlock
//	} else if networkTestnet == network {
//		return testnetCrossChainSupportBscBlock
//	}
//	return mainnetCrossChainSupportBscBlock
//}
//
//func CrossChainSupportPolygonAndTronBlock() int64 {
//	if networkDevnet == network {
//		return devnetCrossChainSupportTronAndPolygonBlock
//	} else if networkTestnet == network {
//		return testnetCrossChainSupportTronAndPolygonBlock
//	}
//	return mainnetCrossChainSupportTronAndPolygonBlock
//}

func EIP155ChainID() *big.Int {
	if networkDevnet == network {
		return big.NewInt(devnetEvmChainID)
	} else if networkTestnet == network {
		return big.NewInt(testnetEvmChainID)
	}
	return big.NewInt(mainnetEvmChainID)
}

//func EvmV0SupportBlock() int64 {
//	//if networkDevnet == network {
//	//	return devnetSupportEvmV0Block
//	//} else if networkTestnet == network {
//	//	return testnetSupportEvmV0Block
//	//}
//	//return mainnetSupportEvmV0Block
//}

func EvmV1SupportBlock() int64 {
	if networkDevnet == network {
		return devnetSupportEvmBlock
	} else if networkTestnet == network {
		return testnetSupportEvmBlock
	}
	return mainnetSupportEvmBlock
}

//func EthSecp256k1MultisignSupportBlock() int64 {
//	if networkDevnet == network {
//		return devnetSupportEthSecp256k1Multisign
//	} else if networkTestnet == network {
//		return testnetSupportEthSecp256k1Multisign
//	} else {
//		return mainnetSupportEthSecp256k1Multisign
//	}
//}

//func RequestBatchBaseFeeBlock() int64 {
//	return EvmV1SupportBlock()
//}
//
//func IsRequestBatchBaseFee(height int64) bool {
//	return height >= RequestBatchBaseFeeBlock()
//}

func ClearKVStores() []string {
	if networkDevnet == network {
		return []string{}
	} else if networkTestnet == network {
		return []string{"evm", "feemarket"}
	}
	return []string{}
}

// ChangeNetworkForTest change network for test
func ChangeNetworkForTest(newNetwork string) {
	if network != networkDevnet && network != networkTestnet && network != networkMainnet {
		panic("Unsupported network:" + newNetwork)
	}
	network = newNetwork
}
