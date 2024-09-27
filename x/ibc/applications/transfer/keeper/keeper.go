package keeper

import (
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/functionx/fx-core/v8/x/ibc/applications/transfer/types"
)

// Keeper defines the IBC fungible transfer keeper
type Keeper struct {
	ibctransferkeeper.Keeper
	cdc codec.Codec

	porttypes.ICS4Wrapper
	erc20Keeper types.Erc20Keeper
	evmKeeper   types.EvmKeeper
	refundHook  types.RefundHook
}

// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(keeper ibctransferkeeper.Keeper, cdc codec.Codec, ics4Wrapper porttypes.ICS4Wrapper) Keeper {
	return Keeper{
		Keeper:      keeper,
		cdc:         cdc,
		ICS4Wrapper: ics4Wrapper,
	}
}

func (k Keeper) SetRefundHook(hook types.RefundHook) Keeper {
	k.refundHook = hook
	return k
}

func (k Keeper) SetErc20Keeper(erc20Keeper types.Erc20Keeper) Keeper {
	k.erc20Keeper = erc20Keeper
	return k
}

func (k Keeper) SetEvmKeeper(evmKeeper types.EvmKeeper) Keeper {
	k.evmKeeper = evmKeeper
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+ibcexported.ModuleName+"-"+types.CompatibleModuleName)
}
