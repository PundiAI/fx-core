package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeUpdateCrossChainOracles defines the type for a UpdateCrossChainOraclesProposal
	ProposalTypeUpdateCrossChainOracles = "UpdateCrossChainOracles"
)

var (
	_ govtypes.Content = &UpdateCrossChainOraclesProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateCrossChainOracles)
}

func (m *UpdateCrossChainOraclesProposal) GetTitle() string { return m.Title }

func (m *UpdateCrossChainOraclesProposal) GetDescription() string { return m.Description }

func (m *UpdateCrossChainOraclesProposal) ProposalRoute() string { return RouterKey }

func (m *UpdateCrossChainOraclesProposal) ProposalType() string {
	return ProposalTypeUpdateCrossChainOracles
}

func (m *UpdateCrossChainOraclesProposal) ValidateBasic() error {
	if err := ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.Oracles) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "oracles")
	}

	oraclesMap := make(map[string]bool)
	for _, addr := range m.Oracles {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return sdkerrors.Wrap(ErrInvalid, "oracle address")
		}
		if oraclesMap[addr] {
			return sdkerrors.Wrap(ErrDuplicate, "oracle address")
		}
		oraclesMap[addr] = true
	}
	return nil
}

func (m *UpdateCrossChainOraclesProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Chain Oracles Proposal:
  Title:       %s
  Description: %s
  ChainName: %s
  Oracles: %v
`, m.Title, m.Description, m.ChainName, m.Oracles))
	return b.String()
}
