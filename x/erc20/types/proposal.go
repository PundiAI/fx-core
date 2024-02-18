package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

// constants
const (
	ProposalTypeRegisterCoin          string = "RegisterCoin"
	ProposalTypeRegisterERC20         string = "RegisterERC20"
	ProposalTypeToggleTokenConversion string = "ToggleTokenConversion" // #nosec G101
	ProposalTypeUpdateDenomAlias      string = "UpdateDenomAlias"
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

// CreateDenomDescription generates a string with the coin description
func CreateDenomDescription(address string) string {
	return fmt.Sprintf("Function X coin token representation of %s", address)
}

// NewRegisterCoinProposal returns new instance of RegisterCoinProposal
func NewRegisterCoinProposal(title, description string, coinMetadata banktypes.Metadata) govv1betal.Content {
	return &RegisterCoinProposal{
		Title:       title,
		Description: description,
		Metadata:    coinMetadata,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterCoinProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterCoinProposal) ProposalType() string {
	return ProposalTypeRegisterCoin
}

// ValidateBasic performs a stateless check of the proposal fields
func (m *RegisterCoinProposal) ValidateBasic() error {
	if err := m.Metadata.Validate(); err != nil {
		return errortypes.ErrInvalidRequest.Wrapf("invalid metadata: %s", err.Error())
	}

	if err := fxtypes.ValidateMetadata(m.Metadata); err != nil {
		return errortypes.ErrInvalidRequest.Wrapf("invalid metadata: %s", err.Error())
	}

	if err := ibctransfertypes.ValidateIBCDenom(m.Metadata.Base); err != nil {
		return errortypes.ErrInvalidRequest.Wrapf("invalid metadata base: %s", err.Error())
	}

	return govv1betal.ValidateAbstract(m)
}

// NewRegisterERC20Proposal returns new instance of RegisterERC20Proposal
func NewRegisterERC20Proposal(title, description, erc20Addr string, aliases []string) govv1betal.Content {
	return &RegisterERC20Proposal{
		Title:        title,
		Description:  description,
		Erc20Address: erc20Addr,
		Aliases:      aliases,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterERC20Proposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterERC20Proposal) ProposalType() string {
	return ProposalTypeRegisterERC20
}

// ValidateBasic performs a stateless check of the proposal fields
func (m *RegisterERC20Proposal) ValidateBasic() error {
	if err := fxtypes.ValidateEthereumAddress(m.Erc20Address); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid ERC20 address: %s", err.Error())
	}
	seenAliases := make(map[string]bool)
	for _, alias := range m.Aliases {
		if seenAliases[alias] {
			return errortypes.ErrInvalidAddress.Wrapf("duplicate denomination unit alias %s", alias)
		}
		if strings.TrimSpace(alias) == "" {
			return errortypes.ErrInvalidAddress.Wrapf("alias for denom unit %s cannot be blank", alias)
		}
		if err := sdk.ValidateDenom(alias); err != nil {
			return errortypes.ErrInvalidRequest.Wrap("invalid alias")
		}
		seenAliases[alias] = true
	}
	return govv1betal.ValidateAbstract(m)
}

// NewToggleTokenConversionProposal returns new instance of ToggleTokenConversionProposal
func NewToggleTokenConversionProposal(title, description string, token string) govv1betal.Content {
	return &ToggleTokenConversionProposal{
		Title:       title,
		Description: description,
		Token:       token,
	}
}

// ProposalRoute returns router key for this proposal
func (*ToggleTokenConversionProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*ToggleTokenConversionProposal) ProposalType() string {
	return ProposalTypeToggleTokenConversion
}

// ValidateBasic performs a stateless check of the proposal fields
func (m *ToggleTokenConversionProposal) ValidateBasic() error {
	// check if the token is a hex address, if not, check if it is a valid SDK
	// denom
	if err := fxtypes.ValidateEthereumAddress(m.Token); err != nil {
		if err := sdk.ValidateDenom(m.Token); err != nil {
			return errortypes.ErrInvalidRequest.Wrap("invalid token")
		}
	}

	return govv1betal.ValidateAbstract(m)
}

// NewUpdateDenomAliasProposal returns new instance of UpdateDenomAliasProposal
func NewUpdateDenomAliasProposal(title, description string, denom, alias string) govv1betal.Content {
	return &UpdateDenomAliasProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
		Alias:       alias,
	}
}

// ProposalRoute returns router key for this proposal
func (*UpdateDenomAliasProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*UpdateDenomAliasProposal) ProposalType() string {
	return ProposalTypeUpdateDenomAlias
}

// ValidateBasic performs a stateless check of the proposal fields
func (m *UpdateDenomAliasProposal) ValidateBasic() error {
	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return errortypes.ErrInvalidRequest.Wrapf("invalid denom: %s", err.Error())
	}
	if err := sdk.ValidateDenom(m.Alias); err != nil {
		return errortypes.ErrInvalidRequest.Wrapf("invalid alias: %s", err.Error())
	}
	return govv1betal.ValidateAbstract(m)
}
