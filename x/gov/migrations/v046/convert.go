package v046

import (
	"encoding/base64"
	"encoding/json"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	fxgovtypes "github.com/functionx/fx-core/v7/x/gov/types"
)

func convertToNewProposal(oldProp v1beta1.Proposal) (v1.Proposal, error) {
	msg, err := v1.NewLegacyContent(oldProp.GetContent(), authtypes.NewModuleAddress(govtypes.ModuleName).String())
	if err != nil {
		return v1.Proposal{}, err
	}
	msgAny, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return v1.Proposal{}, err
	}

	fxMetadata := fxgovtypes.FXMetadata{
		Title:   oldProp.GetContent().GetTitle(),
		Summary: oldProp.GetContent().GetDescription(),
	}
	mdBytes, err := json.Marshal(fxMetadata)
	if err != nil {
		return v1.Proposal{}, errortypes.ErrInvalidRequest.Wrapf("proposal metadata: %s", err.Error())
	}

	return v1.Proposal{
		Id:       oldProp.ProposalId,
		Messages: []*codectypes.Any{msgAny},
		Status:   v1.ProposalStatus(oldProp.Status),
		FinalTallyResult: &v1.TallyResult{
			YesCount:        oldProp.FinalTallyResult.Yes.String(),
			NoCount:         oldProp.FinalTallyResult.No.String(),
			AbstainCount:    oldProp.FinalTallyResult.Abstain.String(),
			NoWithVetoCount: oldProp.FinalTallyResult.NoWithVeto.String(),
		},
		SubmitTime:      &oldProp.SubmitTime,
		DepositEndTime:  &oldProp.DepositEndTime,
		TotalDeposit:    oldProp.TotalDeposit,
		VotingStartTime: &oldProp.VotingStartTime,
		VotingEndTime:   &oldProp.VotingEndTime,
		Metadata:        base64.StdEncoding.EncodeToString(mdBytes),
	}, nil
}
