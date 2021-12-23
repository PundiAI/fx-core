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

// testnet constant
const (
	testnetCrossChainSupportBscBlock     = 1
	testnetCrossChainSupportTronBlock    = 1
	testnetCrossChainSupportPolygonBlock = 1

	testnetGravityPruneValsetAndAttestationBlock = 1
	testnetGravityValsetSlashBlock               = 1
	testnetSupportEvmBlock                       = math.MaxInt
)

// mainnet constant
const (
	mainnetCrossChainSupportBscBlock     = 1354000
	mainnetCrossChainSupportTronBlock    = 2062000
	mainnetCrossChainSupportPolygonBlock = 2062000

	//
	mainnetGravityPruneValsetAndAttestationBlock = 610000
	// gravity not slash no set eth address validator
	mainnetGravityValsetSlashBlock = 1685000
	mainnetSupportEvmBlock         = math.MaxInt
)

// devnet constant
const (
	devnetCrossChainSupportBscBlock     = 1
	devnetCrossChainSupportTronBlock    = 1
	devnetCrossChainSupportPolygonBlock = 1

	devnetGravityPruneValsetAndAttestationBlock = 1
	devnetGravityValsetSlashBlock               = 1
	devnetSupportEvmBlock                       = 300
	devnetSupportIntrarelayerBlock              = math.MaxInt
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

func GravityPruneValsetsAndAttestationBlock() int64 {
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

func CrossChainSupportTronBlock() int64 {
	if networkDevnet == network {
		return devnetCrossChainSupportTronBlock
	} else if networkTestnet == network {
		return testnetCrossChainSupportTronBlock
	}
	return mainnetCrossChainSupportTronBlock
}

func CrossChainSupportPolygonBlock() int64 {
	if networkDevnet == network {
		return devnetCrossChainSupportPolygonBlock
	} else if networkTestnet == network {
		return testnetCrossChainSupportPolygonBlock
	}
	return mainnetCrossChainSupportPolygonBlock
}

func EIP155ChainID() *big.Int {
	if networkDevnet == network {
		return big.NewInt(221)
	} else if networkTestnet == network {
		return big.NewInt(555)
	}
	return big.NewInt(1)
}

func EvmSupportBlock() int64 {
	if networkDevnet == network {
		return devnetSupportEvmBlock
	} else if networkTestnet == network {
		return testnetSupportEvmBlock
	}
	return mainnetSupportEvmBlock
}

func IntrarelayerSupportBlock() int64 {
	if networkDevnet == network {
		return devnetSupportIntrarelayerBlock
	} else if networkTestnet == network {
		return devnetSupportIntrarelayerBlock
	}
	return devnetSupportIntrarelayerBlock
}
