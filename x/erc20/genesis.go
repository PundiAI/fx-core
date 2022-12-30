package erc20

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/keeper"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// InitGenesis import module genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, accountKeeper authkeeper.AccountKeeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)

	// ensure erc20 module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		panic("the erc20 module account has not been set")
	}

	for _, pair := range data.TokenPairs {
		id := pair.GetID()
		k.SetTokenPair(ctx, pair)
		k.SetDenomMap(ctx, pair.Denom, id)
		k.SetERC20Map(ctx, pair.GetERC20Contract(), id)
	}

	// init logic contract
	initContract := []fxtypes.Contract{fxtypes.GetERC20(), fxtypes.GetWFX()}
	for _, contract := range initContract {
		if len(contract.Code) <= 0 || contract.Address == common.HexToAddress(fxtypes.EmptyEvmAddress) {
			panic(fmt.Sprintf("invalid contract: %s/%s", contract.Address.String(), contract.Version))
		}
		if err := k.CreateContractWithCode(ctx, contract.Address, contract.Code); err != nil {
			panic(fmt.Sprintf("create contract %s/%s with code error %s", contract.Address.String(), contract.Version, err.Error()))
		}
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:     k.GetParams(ctx),
		TokenPairs: k.GetAllTokenPairs(ctx),
	}
}
