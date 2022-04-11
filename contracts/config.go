package contracts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Keeper interface {
	CreateContractWithCode(ctx sdk.Context, addr common.Address, code []byte) error
}

type Config struct {
	ContractAddr common.Address
	ABI          abi.ABI
	Bin          []byte
	Code         []byte
	Description  string
}

type Upgrade struct {
	Name    string
	Configs []Config
}

type BlockConfig struct {
	InitUpgradeBlock int64
	TestUpgradeBlock int64
}

func (c BlockConfig) IsOnInitUpgrade(blockHeight int64) bool {
	if c.InitUpgradeBlock == blockHeight {
		return true
	}
	return false
}
func (c BlockConfig) GTEInitUpgrade(blockHeight int64) bool {
	if blockHeight >= c.InitUpgradeBlock {
		return true
	}
	return false
}

func (c BlockConfig) IsOnTestUpgrade(blockHeight int64) bool {
	if c.TestUpgradeBlock == blockHeight {
		return true
	}
	return false
}
func (c BlockConfig) GTETestUpgrade(blockHeight int64) bool {
	if blockHeight >= c.TestUpgradeBlock {
		return true
	}
	return false
}

func GetERC20Config(height int64) Config {
	return Config{}
}

func GetWFXConfig(height int64) Config {
	return Config{}
}

func GetERC1967ProxyConfig(height int64) Config {
	return Config{}
}
