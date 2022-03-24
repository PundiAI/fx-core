package types

import (
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
	"strings"
)

const (
	// ProposalTypeInitEvm defines the type for a InitEvmProposal
	ProposalTypeInitEvm = "InitEvm"
)

var (
	_ govtypes.Content = &InitEvmProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeInitEvm)
}

// Proposal handler

// NewInitEvmProposal returns new instance of InitEvmProposal
func NewInitEvmProposal(
	title, description string,
	evmParams *Params,
	feemarketParams *feemarkettypes.Params,
	intrarelayerParams *IntrarelayerParams,
	metadata []banktypes.Metadata,
) govtypes.Content {
	return &InitEvmProposal{
		Title:              title,
		Description:        description,
		EvmParams:          evmParams,
		FeemarketParams:    feemarketParams,
		IntrarelayerParams: intrarelayerParams,
		Metadata:           metadata,
	}
}

func (m *InitEvmProposal) ProposalRoute() string {
	return RouterKey
}

func (m *InitEvmProposal) ProposalType() string {
	return ProposalTypeInitEvm
}

func (m *InitEvmProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}
	if err := m.EvmParams.Validate(); err != nil {
		return err
	}
	if err := m.FeemarketParams.Validate(); err != nil {
		return err
	}
	if err := m.IntrarelayerParams.Validate(); err != nil {
		return err
	}

	if len(m.Metadata) > 0 {
		for _, metadata := range m.Metadata {
			if err := metadata.Validate(); err != nil {
				return err
			}

			if err := ibctransfertypes.ValidateIBCDenom(metadata.Base); err != nil {
				return err
			}

			if err := ValidateIBC(metadata); err != nil {
				return err
			}
		}
	}

	return govtypes.ValidateAbstract(m)
}

func (ip IntrarelayerParams) Validate() error {
	if ip.IbcTransferTimeoutHeight == 0 {
		return fmt.Errorf("ibc transfer timeout hegith cannot be zero: %d", ip.IbcTransferTimeoutHeight)
	}
	return nil
}

func ValidateIBC(metadata banktypes.Metadata) error {
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
