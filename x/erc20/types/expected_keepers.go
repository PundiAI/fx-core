package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// AccountKeeper defines the expected interface needed to retrieve account info.
type AccountKeeper interface {
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	IsSendEnabledDenom(ctx context.Context, denom string) bool
	BlockedAddr(addr sdk.AccAddress) bool

	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	HasDenomMetaData(ctx context.Context, denom string) bool
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)
}

type ERC20TokenKeeper interface {
	Name(ctx context.Context, contractAddr common.Address) (string, error)
	Symbol(ctx context.Context, contractAddr common.Address) (string, error)
	Decimals(ctx context.Context, contractAddr common.Address) (uint8, error)

	Mint(ctx context.Context, contractAddr, from, receiver common.Address, amount *big.Int) (*evmtypes.MsgEthereumTxResponse, error)
	Transfer(ctx context.Context, contractAddr, from, receiver common.Address, amount *big.Int) (*evmtypes.MsgEthereumTxResponse, error)
}

// EVMKeeper defines the expected EVM keeper interface used on erc20
type EVMKeeper interface {
	GetAccount(ctx sdk.Context, addr common.Address) *statedb.Account
	DeployUpgradableContract(ctx sdk.Context, from, logic common.Address, logicData []byte, initializeAbi *abi.ABI, initializeArgs ...interface{}) (common.Address, error)
	QueryContract(ctx context.Context, from, contract common.Address, abi abi.ABI, method string, res interface{}, args ...interface{}) error
	ApplyContract(ctx context.Context, from, contract common.Address, value *big.Int, abi abi.ABI, method string, constructorData ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
	ExecuteEVM(ctx sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte) (*evmtypes.MsgEthereumTxResponse, error)
}
