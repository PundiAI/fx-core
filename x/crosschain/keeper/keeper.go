package keeper

import (
	"fmt"
	"math"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/x/crosschain/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	moduleName string
	cdc        codec.BinaryCodec // The wire codec for binary encoding/decoding.
	storeKey   sdk.StoreKey      // Unexposed key to access store from sdk.Context
	paramSpace paramtypes.Subspace

	bankKeeper        types.BankKeeper
	ibcTransferKeeper types.IBCTransferKeeper
	ibcChannelKeeper  types.IBCChannelKeeper
	erc20Keeper       types.Erc20Keeper
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(cdc codec.BinaryCodec, moduleName string, storeKey sdk.StoreKey, paramSpace paramtypes.Subspace,
	bankKeeper types.BankKeeper, ibcTransferKeeper types.IBCTransferKeeper, channelKeeper types.IBCChannelKeeper, erc20Keeper types.Erc20Keeper) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	// set KeyTable if it has not already been set
	k := Keeper{
		moduleName:        moduleName,
		cdc:               cdc,
		storeKey:          storeKey,
		paramSpace:        paramSpace,
		bankKeeper:        bankKeeper,
		ibcTransferKeeper: ibcTransferKeeper,
		ibcChannelKeeper:  channelKeeper,
		erc20Keeper:       erc20Keeper,
	}
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+k.moduleName)
}

/////////////////////////////
//       PARAMETERS        //
/////////////////////////////

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, ps types.Params) {
	k.paramSpace.SetParamSet(ctx, &ps)
}

func (k Keeper) ChainHasInit(ctx sdk.Context) bool {
	return k.paramSpace.Has(ctx, types.ParamsStoreKeyGravityID)
}

// GetGravityID returns the GravityID the GravityID is essentially a salt value
// for bridge signatures, provided each chain running Gravity has a unique ID
// it won't be possible to play back signatures from one bridge onto another
// even if they share a oracle set.
//
// The lifecycle of the GravityID is that it is set in the Genesis file
// read from the live chain for the contract deployment, once a Gravity contract
// is deployed the GravityID CAN NOT BE CHANGED. Meaning that it can't just be the
// same as the chain id since the chain id may be changed many times with each
// successive chain in charge of the same bridge
func (k Keeper) GetGravityID(ctx sdk.Context) string {
	var gravityId string
	k.paramSpace.Get(ctx, types.ParamsStoreKeyGravityID, &gravityId)
	return gravityId
}

func (k Keeper) GetOracleDepositThreshold(ctx sdk.Context) sdk.Coin {
	var depositThreshold sdk.Coin
	k.paramSpace.Get(ctx, types.ParamOracleDepositThreshold, &depositThreshold)
	return depositThreshold
}

func (k Keeper) SetChainOracles(ctx sdk.Context, chainOracle *types.ChainOracle) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyChainOracles, k.cdc.MustMarshal(chainOracle))
}

func (k Keeper) GetChainOracles(ctx sdk.Context) (chainOracle types.ChainOracle, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyChainOracles)
	if bz == nil {
		return chainOracle, false
	}
	k.cdc.MustUnmarshal(bz, &chainOracle)
	return chainOracle, true
}

/////////////////////////////
//        Oracle           //
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
func (k Keeper) GetAllOracles(ctx sdk.Context) (oracles types.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oracle types.Oracle
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		oracles = append(oracles, oracle)
	}
	sort.Sort(oracles)
	return oracles
}

func (k Keeper) GetAllActiveOracles(ctx sdk.Context) (oracles types.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oracle types.Oracle
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		if oracle.Jailed {
			continue
		}
		oracles = append(oracles, oracle)
	}
	sort.Sort(oracles)
	return oracles
}

/////////////////////////////
//    ADDRESS DELEGATION   //
/////////////////////////////

// SetOracleByOrchestrator sets the Orchestrator key for a given oracle
func (k Keeper) SetOracleByOrchestrator(ctx sdk.Context, oracleAddr sdk.AccAddress, orchestratorAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	// save  oracle -> orchestrator
	store.Set(types.GetOracleAddressByOrchestratorKey(orchestratorAddr), oracleAddr.Bytes())
}

// GetOracleAddressByOrchestratorKey returns the oracle key associated with an orchestrator key
func (k Keeper) GetOracleAddressByOrchestratorKey(ctx sdk.Context, orchestratorAddr sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	oracle := store.Get(types.GetOracleAddressByOrchestratorKey(orchestratorAddr))
	return oracle, oracle != nil
}

// DelOracleByOrchestrator delete the Orchestrator key for a given oracle
func (k Keeper) DelOracleByOrchestrator(ctx sdk.Context, orchestratorAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleAddressByOrchestratorKey(orchestratorAddr)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

/////////////////////////////
//    External ADDRESS     //
/////////////////////////////

// SetExternalAddressForOracle sets the external address for a given oracle
func (k Keeper) SetExternalAddressForOracle(ctx sdk.Context, oracle sdk.AccAddress, externalAddress string) {
	store := ctx.KVStore(k.storeKey)
	// save external address -> oracle
	store.Set(types.GetOracleAddressByExternalKey(externalAddress), oracle.Bytes())
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
//   ORACLE SET REQUESTS   //
/////////////////////////////

// GetCurrentOracleSet gets powers from the store and normalizes them
// into an integer percentage with a resolution of uint32 Max meaning
// a given validators 'gravity power' is computed as
// Cosmos power / total cosmos power = x / uint32 Max
// where x is the voting power on the gravity contract. This allows us
// to only use integer division which produces a known rounding error
// from truncation equal to the ratio of the validators
// Cosmos power / total cosmos power ratio, leaving us at uint32 Max - 1
// total voting power. This is an acceptable rounding error since floating
// point may cause consensus problems if different floating point unit
// implementations are involved.
func (k Keeper) GetCurrentOracleSet(ctx sdk.Context) *types.OracleSet {
	allOracles := k.GetAllActiveOracles(ctx)
	var bridgeValidators []*types.BridgeValidator
	var totalPower uint64

	for _, oracle := range allOracles {
		power := oracle.GetPower()
		if power.LTE(sdk.ZeroInt()) {
			continue
		}
		totalPower += power.Uint64()
		bridgeValidators = append(bridgeValidators, &types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: oracle.ExternalAddress,
		})
	}
	// normalize power values
	for i := range bridgeValidators {
		bridgeValidators[i].Power = sdk.NewUint(bridgeValidators[i].Power).MulUint64(math.MaxUint32).QuoUint64(totalPower).Uint64()
	}

	oracleSetNonce := k.GetLatestOracleSetNonce(ctx) + 1
	return types.NewOracleSet(oracleSetNonce, uint64(ctx.BlockHeight()), bridgeValidators)
}

// AddOracleSetRequest returns a new instance of the Gravity BridgeValidatorSet
func (k Keeper) AddOracleSetRequest(ctx sdk.Context, currentOracleSet *types.OracleSet, gravityId string) *types.OracleSet {
	// if currentOracleSet member is empty, not store OracleSet.
	if len(currentOracleSet.Members) <= 0 {
		return currentOracleSet
	}
	k.StoreOracleSet(ctx, currentOracleSet)

	k.CommonSetOracleTotalPower(ctx)

	checkpoint := currentOracleSet.GetCheckpoint(gravityId)
	k.SetPastExternalSignatureCheckpoint(ctx, checkpoint)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOracleSetUpdate,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOracleSetNonce, fmt.Sprint(currentOracleSet.Nonce)),
		sdk.NewAttribute(types.AttributeKeyOracleSetLen, fmt.Sprint(len(currentOracleSet.Members))),
	))

	return currentOracleSet
}

// StoreOracleSet is for storing a oracle set at a given height
func (k Keeper) StoreOracleSet(ctx sdk.Context, oracleSet *types.OracleSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleSetKey(oracleSet.Nonce), k.cdc.MustMarshal(oracleSet))
	k.SetLatestOracleSetNonce(ctx, oracleSet.Nonce)
}

// HasOracleSetRequest returns true if a oracleSet defined by a nonce exists
func (k Keeper) HasOracleSetRequest(ctx sdk.Context, nonce uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOracleSetKey(nonce))
}

// DeleteOracleSet deletes the oracleSet at a given nonce from state
func (k Keeper) DeleteOracleSet(ctx sdk.Context, nonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetOracleSetKey(nonce))
}

// SetLatestOracleSetNonce sets the latest oracleSet nonce
func (k Keeper) SetLatestOracleSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestOracleSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLatestOracleSetNonce returns the latest oracleSet nonce
func (k Keeper) GetLatestOracleSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LatestOracleSetNonce)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// GetOracleSet returns a oracleSet by nonce
func (k Keeper) GetOracleSet(ctx sdk.Context, nonce uint64) *types.OracleSet {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOracleSetKey(nonce))
	if bz == nil {
		return nil
	}
	var oracleSet types.OracleSet
	k.cdc.MustUnmarshal(bz, &oracleSet)
	return &oracleSet
}

// IterateOracleSets retruns all oracleSetRequests
func (k Keeper) IterateOracleSets(ctx sdk.Context, cb func(key []byte, val *types.OracleSet) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleSetRequestKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var oracleSet types.OracleSet
		k.cdc.MustUnmarshal(iter.Value(), &oracleSet)
		// cb returns true to stop early
		if cb(iter.Key(), &oracleSet) {
			break
		}
	}
}

// GetOracleSets returns all the oracle sets in state
func (k Keeper) GetOracleSets(ctx sdk.Context) (out []*types.OracleSet) {
	k.IterateOracleSets(ctx, func(_ []byte, val *types.OracleSet) bool {
		out = append(out, val)
		return false
	})
	sort.Sort(types.OracleSets(out))
	return
}

// GetLatestOracleSet returns the latest oracle set in state
func (k Keeper) GetLatestOracleSet(ctx sdk.Context) *types.OracleSet {
	latestOracleSetNonce := k.GetLatestOracleSetNonce(ctx)
	return k.GetOracleSet(ctx, latestOracleSetNonce)
}

// SetLastSlashedOracleSetNonce sets the latest slashed oracleSet nonce
func (k Keeper) SetLastSlashedOracleSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedOracleSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLastSlashedOracleSetNonce returns the latest slashed oracleSet nonce
func (k Keeper) GetLastSlashedOracleSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LastSlashedOracleSetNonce)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// SetLastProposalBlockHeight sets the last proposal block height
func (k Keeper) SetLastProposalBlockHeight(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastProposalBlockHeight, sdk.Uint64ToBigEndian(blockHeight))
}

// GetLastProposalBlockHeight returns the last proposal block height
func (k Keeper) GetLastProposalBlockHeight(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LastProposalBlockHeight)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// SetLastOracleSlashBlockHeight sets the last proposal block height
func (k Keeper) SetLastOracleSlashBlockHeight(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastOracleSlashBlockHeight, sdk.Uint64ToBigEndian(blockHeight))
}

// GetLastOracleSlashBlockHeight returns the last proposal block height
func (k Keeper) GetLastOracleSlashBlockHeight(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LastOracleSlashBlockHeight)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

/////////////////////////////
//   ORACLE SET CONFIRMS   //
/////////////////////////////

// GetUnSlashedOracleSets returns all the unSlashed oracle sets in state
func (k Keeper) GetUnSlashedOracleSets(ctx sdk.Context, maxHeight uint64) (oracleSets types.OracleSets) {
	lastSlashedOracleSetNonce := k.GetLastSlashedOracleSetNonce(ctx)
	k.IterateOracleSetBySlashedOracleSetNonce(ctx, lastSlashedOracleSetNonce, maxHeight, func(_ []byte, oracleSet *types.OracleSet) bool {
		if oracleSet.Nonce > lastSlashedOracleSetNonce && maxHeight > oracleSet.Height {
			oracleSets = append(oracleSets, oracleSet)
		}
		return false
	})
	sort.Sort(oracleSets)
	return
}

// IterateOracleSetBySlashedOracleSetNonce iterates through all oracleSet by last slashed oracleSet nonce in ASC order
func (k Keeper) IterateOracleSetBySlashedOracleSetNonce(ctx sdk.Context, lastSlashedOracleSetNonce uint64, maxHeight uint64, cb func([]byte, *types.OracleSet) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleSetRequestKey)
	iter := prefixStore.Iterator(sdk.Uint64ToBigEndian(lastSlashedOracleSetNonce), sdk.Uint64ToBigEndian(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var oracleSet types.OracleSet
		k.cdc.MustUnmarshal(iter.Value(), &oracleSet)
		// cb returns true to stop early
		if cb(iter.Key(), &oracleSet) {
			break
		}
	}
}

// GetOracleSetConfirm returns a oracleSet confirmation by a nonce and external address
func (k Keeper) GetOracleSetConfirm(ctx sdk.Context, nonce uint64, oracleAddr sdk.AccAddress) *types.MsgOracleSetConfirm {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetOracleSetConfirmKey(nonce, oracleAddr))
	if entity == nil {
		return nil
	}
	confirm := types.MsgOracleSetConfirm{}
	k.cdc.MustUnmarshal(entity, &confirm)
	return &confirm
}

// SetOracleSetConfirm sets a oracleSet confirmation
func (k Keeper) SetOracleSetConfirm(ctx sdk.Context, oracleAddr sdk.AccAddress, oracleSetConfirm types.MsgOracleSetConfirm) []byte {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleSetConfirmKey(oracleSetConfirm.Nonce, oracleAddr)
	store.Set(key, k.cdc.MustMarshal(&oracleSetConfirm))
	return key
}

// GetOracleSetConfirms returns all oracle set confirmations by nonce
func (k Keeper) GetOracleSetConfirms(ctx sdk.Context, nonce uint64) (confirms []*types.MsgOracleSetConfirm) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleSetConfirmKey)
	start, end := prefixRange(sdk.Uint64ToBigEndian(nonce))
	iterator := prefixStore.Iterator(start, end)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		confirm := types.MsgOracleSetConfirm{}
		k.cdc.MustUnmarshal(iterator.Value(), &confirm)
		confirms = append(confirms, &confirm)
	}

	return confirms
}

/////////////////////////////
//      BATCH CONFIRMS     //
/////////////////////////////

// IterateOracleSetConfirmByNonce iterates through all oracleSet confirms by nonce in ASC order
// MARK finish-batches: this is where the key is iterated in the old (presumed working) code
func (k Keeper) IterateOracleSetConfirmByNonce(ctx sdk.Context, nonce uint64, cb func([]byte, types.MsgOracleSetConfirm) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleSetConfirmKey)
	iter := prefixStore.Iterator(prefixRange(sdk.Uint64ToBigEndian(nonce)))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := types.MsgOracleSetConfirm{}
		k.cdc.MustUnmarshal(iter.Value(), &confirm)
		// cb returns true to stop early
		if cb(iter.Key(), confirm) {
			break
		}
	}
}

// GetBatchConfirm returns a batch confirmation given its nonce, the token contract, and a oracle address
func (k Keeper) GetBatchConfirm(ctx sdk.Context, nonce uint64, tokenContract string, oracleAddr sdk.AccAddress) *types.MsgConfirmBatch {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetBatchConfirmKey(tokenContract, nonce, oracleAddr))
	if entity == nil {
		return nil
	}
	confirm := types.MsgConfirmBatch{}
	k.cdc.MustUnmarshal(entity, &confirm)
	return &confirm
}

// SetBatchConfirm sets a batch confirmation by a oracle
func (k Keeper) SetBatchConfirm(ctx sdk.Context, oracleAddr sdk.AccAddress, batch *types.MsgConfirmBatch) []byte {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBatchConfirmKey(batch.TokenContract, batch.Nonce, oracleAddr)
	store.Set(key, k.cdc.MustMarshal(batch))
	return key
}

// IterateBatchConfirmByNonceAndTokenContract iterates through all batch confirmations
// MARK finish-batches: this is where the key is iterated in the old (presumed working) code
func (k Keeper) IterateBatchConfirmByNonceAndTokenContract(ctx sdk.Context, nonce uint64, tokenContract string, cb func([]byte, types.MsgConfirmBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BatchConfirmKey)
	prefixKey := append([]byte(tokenContract), sdk.Uint64ToBigEndian(nonce)...)
	iter := prefixStore.Iterator(prefixRange(prefixKey))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := types.MsgConfirmBatch{}
		k.cdc.MustUnmarshal(iter.Value(), &confirm)
		// cb returns true to stop early
		if cb(iter.Key(), confirm) {
			break
		}
	}
}

// GetBatchConfirmByNonceAndTokenContract returns the batch confirms
func (k Keeper) GetBatchConfirmByNonceAndTokenContract(ctx sdk.Context, nonce uint64, tokenContract string) (out []types.MsgConfirmBatch) {
	k.IterateBatchConfirmByNonceAndTokenContract(ctx, nonce, tokenContract, func(_ []byte, msg types.MsgConfirmBatch) bool {
		out = append(out, msg)
		return false
	})
	return
}

/////////////////////////////
//     ORACLE DEPOSIT      //
/////////////////////////////

func (k Keeper) SetTotalDeposit(ctx sdk.Context, totalDeposit sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OracleTotalDepositKey, []byte(totalDeposit.String()))
}

func (k Keeper) GetTotalDeposit(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	deposit := store.Get(types.OracleTotalDepositKey)
	if deposit == nil {
		return sdk.Coin{}
	}
	depositCoin, err := sdk.ParseCoinNormalized(string(deposit))
	if err != nil {
		panic("invalid oracle total deposit" + err.Error())
	}
	return depositCoin
}

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

// SetLastTotalPower Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: power})
	store.Set(types.LastTotalPowerKey, bz)
}

func (k Keeper) CommonSetOracleTotalPower(ctx sdk.Context) {
	oracles := k.GetAllActiveOracles(ctx)
	totalPower := sdk.ZeroInt()
	for _, oracle := range oracles {
		totalPower = totalPower.Add(oracle.GetPower())
	}
	k.SetLastTotalPower(ctx, totalPower)
}

func (k Keeper) UnpackAttestationClaim(att *types.Attestation) (types.ExternalClaim, error) {
	var msg types.ExternalClaim
	err := k.cdc.UnpackAny(att.Claim, &msg)
	return msg, err
}

func (k Keeper) IsOracle(ctx sdk.Context, oracleAddr string) bool {
	chainOracle, found := k.GetChainOracles(ctx)
	if !found {
		return false
	}
	for _, oracle := range chainOracle.Oracles {
		if oracle == oracleAddr {
			return true
		}
	}
	return false
}

// prefixRange turns a prefix into a (start, end) range. The start is the given prefix value and
// the end is calculated by adding 1 bit to the start value. Nil is not allowed as prefix.
// 		Example: []byte{1, 3, 4} becomes []byte{1, 3, 5}
// 				 []byte{15, 42, 255, 255} becomes []byte{15, 43, 0, 0}
//
// In case of an overflow the end is set to nil.
//		Example: []byte{255, 255, 255, 255} becomes nil
// MARK finish-batches: this is where some crazy shit happens
func prefixRange(prefix []byte) ([]byte, []byte) {
	if prefix == nil {
		panic("nil key not allowed")
	}
	// special case: no prefix is whole range
	if len(prefix) == 0 {
		return nil, nil
	}

	// copy the prefix and update last byte
	end := make([]byte, len(prefix))
	copy(end, prefix)
	l := len(end) - 1
	end[l]++

	// wait, what if that overflowed?....
	for end[l] == 0 && l > 0 {
		l--
		end[l]++
	}

	// okay, funny guy, you gave us FFF, no end to this range...
	if l == 0 && end[0] == 0 {
		end = nil
	}
	return prefix, end
}

//SetIbcSequenceHeight set gravity -> ibc sequence block height.
func (k Keeper) SetIbcSequenceHeight(ctx sdk.Context, sourcePort, sourceChannel string, sequence, height uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIbcSequenceHeightKey(sourcePort, sourceChannel, sequence), sdk.Uint64ToBigEndian(height))
}

//GetIbcSequenceHeight get gravity -> ibc sequence block height.
func (k Keeper) GetIbcSequenceHeight(ctx sdk.Context, sourcePort, sourceChannel string, sequence uint64) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIbcSequenceHeightKey(sourcePort, sourceChannel, sequence)
	if !store.Has(key) {
		return 0, false
	}
	value := store.Get(key)
	return sdk.BigEndianToUint64(value), true
}

//setLastEventBlockHeightByOracle set the latest event blockHeight for a give oracle
func (k Keeper) setLastEventBlockHeightByOracle(ctx sdk.Context, oracle sdk.AccAddress, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventBlockHeightByOracleKey(oracle), sdk.Uint64ToBigEndian(blockHeight))
}

//getLastEventBlockHeightByOracle get the latest event blockHeight for a give oracle
func (k Keeper) getLastEventBlockHeightByOracle(ctx sdk.Context, oracle sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLastEventBlockHeightByOracleKey(oracle)
	if !store.Has(key) {
		return 0
	}
	data := store.Get(key)
	return sdk.BigEndianToUint64(data)
}

func (k Keeper) GetModuleName() string {
	return k.moduleName
}

func (k Keeper) GetParamSpace() paramtypes.Subspace {
	return k.paramSpace
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}
