package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	k Keeper
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{
		k: k,
	}
}

func (s msgServer) ConvertCoin(c context.Context, msg *types.MsgConvertCoin) (*types.MsgConvertCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	receiver := common.HexToAddress(msg.Receiver)
	_, err := s.k.ConvertCoin(ctx, sender, receiver, msg.Coin)
	return &types.MsgConvertCoinResponse{}, err
}

func (s msgServer) UpdateParams(c context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if s.k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", s.k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.k.Params.Set(ctx, req.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

func (s msgServer) ToggleTokenConversion(c context.Context, req *types.MsgToggleTokenConversion) (*types.MsgToggleTokenConversionResponse, error) {
	if s.k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", s.k.authority, req.Authority)
	}
	erc20Token, err := s.k.ToggleTokenConvert(c, req.Token)
	if err != nil {
		return nil, err
	}
	return &types.MsgToggleTokenConversionResponse{Erc20Token: erc20Token}, nil
}
