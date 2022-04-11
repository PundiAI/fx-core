package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/x/migrate/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	// Protobuf codec
	cdc codec.BinaryCodec
	// Store key required for the Fee Market Prefix KVStore.
	storeKey sdk.StoreKey
	// Migrate handlers
	migrateI []MigrateI
}

// NewKeeper generates new fee market module keeper
func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	ctx.KVStore(k.storeKey)
	return ctx.Logger().With("module", types.ModuleName)
}

// SetMigrateI set migrate handlers
func (k *Keeper) SetMigrateI(migrate ...MigrateI) {
	k.migrateI = migrate
}

// GetMigrateI get all migrate handlers
func (k *Keeper) GetMigrateI() []MigrateI {
	return k.migrateI
}

// SetMigrateRecord set from and to migrate record
func (k *Keeper) SetMigrateRecord(ctx sdk.Context, from, to sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	mr := types.MigrateRecord{
		From:   from.String(),
		To:     to.String(),
		Height: ctx.BlockHeight(),
	}
	bz := k.cdc.MustMarshal(&mr)
	store.Set(types.GetMigratedRecordKey(from), bz)
	store.Set(types.GetMigratedRecordKey(to), bz)
}

// GetMigrateRecord get address migrate record
func (k *Keeper) GetMigrateRecord(ctx sdk.Context, addr sdk.AccAddress) (mr types.MigrateRecord, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMigratedRecordKey(addr))
	if bz == nil {
		return mr, false
	}
	k.cdc.MustUnmarshal(bz, &mr)
	return mr, true
}

func (k *Keeper) HasMigrateRecord(ctx sdk.Context, addr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetMigratedRecordKey(addr))
}
