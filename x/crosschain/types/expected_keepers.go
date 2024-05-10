package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tranfsertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	GetUnbondingDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd stakingtypes.UnbondingDelegation, found bool)

	RemoveDelegation(ctx sdk.Context, delegation stakingtypes.Delegation) error
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator

	BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error
}

type StakingMsgServer interface {
	Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error)
	BeginRedelegate(goCtx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error)
	Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error)
}

type DistributionKeeper interface {
	GetDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) (period distributiontypes.DelegatorStartingInfo)
}

type DistributionMsgServer interface {
	WithdrawDelegatorReward(goCtx context.Context, msg *distributiontypes.MsgWithdrawDelegatorReward) (*distributiontypes.MsgWithdrawDelegatorRewardResponse, error)
}

// BankKeeper defines the expected bank keeper methods
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	HasDenomMetaData(ctx sdk.Context, denom string) bool
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	IterateAllDenomMetaData(ctx sdk.Context, cb func(bank.Metadata) bool)
}

type Erc20Keeper interface {
	TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, coin, fee sdk.Coin, _, _ bool) error
	ConvertCoin(goCtx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
	ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error)
	HookOutgoingRefund(ctx sdk.Context, moduleName string, txID uint64, sender sdk.AccAddress, totalCoin sdk.Coin) error
	SetOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64)
	HasOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64) bool
	DeleteOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64)
	IsOriginOrConvertedDenom(ctx sdk.Context, denom string) bool
	ToTargetDenom(ctx sdk.Context, denom, base string, aliases []string, fxTarget fxtypes.FxTarget) string
}

// EVMKeeper defines the expected EVM keeper interface used on crosschain
type EVMKeeper interface {
	CallEVM(ctx sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte, commit bool) (*types.MsgEthereumTxResponse, error)
}

type IBCTransferKeeper interface {
	Transfer(goCtx context.Context, msg *tranfsertypes.MsgTransfer) (*tranfsertypes.MsgTransferResponse, error)

	SetDenomTrace(ctx sdk.Context, denomTrace tranfsertypes.DenomTrace)
}

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

type (
	ParamSet = paramtypes.ParamSet
	// Subspace defines an interface that implements the legacy x/params Subspace
	// type.
	//
	// NOTE: This is used solely for migration of x/params managed parameters.
	Subspace interface {
		GetParamSet(ctx sdk.Context, ps ParamSet)
		HasKeyTable() bool
		WithKeyTable(table paramtypes.KeyTable) paramtypes.Subspace
	}
)
