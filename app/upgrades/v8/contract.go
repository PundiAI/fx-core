package v8

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/app/keepers"
	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
)

func updateContract(ctx sdk.Context, app *keepers.AppKeepers) error {
	acc := app.AccountKeeper.GetModuleAddress(evmtypes.ModuleName)
	moduleAddress := common.BytesToAddress(acc.Bytes())

	if err := deployBridgeFeeContract(
		ctx,
		app.EvmKeeper,
		app.Erc20Keeper,
		app.CrosschainKeepers.EthKeeper,
		moduleAddress,
	); err != nil {
		return err
	}

	if err := deployAccessControlContract(ctx, app.EvmKeeper, moduleAddress); err != nil {
		return err
	}

	if err := updateWPUNDIAILogicCode(ctx, app.EvmKeeper); err != nil {
		return err
	}
	return updateERC20LogicCode(ctx, app.EvmKeeper)
}

func deployBridgeFeeContract(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, crosschainKeeper crosschainkeeper.Keeper, evmModuleAddress common.Address) error {
	chains := fxtypes.GetSupportChains()
	bridgeDenoms := make([]contract.BridgeDenoms, len(chains))
	for index, chain := range chains {
		denoms := make([]common.Hash, 0)
		bridgeTokens, err := erc20Keeper.GetBridgeTokens(ctx, chain)
		if err != nil {
			return err
		}
		for _, token := range bridgeTokens {
			denoms = append(denoms, contract.MustStrToByte32(token.GetDenom()))
		}
		bridgeDenoms[index] = contract.BridgeDenoms{
			ChainName: contract.MustStrToByte32(chain),
			Denoms:    denoms,
		}
	}

	oracles := crosschainKeeper.GetAllOracles(ctx, true)
	if oracles.Len() <= 0 {
		return errors.New("no oracle found")
	}
	defaultOracleAddress := common.HexToAddress(oracles[0].ExternalAddress)
	return contract.DeployBridgeFeeContract(ctx, evmKeeper, bridgeDenoms, evmModuleAddress, getContractOwner(ctx), defaultOracleAddress)
}

func updateWPUNDIAILogicCode(ctx sdk.Context, keeper *fxevmkeeper.Keeper) error {
	wpundiai := contract.GetWPUNDIAI()
	return keeper.UpdateContractCode(ctx, wpundiai.Address, wpundiai.Code)
}

func updateERC20LogicCode(ctx sdk.Context, keeper *fxevmkeeper.Keeper) error {
	erc20 := contract.GetERC20()
	return keeper.UpdateContractCode(ctx, erc20.Address, erc20.Code)
}

func deployAccessControlContract(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, evmModuleAddress common.Address) error {
	return contract.DeployAccessControlContract(ctx, evmKeeper, evmModuleAddress, getContractOwner(ctx))
}
