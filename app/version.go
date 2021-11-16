package app

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
)

// devnet constant
const (
	devnetCrossChainSupportBscBlock     = 1
	devnetCrossChainSupportTronBlock    = 225600
	devnetCrossChainSupportPolygonBlock = 225600

	devnetGravityPruneValsetAndAttestationBlock = 1
	devnetGravityValsetSlashBlock               = 1
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
