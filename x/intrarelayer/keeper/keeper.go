package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/x/intrarelayer/types"
)

// Keeper of this module maintains collections of intrarelayer.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	evmKeeper     types.EVMKeeper

	ibcTransferKeeper types.IBCTransferKeeper
	ibcChannelKeeper  types.IBCChannelKeeper
}

// NewKeeper creates new instances of the intrarelayer Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	evmKeeper types.EVMKeeper,
	ibcTransferKeeper types.IBCTransferKeeper,
	ibcChannelKeeper types.IBCChannelKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	k := &Keeper{
		storeKey:          storeKey,
		cdc:               cdc,
		paramstore:        ps,
		accountKeeper:     ak,
		bankKeeper:        bk,
		evmKeeper:         evmKeeper,
		ibcTransferKeeper: ibcTransferKeeper,
		ibcChannelKeeper:  ibcChannelKeeper,
	}
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) HasInit(ctx sdk.Context) bool {
	return k.paramstore.Has(ctx, types.ParamStoreKeyEnableIntrarelayer)
}

func (k *Keeper) SetIBCTransferKeeper(ibcTransferKeepr types.IBCTransferKeeper) *Keeper {
	k.ibcTransferKeeper = ibcTransferKeepr
	return k
}

func (k *Keeper) SetIBCChannelKeeper(ibcChannelKeeper types.IBCChannelKeeper) {
	k.ibcChannelKeeper = ibcChannelKeeper
}
