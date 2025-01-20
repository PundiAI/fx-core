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
	validators, err := keepers.StakingKeeper.GetValidators(ctx, 1)
	if err != nil {
		return err
	}
	if len(validators) == 0 {
		return errors.New("no validators found")
	}
	defValAddress, err := sdk.ValAddressFromBech32(validators[0].OperatorAddress)
	if err != nil {
		return err
	}

	bridgeDenoms := []contract.BridgeDenoms{
		{
			ChainName: contract.MustStrToByte32(ethtypes.ModuleName),
			Denoms:    []common.Hash{contract.MustStrToByte32(fxtypes.DefaultDenom)},
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
		common.BytesToAddress(defValAddress.Bytes()),
	); err != nil {
		return err
	}

	err = contract.DeployAccessControlContract(ctx, keepers.EvmKeeper, moduleAddress, moduleAddress)
	return err
}
