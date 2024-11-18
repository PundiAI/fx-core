package types

import (
	"context"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type CrosschainKeeper interface {
	ExecuteClaim(ctx sdk.Context, eventNonce uint64) (error, error)
	BridgeCoinSupply(ctx context.Context, token, target string) (sdk.Coin, error)
	CrosschainBaseCoin(ctx sdk.Context, from sdk.AccAddress, receipt string, amount, fee sdk.Coin, fxTarget *crosschaintypes.FxTarget, memo string, originToken bool) error
	BridgeCallBaseCoin(ctx sdk.Context, from, refund, to common.Address, coins sdk.Coins, data, memo []byte, quoteId *big.Int, fxTarget *crosschaintypes.FxTarget, originTokenAmount sdkmath.Int) (uint64, error)
	GetBaseDenomByErc20(ctx sdk.Context, erc20Addr common.Address) (erc20types.ERC20Token, error)

	HasOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string) bool
	GetOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string) (sdk.AccAddress, bool)
	GetOracle(ctx sdk.Context, oracleAddr sdk.AccAddress) (oracle crosschaintypes.Oracle, found bool)
}

type GovKeeper interface {
	CheckDisabledPrecompiles(ctx sdk.Context, contractAddress common.Address, methodId []byte) error
}

type StakingKeeper interface {
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, err error)
	SetDelegation(ctx context.Context, delegation stakingtypes.Delegation) error
	RemoveDelegation(ctx context.Context, delegation stakingtypes.Delegation) error
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	Delegate(
		ctx context.Context, delAddr sdk.AccAddress, bondAmt sdkmath.Int, tokenSrc stakingtypes.BondStatus,
		validator stakingtypes.Validator, subtractAccount bool,
	) (newShares sdkmath.LegacyDec, err error)
	HasMaxUnbondingDelegationEntries(ctx context.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) (bool, error)
	Unbond(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) (amount sdkmath.Int, err error)
	UnbondingTime(ctx context.Context) (time.Duration, error)
	SetUnbondingDelegationEntry(
		ctx context.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
		creationHeight int64, minTime time.Time, balance sdkmath.Int,
	) (stakingtypes.UnbondingDelegation, error)
	InsertUBDQueue(ctx context.Context, ubd stakingtypes.UnbondingDelegation, completionTime time.Time) error
	GetAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress) *big.Int
	SetAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress, shares *big.Int)
	HasReceivingRedelegation(ctx context.Context, delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) (bool, error)
	BeginRedelegation(
		ctx context.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdkmath.LegacyDec,
	) (completionTime time.Time, err error)
	GetLastValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
}

type DistrKeeper interface {
	GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error)
	IncrementValidatorPeriod(ctx context.Context, val stakingtypes.ValidatorI) (uint64, error)
	CalculateDelegationRewards(ctx context.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins, err error)
	GetDelegatorStartingInfo(ctx context.Context, val sdk.ValAddress, del sdk.AccAddress) (period distrtypes.DelegatorStartingInfo, err error)
	SetDelegatorStartingInfo(ctx context.Context, val sdk.ValAddress, del sdk.AccAddress, period distrtypes.DelegatorStartingInfo) error
	DeleteDelegatorStartingInfo(ctx context.Context, val sdk.ValAddress, del sdk.AccAddress) error
	GetValidatorCurrentRewards(ctx context.Context, val sdk.ValAddress) (rewards distrtypes.ValidatorCurrentRewards, err error)
	GetValidatorHistoricalRewards(ctx context.Context, val sdk.ValAddress, period uint64) (rewards distrtypes.ValidatorHistoricalRewards, err error)
	SetValidatorHistoricalRewards(ctx context.Context, val sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) error
	DeleteValidatorHistoricalReward(ctx context.Context, val sdk.ValAddress, period uint64) error
}

type EvmKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
}

type SlashingKeeper interface {
	GetValidatorSigningInfo(ctx context.Context, address sdk.ConsAddress) (slashingtypes.ValidatorSigningInfo, error)
}
