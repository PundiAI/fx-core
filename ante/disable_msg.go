package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethante "github.com/evmos/ethermint/app/ante"
)

type DisableMsgDecorator struct {
	govKeeper Govkeeper
	// disabledMsgs is a set that contains type urls of unauthorized msgs.
	disabledMsgTypes []string
}

func NewDisableMsgDecorator(disabledMsgTypes []string, govKeeper Govkeeper) DisableMsgDecorator {
	return DisableMsgDecorator{disabledMsgTypes: disabledMsgTypes, govKeeper: govKeeper}
}

func (dms DisableMsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	disabledMsgs := dms.govKeeper.GetDisabledMsgs(ctx)
	disabledMsgs = append(disabledMsgs, dms.disabledMsgTypes...)
	if len(disabledMsgs) == 0 {
		return next(ctx, tx, simulate)
	}

	return ethante.NewAuthzLimiterDecorator(disabledMsgs).AnteHandle(ctx, tx, simulate, next)
}
