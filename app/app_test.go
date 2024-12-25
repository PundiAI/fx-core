package app_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/types/legacy"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func Test_MsgServiceRouter(t *testing.T) {
	myApp := helpers.NewApp()

	msgServiceRouter := myApp.MsgServiceRouter()
	deprecated := map[string]struct{}{
		sdk.MsgTypeURL(&crosschaintypes.MsgSetOrchestratorAddress{}): {},
		sdk.MsgTypeURL(&crosschaintypes.MsgAddOracleDeposit{}):       {},
		sdk.MsgTypeURL(&legacy.MsgUpdateParams{}):                    {},
		sdk.MsgTypeURL(&legacy.MsgUpdateFXParams{}):                  {},
		sdk.MsgTypeURL(&legacy.MsgUpdateEGFParams{}):                 {},
		sdk.MsgTypeURL(&crosschaintypes.MsgCancelSendToExternal{}):   {},
		sdk.MsgTypeURL(&crosschaintypes.MsgIncreaseBridgeFee{}):      {},
		sdk.MsgTypeURL(&crosschaintypes.MsgRequestBatch{}):           {},
		sdk.MsgTypeURL(&erc20types.MsgConvertDenom{}):                {},
		sdk.MsgTypeURL(&erc20types.MsgConvertERC20{}):                {},
		sdk.MsgTypeURL(&erc20types.MsgUpdateDenomAlias{}):            {},
		sdk.MsgTypeURL(&erc20types.MsgRegisterERC20{}):               {},
		sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}):                {},
	}
	for _, msg := range myApp.InterfaceRegistry().ListImplementations(sdk.MsgInterfaceProtoName) {
		if _, ok := deprecated[msg]; ok {
			continue
		}
		assert.NotNil(t, msgServiceRouter.HandlerByTypeURL(msg), msg)
	}
}
