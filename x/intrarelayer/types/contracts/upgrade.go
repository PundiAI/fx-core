package contracts

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	fxcoretypes "github.com/functionx/fx-core/types"
	"math"
)

var (
	mainnetConfig = UpgradeBlockConfig{
		InitUpgradeBlock: 0,
		TestUpgradeBlock: math.MaxInt64,
	}

	testnetConfig = UpgradeBlockConfig{
		InitUpgradeBlock: 0,
		TestUpgradeBlock: math.MaxInt64,
	}

	devnetConfig = UpgradeBlockConfig{
		InitUpgradeBlock: 0,
		TestUpgradeBlock: 120,
	}
)

var (
	initUpgrade = make(map[string]*Upgrade)

	testUpgrade = make(map[string]*Upgrade)
)

func init() {
	//init system contract
	initUpgrade[fxcoretypes.NetworkMiannet()] = codeInit
	initUpgrade[fxcoretypes.NetworkTestnet()] = codeInit
	initUpgrade[fxcoretypes.NetworkDevnet()] = codeInit

	//test system contract upgrade
	testUpgrade[fxcoretypes.NetworkDevnet()] = upgradeTest
}

func GetABI(height int64, contractType ContractType) (abi.ABI, bool) {
	abis, err := GetSystemContractABI(height)
	if err != nil {
		return abi.ABI{}, false
	}
	abiJson, ok := abis[contractType]
	return abiJson, ok
}
func MustGetABI(height int64, contractType ContractType) abi.ABI {
	abis, err := GetSystemContractABI(height)
	if err != nil {
		panic(err)
	}
	abiJson, ok := abis[contractType]
	if !ok {
		panic(fmt.Sprintf("abi of type %v not found", contractType))
	}
	return abiJson
}

func GetBin(height int64, contractType ContractType) ([]byte, bool) {
	bins, err := GetSystemContractBin(height)
	if err != nil {
		return nil, false
	}
	bin, ok := bins[contractType]
	return bin, ok
}
func MustGetBin(height int64, contractType ContractType) []byte {
	bins, err := GetSystemContractBin(height)
	if err != nil {
		panic(err)
	}
	bin, ok := bins[contractType]
	if !ok {
		panic(fmt.Sprintf("bin of type %v not found", contractType))
	}
	return bin
}

func GetCode(height int64, contractType ContractType) ([]byte, bool) {
	codes, err := GetSystemContractCode(height)
	if err != nil {
		return nil, false
	}
	code, ok := codes[contractType]
	return code, ok
}
func MustGetCode(height int64, contractType ContractType) []byte {
	codes, err := GetSystemContractCode(height)
	if err != nil {
		panic(err)
	}
	code, ok := codes[contractType]
	if !ok {
		panic(fmt.Sprintf("code of type %v not found", contractType))
	}
	return code
}

func GetSystemContractABI(blockHeight int64) (map[ContractType]abi.ABI, error) {
	network := fxcoretypes.Network()
	abis := make(map[ContractType]abi.ABI)
	err := scopeGTEOptionOfSystemContract(network, blockHeight, func(upgrade *Upgrade) error {
		return systemContractABI(abis, upgrade)
	})
	return abis, err
}
func GetSystemContractBin(blockHeight int64) (map[ContractType][]byte, error) {
	network := fxcoretypes.Network()
	bins := make(map[ContractType][]byte)
	err := scopeGTEOptionOfSystemContract(network, blockHeight, func(upgrade *Upgrade) error {
		return systemContractBin(bins, upgrade)
	})
	return bins, err
}
func GetSystemContractCode(blockHeight int64) (map[ContractType][]byte, error) {
	network := fxcoretypes.Network()
	codes := make(map[ContractType][]byte)
	err := scopeGTEOptionOfSystemContract(network, blockHeight, func(upgrade *Upgrade) error {
		return systemContractCode(codes, upgrade)
	})
	return codes, err
}

func GetUpgradeBlockConfig(network string) UpgradeBlockConfig {
	switch network {
	case fxcoretypes.NetworkTestnet():
		return testnetConfig
	case fxcoretypes.NetworkDevnet():
		return devnetConfig
	default:
		return mainnetConfig
	}
}

func InitSystemContract(ctx sdk.Context, k Keeper) error {
	network := fxcoretypes.Network()
	upgrade, ok := initUpgrade[network]
	if !ok {
		return errors.New("empty system contract")
	}
	return upgradeSystemContract(ctx, k, upgrade)
}
func UpgradeSystemContract(ctx sdk.Context, k Keeper) error {
	network := fxcoretypes.Network()
	blockHeight := ctx.BlockHeight()
	return exactOptionOfSystemContract(network, blockHeight, func(upgrade *Upgrade) error {
		return upgradeSystemContract(ctx, k, upgrade)
	})
}

func exactOptionOfSystemContract(network string, blockHeight int64, fn func(upgrade *Upgrade) error) error {
	bc := GetUpgradeBlockConfig(network)
	if bc.IsOnInitUpgrade(blockHeight) {
		if err := fn(initUpgrade[network]); err != nil {
			return err
		}
	}
	if bc.IsOnTestUpgrade(blockHeight) {
		if err := fn(testUpgrade[network]); err != nil {
			return err
		}
	}
	return nil
}

func scopeGTEOptionOfSystemContract(network string, blockHeight int64, fn func(upgrade *Upgrade) error) error {
	bc := GetUpgradeBlockConfig(network)
	if bc.GTEInitUpgrade(blockHeight) {
		if err := fn(initUpgrade[network]); err != nil {
			return err
		}
	}
	if bc.GTETestUpgrade(blockHeight) {
		if err := fn(testUpgrade[network]); err != nil {
			return err
		}
	}
	return nil
}

func upgradeSystemContract(ctx sdk.Context, k Keeper, upgrade *Upgrade) error {
	if upgrade == nil {
		ctx.Logger().Info("empty upgrade config", "height", ctx.BlockHeight())
		return nil
	}
	ctx.Logger().Info("upgrade system contract", "name", upgrade.Name, "height", ctx.BlockHeight())
	for _, cfg := range upgrade.Configs {
		if cfg.ContractAddr.Hex() == EmptyAddress {
			continue
		}
		ctx.Logger().Info("upgrade contract", "address", cfg.ContractAddr.Hex(), "type", cfg.Type)
		if cfg.BeforeUpgrade != nil {
			if err := cfg.BeforeUpgrade(ctx, k, cfg); err != nil {
				return err
			}
		}

		if err := k.CreateContractWithCode(ctx, cfg.ContractAddr, cfg.Code); err != nil {
			return err
		}

		if cfg.AfterUpgrade != nil {
			if err := cfg.AfterUpgrade(ctx, k, cfg); err != nil {
				return err
			}
		}
		ctx.EventManager().EmitEvents(
			sdk.Events{
				sdk.NewEvent(
					EventTypeUpgradeContract,
					sdk.NewAttribute(AttributeKeyAddress, cfg.ContractAddr.Hex()),
				),
			},
		)
	}
	return nil
}
func systemContractABI(abis map[ContractType]abi.ABI, upgrade *Upgrade) error {
	if upgrade == nil {
		return nil
	}
	for _, cfg := range upgrade.Configs {
		abis[cfg.Type] = cfg.ABI
	}
	return nil
}
func systemContractBin(bins map[ContractType][]byte, upgrade *Upgrade) error {
	if upgrade == nil {
		return nil
	}
	for _, cfg := range upgrade.Configs {
		bins[cfg.Type] = cfg.Bin
	}
	return nil
}
func systemContractCode(codes map[ContractType][]byte, upgrade *Upgrade) error {
	if upgrade == nil {
		return nil
	}
	for _, cfg := range upgrade.Configs {
		codes[cfg.Type] = cfg.Code
	}
	return nil
}
