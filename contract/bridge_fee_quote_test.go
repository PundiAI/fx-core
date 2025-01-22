package contract_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func TestBridgeFeeQuoteKeeper_GetDefaultOracleQuote(t *testing.T) {
	app, ctx := helpers.NewAppWithValNumber(t, 1)
	keeper := contract.NewBridgeFeeQuoteKeeper(app.EvmKeeper)
	quote, err := keeper.GetDefaultOracleQuote(ctx, contract.MustStrToByte32(ethtypes.ModuleName), contract.MustStrToByte32(fxtypes.DefaultDenom))
	require.NoError(t, err)
	require.Empty(t, quote)
}
