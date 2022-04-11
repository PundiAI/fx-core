package keeper

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/core"
	"github.com/functionx/fx-core/x/erc20/types"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	evmKeeper     types.EVMKeeper
	// fetch EIP1559 base fee and parameters
	feeMarketKeeper types.FeeMarketKeeper

	ibcTransferKeeper types.IBCTransferKeeper
	ibcChannelKeeper  types.IBCChannelKeeper
}

func (k Keeper) RefundAfter(ctx sdk.Context, sourcePort, sourceChannel string, sequence uint64, sender sdk.AccAddress, receiver string, amount sdk.Coin) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) TransferAfter(ctx sdk.Context, sender, receive string, coins, fee sdk.Coin) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) HasInit(ctx sdk.Context) bool {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) ConvertDenomToFIP20(ctx sdk.Context, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	logger := ctx.Logger()
	if !k.IsDenomRegistered(ctx, coin.Denom) {
		logger.Error("evm transfer, denom not registered", "denom", coin.Denom)
		return nil
	}
	//TODO implement me
	panic("implement me")
}

func (k Keeper) ModuleInit(ctx sdk.Context, enableErc20, enableEvmHook bool, ibcTransferTimeoutHeight uint64) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	//TODO implement me
	panic("implement me")
}

func (k *Keeper) SetIBCTransferKeeper(ibcTransferKeepr types.IBCTransferKeeper) *Keeper {
	k.ibcTransferKeeper = ibcTransferKeepr
	return k
}

func (k *Keeper) SetIBCChannelKeeper(ibcChannelKeeper types.IBCChannelKeeper) {
	k.ibcChannelKeeper = ibcChannelKeeper
}

func (k Keeper) CreateContractWithCode(ctx sdk.Context, addr common.Address, code []byte) error {
	return k.evmKeeper.CreateContractWithCode(ctx, addr, code)
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	evmKeeper types.EVMKeeper,
	feeMarketKeeper types.FeeMarketKeeper,
	ibcTransferKeeper types.IBCTransferKeeper,
	ibcChannelKeeper types.IBCChannelKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:          storeKey,
		cdc:               cdc,
		paramstore:        ps,
		accountKeeper:     ak,
		bankKeeper:        bk,
		evmKeeper:         evmKeeper,
		feeMarketKeeper:   feeMarketKeeper,
		ibcTransferKeeper: ibcTransferKeeper,
		ibcChannelKeeper:  ibcChannelKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
