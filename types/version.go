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

// testnet constant
const (
	testnetCrossChainSupportBscBlock     = 1
	testnetCrossChainSupportTronBlock    = 1
	testnetCrossChainSupportPolygonBlock = 1

	testnetGravityPruneValsetAndAttestationBlock = 1
	testnetGravityValsetSlashBlock               = 1
	testnetSupportEvmBlock                       = 408000
	testnetEvmChainID                            = 90001
	testnetSupportIntrarelayerBlock              = 999999999
	testnetSupportMigrateBlock                   = 999999999
	testnetSupportEGFProposalBlock               = 999999999
	testnetRequestBaseFee                        = 999999999
)

// mainnet constant
const (
	mainnetCrossChainSupportBscBlock     = 1354000
	mainnetCrossChainSupportTronBlock    = 2062000
	mainnetCrossChainSupportPolygonBlock = 2062000

	//
	mainnetGravityPruneValsetAndAttestationBlock = 610000
	// gravity not slash no set eth address validator
	mainnetGravityValsetSlashBlock  = 1685000
	mainnetSupportEvmBlock          = 999999999
	mainnetEvmChainID               = 1
	mainnetSupportIntrarelayerBlock = 999999999
	mainnetRequestBaseFee           = 999999999
	mainnetSupportMigrateBlock      = 999999999
	mainnetSupportEGFProposalBlock  = 999999999
)

// devnet constant
const (
	devnetCrossChainSupportBscBlock     = 1
	devnetCrossChainSupportTronBlock    = 1
	devnetCrossChainSupportPolygonBlock = 1

	devnetGravityPruneValsetAndAttestationBlock = 1
	devnetGravityValsetSlashBlock               = 1
	devnetSupportEvmBlock                       = 100
	devnetEvmChainID                            = 221
	devnetSupportIntrarelayerBlock              = 100
	devnetSupportMigrateBlock                   = 100
	devnetSupportEGFProposalBlock               = 100
	devnetRequestBaseFee                        = 62000
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

func NetworkMiannet() string {
	return networkMainnet
}
func NetworkTestnet() string {
	return networkTestnet
}
func NetworkDevnet() string {
	return networkDevnet
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
		return big.NewInt(devnetEvmChainID)
	} else if networkTestnet == network {
		return big.NewInt(testnetEvmChainID)
	}
	return big.NewInt(mainnetEvmChainID)
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
		return testnetSupportIntrarelayerBlock
	}
	return mainnetSupportIntrarelayerBlock
}

func MigrateSupportBlock() int64 {
	if networkDevnet == network {
		return devnetSupportMigrateBlock
	} else if networkTestnet == network {
		return testnetSupportMigrateBlock
	}
	return mainnetSupportMigrateBlock
}

func EGFProposalSupportBlock() int64 {
	if networkDevnet == network {
		return devnetSupportEGFProposalBlock
	} else if networkTestnet == network {
		return testnetSupportEGFProposalBlock
	}
	return mainnetSupportEGFProposalBlock
}

func RequestBatchBaseFee() int64 {
	if networkDevnet == network {
		return devnetRequestBaseFee
	} else if networkTestnet == network {
		return testnetRequestBaseFee
	}
	return mainnetRequestBaseFee
}

func IsRequestBatchBaseFee(height int64) bool {
	return height >= RequestBatchBaseFee()
}
