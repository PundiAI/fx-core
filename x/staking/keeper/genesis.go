package keeper

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/v7/x/staking/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) (res []abci.ValidatorUpdate) {
	// fx staking init genesis
	for _, allowance := range data.Allowances {
		valAddr, err := sdk.ValAddressFromBech32(allowance.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		ownerAddr := sdk.MustAccAddressFromBech32(allowance.OwnerAddress)
		spenderAddr := sdk.MustAccAddressFromBech32(allowance.SpenderAddress)
		if allowance.Allowance.IsNegative() {
			panic("allowance must be positive")
		}
		k.SetAllowance(ctx, valAddr, ownerAddr, spenderAddr, allowance.Allowance.BigInt())
	}

	// staking init genesis
	stakingGenesisState := stakingtypes.GenesisState{
		Params:               data.Params,
		LastTotalPower:       data.LastTotalPower,
		LastValidatorPowers:  data.LastValidatorPowers,
		Validators:           data.Validators,
		Delegations:          data.Delegations,
		UnbondingDelegations: data.UnbondingDelegations,
		Redelegations:        data.Redelegations,
		Exported:             data.Exported,
	}
	return k.Keeper.InitGenesis(ctx, &stakingGenesisState)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// fx staking export genesis
	var allownaces []types.Allowance
	k.IterateAllAllowance(ctx, func(valAddr sdk.ValAddress, ownerAddr, spenderAddr sdk.AccAddress, allowance *big.Int) bool {
		allownaces = append(allownaces, types.Allowance{
			ValidatorAddress: valAddr.String(),
			OwnerAddress:     ownerAddr.String(),
			SpenderAddress:   spenderAddr.String(),
			Allowance:        sdkmath.NewIntFromBigInt(allowance),
		})
		return false
	})
	// staking export genesis
	data := k.Keeper.ExportGenesis(ctx)
	return &types.GenesisState{
		Params:               data.Params,
		LastTotalPower:       data.LastTotalPower,
		LastValidatorPowers:  data.LastValidatorPowers,
		Validators:           data.Validators,
		Delegations:          data.Delegations,
		UnbondingDelegations: data.UnbondingDelegations,
		Redelegations:        data.Redelegations,
		Exported:             data.Exported,
		Allowances:           allownaces,
	}
}
