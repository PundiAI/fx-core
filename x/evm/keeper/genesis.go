package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/contract"
)

// InitGenesis initializes genesis state based on exported genesis
func (k *Keeper) InitGenesis(ctx sdk.Context) {
	// ensure evm module account is set
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		panic("the EVM module account has not been set")
	}
	// init logic contract
	initContract := []contract.Contract{contract.GetFIP20(), contract.GetWFX()}
	for _, contractAddr := range initContract {
		if len(contractAddr.Code) == 0 || contract.IsZeroEthAddress(contractAddr.Address) {
			panic(fmt.Sprintf("invalid contract: %s", contractAddr.Address.String()))
		}
		if err := k.CreateContractWithCode(ctx, contractAddr.Address, contractAddr.Code); err != nil {
			panic(fmt.Sprintf("create contract %s with code error %s", contractAddr.Address.String(), err.Error()))
		}
	}
}
