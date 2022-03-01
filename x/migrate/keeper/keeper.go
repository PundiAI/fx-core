package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typescommon "github.com/functionx/fx-core/x/migrate/types/common"
	typesv1 "github.com/functionx/fx-core/x/migrate/types/v1"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	// Protobuf codec
	cdc codec.BinaryMarshaler
	// Store key required for the Fee Market Prefix KVStore.
	storeKey sdk.StoreKey
	// Migrate handlers
	migrateI []MigrateI
}

// NewKeeper generates new fee market module keeper
func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	ctx.KVStore(k.storeKey)
	return ctx.Logger().With("module", typescommon.ModuleName)
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
	mr := typesv1.MigrateRecord{
		From:   from.String(),
		To:     to.String(),
		Height: ctx.BlockHeight(),
	}
	bz := k.cdc.MustMarshalBinaryBare(&mr)
	store.Set(typescommon.GetMigratedRecordKey(from), bz)
	store.Set(typescommon.GetMigratedRecordKey(to), bz)
}

// GetMigrateRecord get address migrate record
func (k *Keeper) GetMigrateRecord(ctx sdk.Context, addr sdk.AccAddress) (mr typesv1.MigrateRecord, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(typescommon.GetMigratedRecordKey(addr))
	if bz == nil {
		return mr, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &mr)
	return mr, true
}

func (k *Keeper) HasMigrateRecord(ctx sdk.Context, addr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(typescommon.GetMigratedRecordKey(addr))
}
