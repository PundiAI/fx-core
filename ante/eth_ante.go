package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmante "github.com/evmos/ethermint/app/ante"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// Just copy, because the pendingTxListener property of the TxListenerDecorator is internal
// https://github.com/pundiai/ethermint/blob/fxcore/v0.22.x/app/ante/tx_listener.go

type TxListenerDecorator struct {
	pendingTxListener evmante.PendingTxListener
}

// newTxListenerDecorator creates a new TxListenerDecorator with the provided PendingTxListener.
// CONTRACT: must be put at the last of the chained decorators
func newTxListenerDecorator(pendingTxListener evmante.PendingTxListener) TxListenerDecorator {
	return TxListenerDecorator{pendingTxListener}
}

func (d TxListenerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}
	if ctx.IsCheckTx() && !simulate && d.pendingTxListener != nil {
		for _, msg := range tx.GetMsgs() {
			if ethTx, ok := msg.(*evmtypes.MsgEthereumTx); ok {
				d.pendingTxListener(ethTx.Hash())
			}
		}
	}
	return next(ctx, tx, simulate)
}
