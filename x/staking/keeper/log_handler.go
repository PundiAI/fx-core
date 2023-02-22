package keeper

import (
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/types/contract"
	"github.com/functionx/fx-core/v3/x/staking/types"
)

type LPTokenTransferHandler struct {
	*Keeper
}

func NewLPTokenTransferHandler(k *Keeper) *LPTokenTransferHandler {
	return &LPTokenTransferHandler{k}
}

func (h LPTokenTransferHandler) EventID() common.Hash {
	return fxtypes.GetLPToken().ABI.Events[types.LPTokenTransferEventName].ID
}

func (h LPTokenTransferHandler) Handle(ctx sdk.Context, _ core.Message, log *ethtypes.Log) error {
	valAddr, found := h.GetLPTokenValidator(ctx, log.Address)
	if !found {
		return sdkerrors.Wrapf(types.ErrLPTokenNotFound, "contract: %s", log.Address.String())
	}

	var res types.FXLPTokenTransfer
	if err := contract.ParseLogEvent(fxtypes.GetLPToken().ABI, log, types.LPTokenTransferEventName, &res); err != nil {
		return sdkerrors.Wrapf(types.ErrUnexpectedEvent, "parse lp token transfer: %s", err.Error())
	}

	shares := sdk.NewDecFromBigIntWithPrec(res.Value, sdk.Precision)
	if !shares.IsZero() {
		err := h.TransferDelegate(ctx, valAddr, res.From.Bytes(), res.To.Bytes(), shares)
		if err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRelayTransfer,
		sdk.NewAttribute(sdk.AttributeKeySender, res.From.String()),
		sdk.NewAttribute(types.AttributeKeyTo, res.To.String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, res.Value.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
		sdk.NewAttribute(types.AttributeKeyLPTokenAddress, log.Address.String()),
	))

	telemetry.IncrCounterWithLabels(
		[]string{stakingtypes.ModuleName, "_transfer"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("lp_token", log.Address.String()),
			telemetry.NewLabel("validator", valAddr.String()),
		},
	)

	return nil
}
