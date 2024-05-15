package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
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
	ak                 types.AccountKeeper
	ibcTransferKeeper  types.IBCTransferKeeper
	erc20Keeper        types.Erc20Keeper
	evmKeeper          types.EVMKeeper

	authority    string
	callbackFrom common.Address
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(cdc codec.BinaryCodec, moduleName string, storeKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper, stakingMsgServer types.StakingMsgServer, distributionKeeper types.DistributionMsgServer,
	bankKeeper types.BankKeeper, ibcTransferKeeper types.IBCTransferKeeper, erc20Keeper types.Erc20Keeper, ak types.AccountKeeper,
	evmKeeper types.EVMKeeper, authority string,
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
		ak:                 ak,
		ibcTransferKeeper:  ibcTransferKeeper,
		erc20Keeper:        erc20Keeper,
		evmKeeper:          evmKeeper,
		authority:          authority,
		callbackFrom:       common.BytesToAddress(crypto.Keccak256([]byte(types.ModuleName))),
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

func (k Keeper) ModuleName() string {
	return k.moduleName
}
