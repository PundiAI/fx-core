package keeper

import (
	// this line is used by starport scaffolding # 1
	"github.com/functionx/fx-core/v2/x/other/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		if len(path) <= 0 {
			return nil, sdkerrors.ErrInvalidRequest
		}
		switch path[0] {
		case "gasPrice":
			var gasPrices sdk.Coins
			for _, coin := range ctx.MinGasPrices() {
				gasPrices = append(gasPrices, sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
			}
			return codec.MarshalJSONIndent(legacyQuerierCdc, gasPrices)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}
