package keeper

import (
	"encoding/binary"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	autytypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
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
		callbackFrom:       common.BytesToAddress(autytypes.NewModuleAddress(types.ModuleName)),
	}
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) GetCallbackFrom() common.Address {
	return k.callbackFrom
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

func (k Keeper) BridgeCallEvm(ctx sdk.Context, sender, refundAddr common.Address, coins sdk.Coins, to common.Address, data, memo []byte, value sdkmath.Int, isMemoSendCallTo bool) error {
	if !k.evmKeeper.IsContract(ctx, to) {
		return nil
	}
	var callEvmSender common.Address
	var args []byte

	if isMemoSendCallTo {
		args = data
		callEvmSender = sender
	} else {
		callTokens, callAmounts := k.CoinsToBridgeCallTokens(ctx, coins)
		var err error
		args, err = types.PackBridgeCallback(sender, refundAddr, callTokens, callAmounts, data, memo)
		if err != nil {
			return err
		}
		callEvmSender = k.GetCallbackFrom()
	}

	gasLimit := k.GetBridgeCallMaxGasLimit(ctx)
	txResp, err := k.evmKeeper.CallEVM(ctx, callEvmSender, &to, value.BigInt(), gasLimit, args, true)
	if err != nil {
		return err
	}
	if txResp.Failed() {
		return errorsmod.Wrap(types.ErrInvalid, txResp.VmError)
	}
	return nil
}
