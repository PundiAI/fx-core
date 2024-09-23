package app_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	govlegacy "github.com/functionx/fx-core/v8/x/gov/legacy"
)

func Test_MsgServiceRouter(t *testing.T) {
	myApp := helpers.NewApp()

	msgServiceRouter := myApp.MsgServiceRouter()
	deprecated := map[string]struct{}{
		sdk.MsgTypeURL(&crosschaintypes.MsgSetOrchestratorAddress{}): {},
		sdk.MsgTypeURL(&crosschaintypes.MsgAddOracleDeposit{}):       {},
		sdk.MsgTypeURL(&govlegacy.MsgUpdateParams{}):                 {},
	}
	for _, msg := range myApp.InterfaceRegistry().ListImplementations(sdk.MsgInterfaceProtoName) {
		if _, ok := deprecated[msg]; ok {
			continue
		}
		assert.NotNil(t, msgServiceRouter.HandlerByTypeURL(msg), msg)
	}
}
