package types

import (
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// constants
const (
	// Deprecated: ProposalTypeRegisterCoin Do not use.
	ProposalTypeRegisterCoin string = "RegisterCoin"
	// Deprecated: ProposalTypeRegisterERC20 Do not use.
	ProposalTypeRegisterERC20 string = "RegisterERC20"
	// Deprecated: ProposalTypeToggleTokenConversion Do not use.
	ProposalTypeToggleTokenConversion string = "ToggleTokenConversion" // #nosec G101
	// Deprecated: ProposalTypeUpdateDenomAlias Do not use.
	ProposalTypeUpdateDenomAlias string = "UpdateDenomAlias"
)

// Implements Proposal Interface
var (
	_ govv1betal.Content = &RegisterCoinProposal{}
	_ govv1betal.Content = &RegisterERC20Proposal{}
	_ govv1betal.Content = &ToggleTokenConversionProposal{}
	_ govv1betal.Content = &UpdateDenomAliasProposal{}
)

func init() {
	govv1betal.RegisterProposalType(ProposalTypeRegisterCoin)
	govv1betal.RegisterProposalType(ProposalTypeRegisterERC20)
	govv1betal.RegisterProposalType(ProposalTypeToggleTokenConversion)
	govv1betal.RegisterProposalType(ProposalTypeUpdateDenomAlias)
}

// ProposalRoute returns router key for this proposal
func (*RegisterCoinProposal) ProposalRoute() string { return ModuleName }

// ProposalType returns proposal type for this proposal
func (*RegisterCoinProposal) ProposalType() string {
	return ProposalTypeRegisterCoin
}

// ValidateBasic performs a stateless check of the proposal fields
func (m *RegisterCoinProposal) ValidateBasic() error {
	return nil
}

// ProposalRoute returns router key for this proposal
func (*RegisterERC20Proposal) ProposalRoute() string { return ModuleName }

// ProposalType returns proposal type for this proposal
func (*RegisterERC20Proposal) ProposalType() string {
	return ProposalTypeRegisterERC20
}

// ValidateBasic performs a stateless check of the proposal fields
func (m *RegisterERC20Proposal) ValidateBasic() error {
	return nil
}

// ProposalRoute returns router key for this proposal
func (*ToggleTokenConversionProposal) ProposalRoute() string { return ModuleName }

// ProposalType returns proposal type for this proposal
func (*ToggleTokenConversionProposal) ProposalType() string {
	return ProposalTypeToggleTokenConversion
}

// ValidateBasic performs a stateless check of the proposal fields
func (m *ToggleTokenConversionProposal) ValidateBasic() error {
	return nil
}

// ProposalRoute returns router key for this proposal
func (*UpdateDenomAliasProposal) ProposalRoute() string { return ModuleName }

// ProposalType returns proposal type for this proposal
func (*UpdateDenomAliasProposal) ProposalType() string {
	return ProposalTypeUpdateDenomAlias
}

// ValidateBasic performs a stateless check of the proposal fields
func (m *UpdateDenomAliasProposal) ValidateBasic() error {
	return nil
}
