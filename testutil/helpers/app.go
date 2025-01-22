package helpers

import (
	"testing"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"

	"github.com/pundiai/fx-core/v8/app"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

type AppOpts struct {
	Logger log.Logger
	Home   string
	DB     dbm.DB
}

func NewApp(opts ...func(*AppOpts)) *app.App {
	defOpts := &AppOpts{
		Logger: log.NewNopLogger(),
		Home:   fxtypes.GetDefaultNodeHome(),
		DB:     dbm.NewMemDB(),
	}
	for _, opt := range opts {
		opt(defOpts)
	}
	return app.New(
		defOpts.Logger,
		defOpts.DB,
		nil,
		true,
		map[int64]bool{},
		defOpts.Home,
		viper.New(),
	)
}

func NewAppWithValNumber(t *testing.T, valNumber int) (*app.App, sdk.Context) {
	t.Helper()

	valSet, valPrivs := generateGenesisValidator(valNumber)
	myApp := setupWithGenesisValSet(t, valSet, valPrivs)
	ctx := myApp.GetContextForFinalizeBlock(nil)
	ctx = ctx.WithProposer(valSet.Proposer.Address.Bytes())
	return myApp, ctx
}
