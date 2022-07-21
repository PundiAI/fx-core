package simapp

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	fxapp "github.com/functionx/fx-core/v2/app"
	"github.com/functionx/fx-core/v2/app/helpers"
)

var (
//_ ibctesting.TestingApp = (*SimApp)(nil)
)

type SimApp struct {
	*fxapp.App
	// make scoped keepers public for test purposes
	//ScopedIBCMockKeeper capabilitykeeper.ScopedKeeper
}

func (app *SimApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

// TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *SimApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetStakingKeeper implements the TestingApp interface.
func (app *SimApp) GetStakingKeeper() stakingkeeper.Keeper {
	return app.StakingKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *SimApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *SimApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *SimApp) GetTxConfig() client.TxConfig {
	return fxapp.MakeEncodingConfig().TxConfig
}

func NewSimApp() *SimApp {
	db := dbm.NewMemDB()
	encCdc := fxapp.MakeEncodingConfig()
	fxApp := fxapp.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, fxapp.DefaultNodeHome, 5, encCdc, helpers.EmptyAppOptions{})

	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// not replicate if you do not need to test core IBC or light clients.
	//scopedIBCMockKeeper := fxApp.CapabilityKeeper.ScopeToModule(ibcmock.ModuleName)
	app := &SimApp{
		App: fxApp,
		//ScopedIBCMockKeeper: scopedIBCMockKeeper,
	}
	return app
}
