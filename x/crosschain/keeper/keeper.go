package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	autytypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
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

	authority          string
	sender             common.Address
	bridgeFeeCollector sdk.AccAddress
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(cdc codec.BinaryCodec, moduleName string, storeKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper, stakingMsgServer types.StakingMsgServer, distributionKeeper types.DistributionMsgServer,
	bankKeeper types.BankKeeper, ibcTransferKeeper types.IBCTransferKeeper, erc20Keeper types.Erc20Keeper, ak types.AccountKeeper,
	evmKeeper types.EVMKeeper, sender common.Address, authority string,
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
		sender:             sender,
		bridgeFeeCollector: autytypes.NewModuleAddress(types.BridgeFeeCollectorName),
	}
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) GetCallSender() common.Address {
	return k.sender
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+k.moduleName)
}

func (k Keeper) ModuleName() string {
	return k.moduleName
}

func (k Keeper) autoIncrementID(ctx sdk.Context, idKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(idKey)
	var id uint64 = 1
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(idKey, bz)
	return id
}
