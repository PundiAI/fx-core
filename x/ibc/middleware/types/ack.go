package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

func NewAckErrorWithErrorEvent(ctx sdk.Context, err error) exported.Acknowledgement {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventTypeReceive,
		sdk.NewAttribute(AttributeKeyError, err.Error()),
	))

	return channeltypes.NewErrorAcknowledgement(err)
}
