package store_test

import (
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
)

func TestLocalStoreInV2(t *testing.T) {
	if os.Getenv("LOCAL_STORE_TEST") != "true" {
		t.Skip("skipping local store test")
	}
	logger := log.NewNopLogger()
	db, err := sdk.NewLevelDB("application", filepath.Join(app.DefaultNodeHome, "data"))
	require.NoError(t, err)

	appEncodingCfg := app.MakeEncodingConfig()
	myApp := app.New(logger, db,
		nil, true, map[int64]bool{}, app.DefaultNodeHome, 0,
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
				vm := myApp.UpgradeKeeper.GetModuleVersionMap(ctx)
				for k, v := range vm {
					t.Log(k, v)
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.testCase)
	}
}
