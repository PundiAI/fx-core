package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	v010 "github.com/functionx/fx-core/v2/x/crosschain/legacy/v1"

	fxtypes "github.com/functionx/fx-core/v2/types"
	"github.com/functionx/fx-core/v2/x/crosschain/types"
)

func MigrateOracle(ctx sdk.Context, cdc codec.BinaryCodec, storeKey sdk.StoreKey, stakingKeeper StakingKeeper) (types.Oracles, sdk.ValAddress, error) {

	validatorsByPower := stakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validatorsByPower) <= 0 {
		panic("no found bonded validator")
	}
	delegateValidator := validatorsByPower[0]

	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	oracles := types.Oracles{}
	for ; iterator.Valid(); iterator.Next() {
		var legacyOracle v010.LegacyOracle
		cdc.MustUnmarshal(iterator.Value(), &legacyOracle)
		if legacyOracle.DepositAmount.Denom != fxtypes.DefaultDenom {
			return nil, nil, sdkerrors.Wrapf(types.ErrInvalid, "delegate denom: %s", legacyOracle.DepositAmount.Denom)
		}

		oracle := types.Oracle{
			OracleAddress:     legacyOracle.OracleAddress,
			BridgerAddress:    legacyOracle.OrchestratorAddress,
			ExternalAddress:   legacyOracle.ExternalAddress,
			DelegateAmount:    legacyOracle.DepositAmount.Amount,
			StartHeight:       legacyOracle.StartHeight,
			Online:            !legacyOracle.Jailed,
			DelegateValidator: delegateValidator.OperatorAddress,
			SlashTimes:        0,
		}
		store.Set(types.GetOracleKey(oracle.GetOracle()), cdc.MustMarshal(&oracle))
		oracles = append(oracles, oracle)
	}

	return oracles, delegateValidator.GetOperator(), nil
}
