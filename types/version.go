package types

import (
	"math"
	"math/big"
	"os"
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
	mainnetEvmChainID = 530

	mainnetSupportDenomManyToOneBlock = math.MaxInt64
)

// testnet constant
const (
	testnetChainId        = "dhobyghaut"
	testnetEvmChainID     = 90001
	testnetIBCRouterBlock = 3433511

	testnetSupportDenomManyToOneBlock = 3918000
)

// SupportDenomManyToOneMsgTypes return msg types
// use method return to prevent accidental modification
func SupportDenomManyToOneMsgTypes() []string {
	return []string{
		"/fx.erc20.v1.MsgConvertDenom",
		"/fx.erc20.v1.UpdateDenomAliasProposal",
	}
}

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

func SetTestingManyToOneBlock(fn func() int64) {
	if os.Getenv("GO_ENV") != "testing" {
		panic("invalid env")
	}
	testingManyToOneBlock = fn
}

var testingManyToOneBlock func() int64

func SupportDenomManyToOneBlock() int64 {
	if os.Getenv("GO_ENV") == "testing" {
		return testingManyToOneBlock()
	}
	if testnetChainId == chainId {
		return testnetSupportDenomManyToOneBlock
	}
	return mainnetSupportDenomManyToOneBlock
}
