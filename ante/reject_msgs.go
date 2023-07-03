package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// RejectExtensionOptionsDecorator is an AnteDecorator that rejects all extension
// options which can optionally be included in protobuf transactions. Users that
// need extension options should create a custom AnteHandler chain that handles
// needed extension options properly and rejects unknown ones.
type RejectExtensionOptionsDecorator struct{}

// NewRejectExtensionOptionsDecorator creates a new RejectExtensionOptionsDecorator
func NewRejectExtensionOptionsDecorator() RejectExtensionOptionsDecorator {
	return RejectExtensionOptionsDecorator{}
}

var _ sdk.AnteDecorator = RejectExtensionOptionsDecorator{}

// AnteHandle implements the AnteDecorator.AnteHandle method
func (r RejectExtensionOptionsDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if hasExtOptsTx, ok := tx.(ante.HasExtensionOptionsTx); ok {
		if len(hasExtOptsTx.GetExtensionOptions()) != 0 {
			return ctx, errortypes.ErrUnknownExtensionOptions
		}
		if len(hasExtOptsTx.GetNonCriticalExtensionOptions()) != 0 {
			return ctx, errorsmod.Wrap(errortypes.ErrUnknownRequest, "unknown non critical extension options")
		}
	}
	return next(ctx, tx, simulate)
}

// RejectValidatorGrantedDecorator is an AnteDecorator that rejects all transactions from validator granted
type RejectValidatorGrantedDecorator struct {
	sk StakingKeeper
}

func NewRejectValidatorGrantedDecorator(sk StakingKeeper) RejectValidatorGrantedDecorator {
	return RejectValidatorGrantedDecorator{
		sk: sk,
	}
}

func (r RejectValidatorGrantedDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		signers := msg.GetSigners()
		for _, signer := range signers {
			if r.sk.HasValidatorOperator(ctx, signer.Bytes()) {
				return ctx, errorsmod.Wrap(errortypes.ErrInvalidAddress, "validator granted")
			}
		}
	}

	return next(ctx, tx, simulate)
}
