package keeper

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/functionx/fx-core/v3/x/gov/types"
)

type msgServer struct {
	govv1.MsgServer
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(m govv1.MsgServer, k Keeper) govv1.MsgServer {
	return &msgServer{MsgServer: m, Keeper: k}
}

var _ govv1.MsgServer = msgServer{}

func (k msgServer) SubmitProposal(goCtx context.Context, msg *govv1.MsgSubmitProposal) (*govv1.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposalMsgs, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, proposalMsgs, msg.Metadata)
	if err != nil {
		return nil, err
	}

	bytes, err := proposal.Marshal()
	if err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(3*ctx.KVGasConfig().WriteCostPerByte*uint64(len(bytes)),
		"submit proposal")

	defer telemetry.IncrCounter(1, govtypes.ModuleName, "proposal")

	if sdk.NewCoins(msg.GetInitialDeposit()...).IsAllLT(types.GetInitialDeposit()) {
		return nil, errorsmod.Wrapf(types.ErrInitialAmountTooLow, "%s is smaller than %s", msg.GetInitialDeposit(), types.GetInitialDeposit())
	}

	proposer, err := sdk.AccAddressFromBech32(msg.GetProposer())
	if err != nil {
		return nil, err
	}
	votingStarted, err := k.Keeper.AddDeposit(ctx, proposal.Id, proposer, msg.GetInitialDeposit())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.GetProposer()),
	))

	if votingStarted {
		submitEvent := sdk.NewEvent(govtypes.EventTypeSubmitProposal,
			sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.Id)),
		)

		ctx.EventManager().EmitEvent(submitEvent)
	}

	return &govv1.MsgSubmitProposalResponse{
		ProposalId: proposal.Id,
	}, nil
}

func (k msgServer) Deposit(goCtx context.Context, msg *govv1.MsgDeposit) (*govv1.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	votingStarted, err := k.Keeper.AddDeposit(ctx, msg.ProposalId, accAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{govtypes.ModuleName, "deposit"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor),
	))

	if votingStarted {
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			govtypes.EventTypeProposalDeposit,
			sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalId)),
		))
	}

	return &govv1.MsgDepositResponse{}, nil
}
