package types

import (
	"context"
	"math/big"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
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
	IsSendEnabledCoin(ctx context.Context, coin sdk.Coin) bool
	BlockedAddr(addr sdk.AccAddress) bool
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	HasDenomMetaData(ctx context.Context, denom string) bool
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// EVMKeeper defines the expected EVM keeper interface used on erc20
type EVMKeeper interface {
	GetAccount(ctx sdk.Context, addr common.Address) *statedb.Account
	QueryContract(ctx sdk.Context, from, contract common.Address, abi abi.ABI, method string, res interface{}, constructorData ...interface{}) error
	ApplyContract(ctx sdk.Context, from, contract common.Address, value *big.Int, abi abi.ABI, method string, constructorData ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
	DeployUpgradableContract(ctx sdk.Context, from, logic common.Address, logicData []byte, initializeAbi *abi.ABI, initializeArgs ...interface{}) (common.Address, error)
}

type IBCTransferKeeper interface {
	Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error)
	GetDenomTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) (types.DenomTrace, bool)
}
