package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/v6/x/crosschain/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	moduleName string
	cdc        codec.BinaryCodec   // The wire codec for binary encoding/decoding.
	storeKey   storetypes.StoreKey // Unexposed key to access store from sdk.Context

	stakingKeeper      types.StakingKeeper
	stakingMsgServer   types.StakingMsgServer
	distributionKeeper types.DistributionMsgServer
	bankKeeper         types.BankKeeper
	ibcTransferKeeper  types.IBCTransferKeeper
	erc20Keeper        types.Erc20Keeper

	authority string
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(cdc codec.BinaryCodec, moduleName string, storeKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper, stakingMsgServer types.StakingMsgServer, distributionKeeper types.DistributionMsgServer,
	bankKeeper types.BankKeeper, ibcTransferKeeper types.IBCTransferKeeper, erc20Keeper types.Erc20Keeper, ak types.AccountKeeper,
	authority string,
) Keeper {
	if addr := ak.GetModuleAddress(moduleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", moduleName))
	}

	return Keeper{
		moduleName: moduleName,
		cdc:        cdc,
		storeKey:   storeKey,

		stakingKeeper:      stakingKeeper,
		stakingMsgServer:   stakingMsgServer,
		distributionKeeper: distributionKeeper,
		bankKeeper:         bankKeeper,
		ibcTransferKeeper:  ibcTransferKeeper,
		erc20Keeper:        erc20Keeper,
		authority:          authority,
	}
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+k.moduleName)
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

// SetLastEventBlockHeightByOracle set the latest event blockHeight for a give oracle
func (k Keeper) SetLastEventBlockHeightByOracle(ctx sdk.Context, oracleAddr sdk.AccAddress, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventBlockHeightByOracleKey(oracleAddr), sdk.Uint64ToBigEndian(blockHeight))
}

// GetLastEventBlockHeightByOracle get the latest event blockHeight for a give oracle
func (k Keeper) GetLastEventBlockHeightByOracle(ctx sdk.Context, oracleAddr sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLastEventBlockHeightByOracleKey(oracleAddr)
	if !store.Has(key) {
		return 0
	}
	data := store.Get(key)
	return sdk.BigEndianToUint64(data)
}

func (k Keeper) UpdateChainOracles(ctx sdk.Context, oracles []string) error {
	if len(oracles) > types.MaxOracleSize {
		return errorsmod.Wrapf(types.ErrInvalid,
			fmt.Sprintf("oracle length must be less than or equal: %d", types.MaxOracleSize))
	}

	newOracleMap := make(map[string]bool, len(oracles))
	for _, oracle := range oracles {
		newOracleMap[oracle] = true
	}

	var unbondedOracleList []types.Oracle
	totalPower, deleteTotalPower := sdkmath.ZeroInt(), sdkmath.ZeroInt()

	allOracles := k.GetAllOracles(ctx, false)
	proposalOracle, _ := k.GetProposalOracle(ctx)
	oldOracleMap := make(map[string]bool, len(oracles))
	for _, oracle := range proposalOracle.Oracles {
		oldOracleMap[oracle] = true
	}

	for _, oracle := range allOracles {
		if oracle.Online {
			totalPower = totalPower.Add(oracle.GetPower())
		}
		// oracle in new proposal
		if _, ok := newOracleMap[oracle.OracleAddress]; ok {
			continue
		}
		// oracle not in new proposal and oracle in old proposal
		if _, ok := oldOracleMap[oracle.OracleAddress]; ok {
			unbondedOracleList = append(unbondedOracleList, oracle)
			if oracle.Online {
				deleteTotalPower = deleteTotalPower.Add(oracle.GetPower())
			}
		}
	}

	maxChangePowerThreshold := types.AttestationProposalOracleChangePowerThreshold.Mul(totalPower).Quo(sdkmath.NewInt(100))
	k.Logger(ctx).Info("update chain oracles proposal",
		"maxChangePowerThreshold", maxChangePowerThreshold.String(), "deleteTotalPower", deleteTotalPower.String())
	if deleteTotalPower.GT(sdkmath.ZeroInt()) && deleteTotalPower.GTE(maxChangePowerThreshold) {
		return errorsmod.Wrapf(types.ErrInvalid, "max change power, "+
			"maxChangePowerThreshold: %s, deleteTotalPower: %s", maxChangePowerThreshold.String(), deleteTotalPower.String())
	}

	// update proposal oracle
	k.SetProposalOracle(ctx, &types.ProposalOracle{Oracles: oracles})

	for _, unbondedOracle := range unbondedOracleList {
		if err := k.UnbondedOracleFromProposal(ctx, unbondedOracle); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) UnbondedOracleFromProposal(ctx sdk.Context, oracle types.Oracle) error {
	delegateAddr := oracle.GetDelegateAddress(k.moduleName)
	valAddr := oracle.GetValidator()
	getOracleDelegateToken, err := k.GetOracleDelegateToken(ctx, delegateAddr, valAddr)
	if err != nil {
		return err
	}
	msgUndelegate := stakingtypes.NewMsgUndelegate(delegateAddr, valAddr, types.NewDelegateAmount(getOracleDelegateToken))
	if _, err = k.stakingMsgServer.Undelegate(sdk.WrapSDKContext(ctx), msgUndelegate); err != nil {
		return err
	}

	oracle.Online = false
	k.SetOracle(ctx, oracle)

	return nil
}

func (k Keeper) ModuleName() string {
	return k.moduleName
}
