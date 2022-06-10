package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	fxtypes "github.com/functionx/fx-core/types"
	v042 "github.com/functionx/fx-core/x/crosschain/legacy/v042"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func MigrateOracle(ctx sdk.Context, cdc codec.BinaryCodec, storeKey sdk.StoreKey, stakingKeeper StakingKeeper) (types.Oracles, error) {

	validatorsByPower := stakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validatorsByPower) <= 0 {
		panic("no found bonded validator")
	}
	delegateValidator := validatorsByPower[0].OperatorAddress

	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	oracles := types.Oracles{}
	for ; iterator.Valid(); iterator.Next() {
		var legacyOracle v042.LegacyOracle
		cdc.MustUnmarshal(iterator.Value(), &legacyOracle)
		if legacyOracle.DelegateAmount.Denom != fxtypes.DefaultDenom {
			return nil, sdkerrors.Wrapf(types.ErrInvalid, "delegate denom: %s", legacyOracle.DelegateAmount.Denom)
		}
		oracleAddr, err := sdk.AccAddressFromBech32(legacyOracle.OracleAddress)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
		}

		oracle := types.Oracle{
			OracleAddress:     legacyOracle.OracleAddress,
			BridgerAddress:    legacyOracle.BridgerAddress,
			ExternalAddress:   legacyOracle.ExternalAddress,
			DelegateAmount:    legacyOracle.DelegateAmount.Amount,
			StartHeight:       legacyOracle.StartHeight,
			Jailed:            legacyOracle.Jailed,
			JailedHeight:      legacyOracle.JailedHeight,
			DelegateValidator: delegateValidator,
			IsValidator:       false,
		}
		validator, found := stakingKeeper.GetValidator(ctx, oracleAddr.Bytes())
		if found {
			oracle.IsValidator = true
			oracle.DelegateValidator = validator.OperatorAddress
		}
		store.Set(types.GetOracleKey(oracle.GetOracle()), cdc.MustMarshal(&oracle))
		oracles = append(oracles, oracle)
	}
	return oracles, nil
}
