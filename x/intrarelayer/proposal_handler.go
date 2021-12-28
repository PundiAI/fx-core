package intrarelayer

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	fxtype "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/intrarelayer/keeper"
	"github.com/functionx/fx-core/x/intrarelayer/types"
)

// NewIntrarelayerProposalHandler creates a governance handler to manage new proposal types.
// It enables RegisterTokenPairProposal to propose a registration of token mapping
func NewIntrarelayerProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.InitIntrarelayerProposal:
			return handleInitIntrarelayerProposalHandler(ctx, k, c)
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

func handleInitIntrarelayerProposalHandler(ctx sdk.Context, k *keeper.Keeper, p *types.InitIntrarelayerProposal) error {
	if ctx.BlockHeight() < fxtype.IntrarelayerSupportBlock() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("evm module not enable"))
	}
	return k.InitIntrarelayer(ctx, p)
}

func handleRegisterCoinProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterCoinProposal) error {
	if !k.HasInit(ctx) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("intrarelayer module not enable"))
	}
	pair, err := k.RegisterCoin(ctx, p.Metadata)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterCoin,
			sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
		),
	)

	return nil
}

func handleRegisterERC20Proposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterERC20Proposal) error {
	if !k.HasInit(ctx) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("intrarelayer module not enable"))
	}
	pair, err := k.RegisterERC20(ctx, common.HexToAddress(p.Erc20Address))
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterERC20,
			sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
		),
	)

	return nil
}

func handleToggleRelayProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ToggleTokenRelayProposal) error {
	if !k.HasInit(ctx) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("intrarelayer module not enable"))
	}
	pair, err := k.ToggleRelay(ctx, p.Token)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeToggleTokenRelay,
			sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
		),
	)

	return nil
}
