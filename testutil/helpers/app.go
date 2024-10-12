package helpers

import (
	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/viper"

	"github.com/functionx/fx-core/v8/app"
	fxtypes "github.com/functionx/fx-core/v8/types"
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
