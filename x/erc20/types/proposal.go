package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

// constants
const (
	ProposalTypeRegisterCoin          string = "RegisterCoin"
	ProposalTypeRegisterERC20         string = "RegisterERC20"
	ProposalTypeToggleTokenConversion string = "ToggleTokenConversion" // #nosec
	ProposalTypeUpdateDenomAlias      string = "UpdateDenomAlias"
)

// Implements Proposal Interface
var (
	_ govtypes.Content = &RegisterCoinProposal{}
	_ govtypes.Content = &RegisterERC20Proposal{}
	_ govtypes.Content = &ToggleTokenConversionProposal{}
	_ govtypes.Content = &UpdateDenomAliasProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeRegisterCoin)
	govtypes.RegisterProposalType(ProposalTypeRegisterERC20)
	govtypes.RegisterProposalType(ProposalTypeToggleTokenConversion)
	govtypes.RegisterProposalType(ProposalTypeUpdateDenomAlias)
	govtypes.RegisterProposalTypeCodec(&RegisterCoinProposal{}, "erc20/RegisterCoinProposal")
	govtypes.RegisterProposalTypeCodec(&RegisterERC20Proposal{}, "erc20/RegisterERC20Proposal")
	govtypes.RegisterProposalTypeCodec(&ToggleTokenConversionProposal{}, "erc20/ToggleTokenConversionProposal")
	govtypes.RegisterProposalTypeCodec(&UpdateDenomAliasProposal{}, "erc20/UpdateDenomAliasProposal")
}

// CreateDenomDescription generates a string with the coin description
func CreateDenomDescription(address string) string {
	return fmt.Sprintf("Function X coin token representation of %s", address)
}

// NewRegisterCoinProposal returns new instance of RegisterCoinProposal
func NewRegisterCoinProposal(title, description string, coinMetadata banktypes.Metadata) govtypes.Content {
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
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid metadata: %s", err.Error())
	}

	if err := fxtypes.ValidateMetadata(m.Metadata); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid metadata: %s", err.Error())
	}

	if err := ibctransfertypes.ValidateIBCDenom(m.Metadata.Base); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid metadata base: %s", err.Error())
	}

	return govtypes.ValidateAbstract(m)
}

// NewRegisterERC20Proposal returns new instance of RegisterERC20Proposal
func NewRegisterERC20Proposal(title, description, erc20Addr string) govtypes.Content {
	return &RegisterERC20Proposal{
		Title:        title,
		Description:  description,
		Erc20Address: erc20Addr,
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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid ERC20 address: %s", err.Error())
	}
	return govtypes.ValidateAbstract(m)
}

// NewToggleTokenConversionProposal returns new instance of ToggleTokenConversionProposal
func NewToggleTokenConversionProposal(title, description string, token string) govtypes.Content {
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
			return sdkerrors.ErrInvalidRequest.Wrap("invalid token")
		}
	}

	return govtypes.ValidateAbstract(m)
}

// NewUpdateDenomAliasProposal returns new instance of UpdateDenomAliasProposal
func NewUpdateDenomAliasProposal(title, description string, denom, alias string) govtypes.Content {
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
		return sdkerrors.ErrInvalidRequest.Wrap("invalid denom")
	}
	if err := sdk.ValidateDenom(m.Alias); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid alias")
	}
	return govtypes.ValidateAbstract(m)
}
