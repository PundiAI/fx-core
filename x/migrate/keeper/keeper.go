package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/x/migrate/types"
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

	bzFrom := make([]byte, 1+sdk.AddrLen+8)
	bzTo := make([]byte, 1+sdk.AddrLen+8)

	height := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))

	copy(bzFrom[:], types.PrefixMigrateFromFlag)
	copy(bzFrom[1:], to.Bytes())
	copy(bzFrom[1+sdk.AddrLen:], height)

	copy(bzTo[:], types.PrefixMigrateToFlag)
	copy(bzTo[1:], from.Bytes())
	copy(bzTo[1+sdk.AddrLen:], height)

	store.Set(types.GetMigratedRecordKey(from), bzFrom)
	store.Set(types.GetMigratedRecordKey(to), bzTo)
}

// GetMigrateRecord get address migrate record
func (k *Keeper) GetMigrateRecord(ctx sdk.Context, addr sdk.AccAddress) (mr types.MigrateRecord, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMigratedRecordKey(addr))
	if len(bz) < sdk.AddrLen+9 {
		return mr, false
	}
	mr.Height = int64(sdk.BigEndianToUint64(bz[sdk.AddrLen+1:]))
	if bytes.Equal(bz[:1], types.PrefixMigrateFromFlag) {
		mr.From = addr.String()
		mr.To = sdk.AccAddress(bz[1 : sdk.AddrLen+1]).String()
	} else {
		mr.From = sdk.AccAddress(bz[1 : sdk.AddrLen+1]).String()
		mr.To = addr.String()
	}
	return mr, true
}

func (k *Keeper) HasMigrateRecord(ctx sdk.Context, addr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetMigratedRecordKey(addr))
}
