package types

import (
	context "context"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/x/evm/statedb"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
)

// AccountKeeper defines the expected interface needed to retrieve account info.
type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(sdk.Context, authtypes.AccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
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

// EVMKeeper defines the expected EVM keeper interface used on erc20
type EVMKeeper interface {
	SetParams(ctx sdk.Context, params evmtypes.Params)
	GetParams(ctx sdk.Context) evmtypes.Params
	GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account
	EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
	CreateContractWithCode(ctx sdk.Context, addr common.Address, code []byte) error
}

type FeeMarketKeeper interface {
	GetBaseFee(ctx sdk.Context) *big.Int
	GetParams(ctx sdk.Context) feemarkettypes.Params
	SetParams(ctx sdk.Context, params feemarkettypes.Params)
	SetBaseFee(ctx sdk.Context, baseFee *big.Int)
	SetBlockGasUsed(ctx sdk.Context, gas uint64)
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

// GovKeeper defines the expected governance keeper interface used on erc20
type GovKeeper interface {
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams
	GetProposal(ctx sdk.Context, proposalID uint64) (govtypes.Proposal, bool)
	InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, timestamp time.Time)
	RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, timestamp time.Time)
	SetProposal(ctx sdk.Context, proposal govtypes.Proposal)
}
