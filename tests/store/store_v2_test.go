package store_test

import (
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	v2 "github.com/functionx/fx-core/v3/x/gravity/legacy/v2"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
)

func TestLocalStoreInV2(t *testing.T) {
	if !helpers.IsLocalTest() {
		t.Skip("skipping local test")
	}
	logger := log.NewNopLogger()
	db, err := sdk.NewLevelDB("application", filepath.Join(fxtypes.GetDefaultNodeHome(), "data"))
	require.NoError(t, err)

	appEncodingCfg := app.MakeEncodingConfig()
	myApp := app.New(logger, db,
		nil, true, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		appEncodingCfg, app.EmptyAppOptions{},
	)
	ctx := myApp.NewUncachedContext(false, tmproto.Header{Height: myApp.LastBlockHeight()})

	tests := []struct {
		name     string
		testCase func(t *testing.T)
	}{
		{
			name: "ModuleConsensusVersion",
			testCase: func(t *testing.T) {
				moduleVersion := module.VersionMap{
					"polygon":      2,
					"genutil":      1,
					"gov":          2,
					"params":       1,
					"slashing":     2,
					"transfer":     1,
					"authz":        1,
					"crisis":       1,
					"bsc":          2,
					"capability":   1,
					"crosschain":   1,
					"bank":         2,
					"erc20":        1,
					"feemarket":    3,
					"migrate":      1,
					"mint":         1,
					"staking":      2,
					"tron":         2,
					"vesting":      1,
					"feegrant":     1,
					"evidence":     1,
					"evm":          2,
					"ibc":          2,
					"other":        1,
					"upgrade":      1,
					"distribution": 2,
					"gravity":      1,
					"auth":         2,
				}
				vm := myApp.UpgradeKeeper.GetModuleVersionMap(ctx)
				for k, v := range vm {
					require.Equal(t, moduleVersion[k], v)
				}
			},
		},
		{
			name: "Iterator gravity module store",
			testCase: func(t *testing.T) {
				store := ctx.MultiStore().GetKVStore(myApp.GetKey(gravitytypes.ModuleName))
				v2.ParseStore(t, store)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.testCase)
	}
}
