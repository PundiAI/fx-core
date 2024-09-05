package app_test

import (
	"os"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v8/app"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	govlegacy "github.com/functionx/fx-core/v8/x/gov/legacy"
	gravitytypes "github.com/functionx/fx-core/v8/x/gravity/types"
)

func Test_MsgServiceRouter(t *testing.T) {
	home := t.TempDir()
	chainId := fxtypes.TestnetChainId
	encodingConfig := app.MakeEncodingConfig()

	myApp := app.New(log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowAll()),
		dbm.NewMemDB(), nil, true, map[int64]bool{}, home, 0,
		encodingConfig, app.EmptyAppOptions{}, baseapp.SetChainID(chainId))

	msgServiceRouter := myApp.MsgServiceRouter()
	// nolint:staticcheck
	deprecated := map[string]struct{}{
		sdk.MsgTypeURL(&crosschaintypes.MsgSetOrchestratorAddress{}): {},
		sdk.MsgTypeURL(&crosschaintypes.MsgAddOracleDeposit{}):       {},
		sdk.MsgTypeURL(&gravitytypes.MsgSetOrchestratorAddress{}):    {},
		sdk.MsgTypeURL(&gravitytypes.MsgFxOriginatedTokenClaim{}):    {},
		sdk.MsgTypeURL(&govlegacy.MsgUpdateParams{}):                 {},
	}
	for _, msg := range encodingConfig.InterfaceRegistry.ListImplementations(sdk.MsgInterfaceProtoName) {
		if _, ok := deprecated[msg]; ok {
			continue
		}
		assert.NotNil(t, msgServiceRouter.HandlerByTypeURL(msg), msg)
	}
}
