package keeper

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	fxcoretypes "github.com/functionx/fx-core/types"
	"github.com/tendermint/tendermint/libs/log"

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
	if ctx.BlockHeight() < fxcoretypes.EvmSupportBlock() || !k.HasInit(ctx) {
		ctx.Logger().Info("ignore refund, module not enable", "module", types.ModuleName)
		return nil
	}
	//check tx
	if !k.HashIBCTransferHash(ctx, sourcePort, sourceChannel, sequence) {
		ctx.Logger().Info("ignore refund, transaction not belong to evm ibc transfer", "module", types.ModuleName)
		return nil
	}
	return k.RelayConvertCoin(ctx, sender, common.BytesToAddress(sender.Bytes()), amount)
}

func (k Keeper) TransferAfter(ctx sdk.Context, sender, receive string, coin, fee sdk.Coin) error {
	sendAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return fmt.Errorf("invalid sender address %s, error %s", sender, err.Error())
	}
	if !common.IsHexAddress(receive) {
		return fmt.Errorf("invalid receiver address %s", receive)
	}
	if ctx.BlockHeight() < fxcoretypes.EvmSupportBlock() || !k.HasInit(ctx) {
		return errors.New("erc20 module not enable")
	}
	if !k.IsDenomRegistered(ctx, coin.Denom) {
		return fmt.Errorf("denom %s not resgister", coin.Denom)
	}
	return k.RelayConvertCoin(ctx, sendAddr, common.HexToAddress(receive), coin.Add(fee))
}

func (k Keeper) HasInit(ctx sdk.Context) bool {
	return k.paramstore.Has(ctx, types.ParamStoreKeyEnableErc20)
}

func (k Keeper) RelayConvertCoin(ctx sdk.Context, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	if ctx.BlockHeight() < fxcoretypes.EvmSupportBlock() || !k.HasInit(ctx) {
		return errors.New("erc20 module not enable")
	}
	if !k.IsDenomRegistered(ctx, coin.Denom) {
		return fmt.Errorf("denom(%s) not registered", coin.Denom)
	}
	msg := &types.MsgConvertCoin{
		Coin:     coin,
		Receiver: receiver.Hex(),
		Sender:   sender.String(),
	}
	_, err := k.ConvertCoin(sdk.WrapSDKContext(ctx), msg)
	return err
}

func (k Keeper) ModuleInit(ctx sdk.Context, enableErc20, enableEvmHook bool, ibcTimeout uint64) error {
	k.SetParams(ctx, types.Params{
		EnableErc20:   enableErc20,
		EnableEVMHook: enableEvmHook,
		IbcTimeout:    ibcTimeout,
	})
	// ensure erc20 module account is set on genesis
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		return errors.New("the erc20 module account has not been set")
	}
	//init system contract
	if err := k.InitSystemContract(ctx); err != nil {
		return fmt.Errorf("intrarelayer module init system contract error: %v", err)
	}
	return nil
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
