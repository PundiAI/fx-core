package keeper_test

import (
	"math/big"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/staking/types"
)

func (suite *KeeperTestSuite) TestInitGenesis() {
	allowances := []types.Allowance{
		{
			ValidatorAddress: sdk.ValAddress(helpers.GenHexAddress().Bytes()).String(),
			OwnerAddress:     helpers.GenAccAddress().String(),
			SpenderAddress:   helpers.GenAccAddress().String(),
			Allowance:        sdkmath.NewIntFromUint64(tmrand.Uint64()),
		},
		{
			ValidatorAddress: sdk.ValAddress(helpers.GenHexAddress().Bytes()).String(),
			OwnerAddress:     helpers.GenAccAddress().String(),
			SpenderAddress:   helpers.GenAccAddress().String(),
			Allowance:        sdkmath.NewIntFromUint64(tmrand.Uint64()),
		},
		{
			ValidatorAddress: sdk.ValAddress(helpers.GenHexAddress().Bytes()).String(),
			OwnerAddress:     helpers.GenAccAddress().String(),
			SpenderAddress:   helpers.GenAccAddress().String(),
			Allowance:        sdkmath.NewIntFromUint64(tmrand.Uint64()),
		},
	}
	sort.Sort(types.Allowances(allowances))

	stakingGenesisState := suite.app.StakingKeeper.Keeper.ExportGenesis(suite.ctx)
	fxStakingGenesisState := types.GenesisState{
		Params:               stakingGenesisState.Params,
		LastTotalPower:       stakingGenesisState.LastTotalPower,
		LastValidatorPowers:  stakingGenesisState.LastValidatorPowers,
		Validators:           stakingGenesisState.Validators,
		Delegations:          stakingGenesisState.Delegations,
		UnbondingDelegations: stakingGenesisState.UnbondingDelegations,
		Redelegations:        stakingGenesisState.Redelegations,
		Exported:             stakingGenesisState.Exported,
		Allowances:           allowances,
	}

	res := suite.app.StakingKeeper.InitGenesis(suite.ctx, &fxStakingGenesisState)
	suite.Require().Equal(len(stakingGenesisState.LastValidatorPowers), len(res))

	allAllowance := make([]types.Allowance, 0)
	suite.app.StakingKeeper.IterateAllAllowance(suite.ctx, func(valAddr sdk.ValAddress, owner, spender sdk.AccAddress, allowance *big.Int) (stop bool) {
		allAllowance = append(allAllowance, types.Allowance{
			ValidatorAddress: valAddr.String(),
			OwnerAddress:     owner.String(),
			SpenderAddress:   spender.String(),
			Allowance:        sdkmath.NewIntFromBigInt(allowance),
		})
		return false
	})
	sort.Sort(types.Allowances(allAllowance))
	suite.Require().EqualValues(allowances, allAllowance)
}

func (suite *KeeperTestSuite) TestExpGenesis() {
	allowances := []types.Allowance{
		{
			ValidatorAddress: sdk.ValAddress(helpers.GenHexAddress().Bytes()).String(),
			OwnerAddress:     helpers.GenAccAddress().String(),
			SpenderAddress:   helpers.GenAccAddress().String(),
			Allowance:        sdkmath.NewIntFromUint64(tmrand.Uint64()),
		},
		{
			ValidatorAddress: sdk.ValAddress(helpers.GenHexAddress().Bytes()).String(),
			OwnerAddress:     helpers.GenAccAddress().String(),
			SpenderAddress:   helpers.GenAccAddress().String(),
			Allowance:        sdkmath.NewIntFromUint64(tmrand.Uint64()),
		},
		{
			ValidatorAddress: sdk.ValAddress(helpers.GenHexAddress().Bytes()).String(),
			OwnerAddress:     helpers.GenAccAddress().String(),
			SpenderAddress:   helpers.GenAccAddress().String(),
			Allowance:        sdkmath.NewIntFromUint64(tmrand.Uint64()),
		},
	}
	sort.Sort(types.Allowances(allowances))
	stakingGenesisState := suite.app.StakingKeeper.Keeper.ExportGenesis(suite.ctx)
	fxStakingGenesisState := types.GenesisState{
		Params:               stakingGenesisState.Params,
		LastTotalPower:       stakingGenesisState.LastTotalPower,
		LastValidatorPowers:  stakingGenesisState.LastValidatorPowers,
		Validators:           stakingGenesisState.Validators,
		Delegations:          stakingGenesisState.Delegations,
		UnbondingDelegations: stakingGenesisState.UnbondingDelegations,
		Redelegations:        stakingGenesisState.Redelegations,
		Exported:             stakingGenesisState.Exported,
		Allowances:           allowances,
	}

	res := suite.app.StakingKeeper.InitGenesis(suite.ctx, &fxStakingGenesisState)
	suite.Require().Equal(len(stakingGenesisState.LastValidatorPowers), len(res))

	genesisState := suite.app.StakingKeeper.ExportGenesis(suite.ctx)
	allAllowance := genesisState.Allowances
	sort.Sort(types.Allowances(allAllowance))
	genesisState.Allowances = allAllowance
	suite.Require().EqualValues(&fxStakingGenesisState, genesisState)
}
