// nolint:staticcheck
package erc20

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/x/erc20/keeper"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

// NewErc20ProposalHandler creates a governance handler to manage new proposal types.
// It enables RegisterTokenPairProposal to propose a registration of token mapping
// Deprecated: instead of defining gov router proposal handlers, the v0.46 gov execution models is based on sdk.Msgs
func NewErc20ProposalHandler(k keeper.Keeper) govv1betal.Handler {
	return func(ctx sdk.Context, content govv1betal.Content) error {
		switch c := content.(type) {
		case *types.RegisterCoinProposal:
			return handleRegisterCoinProposal(ctx, k, c)
		case *types.RegisterERC20Proposal:
			return handleRegisterERC20Proposal(ctx, k, c)
		case *types.ToggleTokenConversionProposal:
			return handleToggleConversionProposal(ctx, k, c)
		case *types.UpdateDenomAliasProposal:
			return handleUpdateDenomAliasProposal(ctx, k, c)
		default:
			return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}

// Deprecated
func handleRegisterCoinProposal(ctx sdk.Context, k keeper.Keeper, p *types.RegisterCoinProposal) error {
	_, err := k.RegisterNativeCoin(ctx, p.Metadata)
	return err
}

// Deprecated
func handleRegisterERC20Proposal(ctx sdk.Context, k keeper.Keeper, p *types.RegisterERC20Proposal) error {
	_, err := k.RegisterNativeERC20(ctx, common.HexToAddress(p.Erc20Address), p.Aliases...)
	return err
}

// Deprecated
func handleToggleConversionProposal(ctx sdk.Context, k keeper.Keeper, p *types.ToggleTokenConversionProposal) error {
	_, err := k.ToggleTokenConvert(ctx, p.Token)
	return err
}

// Deprecated
func handleUpdateDenomAliasProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateDenomAliasProposal) error {
	_, err := k.UpdateDenomAliases(ctx, p.Denom, p.Alias)
	return err
}
