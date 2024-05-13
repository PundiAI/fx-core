package keepers

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/ethereum/go-ethereum/core/vm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	arbitrumtypes "github.com/functionx/fx-core/v7/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v7/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	precompilescrosschain "github.com/functionx/fx-core/v7/x/evm/precompiles/crosschain"
	precompilesstaking "github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
	layer2types "github.com/functionx/fx-core/v7/x/layer2/types"
	migratetypes "github.com/functionx/fx-core/v7/x/migrate/types"
	optimismtypes "github.com/functionx/fx-core/v7/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v7/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

func (appKeepers *AppKeepers) generateKeys() {
	// Define what keys will be used in the cosmos-sdk key/value store.
	// Cosmos-SDK modules each have a "key" that allows the application to reference what they've stored on the chain.
	appKeepers.keys = sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey,
		bsctypes.StoreKey, polygontypes.StoreKey, avalanchetypes.StoreKey, ethtypes.StoreKey, trontypes.StoreKey,
		arbitrumtypes.ModuleName, optimismtypes.ModuleName, layer2types.ModuleName,
		evmtypes.StoreKey, feemarkettypes.StoreKey,
		erc20types.StoreKey, migratetypes.StoreKey,
	)

	// Define transient store keys
	appKeepers.tkeys = sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)

	// MemKeys are for information that is stored only in RAM.
	appKeepers.memKeys = sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
}

func (appKeepers *AppKeepers) GetKVStoreKey() map[string]*storetypes.KVStoreKey {
	return appKeepers.keys
}

func (appKeepers *AppKeepers) GetTransientStoreKey() map[string]*storetypes.TransientStoreKey {
	return appKeepers.tkeys
}

func (appKeepers *AppKeepers) GetMemoryStoreKey() map[string]*storetypes.MemoryStoreKey {
	return appKeepers.memKeys
}

// EvmPrecompiled  set evm precompiled contracts
func (appKeepers *AppKeepers) EvmPrecompiled() {
	precompiled := evmkeeper.BerlinPrecompiled()

	// staking precompile
	precompiled[precompilesstaking.GetAddress()] = func(ctx sdk.Context) vm.PrecompiledContract {
		return precompilesstaking.NewPrecompiledContract(
			ctx,
			appKeepers.BankKeeper,
			appKeepers.StakingKeeper,
			appKeepers.DistrKeeper,
			appKeepers.EvmKeeper,
		)
	}

	// cross chain precompile
	crosschainRouter := precompilescrosschain.NewRouter().
		AddRoute(ethtypes.ModuleName, appKeepers.EthKeeper).
		AddRoute(bsctypes.ModuleName, appKeepers.BscKeeper).
		AddRoute(polygontypes.ModuleName, appKeepers.PolygonKeeper).
		AddRoute(trontypes.ModuleName, appKeepers.TronKeeper).
		AddRoute(avalanchetypes.ModuleName, appKeepers.AvalancheKeeper).
		AddRoute(arbitrumtypes.ModuleName, appKeepers.ArbitrumKeeper).
		AddRoute(optimismtypes.ModuleName, appKeepers.OptimismKeeper).
		AddRoute(layer2types.ModuleName, appKeepers.Layer2Keeper)
	precompiled[precompilescrosschain.GetAddress()] = func(ctx sdk.Context) vm.PrecompiledContract {
		return precompilescrosschain.NewPrecompiledContract(
			ctx,
			appKeepers.BankKeeper,
			appKeepers.Erc20Keeper,
			appKeepers.IBCTransferKeeper,
			appKeepers.AccountKeeper,
			crosschainRouter,
		)
	}

	// set precompiled contracts
	appKeepers.EvmKeeper.WithPrecompiled(precompiled)
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetKey(storeKey string) *storetypes.KVStoreKey {
	return appKeepers.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return appKeepers.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (appKeepers *AppKeepers) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return appKeepers.memKeys[storeKey]
}
