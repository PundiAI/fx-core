package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/tendermint/tendermint/libs/log"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
)

// Keeper defines the IBC fungible transfer keeper
type Keeper struct {
	ibctransferkeeper.Keeper
	storeKey   storetypes.StoreKey
	cdc        codec.Codec
	paramSpace paramtypes.Subspace

	ics4Wrapper   porttypes.ICS4Wrapper
	channelKeeper transfertypes.ChannelKeeper
	portKeeper    transfertypes.PortKeeper
	authKeeper    transfertypes.AccountKeeper
	bankKeeper    transfertypes.BankKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
	erc20Keeper   types.Erc20Keeper
	evmKeeper     types.EvmKeeper
	router        *fxtypes.Router
	refundHook    types.RefundHook
}

// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(keeper ibctransferkeeper.Keeper,
	cdc codec.Codec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ics4Wrapper porttypes.ICS4Wrapper, channelKeeper transfertypes.ChannelKeeper, portKeeper transfertypes.PortKeeper,
	authKeeper transfertypes.AccountKeeper, bankKeeper transfertypes.BankKeeper, scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	// ensure ibc transfer module account is set
	if addr := authKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the FX IBC transfer module account has not been set")
	}

	return Keeper{
		Keeper:        keeper,
		cdc:           cdc,
		storeKey:      key,
		paramSpace:    paramSpace,
		ics4Wrapper:   ics4Wrapper,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

// SetRouter sets the Router in IBC Transfer Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k Keeper) SetRouter(rtr fxtypes.Router) Keeper {
	if k.router != nil && k.router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.router = &rtr
	k.router.Seal()
	return k
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
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+types.CompatibleModuleName)
}
