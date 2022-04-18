package erc20

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/functionx/fx-core/x/erc20/keeper"
	"github.com/functionx/fx-core/x/erc20/types"
)

// NewErc20ProposalHandler creates a governance handler to manage new proposal types.
// It enables RegisterTokenPairProposal to propose a registration of token mapping
func NewErc20ProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		if ctx.BlockHeight() < fxtypes.EvmSupportBlock() {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "erc20 module not enable")
		}
		switch c := content.(type) {
		case *types.RegisterCoinProposal:
			return handleRegisterCoinProposal(ctx, k, c)
		case *types.RegisterERC20Proposal:
			return handleRegisterERC20Proposal(ctx, k, c)
		case *types.ToggleTokenRelayProposal:
			return handleToggleRelayProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}

func handleRegisterCoinProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterCoinProposal) error {
	pair, err := k.RegisterCoin(ctx, p.Metadata)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterCoin,
			sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
		),
	)

	return nil
}

// nolint
func handleRegisterERC20Proposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterERC20Proposal) error {
	pair, err := k.RegisterERC20(ctx, common.HexToAddress(p.Erc20Address))
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterERC20,
			sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
		),
	)

	return nil
}

func handleToggleRelayProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ToggleTokenRelayProposal) error {
	pair, err := k.ToggleRelay(ctx, p.Token)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeToggleTokenRelay,
			sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
		),
	)

	return nil
}
