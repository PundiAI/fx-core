package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// Deprecated: ProposalTypeInitCrossChainParams
	ProposalTypeInitCrossChainParams = "InitCrossChainParams"
	// ProposalTypeUpdateChainOracles defines the type for a UpdateChainOraclesProposal
	ProposalTypeUpdateChainOracles = "UpdateChainOracles"
)

var (
	_ govtypes.Content = &InitCrossChainParamsProposal{}
	_ govtypes.Content = &UpdateChainOraclesProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeInitCrossChainParams)
	govtypes.RegisterProposalTypeCodec(&InitCrossChainParamsProposal{}, "crosschain/InitCrossChainParamsProposal")
	govtypes.RegisterProposalType(ProposalTypeUpdateChainOracles)
	govtypes.RegisterProposalTypeCodec(&UpdateChainOraclesProposal{}, "crosschain/UpdateChainOraclesProposal")
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
	if err := ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid chain name")
	}
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.Oracles) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty oracles")
	}

	oraclesMap := make(map[string]bool)
	for _, addr := range m.Oracles {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
		}
		if oraclesMap[addr] {
			return ErrDuplicate.Wrapf("oracle address: %s", addr)
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
