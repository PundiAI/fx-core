package contract_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
)

func TestAccessControlCaller_HasRole(t *testing.T) {
	app, ctx := helpers.NewAppWithValNumber(t, 1)

	keeper := contract.NewAccessControlKeeper(app.EvmKeeper, contract.AccessControlAddress)

	has, err := keeper.HasRole(ctx, common.HexToHash(contract.TransferModuleRole), helpers.GenHexAddress())
	require.NoError(t, err)

	require.False(t, has)
}
