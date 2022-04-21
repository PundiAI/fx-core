package forks

import (
	"errors"
	"github.com/ethereum/go-ethereum/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/types"
	erc20keeper "github.com/functionx/fx-core/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/x/erc20/types"
	evmkeeper "github.com/functionx/fx-core/x/evm/keeper"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	feemarketkeeper "github.com/functionx/fx-core/x/feemarket/keeper"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
)

func InitSupportEvm(ctx sdk.Context, accountKeeper authkeeper.AccountKeeper,
	feeMarketKeeper feemarketkeeper.Keeper, feemarketParams feemarkettypes.Params,
	evmKeeper *evmkeeper.Keeper, evmParams evmtypes.Params,
	erc20Keeper erc20keeper.Keeper, erc20Params erc20types.Params,
) error {
	logger := ctx.Logger()
	// init fee market
	logger.Info("init fee market", "params", feemarketParams.String())
	// set feeMarket blockGasUsed
	feeMarketKeeper.SetBlockGasUsed(ctx, 0)
	// init feeMarket module erc20Params
	feeMarketKeeper.SetParams(ctx, feemarketParams)

	// init evm
	logger.Info("init evm", "params", evmParams.String())
	evmKeeper.SetParams(ctx, evmParams)

	// init erc20
	logger.Info("init erc20", "params", erc20Params.String())
	erc20Keeper.SetParams(ctx, erc20Params)

	// init contract
	// ensure erc20 module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, erc20types.ModuleName); acc == nil {
		return errors.New("the erc20 module account has not been set")
	}
	for _, contract := range fxtypes.GetInitContracts() {
		if len(contract.Code) <= 0 || contract.Address == common.HexToAddress(fxtypes.EmptyEvmAddress) {
			return errors.New("invalid contract")
		}
		if err := evmKeeper.CreateContractWithCode(ctx, contract.Address, contract.Code); err != nil {
			return err
		}
	}

	// init coin
	for _, metadata := range fxtypes.GetMetadata() {
		logger.Info("add metadata", "coin", metadata.String())
		pair, err := erc20Keeper.RegisterCoin(ctx, metadata)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			erc20types.EventTypeRegisterCoin,
			sdk.NewAttribute(erc20types.AttributeKeyDenom, pair.Denom),
			sdk.NewAttribute(erc20types.AttributeKeyTokenAddress, pair.Erc20Address),
		))
	}
	return nil
}

func DefaultFeeMarket() feemarkettypes.Params {
	var (
		initialBaseFee           = 500000000000                            //500,000,000,000
		baseFeeChangeDenominator = uint32(params.BaseFeeChangeDenominator) //8
		elasticityMultiplier     = uint32(params.ElasticityMultiplier)     //2
		baseFee                  = uint64(initialBaseFee)                  //500Gwei
		minBaseFee               = sdk.NewInt(int64(initialBaseFee))       //500Gwei
		maxBaseFee               = feemarkettypes.MaxBaseFee               //MaxUint64 - 1
		maxGas                   = uint64(3e7)                             //30,000,000
	)
	return feemarkettypes.NewParams(baseFeeChangeDenominator, elasticityMultiplier, baseFee, minBaseFee, maxBaseFee, maxGas)
}
