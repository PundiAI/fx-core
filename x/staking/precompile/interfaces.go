package precompile

import (
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, recipientAddr sdk.AccAddress, senderModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type StakingKeeper interface {
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	SetDelegation(ctx sdk.Context, delegation stakingtypes.Delegation)
	RemoveDelegation(ctx sdk.Context, delegation stakingtypes.Delegation) error
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdkmath.Int, tokenSrc stakingtypes.BondStatus,
		validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error)
	HasMaxUnbondingDelegationEntries(ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) bool
	Unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) (amount sdkmath.Int, err error)
	UnbondingTime(ctx sdk.Context) (res time.Duration)
	SetUnbondingDelegationEntry(ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
		creationHeight int64, minTime time.Time, balance sdkmath.Int) stakingtypes.UnbondingDelegation
	InsertUBDQueue(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation, completionTime time.Time)
	GetAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress) *big.Int
	SetAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress, shares *big.Int)
	HasReceivingRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) bool
	BeginRedelegation(
		ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdk.Dec,
	) (completionTime time.Time, err error)
	GetLastValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
}

type DistrKeeper interface {
	GetDelegatorWithdrawAddr(ctx sdk.Context, delAddr sdk.AccAddress) sdk.AccAddress
	IncrementValidatorPeriod(ctx sdk.Context, val stakingtypes.ValidatorI) uint64
	CalculateDelegationRewards(ctx sdk.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins)
	GetDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) (period distrtypes.DelegatorStartingInfo)
	SetDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress, period distrtypes.DelegatorStartingInfo)
	DeleteDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress)
	GetValidatorCurrentRewards(ctx sdk.Context, val sdk.ValAddress) (rewards distrtypes.ValidatorCurrentRewards)
	GetValidatorHistoricalRewards(ctx sdk.Context, val sdk.ValAddress, period uint64) (rewards distrtypes.ValidatorHistoricalRewards)
	SetValidatorHistoricalRewards(ctx sdk.Context, val sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards)
	DeleteValidatorHistoricalReward(ctx sdk.Context, val sdk.ValAddress, period uint64)
}

type EvmKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
}

type GovKeeper interface {
	CheckDisabledPrecompiles(ctx sdk.Context, contractAddress common.Address, methodId []byte) error
}

type SlashingKeeper interface {
	GetValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress) (info slashingtypes.ValidatorSigningInfo, found bool)
}
