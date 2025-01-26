package app_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/types/legacy"
	"github.com/pundiai/fx-core/v8/types/legacy/gravity"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func Test_MsgServiceRouter(t *testing.T) {
	myApp := helpers.NewApp()

	msgServiceRouter := myApp.MsgServiceRouter()
	deprecated := map[string]struct{}{
		sdk.MsgTypeURL(&legacy.MsgUpdateParams{}):    {},
		sdk.MsgTypeURL(&legacy.MsgUpdateFXParams{}):  {},
		sdk.MsgTypeURL(&legacy.MsgUpdateEGFParams{}): {},

		sdk.MsgTypeURL(&legacy.MsgConvertDenom{}):     {},
		sdk.MsgTypeURL(&legacy.MsgConvertERC20{}):     {},
		sdk.MsgTypeURL(&legacy.MsgUpdateDenomAlias{}): {},
		sdk.MsgTypeURL(&legacy.MsgRegisterERC20{}):    {},
		sdk.MsgTypeURL(&legacy.MsgRegisterCoin{}):     {},

		sdk.MsgTypeURL(&legacy.MsgValsetConfirm{}):           {},
		sdk.MsgTypeURL(&legacy.MsgConfirmBatch{}):            {},
		sdk.MsgTypeURL(&gravity.MsgSetOrchestratorAddress{}): {},
		sdk.MsgTypeURL(&legacy.MsgFxOriginatedTokenClaim{}):  {},
		sdk.MsgTypeURL(&gravity.MsgRequestBatch{}):           {},
		sdk.MsgTypeURL(&legacy.MsgWithdrawClaim{}):           {},
		sdk.MsgTypeURL(&legacy.MsgSendToEth{}):               {},
		sdk.MsgTypeURL(&legacy.MsgCancelSendToEth{}):         {},
		sdk.MsgTypeURL(&legacy.MsgValsetUpdatedClaim{}):      {},
		sdk.MsgTypeURL(&legacy.MsgDepositClaim{}):            {},

		sdk.MsgTypeURL(&legacy.MsgGrantPrivilege{}):      {},
		sdk.MsgTypeURL(&legacy.MsgEditConsensusPubKey{}): {},

		sdk.MsgTypeURL(&legacy.MsgSetOrchestratorAddress{}): {},
		sdk.MsgTypeURL(&legacy.MsgAddOracleDeposit{}):       {},
		sdk.MsgTypeURL(&legacy.MsgCancelSendToExternal{}):   {},
		sdk.MsgTypeURL(&legacy.MsgIncreaseBridgeFee{}):      {},
		sdk.MsgTypeURL(&legacy.MsgRequestBatch{}):           {},

		sdk.MsgTypeURL(&legacy.MsgTransfer{}): {},

		// MsgClaim
		sdk.MsgTypeURL(&crosschaintypes.MsgSendToExternalClaim{}):   {},
		sdk.MsgTypeURL(&crosschaintypes.MsgSendToFxClaim{}):         {},
		sdk.MsgTypeURL(&crosschaintypes.MsgBridgeCallClaim{}):       {},
		sdk.MsgTypeURL(&crosschaintypes.MsgBridgeTokenClaim{}):      {},
		sdk.MsgTypeURL(&crosschaintypes.MsgOracleSetUpdatedClaim{}): {},
		sdk.MsgTypeURL(&crosschaintypes.MsgBridgeCallResultClaim{}): {},

		// MsgConfirm
		sdk.MsgTypeURL(&crosschaintypes.MsgConfirmBatch{}):      {},
		sdk.MsgTypeURL(&crosschaintypes.MsgOracleSetConfirm{}):  {},
		sdk.MsgTypeURL(&crosschaintypes.MsgBridgeCallConfirm{}): {},
	}
	for _, msg := range myApp.InterfaceRegistry().ListImplementations(sdk.MsgInterfaceProtoName) {
		if _, ok := deprecated[msg]; ok {
			continue
		}
		assert.NotNil(t, msgServiceRouter.HandlerByTypeURL(msg), msg)
	}
}
