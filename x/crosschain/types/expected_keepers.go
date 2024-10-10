package types

import (
	"context"
	"math/big"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tranfsertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v8/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

type StakingKeeper interface {
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, err error)
	GetUnbondingDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd stakingtypes.UnbondingDelegation, err error)
}

type StakingMsgServer interface {
	Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error)
	BeginRedelegate(goCtx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error)
	Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error)
}

type DistributionMsgServer interface {
	WithdrawDelegatorReward(goCtx context.Context, msg *distributiontypes.MsgWithdrawDelegatorReward) (*distributiontypes.MsgWithdrawDelegatorRewardResponse, error)
}

// BankKeeper defines the expected bank keeper methods
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx context.Context, denom string) sdk.Coin

	GetDenomMetaData(ctx context.Context, denom string) (bank.Metadata, bool)
	HasDenomMetaData(ctx context.Context, denom string) bool
	SetDenomMetaData(ctx context.Context, denomMetaData bank.Metadata)
	IterateAllDenomMetaData(ctx context.Context, cb func(bank.Metadata) bool)
}

type Erc20Keeper interface {
	TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, coin, fee sdk.Coin, _ bool) error
	ConvertCoin(ctx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
	ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error)
	HookOutgoingRefund(ctx sdk.Context, moduleName string, txID uint64, sender sdk.AccAddress, totalCoin sdk.Coin) error
	SetOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64)
	HasOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64) bool
	DeleteOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64)
	IsOriginOrConvertedDenom(ctx sdk.Context, denom string) bool
	ToTargetDenom(ctx sdk.Context, denom, base string, aliases []string, fxTarget fxtypes.FxTarget) string
	GetTokenPair(ctx sdk.Context, tokenOrDenom string) (erc20types.TokenPair, bool)
	IbcRefund(ctx sdk.Context, channel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error
	DeleteIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64) bool
}

// EVMKeeper defines the expected EVM keeper interface used on crosschain
type EVMKeeper interface {
	CallEVM(ctx sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte, commit bool) (*types.MsgEthereumTxResponse, error)
	IsContract(ctx sdk.Context, account common.Address) bool
}

type IBCTransferKeeper interface {
	Transfer(ctx context.Context, msg *tranfsertypes.MsgTransfer) (*tranfsertypes.MsgTransferResponse, error)
	SetDenomTrace(ctx sdk.Context, denomTrace tranfsertypes.DenomTrace)
	GetDenomTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) (tranfsertypes.DenomTrace, bool)
}

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

type BridgeTokenKeeper interface {
	HasToken(ctx context.Context, denom string) (bool, error)
	GetBridgeDenom(ctx context.Context, denom string) ([]string, error)
	GetBaseDenom(ctx context.Context, alias string) (string, error)
	GetAllTokens(ctx context.Context) ([]string, error)
	UpdateBridgeDenom(ctx context.Context, denom string, bridgeDenoms ...string) error
	SetToken(ctx context.Context, name, symbol string, decimals uint32, bridgeDenoms ...string) error
}
