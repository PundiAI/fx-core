package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ethermint "github.com/functionx/fx-core/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

// constants
const (
	ProposalTypeInitIntrarelayerParams string = "InitIntrarelayerParams"
	ProposalTypeRegisterCoin           string = "RegisterCoin"
	ProposalTypeRegisterFIP20          string = "RegisterFIP20"
	ProposalTypeToggleTokenRelay       string = "ToggleTokenRelay" // #nosec
	ProposalTypeUpdateTokenPairFIP20   string = "UpdateTokenPairFIP20"
)

// Implements Proposal Interface
var (
	_ govtypes.Content = &InitIntrarelayerParamsProposal{}
	_ govtypes.Content = &RegisterCoinProposal{}
	_ govtypes.Content = &RegisterFIP20Proposal{}
	_ govtypes.Content = &ToggleTokenRelayProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeInitIntrarelayerParams)
	govtypes.RegisterProposalType(ProposalTypeRegisterCoin)
	govtypes.RegisterProposalType(ProposalTypeRegisterFIP20)
	govtypes.RegisterProposalType(ProposalTypeToggleTokenRelay)
	govtypes.RegisterProposalType(ProposalTypeUpdateTokenPairFIP20)
	govtypes.RegisterProposalTypeCodec(&InitIntrarelayerParamsProposal{}, "intrarelayer/InitIntrarelayerParamsProposal")
	govtypes.RegisterProposalTypeCodec(&RegisterCoinProposal{}, "intrarelayer/RegisterCoinProposal")
	govtypes.RegisterProposalTypeCodec(&RegisterFIP20Proposal{}, "intrarelayer/RegisterFIP20Proposal")
	govtypes.RegisterProposalTypeCodec(&ToggleTokenRelayProposal{}, "intrarelayer/ToggleTokenRelayProposal")
}

// CreateDenom generates a string the module name plus the address to avoid conflicts with names staring with a number
func CreateDenom(address string) string {
	return fmt.Sprintf("%s/%s", ModuleName, address)
}

// NewInitIntrarelayerParamsProposal returns new instance of InitIntrarelayerParamsProposal
func NewInitIntrarelayerParamsProposal(title, description string, params *Params, metadata []banktypes.Metadata) govtypes.Content {
	return &InitIntrarelayerParamsProposal{
		Title:       title,
		Description: description,
		Params:      params,
		Metadata:    metadata,
	}
}

// ProposalRoute returns router key for this proposal
func (*InitIntrarelayerParamsProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns proposal type for this proposal
func (*InitIntrarelayerParamsProposal) ProposalType() string {
	return ProposalTypeInitIntrarelayerParams
}

// ValidateBasic performs a stateless check of the proposal fields
func (iip *InitIntrarelayerParamsProposal) ValidateBasic() error {

	if iip.Params == nil {
		return errors.New("empty params")
	}
	if err := iip.Params.Validate(); err != nil {
		return err
	}

	if len(iip.Metadata) > 0 {
		for _, metadata := range iip.Metadata {
			if err := metadata.Validate(); err != nil {
				return err
			}

			if err := ibctransfertypes.ValidateIBCDenom(metadata.Base); err != nil {
				return err
			}

			if err := validateIBC(metadata); err != nil {
				return err
			}
		}
	}

	return govtypes.ValidateAbstract(iip)
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
func (rtbp *RegisterCoinProposal) ValidateBasic() error {
	if err := rtbp.Metadata.Validate(); err != nil {
		return err
	}

	if err := ibctransfertypes.ValidateIBCDenom(rtbp.Metadata.Base); err != nil {
		return err
	}

	if err := validateIBC(rtbp.Metadata); err != nil {
		return err
	}

	return govtypes.ValidateAbstract(rtbp)
}

func validateIBC(metadata banktypes.Metadata) error {
	// Check ibc/ denom
	denomSplit := strings.SplitN(metadata.Base, "/", 2)

	if denomSplit[0] == metadata.Base && strings.TrimSpace(metadata.Base) != "" {
		// Not IBC
		return nil
	}

	if len(denomSplit) != 2 || denomSplit[0] != ibctransfertypes.DenomPrefix {
		// NOTE: should be unaccessible (covered on ValidateIBCDenom)
		return fmt.Errorf("invalid metadata. %s denomination should be prefixed with the format 'ibc/", metadata.Base)
	}
	return nil
}

// ValidateIntrarelayerDenom checks if a denom is a valid intrarelayer/
// denomination
func ValidateIntrarelayerDenom(denom string) error {
	denomSplit := strings.SplitN(denom, "/", 2)

	if len(denomSplit) != 2 || denomSplit[0] != ModuleName {
		return fmt.Errorf("invalid denom. %s denomination should be prefixed with the format 'intrarelayer/", denom)
	}

	return ethermint.ValidateAddress(denomSplit[1])
}

// NewRegisterFIP20Proposal returns new instance of RegisterFIP20Proposal
func NewRegisterFIP20Proposal(title, description, fip20Addr string) govtypes.Content {
	return &RegisterFIP20Proposal{
		Title:        title,
		Description:  description,
		Fip20Address: fip20Addr,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterFIP20Proposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterFIP20Proposal) ProposalType() string {
	return ProposalTypeRegisterFIP20
}

// ValidateBasic performs a stateless check of the proposal fields
func (rtbp *RegisterFIP20Proposal) ValidateBasic() error {
	if err := ethermint.ValidateAddress(rtbp.Fip20Address); err != nil {
		return sdkerrors.Wrap(err, "FIP20 address")
	}
	return govtypes.ValidateAbstract(rtbp)
}

// NewToggleTokenRelayProposal returns new instance of ToggleTokenRelayProposal
func NewToggleTokenRelayProposal(title, description string, token string) govtypes.Content {
	return &ToggleTokenRelayProposal{
		Title:       title,
		Description: description,
		Token:       token,
	}
}

// ProposalRoute returns router key for this proposal
func (*ToggleTokenRelayProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*ToggleTokenRelayProposal) ProposalType() string {
	return ProposalTypeToggleTokenRelay
}

// ValidateBasic performs a stateless check of the proposal fields
func (etrp *ToggleTokenRelayProposal) ValidateBasic() error {
	// check if the token is a hex address, if not, check if it is a valid SDK
	// denom
	if err := ethermint.ValidateAddress(etrp.Token); err != nil {
		if err := sdk.ValidateDenom(etrp.Token); err != nil {
			return err
		}
	}

	return govtypes.ValidateAbstract(etrp)
}
