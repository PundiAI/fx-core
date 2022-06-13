package migrate

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/tharsis/ethermint/types"

	"github.com/functionx/fx-core/x/migrate/keeper"
	"github.com/functionx/fx-core/x/migrate/types"
)

var emptyCodeHash = crypto.Keccak256(nil)

// InitGenesis import module genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, _ types.GenesisState) {
	// migrate base account to eth account
	k.AccountKeeper.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {
		if _, ok := account.(ethermint.EthAccountI); ok {
			return false
		}
		baseAccount, ok := account.(*authtypes.BaseAccount)
		if !ok {
			k.Logger(ctx).Info("ignore account", "address", account.GetAddress(), "type", fmt.Sprintf("%T", account))
			return false
		}
		ethAccount := &ethermint.EthAccount{
			BaseAccount: baseAccount,
			CodeHash:    common.BytesToHash(emptyCodeHash).String(),
		}
		k.AccountKeeper.SetAccount(ctx, ethAccount)
		k.Logger(ctx).Info("migrate account", "address", account.GetAddress())
		return false
	})
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
