package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
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

// --- ADDRESS Bridger --- //

// SetOracleByBridger sets the bridger key for a given oracle
func (k Keeper) SetOracleByBridger(ctx sdk.Context, bridgerAddr, oracleAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleAddressByBridgerKey(bridgerAddr), oracleAddr.Bytes())
}

// GetOracleAddressByBridgerKey returns the oracle key associated with an bridger key
func (k Keeper) GetOracleAddressByBridgerKey(ctx sdk.Context, bridgerAddr sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	oracle := store.Get(types.GetOracleAddressByBridgerKey(bridgerAddr))
	return oracle, oracle != nil
}

// DelOracleByBridger delete the bridger key for a given oracle
func (k Keeper) DelOracleByBridger(ctx sdk.Context, bridgerAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleAddressByBridgerKey(bridgerAddr)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

// --- External ADDRESS --- //

// SetOracleByExternalAddress sets the external address for a given oracle
func (k Keeper) SetOracleByExternalAddress(ctx sdk.Context, externalAddress string, oracleAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleAddressByExternalKey(externalAddress), oracleAddr.Bytes())
}

// GetOracleByExternalAddress returns the external address for a given gravity oracle
func (k Keeper) GetOracleByExternalAddress(ctx sdk.Context, externalAddress string) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOracleAddressByExternalKey(externalAddress))
	return bz, bz != nil
}

// DelOracleByExternalAddress delete the external address for a give oracle
func (k Keeper) DelOracleByExternalAddress(ctk sdk.Context, externalAddress string) {
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
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdkmath.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastTotalPowerKey, k.cdc.MustMarshal(&sdk.IntProto{Int: power}))
}

func (k Keeper) CommonSetOracleTotalPower(ctx sdk.Context) {
	oracles := k.GetAllOracles(ctx, true)
	totalPower := sdkmath.ZeroInt()
	for _, oracle := range oracles {
		totalPower = totalPower.Add(oracle.GetPower())
	}
	k.SetLastTotalPower(ctx, totalPower)
}

// --- ORACLES --- //

func (k Keeper) IterateOracle(ctx sdk.Context, cb func(oracle types.Oracle) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		oracle := types.Oracle{}
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		if cb(oracle) {
			break
		}
	}
}

// SetOracle save Oracle data
func (k Keeper) SetOracle(ctx sdk.Context, oracle types.Oracle) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleKey(oracle.GetOracle()), k.cdc.MustMarshal(&oracle))
}

func (k Keeper) HasOracle(ctx sdk.Context, addr sdk.AccAddress) (found bool) {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOracleKey(addr))
}

// GetOracle get Oracle data
func (k Keeper) GetOracle(ctx sdk.Context, addr sdk.AccAddress) (oracle types.Oracle, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetOracleKey(addr))
	if value == nil {
		return oracle, false
	}
	k.cdc.MustUnmarshal(value, &oracle)
	return oracle, true
}

func (k Keeper) DelOracle(ctx sdk.Context, oracle sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleKey(oracle)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

func (k Keeper) GetAllOracles(ctx sdk.Context, isOnline bool) (oracles types.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
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

func (k Keeper) SlashOracle(ctx sdk.Context, oracleAddrStr string) {
	oracleAddr := sdk.MustAccAddressFromBech32(oracleAddrStr)
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		panic(types.ErrNoFoundOracle)
	}
	if !oracle.Online {
		return
	}

	oracle.Online = false
	oracle.SlashTimes += 1
	k.SetOracle(ctx, oracle)
	k.SetLastOracleSlashBlockHeight(ctx, uint64(ctx.BlockHeight()))
}
