package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) AddERC20Token(ctx context.Context, name, symbol string, decimals uint8, erc20Addr common.Address, contractOwner types.Owner) (types.ERC20Token, error) {
	var metadata banktypes.Metadata
	if symbol == fxtypes.DefaultDenom {
		metadata = fxtypes.NewFXMetaData()
	} else {
		metadata = fxtypes.NewMetadata(name, symbol, uint32(decimals))
	}
	if err := metadata.Validate(); err != nil {
		return types.ERC20Token{}, sdkerrors.ErrInvalidRequest.Wrapf("metadata: %s", err.Error())
	}
	if !k.bankKeeper.HasDenomMetaData(ctx, metadata.Base) {
		k.bankKeeper.SetDenomMetaData(ctx, metadata)
	}

	if has, err := k.ERC20Token.Has(ctx, metadata.Base); err != nil {
		return types.ERC20Token{}, err
	} else if has {
		return types.ERC20Token{}, types.ErrExists.Wrapf("denom %s is already registered", metadata.Base)
	}

	erc20Token := types.ERC20Token{
		Erc20Address:  erc20Addr.String(),
		Denom:         metadata.Base,
		Enabled:       true,
		ContractOwner: contractOwner,
	}
	if err := k.DenomIndex.Set(ctx, erc20Token.Erc20Address, erc20Token.Denom); err != nil {
		return types.ERC20Token{}, err
	}
	if err := k.ERC20Token.Set(ctx, erc20Token.Denom, erc20Token); err != nil {
		return types.ERC20Token{}, err
	}
	return erc20Token, nil
}

func (k Keeper) GetERC20Token(ctx context.Context, baseDenom string) (types.ERC20Token, error) {
	return k.ERC20Token.Get(ctx, baseDenom)
}

func (k Keeper) ToggleTokenConvert(ctx context.Context, token string) (types.ERC20Token, error) {
	baseDenom, err := k.DenomIndex.Get(ctx, token)
	if err != nil {
		baseDenom = token
	}
	erc20Token, err := k.ERC20Token.Get(ctx, baseDenom)
	if err != nil {
		return types.ERC20Token{}, sdkerrors.ErrNotFound.Wrapf("token %s not found", token)
	}
	erc20Token.Enabled = !erc20Token.Enabled

	if err = k.ERC20Token.Set(ctx, baseDenom, erc20Token); err != nil {
		return types.ERC20Token{}, err
	}

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeToggleTokenRelay,
		sdk.NewAttribute(types.AttributeKeyDenom, erc20Token.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, erc20Token.Erc20Address),
	))
	return erc20Token, nil
}
