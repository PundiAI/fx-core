package keeper

import (
	"sort"
	"time"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/crosschain/types"
)

/////////////////////////////
//    PROPOSAL ORACLE      //
/////////////////////////////

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

/////////////////////////////
//     ADDRESS Bridger     //
/////////////////////////////

// SetOracleByBridger sets the bridger key for a given oracle
func (k Keeper) SetOracleByBridger(ctx sdk.Context, oracleAddr sdk.AccAddress, bridgerAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	// save external oracleAddr -> bridgerAddr
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

/////////////////////////////
//    External ADDRESS     //
/////////////////////////////

// SetExternalAddressForOracle sets the external address for a given oracle
func (k Keeper) SetExternalAddressForOracle(ctx sdk.Context, oracleAddr sdk.AccAddress, externalAddress string) {
	store := ctx.KVStore(k.storeKey)
	// save external address -> oracleAddr
	store.Set(types.GetOracleAddressByExternalKey(externalAddress), oracleAddr.Bytes())
}

// GetOracleByExternalAddress returns the external address for a given gravity oracle
func (k Keeper) GetOracleByExternalAddress(ctx sdk.Context, externalAddress string) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOracleAddressByExternalKey(externalAddress))
	return bz, bz != nil
}

// DelExternalAddressForOracle delete the external address for a give oracle
func (k Keeper) DelExternalAddressForOracle(ctk sdk.Context, externalAddress string) {
	store := ctk.KVStore(k.storeKey)
	oracleAddr := types.GetOracleAddressByExternalKey(externalAddress)
	if !store.Has(oracleAddr) {
		return
	}
	store.Delete(oracleAddr)
}

/////////////////////////////
//   ORACLE TOTAL POWER    //
/////////////////////////////

// GetLastTotalPower Load the last total oracle power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastTotalPowerKey)

	if bz == nil {
		return sdk.ZeroInt()
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshal(bz, &ip)

	return ip.Int
}

// setLastTotalPower Set the last total validator power.
func (k Keeper) setLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: power})
	store.Set(types.LastTotalPowerKey, bz)
}

func (k Keeper) CommonSetOracleTotalPower(ctx sdk.Context) {
	oracles := k.GetAllOracles(ctx, true)
	totalPower := sdk.ZeroInt()
	for _, oracle := range oracles {
		totalPower = totalPower.Add(oracle.GetPower())
	}
	k.setLastTotalPower(ctx, totalPower)
}

/////////////////////////////
//        ORACLES          //
/////////////////////////////

// SetOracle save Oracle data
func (k Keeper) SetOracle(ctx sdk.Context, oracle types.Oracle) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&oracle)
	store.Set(types.GetOracleKey(oracle.GetOracle()), bz)
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

// GetAllOracles
func (k Keeper) GetAllOracles(ctx sdk.Context, isActive bool) (oracles types.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oracle types.Oracle
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		if isActive && oracle.Jailed {
			continue
		}
		oracles = append(oracles, oracle)
	}
	sort.Sort(oracles)
	return oracles
}

func (k Keeper) SlashOracle(ctx sdk.Context, oracle types.Oracle, slashFraction sdk.Dec) {
	if oracle.Jailed {
		return
	}

	if oracle.IsValidator {
		valAddr, err := sdk.ValAddressFromBech32(oracle.DelegateValidator)
		if err != nil {
			panic(err)
		}
		validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
		if found && !validator.IsJailed() {
			consAddr, _ := validator.GetConsAddr()
			power := validator.ConsensusPower(sdk.DefaultPowerReduction)
			k.stakingKeeper.Slash(ctx, consAddr, ctx.BlockHeight(), power, slashFraction)
			k.stakingKeeper.Jail(ctx, consAddr)
		}
	} else {
		slashAmount := oracle.DelegateAmount.ToDec().Mul(slashFraction).TruncateInt()
		slashAmount = sdk.MinInt(slashAmount, oracle.DelegateAmount)
		slashAmount = sdk.MaxInt(slashAmount, sdk.ZeroInt())
		if slashAmount.IsPositive() {
			oracle.DelegateAmount = oracle.DelegateAmount.Sub(slashAmount)
			oracleAddr, err := sdk.AccAddressFromBech32(oracle.OracleAddress)
			if err != nil {
				panic(err)
			}
			delegateAddr := types.GetOracleDelegateAddress(k.moduleName, oracleAddr)
			valAddr, err := sdk.ValAddressFromBech32(oracle.DelegateValidator)
			if err != nil {
				panic(err)
			}
			sharesAmount, err := k.stakingKeeper.ValidateUnbondAmount(ctx, delegateAddr, valAddr, slashAmount)
			if err != nil {
				panic(err)
			}
			completionTime, err := k.stakingKeeper.Undelegate(ctx, delegateAddr, valAddr, sharesAmount)
			if err != nil {
				panic(err)
			}
			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					stakingtypes.EventTypeUnbond,
					sdk.NewAttribute(stakingtypes.AttributeKeyValidator, oracle.DelegateValidator),
					sdk.NewAttribute(sdk.AttributeKeyAmount, slashAmount.String()),
					sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
				),
			})
		}
	}

	oracle.Jailed = true
	oracle.JailedHeight = ctx.BlockHeight()
	k.SetOracle(ctx, oracle)
	k.SetLastOracleSlashBlockHeight(ctx, uint64(ctx.BlockHeight()))
}
