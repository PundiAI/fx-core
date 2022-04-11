package contracts

import (
	fxcoretypes "github.com/functionx/fx-core/types"
	"math"
)

var (
	mainnetConfig = BlockConfig{
		InitUpgradeBlock: 0,
		TestUpgradeBlock: math.MaxInt64,
	}

	testnetConfig = BlockConfig{
		InitUpgradeBlock: 0,
		TestUpgradeBlock: math.MaxInt64,
	}

	devnetConfig = BlockConfig{
		InitUpgradeBlock: 0,
		TestUpgradeBlock: math.MaxInt64,
	}
)

var (
	codeInit    = &Upgrade{Name: "code init", Configs: []Config{erc1967ProxyConfig, fip20InitConfig, wfxInitConfig}}
	upgradeTest = &Upgrade{Name: "upgrade test", Configs: []Config{fip20TestConfig}}

	initUpgrade = map[string]*Upgrade{
		fxcoretypes.NetworkMiannet(): codeInit,
		fxcoretypes.NetworkTestnet(): codeInit,
		fxcoretypes.NetworkDevnet():  codeInit,
	}

	testUpgrade = map[string]*Upgrade{
		fxcoretypes.NetworkDevnet(): upgradeTest,
	}
)

func GetUpgradeBlockConfig(network string) BlockConfig {
	switch network {
	case fxcoretypes.NetworkTestnet():
		return testnetConfig
	case fxcoretypes.NetworkDevnet():
		return devnetConfig
	default:
		return mainnetConfig
	}
}

func GetERC20Config(height int64) Config {
	network := fxcoretypes.Network()
	bc := GetUpgradeBlockConfig(network)

	if height >= bc.InitUpgradeBlock && height < bc.TestUpgradeBlock {
		return fip20InitConfig
	} else if height >= bc.TestUpgradeBlock {
		return fip20TestConfig
	}

	return fip20InitConfig
}

func GetWFXConfig(height int64) Config {
	network := fxcoretypes.Network()
	bc := GetUpgradeBlockConfig(network)

	if height > bc.InitUpgradeBlock {
		return wfxInitConfig
	}
	return wfxInitConfig
}

func GetERC1967ProxyConfig(height int64) Config {
	network := fxcoretypes.Network()
	bc := GetUpgradeBlockConfig(network)

	if height > bc.InitUpgradeBlock {
		return erc1967ProxyConfig
	}
	return erc1967ProxyConfig
}

// GetInitConfig get a copy of init config
func GetInitConfig(network string) *Upgrade {
	upgrade, ok := initUpgrade[network]
	if ok {
		upgradeCopy := *upgrade
		return &upgradeCopy
	}
	return nil
}

// GetTestConfig get a copy of test config
func GetTestConfig(network string) *Upgrade {
	upgrade, ok := testUpgrade[network]
	if ok {
		upgradeCopy := *upgrade
		return &upgradeCopy
	}
	return nil
}
