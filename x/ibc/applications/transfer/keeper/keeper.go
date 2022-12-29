package keeper

import (
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/v3/x/ibc/applications/transfer/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
)

// Keeper defines the IBC fungible transfer keeper
type Keeper struct {
	ibctransferkeeper.Keeper
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace

	ics4Wrapper   porttypes.ICS4Wrapper
	channelKeeper transfertypes.ChannelKeeper
	portKeeper    transfertypes.PortKeeper
	authKeeper    transfertypes.AccountKeeper
	bankKeeper    transfertypes.BankKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
	Router        *types.Router
	RefundHook    types.RefundHook
	AckHook       types.AckHook
}

// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(keeper ibctransferkeeper.Keeper,
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
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
func (k Keeper) SetRouter(rtr *types.Router) Keeper {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
	return k
}

func (k Keeper) GetRouter() *types.Router {
	return k.Router
}

func (k Keeper) SetRefundHook(hook types.RefundHook) Keeper {
	k.RefundHook = hook
	return k
}

func (k Keeper) SetAckHook(hook types.AckHook) Keeper {
	k.AckHook = hook
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+types.CompatibleModuleName)
}
