package helpers

import (
	"encoding/json"
	"os"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	ibctestingtypes "github.com/cosmos/ibc-go/v7/testing/types"

	"github.com/functionx/fx-core/v8/app"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

type TestingApp struct {
	*app.App
	TxConfig client.TxConfig
}

func (app *TestingApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *TestingApp) GetTxConfig() client.TxConfig {
	return app.TxConfig
}

func (app *TestingApp) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.StakingKeeper
}

func (app *TestingApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *TestingApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// SetupTestingApp initializes the IBC-go testing application
func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	cfg := app.MakeEncodingConfig()
	myApp := app.New(log.NewNopLogger(), dbm.NewMemDB(),
		nil, true, map[int64]bool{}, os.TempDir(), 5, cfg, app.EmptyAppOptions{},
		baseapp.SetChainID(fxtypes.ChainId()))
	testingApp := &TestingApp{App: myApp, TxConfig: cfg.TxConfig}
	return testingApp, app.NewDefAppGenesisByDenom(fxtypes.DefaultDenom, cfg.Codec)
}
