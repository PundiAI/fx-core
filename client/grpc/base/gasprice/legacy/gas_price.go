package legacy

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Deprecated: NewQuerier
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
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query endpoint: %s", path[0])
		}
	}
}

// Deprecated:
type Querier struct{}

var _ QueryServer = Querier{}

func (q Querier) GasPrice(c context.Context, _ *GasPriceRequest) (*GasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var gasPrices sdk.Coins
	for _, coin := range ctx.MinGasPrices() {
		gasPrices = append(gasPrices, sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
	}
	return &GasPriceResponse{GasPrices: gasPrices}, nil
}
