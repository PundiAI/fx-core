package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/x/gov/types"
)

type Keeper struct {
	*govkeeper.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey storetypes.StoreKey

	authKeeper govtypes.AccountKeeper
	bankKeeper govtypes.BankKeeper
	sk         govtypes.StakingKeeper

	cdc codec.BinaryCodec

	authority string

	storeKeys      map[string]*storetypes.KVStoreKey
	CustomerParams collections.Map[string, types.CustomParams]
}

func NewKeeper(storeService corestoretypes.KVStoreService, ak govtypes.AccountKeeper, bk govtypes.BankKeeper, sk govtypes.StakingKeeper, keys map[string]*storetypes.KVStoreKey, gk *govkeeper.Keeper, cdc codec.BinaryCodec, authority string) *Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	sb := collections.NewSchemaBuilder(storeService)
	return &Keeper{
		storeKey:       keys[govtypes.StoreKey],
		authKeeper:     ak,
		bankKeeper:     bk,
		sk:             sk,
		Keeper:         gk,
		cdc:            cdc,
		authority:      authority,
		storeKeys:      keys,
		CustomerParams: collections.NewMap(sb, types.CustomParamsKey, "customParams", collections.StringKey, codec.CollValue[types.CustomParams](cdc)),
	}
}

func (keeper Keeper) InitCustomParams(ctx sdk.Context) error {
	customParamsList := types.DefaultInitGenesisCustomParams()
	for _, customParams := range customParamsList {
		if err := keeper.CustomerParams.Set(ctx, customParams.MsgType, customParams.Params); err != nil {
			return err
		}
	}
	return nil
}

func (keeper Keeper) CheckDisabledPrecompiles(ctx sdk.Context, contractAddress common.Address, methodId []byte) error {
	switchParams := keeper.GetSwitchParams(ctx)
	return CheckContractAddressIsDisabled(switchParams.DisablePrecompiles, contractAddress, methodId)
}

func CheckContractAddressIsDisabled(disabledPrecompiles []string, addr common.Address, methodId []byte) error {
	if len(disabledPrecompiles) == 0 {
		return nil
	}

	addrStr := strings.ToLower(addr.String())
	methodIdStr := hex.EncodeToString(methodId)
	addrMethodId := fmt.Sprintf("%s/%s", addrStr, methodIdStr)
	for _, disabledPrecompile := range disabledPrecompiles {
		disabledPrecompile = strings.ToLower(disabledPrecompile)
		if disabledPrecompile == addrStr {
			return errors.New("precompile address is disabled")
		}

		if disabledPrecompile == addrMethodId {
			return fmt.Errorf("precompile method %s is disabled", methodIdStr)
		}
	}

	return nil
}
