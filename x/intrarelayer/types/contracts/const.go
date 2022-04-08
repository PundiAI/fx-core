package contracts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ContractType int

type Keeper interface {
	CreateContractWithCode(ctx sdk.Context, addr common.Address, code []byte) error
}

type UpgradeHook func(ctx sdk.Context, k Keeper, uc UpgradeConfig) error

type UpgradeConfig struct {
	BeforeUpgrade UpgradeHook
	AfterUpgrade  UpgradeHook
	ContractAddr  common.Address
	ABI           abi.ABI
	Bin           []byte
	Code          []byte
	Type          ContractType
	Description   string
}

type Upgrade struct {
	Name    string
	Configs []UpgradeConfig
}

type CompileContract struct {
	ABI  abi.ABI
	Bin  []byte
	Code []byte
}

type UpgradeBlockConfig struct {
	InitUpgradeBlock int64
	TestUpgradeBlock int64
}

func (c UpgradeBlockConfig) IsOnInitUpgrade(blockHeight int64) bool {
	if c.InitUpgradeBlock == blockHeight {
		return true
	}
	return false
}
func (c UpgradeBlockConfig) GTEInitUpgrade(blockHeight int64) bool {
	if blockHeight >= c.InitUpgradeBlock {
		return true
	}
	return false
}

func (c UpgradeBlockConfig) IsOnTestUpgrade(blockHeight int64) bool {
	if c.TestUpgradeBlock == blockHeight {
		return true
	}
	return false
}
func (c UpgradeBlockConfig) GTETestUpgrade(blockHeight int64) bool {
	if blockHeight >= c.TestUpgradeBlock {
		return true
	}
	return false
}

const (
	EventTypeUpgradeContract = "upgrade_contract"
	AttributeKeyAddress      = "address"
)

const (
	ERC1967ProxyType ContractType = iota + 1
	FIP20UpgradeType
	WFXUpgradeType
)

const (
	EmptyAddress            = "0x0000000000000000000000000000000000000000"
	FIP20UpgradeCodeAddress = "0x0000000000000000000000000000000000001001"
	WFXUpgradeCodeAddress   = "0x0000000000000000000000000000000000001002"
)
