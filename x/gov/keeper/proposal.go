package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	fxgovtypes "github.com/functionx/fx-core/v7/x/gov/types"
)

// SubmitProposal creates a new proposal given an array of messages
//
//gocyclo:ignore
func (keeper Keeper) SubmitProposal(ctx sdk.Context, messages []sdk.Msg, fxMetadata string) (v1.Proposal, error) {
	// TODO proposal metadata contain title, summary and metadata, for compatibility with cosmos v0.46.x,
	//  when upgrade to v0.47.x, will migrate to new proposal struct
	fxMD, err := fxgovtypes.ParseFXMetadata(fxMetadata)
	if err != nil {
		return v1.Proposal{}, errortypes.ErrInvalidRequest.Wrapf("invalid fx metadata content: %s", err)
	}
	if err := keeper.AssertFXMetadata(fxMD); err != nil {
		if err != nil {
			return v1.Proposal{}, err
		}
	}

	// Will hold a comma-separated string of all Msg type URLs.
	msgsStr := ""

	// record MsgTypeURL messages Is it the same
	msgType := ""
	// Loop through all messages and confirm that each has a handler and the gov module account
	// as the only signer
	for _, msg := range messages {
		msgsStr += fmt.Sprintf(",%s", sdk.MsgTypeURL(msg))

		if msgType != "" && !strings.EqualFold(msgType, sdk.MsgTypeURL(msg)) {
			return v1.Proposal{}, errorsmod.Wrap(types.ErrInvalidProposalContent, "proposal MsgTypeURL is different")
		}
		msgType = sdk.MsgTypeURL(msg)
		// perform a basic validation of the message
		if err := msg.ValidateBasic(); err != nil {
			return v1.Proposal{}, errorsmod.Wrap(types.ErrInvalidProposalMsg, err.Error())
		}

		signers := msg.GetSigners()
		if len(signers) != 1 {
			return v1.Proposal{}, types.ErrInvalidSigner
		}

		// assert that the governance module account is the only signer of the messages
		if !signers[0].Equals(keeper.GetGovernanceAccount(ctx).GetAddress()) {
			return v1.Proposal{}, errorsmod.Wrapf(types.ErrInvalidSigner, signers[0].String())
		}

		// use the msg service router to see that there is a valid route for that message.
		handler := keeper.Router().Handler(msg)
		if handler == nil {
			return v1.Proposal{}, errorsmod.Wrap(types.ErrUnroutableProposalMsg, sdk.MsgTypeURL(msg))
		}

		// Only if it's a MsgExecLegacyContent do we try to execute the
		// proposal in a cached context.
		// For other Msgs, we do not verify the proposal messages any further.
		// They may fail upon execution.
		// ref: https://github.com/cosmos/cosmos-sdk/pull/10868#discussion_r784872842
		if msg, ok := msg.(*v1.MsgExecLegacyContent); ok {
			cacheCtx, _ := ctx.CacheContext()
			if _, err := handler(cacheCtx, msg); err != nil {
				return v1.Proposal{}, errorsmod.Wrap(types.ErrNoProposalHandlerExists, err.Error())
			}
		}
	}

	proposalID, err := keeper.GetProposalID(ctx)
	if err != nil {
		return v1.Proposal{}, err
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := keeper.GetMaxDepositPeriod(ctx, msgType)

	proposal, err := v1.NewProposal(messages, proposalID, fxMetadata, submitTime, submitTime.Add(*depositPeriod))
	if err != nil {
		return v1.Proposal{}, err
	}

	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposalID, *proposal.DepositEndTime)
	keeper.SetProposalID(ctx, proposalID+1)

	// called right after a proposal is submitted
	keeper.AfterProposalSubmission(ctx, proposalID)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposalMessages, msgsStr),
		),
	)

	return proposal, nil
}

// AssertFXMetadata returns an error if given metadata invalid
func (keeper Keeper) AssertFXMetadata(pm fxgovtypes.FXMetadata) error {
	if len(strings.TrimSpace(pm.Title)) == 0 {
		return errorsmod.Wrap(types.ErrInvalidProposalContent, "proposal title cannot be blank")
	}
	if uint64(len(pm.Title)) > keeper.config.MaxTitleLen {
		return errorsmod.Wrapf(types.ErrInvalidProposalContent, "proposal title is longer than max length of %d", keeper.config.MaxTitleLen)
	}

	if len(strings.TrimSpace(pm.Summary)) == 0 {
		return errorsmod.Wrap(types.ErrInvalidProposalContent, "proposal summary cannot be blank")
	}
	if uint64(len(pm.Summary)) > keeper.config.MaxSummaryLen {
		return errorsmod.Wrapf(types.ErrInvalidProposalContent, "proposal summary is longer than max length of %d", keeper.config.MaxSummaryLen)
	}

	if pm.Metadata != "" && uint64(len(pm.Metadata)) > keeper.config.MaxMetadataLen {
		return types.ErrInvalidProposalContent.Wrapf("proposal metadata is longer than max length of %d", keeper.config.MaxMetadataLen)
	}
	return nil
}
