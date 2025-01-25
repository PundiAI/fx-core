package v8

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
)

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

func updateWPUNDIAILogicCode(ctx sdk.Context, keeper *fxevmkeeper.Keeper) {
	wpundiai := contract.GetWPUNDIAI()
	if err := keeper.UpdateContractCode(ctx, wpundiai.Address, wpundiai.Code); err != nil {
		ctx.Logger().Error("update WPUNDIAI contract", "module", "upgrade", "err", err.Error())
	} else {
		ctx.Logger().Info("update WPUNDIAI contract", "module", "upgrade", "codeHash", wpundiai.CodeHash())
	}
}

func updateERC20LogicCode(ctx sdk.Context, keeper *fxevmkeeper.Keeper) {
	erc20 := contract.GetERC20()
	if err := keeper.UpdateContractCode(ctx, erc20.Address, erc20.Code); err != nil {
		ctx.Logger().Error("update ERC20 contract", "module", "upgrade", "err", err.Error())
	} else {
		ctx.Logger().Info("update ERC20 contract", "module", "upgrade", "codeHash", erc20.CodeHash())
	}
}

func deployAccessControlContract(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, evmModuleAddress common.Address) error {
	return contract.DeployAccessControlContract(ctx, evmKeeper, evmModuleAddress, getContractOwner(ctx))
}
