package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/v2/x/migrate/types"
)

const (
	QueryMigrateRecord = "migrateRecord"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		if len(path) <= 0 {
			return nil, sdkerrors.ErrInvalidRequest
		}
		switch path[0] {
		case QueryMigrateRecord:
			return queryMigrateRecord(ctx, path[1], keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint", types.ModuleName)
		}
	}
}

func queryMigrateRecord(ctx sdk.Context, address string, keeper Keeper) ([]byte, error) {
	bech32, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	record, _ := keeper.GetMigrateRecord(ctx, bech32)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, record)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
