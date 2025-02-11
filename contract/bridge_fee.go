package contract

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const DefaultMaxQuoteCap = 3

func DeployBridgeFeeContract(
	ctx sdk.Context,
	evmKeeper EvmKeeper,
	bridgeDenoms []BridgeDenoms,
	evmModuleAddress, owner, defaultOracleAddress common.Address,
) error {
	if err := DeployProxy(
		ctx,
		evmKeeper,
		GetBridgeFeeQuote().ABI,
		GetBridgeFeeQuote().Bin,
		common.HexToAddress(BridgeFeeAddress),
		evmModuleAddress,
	); err != nil {
		return err
	}
	if err := DeployProxy(
		ctx,
		evmKeeper,
		GetBridgeFeeOracle().ABI,
		GetBridgeFeeOracle().Bin,
		common.HexToAddress(BridgeFeeOracleAddress),
		evmModuleAddress,
	); err != nil {
		return err
	}

	bridgeFeeOracleKeeper := NewBridgeFeeOracleKeeper(evmKeeper)
	if err := initBridgeFeeOracle(ctx, bridgeFeeOracleKeeper, owner, defaultOracleAddress); err != nil {
		return err
	}
	bridgeFeeQuoteKeeper := NewBridgeFeeQuoteKeeper(evmKeeper)
	return initBridgeFeeQuote(ctx, bridgeFeeQuoteKeeper, bridgeDenoms, owner)
}

func DeployProxy(
	ctx sdk.Context,
	evmKeeper EvmKeeper,
	logicABI abi.ABI,
	logicBin []byte,
	proxyAddress, evmModuleAddress common.Address,
) error {
	logicContract, err := evmKeeper.DeployContract(ctx, evmModuleAddress, logicABI, logicBin)
	if err != nil {
		return err
	}
	if err = evmKeeper.CreateContractWithCode(ctx, proxyAddress, GetBridgeProxy().Code); err != nil {
		return err
	}
	if _, err = evmKeeper.ApplyContract(ctx, evmModuleAddress, proxyAddress, nil, GetBridgeProxy().ABI, "init", logicContract); err != nil {
		return err
	}
	return nil
}

func initBridgeFeeOracle(
	ctx sdk.Context,
	bridgeFeeOracleKeeper BridgeFeeOracleKeeper,
	owner, defaultOracleAddress common.Address,
) error {
	if _, err := bridgeFeeOracleKeeper.Initialize(ctx); err != nil {
		return err
	}
	role, err := bridgeFeeOracleKeeper.GetQuoteRole(ctx)
	if err != nil {
		return err
	}
	if _, err = bridgeFeeOracleKeeper.GrantRole(ctx, role, common.HexToAddress(BridgeFeeAddress)); err != nil {
		return err
	}
	ownerRole, err := bridgeFeeOracleKeeper.GetOwnerRole(ctx)
	if err != nil {
		return err
	}
	if _, err = bridgeFeeOracleKeeper.GrantRole(ctx, ownerRole, owner); err != nil {
		return err
	}
	upgradeRole, err := bridgeFeeOracleKeeper.GetUpgradeRole(ctx)
	if err != nil {
		return err
	}
	if _, err = bridgeFeeOracleKeeper.GrantRole(ctx, upgradeRole, owner); err != nil {
		return err
	}
	if _, err = bridgeFeeOracleKeeper.SetDefaultOracle(ctx, defaultOracleAddress); err != nil {
		return err
	}
	return nil
}

func initBridgeFeeQuote(
	ctx sdk.Context,
	bridgeFeeQuoteKeeper BridgeFeeQuoteKeeper,
	bridgeDenoms []BridgeDenoms,
	owner common.Address,
) error {
	if _, err := bridgeFeeQuoteKeeper.Initialize(ctx, common.HexToAddress(BridgeFeeOracleAddress), DefaultMaxQuoteCap); err != nil {
		return err
	}
	ownerRole, err := bridgeFeeQuoteKeeper.GetOwnerRole(ctx)
	if err != nil {
		return err
	}
	if _, err = bridgeFeeQuoteKeeper.GrantRole(ctx, ownerRole, owner); err != nil {
		return err
	}
	upgradeRole, err := bridgeFeeQuoteKeeper.GetUpgradeRole(ctx)
	if err != nil {
		return err
	}
	if _, err = bridgeFeeQuoteKeeper.GrantRole(ctx, upgradeRole, owner); err != nil {
		return err
	}
	for _, bridgeDenom := range bridgeDenoms {
		if _, err = bridgeFeeQuoteKeeper.RegisterChain(ctx, bridgeDenom.ChainName, bridgeDenom.Denoms...); err != nil {
			return err
		}
	}
	return nil
}

func UpgradeBridgeFeeContract(ctx sdk.Context, evmKeeper EvmKeeper, evmModuleAddress common.Address) error {
	logicContract, err := evmKeeper.DeployContract(ctx, evmModuleAddress, GetBridgeFeeQuote().ABI, GetBridgeFeeQuote().Bin)
	if err != nil {
		return err
	}

	bridgeFeeQuoteKeeper := NewBridgeFeeQuoteKeeper(evmKeeper)
	_, err = bridgeFeeQuoteKeeper.UpgradeTo(ctx, logicContract)
	return err
}
