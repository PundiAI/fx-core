package types

import (
	"context"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/functionx/fx-core/x/evm/statedb"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// AccountKeeper defines the expected interface needed to retrieve account info.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)
	NewAccountWithAddress(sdk.Context, sdk.AccAddress) types.AccountI
	SetAccount(sdk.Context, types.AccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendEnabledCoin(ctx sdk.Context, coin sdk.Coin) bool
	BlockedAddr(addr sdk.AccAddress) bool
	GetDenomMetaData(ctx sdk.Context, denom string) banktypes.Metadata
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// EVMKeeper defines the expected EVM keeper interface used on intrarelayer
type EVMKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
	GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account
	EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
}

type GravityKeeper interface {
	AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, counterpartReceiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error)
}

type CrossChainKeeper interface {
	GetModuleName() string
	AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, receiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error)
}

type IBCTransferKeeper interface {
	GetRouter() *ibctransfertypes.Router
	Transfer(goCtx context.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error)
}

type IBCChannelKeeper interface {
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
}

// GovKeeper defines the expected governance keeper interface used on intrarelayer
type GovKeeper interface {
	Logger(sdk.Context) log.Logger
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams
	GetProposal(ctx sdk.Context, proposalID uint64) (govtypes.Proposal, bool)
	InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, timestamp time.Time)
	RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, timestamp time.Time)
	SetProposal(ctx sdk.Context, proposal govtypes.Proposal)
}
