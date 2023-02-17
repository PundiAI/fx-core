package v4

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/x/staking"
	"github.com/functionx/fx-core/v3/x/staking/keeper"
	types "github.com/functionx/fx-core/v3/x/staking/types"
)

func Migrate(ctx sdk.Context, cdc codec.BinaryCodec, k keeper.Keeper, ak types.AccountKeeper) (err error) {
	// 1. create lpToken module address if not exist
	staking.CreateLPTokenModuleAccount(ctx, types.LPTokenOwnerModuleName, ak)

	// 2. create lpToken contract for all validators
	validatorLPTokenMap := map[string]common.Address{}
	iterator := k.IteratorValidators(ctx)

	for ; iterator.Valid(); iterator.Next() {
		validator := stakingtypes.MustUnmarshalValidator(cdc, iterator.Value())
		lpTokenContract, err := k.DeployLPToken(ctx, validator.GetOperator())
		if err != nil {
			return err
		}
		validatorLPTokenMap[validator.GetOperator().String()] = lpTokenContract
	}

	// 3. mint tokenToken for all delegations
	k.IterateAllDelegations(ctx, func(delegation stakingtypes.Delegation) (stop bool) {
		lpTokenContract, found := validatorLPTokenMap[delegation.GetValidatorAddr().String()]
		if !found {
			return false
		}
		err = k.MintLPToken(ctx, lpTokenContract, delegation.GetDelegatorAddr(), delegation.GetShares())
		return err != nil
	})
	return err
}
