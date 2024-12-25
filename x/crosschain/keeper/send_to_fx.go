package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	fxtelemetry "github.com/pundiai/fx-core/v8/telemetry"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (k Keeper) SendToFxExecuted(ctx sdk.Context, claim *types.MsgSendToFxClaim) error {
	if !ctx.IsCheckTx() {
		defer func() {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "send_to_fx"},
				float32(1),
				[]metrics.Label{
					telemetry.NewLabel("module", k.moduleName),
				},
			)
			fxtelemetry.SetGaugeLabelsWithDenom(
				[]string{types.ModuleName, "send_to_fx_amount"},
				claim.TokenContract, claim.Amount.BigInt(),
				telemetry.NewLabel("module", k.moduleName),
			)
		}()
	}

	receiveAddr, err := sdk.AccAddressFromBech32(claim.Receiver)
	if err != nil {
		return types.ErrInvalid.Wrapf("receiver address")
	}

	bridgeToken, err := k.DepositBridgeToken(ctx, receiveAddr, claim.Amount, claim.TokenContract)
	if err != nil {
		return err
	}
	baseCoin, err := k.BridgeTokenToBaseCoin(ctx, receiveAddr, claim.Amount, bridgeToken)
	if err != nil {
		return err
	}
	if !bridgeToken.IsOrigin() {
		_, err = k.erc20Keeper.BaseCoinToEvm(ctx, common.BytesToAddress(receiveAddr.Bytes()), baseCoin)
		if err != nil {
			return err
		}
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEvmTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
	))
	return nil
}
