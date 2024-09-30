package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/x/migrate/types"
)

type GovMigrate struct {
	govKeeper     types.GovKeeper
	accountKeeper govtypes.AccountKeeper
}

func NewGovMigrate(govKeeper types.GovKeeper, accountKeeper govtypes.AccountKeeper) MigrateI {
	return &GovMigrate{
		govKeeper:     govKeeper,
		accountKeeper: accountKeeper,
	}
}

func (m *GovMigrate) Validate(ctx sdk.Context, _ codec.BinaryCodec, from sdk.AccAddress, to common.Address) error {
	if err := m.govKeeper.IteratorInactiveProposal(ctx, ctx.BlockTime(), m.DepositPeriodCallback(ctx, from, to)); err != nil {
		return err
	}
	return m.govKeeper.IteratorActiveProposal(ctx, ctx.BlockTime(), m.VotePeriodCallback(ctx, from, to))
}

func (m *GovMigrate) Execute(_ sdk.Context, _ codec.BinaryCodec, _ sdk.AccAddress, _ common.Address) error {
	return nil
}

func (m *GovMigrate) DepositPeriodCallback(ctx sdk.Context, from sdk.AccAddress, to common.Address) func(proposal govv1.Proposal) (bool, error) {
	return func(proposal govv1.Proposal) (bool, error) {
		proposer, err := m.accountKeeper.AddressCodec().StringToBytes(proposal.GetProposer())
		if err != nil {
			return false, err
		}

		if from.Equals(sdk.AccAddress(proposer)) {
			return false, sdkerrors.ErrInvalidRequest.Wrapf("can not migrate, %s is proposer of %d", from.String(), proposal.Id)
		}
		if sdk.AccAddress(to.Bytes()).Equals(sdk.AccAddress(proposer)) {
			return false, sdkerrors.ErrInvalidRequest.Wrapf("can not migrate, %s is proposer of %d", to.String(), proposal.Id)
		}
		hasDeposit, err := m.govKeeper.HasDeposit(ctx, proposal.Id, from)
		if err != nil {
			return false, err
		}
		if hasDeposit {
			return false, sdkerrors.ErrInvalidRequest.Wrapf("can not migrate, %s have deposit of proposal %d", from.String(), proposal.Id)
		}

		hasDeposit, err = m.govKeeper.HasDeposit(ctx, proposal.Id, to.Bytes())
		if err != nil {
			return false, err
		}
		if hasDeposit {
			return false, sdkerrors.ErrInvalidRequest.Wrapf("can not migrate, %s have deposit of proposal %d", to.String(), proposal.Id)
		}
		return false, nil
	}
}

func (m *GovMigrate) VotePeriodCallback(ctx sdk.Context, from sdk.AccAddress, to common.Address) func(proposal govv1.Proposal) (bool, error) {
	return func(proposal govv1.Proposal) (bool, error) {
		b, err := m.DepositPeriodCallback(ctx, from, to)(proposal)
		if err != nil {
			return b, err
		}
		hasVote, err := m.govKeeper.HasVote(ctx, proposal.Id, from)
		if err != nil {
			return false, err
		}
		if hasVote {
			return false, sdkerrors.ErrInvalidRequest.Wrapf("can not migrate, %s have vote of proposal %d", from.String(), proposal.Id)
		}

		hasVote, err = m.govKeeper.HasVote(ctx, proposal.Id, to.Bytes())
		if err != nil {
			return false, err
		}
		if hasVote {
			return false, sdkerrors.ErrInvalidRequest.Wrapf("can not migrate, %s have vote of proposal %d", to.String(), proposal.Id)
		}
		return false, nil
	}
}
