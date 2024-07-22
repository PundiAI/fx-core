package types

import (
	"context"

	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type MsgServerPro interface {
	MsgServer
	// SubmitProposal defines a method to create new proposal given the messages.
	SubmitProposal(context.Context, *v1.MsgSubmitProposal) (*v1.MsgSubmitProposalResponse, error)
	// ExecLegacyContent defines a Msg to be in included in a MsgSubmitProposal
	// to execute a legacy content-based proposal.
	ExecLegacyContent(context.Context, *v1.MsgExecLegacyContent) (*v1.MsgExecLegacyContentResponse, error)
	// Vote defines a method to add a vote on a specific proposal.
	Vote(context.Context, *v1.MsgVote) (*v1.MsgVoteResponse, error)
	// VoteWeighted defines a method to add a weighted vote on a specific proposal.
	VoteWeighted(context.Context, *v1.MsgVoteWeighted) (*v1.MsgVoteWeightedResponse, error)
	// Deposit defines a method to add deposit on a specific proposal.
	Deposit(context.Context, *v1.MsgDeposit) (*v1.MsgDepositResponse, error)

	UpdateParams(context.Context, *v1.MsgUpdateParams) (*v1.MsgUpdateParamsResponse, error)
}
