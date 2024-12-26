package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
)

func DeployBridgeFeeContract(
	ctx sdk.Context,
	evmKeeper *fxevmkeeper.Keeper,
	bridgeDenoms []contract.BridgeDenoms,
	evmModuleAddress, ownerAddress, defaultOracleAddress common.Address,
) error {
	bridgeFeeQuoteKeeper := contract.NewBridgeFeeQuoteKeeper(evmKeeper, contract.BridgeFeeAddress)
	bridgeFeeOracleKeeper := contract.NewBridgeFeeOracleKeeper(evmKeeper, contract.BridgeFeeOracleAddress)

	return contract.DeployBridgeFeeContract(
		ctx,
		evmKeeper,
		bridgeFeeQuoteKeeper,
		bridgeFeeOracleKeeper,
		bridgeDenoms,
		evmModuleAddress,
		ownerAddress,
		defaultOracleAddress,
	)
}

func DeployAccessControlContract(
	ctx sdk.Context,
	evmKeeper *fxevmkeeper.Keeper,
	evmModuleAddress, adminAddress common.Address,
) error {
	accessControl := contract.NewAccessControlKeeper(evmKeeper, contract.AccessControlAddress)
	return contract.DeployAccessControlContract(ctx, evmKeeper, accessControl, evmModuleAddress, adminAddress)
}
