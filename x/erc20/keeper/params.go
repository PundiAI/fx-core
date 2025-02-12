package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) CheckEnableErc20(ctx context.Context) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	if !params.EnableErc20 {
		return sdkerrors.ErrInvalidRequest.Wrapf("ERC20 is not enabled")
	}
	return nil
}
