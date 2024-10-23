package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/x/erc20/types"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	evmKeeper      types.EVMKeeper
	evmErc20Keeper types.ERC20TokenKeeper

	contractOwner common.Address

	authority string

	Schema      collections.Schema
	Params      collections.Item[types.Params]
	ERC20Token  collections.Map[string, types.ERC20Token]                            // baseDenom -> ERC20Token
	BridgeToken collections.Map[collections.Pair[string, string], types.BridgeToken] // baseDenom -> BridgeToken
	IBCToken    collections.Map[collections.Pair[string, string], types.IBCToken]    // baseDenom -> IBCToken
	DenomIndex  collections.Map[string, string]                                      // bridgeDenom/erc20_contract/ibc_denom -> baseDenom
	Cache       collections.Map[string, sdkmath.Int]                                 // crosschain cache
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	evmKeeper types.EVMKeeper,
	evmErc20Keeper types.ERC20TokenKeeper,
	authority string,
) Keeper {
	moduleAddress := ak.GetModuleAddress(types.ModuleName)
	if moduleAddress == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		accountKeeper:  ak,
		bankKeeper:     bk,
		evmKeeper:      evmKeeper,
		evmErc20Keeper: evmErc20Keeper,
		contractOwner:  common.BytesToAddress(moduleAddress),
		authority:      authority,
		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		ERC20Token:     collections.NewMap(sb, types.ERC20TokenKey, "erc20_token", collections.StringKey, codec.CollValue[types.ERC20Token](cdc)),
		BridgeToken:    collections.NewMap(sb, types.BridgeTokenKey, "bridge_token", collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.BridgeToken](cdc)),
		IBCToken:       collections.NewMap(sb, types.IBCTokenKey, "ibc_token", collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.IBCToken](cdc)),
		DenomIndex:     collections.NewMap(sb, types.DenomIndexKey, "denom_index", collections.StringKey, collections.StringValue),
		Cache:          collections.NewMap(sb, types.CacheKey, "cache", collections.StringKey, sdk.IntValue),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) MintingEnabled(ctx context.Context, receiver sdk.AccAddress, isBaseDenom bool, tokenName string) (erc20Token types.ERC20Token, err error) {
	if err = k.CheckEnableErc20(ctx); err != nil {
		return types.ERC20Token{}, err
	}

	if isBaseDenom {
		erc20Token, err = k.ERC20Token.Get(ctx, tokenName)
		if err != nil {
			return types.ERC20Token{}, err
		}
	} else {
		erc20Token, err = k.GetERC20Token(ctx, tokenName)
		if err != nil {
			return types.ERC20Token{}, err
		}
	}

	if !erc20Token.Enabled {
		return erc20Token, types.ErrDisabled.Wrapf("token %s is disabled", erc20Token.Denom)
	}

	if k.bankKeeper.BlockedAddr(receiver.Bytes()) {
		return erc20Token, sdkerrors.ErrUnauthorized.Wrapf("%s is not allowed to receive transactions", receiver)
	}

	if !k.bankKeeper.IsSendEnabledDenom(ctx, erc20Token.Denom) {
		return erc20Token, banktypes.ErrSendDisabled.Wrapf("minting '%s' denom is currently disabled", tokenName)
	}
	return erc20Token, nil
}
