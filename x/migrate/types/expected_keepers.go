package types

import (
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the expected staking keeper methods
type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation)
	GetUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (unbondingDelegations []stakingtypes.UnbondingDelegation)
	GetRedelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (redelegations []stakingtypes.Redelegation)
	GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []stakingtypes.DVPair)
	GetRedelegationQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvvTriplets []stakingtypes.DVVTriplet)
}

// AccountKeeper defines the expected account keeper methods
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	IterateAccounts(ctx sdk.Context, cb func(account authtypes.AccountI) (stop bool))
}

// BankKeeper defines the expected bank keeper methods
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type GovKeeper interface {
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams

	ActiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator
	InactiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator

	GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (deposit govtypes.Deposit, found bool)
	GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (vote govtypes.Vote, found bool)
}
