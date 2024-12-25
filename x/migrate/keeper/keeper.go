package keeper

import (
	"bytes"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/x/migrate/types"
)

// secp256k1
const addressLen = 20

type Keeper struct {
	// Protobuf codec
	cdc codec.BinaryCodec
	// Store key required for the Fee Market Prefix KVStore.
	storeKey storetypes.StoreKey
	// account keeper
	accountKeeper types.AccountKeeper
	// Migrate handlers
	migrates []MigrateI
}

// NewKeeper generates new fee market module keeper
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, ak types.AccountKeeper, migrates ...MigrateI) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		accountKeeper: ak,
		migrates:      migrates,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetMigrateI get all migrate handlers
func (k Keeper) GetMigrateI() []MigrateI {
	return k.migrates
}

// SetMigrateRecord set from and to migrate record
func (k Keeper) SetMigrateRecord(ctx sdk.Context, from sdk.AccAddress, to common.Address) {
	store := ctx.KVStore(k.storeKey)

	bzFrom := make([]byte, 1+addressLen+8)
	bzTo := make([]byte, 1+addressLen+8)

	height := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))

	copy(bzFrom, types.ValuePrefixMigrateFromFlag)
	copy(bzFrom[1:], to.Bytes())
	copy(bzFrom[1+addressLen:], height)

	copy(bzTo, types.ValuePrefixMigrateToFlag)
	copy(bzTo[1:], from.Bytes())
	copy(bzTo[1+addressLen:], height)

	store.Set(types.GetMigratedRecordKey(from.Bytes()), bzFrom)
	store.Set(types.GetMigratedRecordKey(to.Bytes()), bzTo)

	store.Set(types.GetMigratedDirectionFrom(from), []byte{1})
	store.Set(types.GetMigratedDirectionTo(to), []byte{1})
}

// GetMigrateRecord get address migrate record
func (k Keeper) GetMigrateRecord(ctx sdk.Context, addr []byte) (mr types.MigrateRecord, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMigratedRecordKey(addr))
	if len(bz) < addressLen+9 {
		return mr, false
	}
	mr.Height = int64(sdk.BigEndianToUint64(bz[addressLen+1:]))
	if bytes.Equal(bz[:1], types.ValuePrefixMigrateFromFlag) {
		mr.From = sdk.AccAddress(addr).String()
		mr.To = common.BytesToAddress(bz[1 : addressLen+1]).String()
	} else {
		mr.From = sdk.AccAddress(bz[1 : addressLen+1]).String()
		mr.To = common.BytesToAddress(addr).String()
	}
	return mr, true
}

func (k Keeper) HasMigrateRecord(ctx sdk.Context, addr []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetMigratedRecordKey(addr))
}

func (k Keeper) HasMigratedDirectionFrom(ctx sdk.Context, addr []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetMigratedDirectionFrom(addr))
}

func (k Keeper) HasMigratedDirectionTo(ctx sdk.Context, addr common.Address) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetMigratedDirectionTo(addr))
}

// checkMigrateFrom check migrate from address
func (k Keeper) checkMigrateFrom(ctx sdk.Context, addr sdk.AccAddress) (sdk.AccountI, error) {
	fromAccount := k.accountKeeper.GetAccount(ctx, addr)
	if fromAccount == nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("empty account: %s", addr.String())
	}
	fromPubKey := fromAccount.GetPubKey()
	if fromPubKey == nil {
		return nil, sdkerrors.ErrInvalidPubKey.Wrapf("empty public key: %s", addr.String())
	}
	if fromPubKey.Type() != new(secp256k1.PubKey).Type() {
		return nil, sdkerrors.ErrInvalidPubKey.Wrapf("account type not support: %s(%s)", addr.String(), fromPubKey.Type())
	}
	return fromAccount, nil
}

func (k Keeper) IterateMigrateRecords(ctx sdk.Context, cb func(types.MigrateRecord) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.KeyPrefixMigratedRecord)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		if bytes.Equal(iter.Value()[:1], types.ValuePrefixMigrateToFlag) {
			continue
		}
		if cb(types.MigrateRecord{
			From:   sdk.AccAddress(iter.Key()[1:]).String(),
			To:     common.BytesToAddress(iter.Value()[1 : addressLen+1]).String(),
			Height: int64(sdk.BigEndianToUint64(iter.Value()[1+addressLen:])),
		}) {
			break
		}
	}
}
