package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeInitCrossChainParams defines the type for a InitCrossChainParamsProposal
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
	govtypes.RegisterProposalType(ProposalTypeUpdateChainOracles)
}

// Proposal handler

func (m *InitCrossChainParamsProposal) GetTitle() string {
	return m.Title
}

func (m *InitCrossChainParamsProposal) GetDescription() string {
	return m.Description
}

func (m *InitCrossChainParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *InitCrossChainParamsProposal) ProposalType() string {
	return ProposalTypeInitCrossChainParams
}

func (m *InitCrossChainParamsProposal) ValidateBasic() error {
	if err := ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	if len(m.Params.Oracles) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "oracles cannot be empty")
	}
	return nil
}

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

func (m *UpdateChainOraclesProposal) GetTitle() string {
	return m.Title
}

func (m *UpdateChainOraclesProposal) GetDescription() string {
	return m.Description
}

func (m *UpdateChainOraclesProposal) ProposalRoute() string {
	return RouterKey
}

func (m *UpdateChainOraclesProposal) ProposalType() string {
	return ProposalTypeUpdateChainOracles
}

func (m *UpdateChainOraclesProposal) ValidateBasic() error {
	if err := ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.Oracles) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "oracles cannot be empty")
	}

	oraclesMap := make(map[string]bool)
	for _, addr := range m.Oracles {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return sdkerrors.Wrap(ErrOracleAddress, addr)
		}
		if oraclesMap[addr] {
			return sdkerrors.Wrap(ErrInvalid, fmt.Sprintf("duplicate oracle %s", addr))
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
