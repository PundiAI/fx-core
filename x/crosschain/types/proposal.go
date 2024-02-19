package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	// Deprecated: ProposalTypeInitCrossChainParams
	ProposalTypeInitCrossChainParams = "InitCrossChainParams"
	// ProposalTypeUpdateChainOracles defines the type for a UpdateChainOraclesProposal
	ProposalTypeUpdateChainOracles = "UpdateChainOracles"
)

var (
	_ govv1betal.Content = &InitCrossChainParamsProposal{}
	_ govv1betal.Content = &UpdateChainOraclesProposal{}
)

func init() {
	govv1betal.RegisterProposalType(ProposalTypeInitCrossChainParams)
	govv1betal.RegisterProposalType(ProposalTypeUpdateChainOracles)
}

func (m *InitCrossChainParamsProposal) GetTitle() string { return m.Title }

func (m *InitCrossChainParamsProposal) GetDescription() string { return m.Description }

func (m *InitCrossChainParamsProposal) ProposalRoute() string { return RouterKey }

func (m *InitCrossChainParamsProposal) ProposalType() string { return ProposalTypeInitCrossChainParams }

func (m *InitCrossChainParamsProposal) ValidateBasic() error { return nil }

func (m *InitCrossChainParamsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Init CrossChain Params Proposal:
  Title:       %s
  Description: %s
  ChainName: %s
  Params: %v
`, m.Title, m.Description, m.ChainName, m.Params))
	return b.String()
}

func (m *UpdateChainOraclesProposal) GetTitle() string { return m.Title }

func (m *UpdateChainOraclesProposal) GetDescription() string { return m.Description }

func (m *UpdateChainOraclesProposal) ProposalRoute() string { return RouterKey }

func (m *UpdateChainOraclesProposal) ProposalType() string {
	return ProposalTypeUpdateChainOracles
}

func (m *UpdateChainOraclesProposal) ValidateBasic() error {
	if _, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if err := govv1betal.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.Oracles) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty oracles")
	}

	oraclesMap := make(map[string]bool)
	for _, addr := range m.Oracles {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errortypes.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err.Error())
		}
		if oraclesMap[addr] {
			return errortypes.ErrInvalidAddress.Wrapf("duplicate oracle address: %s", addr)
		}
		oraclesMap[addr] = true
	}
	return nil
}

func (m *UpdateChainOraclesProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Chain Oracles Proposal:
  Title:       %s
  Description: %s
  ChainName: %s
  Oracles: %v
`, m.Title, m.Description, m.ChainName, m.Oracles))
	return b.String()
}
