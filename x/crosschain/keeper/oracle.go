package keeper

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashicorp/go-metrics"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

// --- PROPOSAL ORACLE --- //

func (k Keeper) SetProposalOracle(ctx sdk.Context, proposalOracle *types.ProposalOracle) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposalOracleKey, k.cdc.MustMarshal(proposalOracle))
}

func (k Keeper) GetProposalOracle(ctx sdk.Context) (proposalOracle types.ProposalOracle, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalOracleKey)
	if bz == nil {
		return proposalOracle, false
	}
	k.cdc.MustUnmarshal(bz, &proposalOracle)
	return proposalOracle, true
}

func (k Keeper) IsProposalOracle(ctx sdk.Context, oracleAddr string) bool {
	proposalOracle, found := k.GetProposalOracle(ctx)
	if !found {
		return false
	}
	for _, oracle := range proposalOracle.Oracles {
		if oracle == oracleAddr {
			return true
		}
	}
	return false
}

// --- Bridger ADDRESS --- //

// SetOracleAddrByBridgerAddr sets the bridger key for a given oracle
func (k Keeper) SetOracleAddrByBridgerAddr(ctx sdk.Context, bridgerAddr, oracleAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleAddressByBridgerKey(bridgerAddr), oracleAddr.Bytes())
}

// GetOracleAddrByBridgerAddr returns the oracle key associated with an bridger key
func (k Keeper) GetOracleAddrByBridgerAddr(ctx sdk.Context, bridgerAddr sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	oracle := store.Get(types.GetOracleAddressByBridgerKey(bridgerAddr))
	return oracle, oracle != nil
}

func (k Keeper) HasOracleAddrByBridgerAddr(ctx sdk.Context, bridgerAddr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOracleAddressByBridgerKey(bridgerAddr))
}

// DelOracleAddrByBridgerAddr delete the bridger key for a given oracle
func (k Keeper) DelOracleAddrByBridgerAddr(ctx sdk.Context, bridgerAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleAddressByBridgerKey(bridgerAddr)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

// --- External ADDRESS --- //

// SetOracleAddrByExternalAddr sets the external address for a given oracle
func (k Keeper) SetOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string, oracleAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleAddressByExternalKey(externalAddress), oracleAddr.Bytes())
}

// GetOracleAddrByExternalAddr returns the external address for a given gravity oracle
func (k Keeper) GetOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOracleAddressByExternalKey(externalAddress))
	return bz, bz != nil
}

func (k Keeper) HasOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOracleAddressByExternalKey(externalAddress))
}

// DelOracleAddrByExternalAddr delete the external address for a give oracle
func (k Keeper) DelOracleAddrByExternalAddr(ctk sdk.Context, externalAddress string) {
	store := ctk.KVStore(k.storeKey)
	oracleAddr := types.GetOracleAddressByExternalKey(externalAddress)
	if !store.Has(oracleAddr) {
		return
	}
	store.Delete(oracleAddr)
}

// --- ORACLE TOTAL POWER --- //

// GetLastTotalPower Load the last total oracle power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdkmath.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastTotalPowerKey)

	if bz == nil {
		return sdkmath.ZeroInt()
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshal(bz, &ip)

	return ip.Int
}

// SetLastTotalPower Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context) {
	oracles := k.GetAllOracles(ctx, true)
	totalPower := sdkmath.ZeroInt()
	for _, oracle := range oracles {
		totalPower = totalPower.Add(oracle.GetPower())
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastTotalPowerKey, k.cdc.MustMarshal(&sdk.IntProto{Int: totalPower}))
}

// --- ORACLE --- //

func (k Keeper) IterateOracle(ctx sdk.Context, cb func(oracle types.Oracle) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		oracle := types.Oracle{}
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		if cb(oracle) {
			break
		}
	}
}

func (k Keeper) SetOracle(ctx sdk.Context, oracle types.Oracle) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleKey(oracle.GetOracle()), k.cdc.MustMarshal(&oracle))
}

func (k Keeper) HasOracle(ctx sdk.Context, oracleAddr sdk.AccAddress) (found bool) {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOracleKey(oracleAddr))
}

func (k Keeper) GetOracle(ctx sdk.Context, oracleAddr sdk.AccAddress) (oracle types.Oracle, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetOracleKey(oracleAddr))
	if value == nil {
		return oracle, false
	}
	k.cdc.MustUnmarshal(value, &oracle)
	return oracle, true
}

func (k Keeper) DelOracle(ctx sdk.Context, oracleAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleKey(oracleAddr)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

func (k Keeper) GetAllOracles(ctx sdk.Context, isOnline bool) (oracles types.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oracle types.Oracle
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		if isOnline && !oracle.Online {
			continue
		}
		oracles = append(oracles, oracle)
	}
	return oracles
}

// --- LAST ORACLE SLASH BLOCK HEIGHT --- //

func (k Keeper) SlashOracle(ctx sdk.Context, oracleAddrStr string) {
	oracleAddr := sdk.MustAccAddressFromBech32(oracleAddrStr)
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		panic(types.ErrNoFoundOracle)
	}
	if !oracle.Online {
		return
	}
	if !ctx.IsCheckTx() {
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, "oracle_status"},
			float32(1),
			[]metrics.Label{
				telemetry.NewLabel("module", k.moduleName),
				telemetry.NewLabel("address", oracle.OracleAddress),
			},
		)
	}

	oracle.Online = false
	oracle.SlashTimes++
	k.SetOracle(ctx, oracle)
	k.SetLastOracleSlashBlockHeight(ctx, uint64(ctx.BlockHeight()))
}

// SetLastOracleSlashBlockHeight sets the last proposal block height
func (k Keeper) SetLastOracleSlashBlockHeight(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastOracleSlashBlockHeight, sdk.Uint64ToBigEndian(blockHeight))
}

// GetLastOracleSlashBlockHeight returns the last proposal block height
func (k Keeper) GetLastOracleSlashBlockHeight(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	return sdk.BigEndianToUint64(store.Get(types.LastOracleSlashBlockHeight))
}
