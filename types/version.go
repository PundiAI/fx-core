package types

import "sync"

// testnet constant
const (
	TestnetChainId                       = "dhobyghaut"
	testnetCrossChainSupportBscBlock     = 1
	testnetCrossChainSupportTronBlock    = 1
	testnetCrossChainSupportPolygonBlock = 1

	testnetGravityPruneValsetAndAttestationBlock = 1
	testnetGravityValsetSlashBlock               = 1
)

// mainnet constant
const (
	MainnetChainId                       = "fxcore"
	mainnetCrossChainSupportBscBlock     = 1354000
	mainnetCrossChainSupportTronBlock    = 2062000
	mainnetCrossChainSupportPolygonBlock = 2062000

	// gravity prune validator set and attestation
	mainnetGravityPruneValsetAndAttestationBlock = 610000
	// gravity not slash no set eth address validator
	mainnetGravityValsetSlashBlock = 1685000
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

func GravityPruneValsetsAndAttestationBlock() int64 {
	if TestnetChainId == chainId {
		return testnetGravityPruneValsetAndAttestationBlock
	}
	return mainnetGravityPruneValsetAndAttestationBlock
}

func GravityValsetSlashBlock() int64 {
	if TestnetChainId == chainId {
		return testnetGravityValsetSlashBlock
	}
	return mainnetGravityValsetSlashBlock
}

func CrossChainSupportBscBlock() int64 {
	if TestnetChainId == chainId {
		return testnetCrossChainSupportBscBlock
	}
	return mainnetCrossChainSupportBscBlock
}

func CrossChainSupportTronBlock() int64 {
	if TestnetChainId == chainId {
		return testnetCrossChainSupportTronBlock
	}
	return mainnetCrossChainSupportTronBlock
}

func CrossChainSupportPolygonBlock() int64 {
	if TestnetChainId == chainId {
		return testnetCrossChainSupportPolygonBlock
	}
	return mainnetCrossChainSupportPolygonBlock
}
