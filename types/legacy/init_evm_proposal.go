package legacy

import (
	"fmt"
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*govv1betal.Content)(nil), &InitEvmParamsProposal{})
}

const ProposalTypeInitEvmParams = "InitEvmParams"

var _ govv1betal.Content = &InitEvmParamsProposal{}

func init() {
	govv1betal.RegisterProposalType(ProposalTypeInitEvmParams)
}

func (m *InitEvmParamsProposal) GetTitle() string {
	return m.Title
}

func (m *InitEvmParamsProposal) GetDescription() string {
	return m.Description
}

func (m *InitEvmParamsProposal) ProposalRoute() string {
	return "evm"
}

func (m *InitEvmParamsProposal) ProposalType() string {
	return ProposalTypeInitEvmParams
}

func (m *InitEvmParamsProposal) ValidateBasic() error {
	return nil
}

func (m *InitEvmParamsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Init Evm Params Proposal:
  Title:       %s
  Description: %s
  EvmParams: %v
  FeeMarketParams: %v
`, m.Title, m.Description, m.EvmParams, m.FeemarketParams))
	return b.String()
}
