package app

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/app/keepers"
	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func initChainer(ctx sdk.Context, keepers keepers.AppKeepers) error {
	delegations, err := keepers.StakingKeeper.GetAllDelegations(ctx)
	if err != nil {
		return err
	}
	if len(delegations) == 0 {
		return errors.New("no delegations found")
	}

	bridgeDenoms := []contract.BridgeDenoms{
		{
			ChainName: ethtypes.ModuleName,
			Denoms:    []string{fxtypes.DefaultDenom},
		},
	}

	acc := keepers.AccountKeeper.GetModuleAddress(evmtypes.ModuleName)
	moduleAddress := common.BytesToAddress(acc.Bytes())

	if err = contract.DeployBridgeFeeContract(
		ctx,
		keepers.EvmKeeper,
		bridgeDenoms,
		moduleAddress,
		moduleAddress,
		common.BytesToAddress(sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress).Bytes()),
	); err != nil {
		return err
	}

	err = contract.DeployAccessControlContract(ctx, keepers.EvmKeeper, moduleAddress, moduleAddress)
	return err
}
