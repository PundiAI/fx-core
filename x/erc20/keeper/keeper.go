package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	accountKeeper     types.AccountKeeper
	bankKeeper        types.BankKeeper
	evmKeeper         types.EVMKeeper
	ibcTransferKeeper types.IBCTransferKeeper

	router *fxtypes.Router

	moduleAddress common.Address
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	evmKeeper types.EVMKeeper,
	ibcTransferKeeper types.IBCTransferKeeper,
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
		ibcTransferKeeper: ibcTransferKeeper,
		moduleAddress:     common.BytesToAddress(ak.GetModuleAddress(types.ModuleName)),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetRouter(rtr fxtypes.Router) Keeper {
	if k.router != nil && k.router.Sealed() {
		panic("cannot reset a sealed router")
	}
	if _, found := rtr.GetRoute(types.ModuleName); found {
		panic("cannot set current module")
	}
	k.router = &rtr
	k.router.Seal()
	return k
}

// TransferAfter ibc transfer after
func (k Keeper) TransferAfter(ctx sdk.Context, sender, receive string, coin, fee sdk.Coin) error {
	sendAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender address %s, error %s", sender, err.Error())
	}
	if err = fxtypes.ValidateEthereumAddress(receive); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "receiver address %s, error %s", receive, err.Error())
	}
	return k.RelayConvertCoin(ctx, sendAddr, common.HexToAddress(receive), coin.Add(fee))
}

func (k Keeper) RelayConvertCoin(ctx sdk.Context, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	if !k.IsDenomRegistered(ctx, coin.Denom) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom(%s) not registered", coin.Denom)
	}
	msg := &types.MsgConvertCoin{
		Coin:     coin,
		Receiver: receiver.Hex(),
		Sender:   sender.String(),
	}
	_, err := k.ConvertCoin(sdk.WrapSDKContext(ctx), msg)
	return err
}

func (k Keeper) RelayConvertDenomToOne(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin) (sdk.Coin, error) {
	return k.ConvertDenomToOne(ctx, from, coin)
}

func (k Keeper) RelayConvertDenomToMany(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, target string) (sdk.Coin, error) {
	// convert denom
	cacheCtx, commit := ctx.CacheContext()
	targetCoin, err := k.ConvertDenomToMany(cacheCtx, from, coin, target)
	if err != nil {
		return sdk.Coin{}, err
	}
	commit()
	return targetCoin, nil
}

func (k Keeper) HasDenomAlias(ctx sdk.Context, denom string) (banktypes.Metadata, bool) {
	md, found := k.bankKeeper.GetDenomMetaData(ctx, denom)
	// not register metadata
	if !found {
		return banktypes.Metadata{}, false
	}
	// not have denom units
	if len(md.DenomUnits) == 0 {
		return banktypes.Metadata{}, false
	}
	//not have alias
	if len(md.DenomUnits[0].Aliases) == 0 {
		return banktypes.Metadata{}, false
	}
	return md, true
}
