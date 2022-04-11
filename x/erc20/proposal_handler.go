package erc20

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	fxtype "github.com/functionx/fx-core/types"

	"github.com/functionx/fx-core/x/erc20/keeper"
	"github.com/functionx/fx-core/x/erc20/types"
)

// NewErc20ProposalHandler creates a governance handler to manage new proposal types.
// It enables RegisterTokenPairProposal to propose a registration of token mapping
func NewErc20ProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		if ctx.BlockHeight() < fxtype.EvmSupportBlock() {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("erc20 module not enable"))
		}
		switch c := content.(type) {
		case *types.InitEvmProposal:
			return k.HandleInitEvmProposal(ctx, c.Erc20Params, c.FeemarketParams, c.EvmParams, c.Metadatas)
		case *types.RegisterCoinProposal:
			return handleRegisterCoinProposal(ctx, k, c)
		//case *types.RegisterERC20Proposal:
		//	return handleRegisterERC20Proposal(ctx, k, c)
		case *types.ToggleTokenRelayProposal:
			return handleToggleRelayProposal(ctx, k, c)
		//case *types.UpdateTokenPairERC20Proposal:
		//	return handleUpdateTokenPairERC20Proposal(ctx, k, c)
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
			sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
		),
	)

	return nil
}

func handleRegisterERC20Proposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterERC20Proposal) error {
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

func handleUpdateTokenPairERC20Proposal(ctx sdk.Context, k *keeper.Keeper, p *types.UpdateTokenPairERC20Proposal) error {
	pair, err := k.UpdateTokenPairERC20(ctx, p.GetERC20Address(), p.GetNewERC20Address())
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateTokenPairERC20,
			sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
		),
	)

	return nil
}
