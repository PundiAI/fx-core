package types

import (
	"context"
	"time"

	addresscodec "cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the expected staking keeper methods
type StakingKeeper interface {
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	GetDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation, err error)
	GetUnbondingDelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) (unbondingDelegations []stakingtypes.UnbondingDelegation, err error)
	GetRedelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) (redelegations []stakingtypes.Redelegation, err error)
	GetUBDQueueTimeSlice(ctx context.Context, timestamp time.Time) (dvPairs []stakingtypes.DVPair, err error)
	GetRedelegationQueueTimeSlice(ctx context.Context, timestamp time.Time) (dvvTriplets []stakingtypes.DVVTriplet, err error)
	ValidatorAddressCodec() addresscodec.Codec
}

// AccountKeeper defines the expected account keeper methods
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	IterateAccounts(ctx context.Context, cb func(account sdk.AccountI) (stop bool))
}

// BankKeeper defines the expected bank keeper methods
type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type GovKeeper interface {
	IteratorInactiveProposal(ctx sdk.Context, t time.Time, fn func(proposal v1.Proposal) (bool, error)) error
	IteratorActiveProposal(ctx sdk.Context, t time.Time, fn func(proposal v1.Proposal) (bool, error)) error
	HasDeposit(ctx sdk.Context, proposalId uint64, depositor sdk.AccAddress) (bool, error)
	HasVote(ctx sdk.Context, proposalId uint64, voter sdk.AccAddress) (bool, error)
}
