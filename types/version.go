package types

import (
	"math/big"
	"os"
	"sync"
)

// mainnet
const (
	MainnetChainId    = "fxcore"
	mainnetEvmChainID = 530

	// v2.2.x upgrade code-named is Exponential
	mainnetExponentialBlock = 5940000

	// v2.3.x upgrade code-named is Trigonometric
	mainnetTrigonometricBlock = 6807000
)

// testnet
const (
	TestnetChainId    = "dhobyghaut"
	testnetEvmChainID = 90001

	testnetIBCRouterBlock = 3433511

	testnetExponential1Block = 3918000
	testnetExponential2Block = 4028000

	testnetTrigonometric1Block = 4317959
	testnetTrigonometric2Block = 4714000
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

func IBCRouteBlock() int64 {
	if TestnetChainId == ChainId() {
		return testnetIBCRouterBlock
	}
	return 0
}

func UpgradeExponential1Block() int64 {
	if os.Getenv("GO_ENV") == "testing" {
		return 0
	}
	if TestnetChainId == ChainId() {
		return testnetExponential1Block
	}
	return mainnetExponentialBlock
}

func UpgradeExponential2Block() int64 {
	if os.Getenv("GO_ENV") == "testing" {
		return 0
	}
	if TestnetChainId == ChainId() {
		return testnetExponential2Block
	}
	return mainnetExponentialBlock
}

func UpgradeTrigonometric1Block() int64 {
	if TestnetChainId == ChainId() {
		return testnetTrigonometric1Block
	}
	return mainnetTrigonometricBlock
}

func UpgradeTrigonometric2Block() int64 {
	if TestnetChainId == ChainId() {
		return testnetTrigonometric2Block
	}
	return mainnetTrigonometricBlock
}
