package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/functionx/fx-core/x/evm/types"
	intrarelayertypes "github.com/functionx/fx-core/x/intrarelayer/types"
)

type IntrarelayerKeeperI interface {
	ModuleInit(ctx sdk.Context, enableIntrarelayer, enableEvmHook bool, ibcTransferTimeoutHeight uint64) error
	RegisterCoin(ctx sdk.Context, coinMetadata banktypes.Metadata) (*intrarelayertypes.TokenPair, error)
}

func (k Keeper) HandleInitEvmProposal(ctx sdk.Context, p *types.InitEvmProposal) error {
	// check duplicate init params.
	if k.HasInit(ctx) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate init evm")
	}
	//init fee market
	k.Logger(ctx).Info("init fee market", "params", p.FeemarketParams.String())
	if p.FeemarketParams.BaseFee.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "base fee cannot be negative")
	}
	// set feeMarket baseFee
	k.feeMarketKeeper.SetBaseFee(ctx, p.FeemarketParams.BaseFee.BigInt())
	// set feeMarket blockGasUsed
	k.feeMarketKeeper.SetBlockGasUsed(ctx, 0)
	// init feeMarket module params
	k.feeMarketKeeper.SetParams(ctx, *p.FeemarketParams)

	//init evm
	k.Logger(ctx).Info("init evm", "params", p.EvmParams.String())
	k.SetParams(ctx, *p.EvmParams)

	//init intrarelayer
	k.Logger(ctx).Info("init intrarelayer", "params", p.IntrarelayerParams.String())

	if err := k.intrarelayerKeeper.ModuleInit(ctx, p.IntrarelayerParams.EnableIntrarelayer,
		p.IntrarelayerParams.EnableEVMHook, p.IntrarelayerParams.IbcTransferTimeoutHeight); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	//init register coin
	events := make([]sdk.Event, 0, len(p.Metadata))
	for _, metadata := range p.Metadata {
		k.Logger(ctx).Info("register coin", "coin", metadata.String())
		pair, err := k.intrarelayerKeeper.RegisterCoin(ctx, metadata)
		if err != nil {
			return sdkerrors.Wrapf(intrarelayertypes.ErrInvalidMetadata, fmt.Sprintf("base %s, display %s, error %s",
				metadata.Base, metadata.Display, err.Error()))
		}
		event := sdk.NewEvent(
			intrarelayertypes.EventTypeRegisterCoin,
			sdk.NewAttribute(intrarelayertypes.AttributeKeyCosmosCoin, pair.Denom),
			sdk.NewAttribute(intrarelayertypes.AttributeKeyFIP20Token, pair.Fip20Address),
		)
		events = append(events, event)
	}
	ctx.EventManager().EmitEvents(events)
	return nil
}
