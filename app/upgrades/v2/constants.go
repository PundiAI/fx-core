package v2

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	erc20types "github.com/functionx/fx-core/x/erc20/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
	migratetypes "github.com/functionx/fx-core/x/migrate/types"
)

const (
	UpgradeName = "v2"
)

var (
	initGenesis = map[string]bool{
		feemarkettypes.ModuleName: true,
		evmtypes.ModuleName:       true,
		erc20types.ModuleName:     true,
		migratetypes.ModuleName:   true,
	}

	runMigrates = map[string]uint64{
		authtypes.ModuleName:         1,
		banktypes.ModuleName:         1,
		distributiontypes.ModuleName: 1,
		govtypes.ModuleName:          1,
		slashingtypes.ModuleName:     1,
		stakingtypes.ModuleName:      1,
	}

	storeUpgrades = &store.StoreUpgrades{
		Added: []string{
			feemarkettypes.StoreKey,
			evmtypes.StoreKey,
			erc20types.StoreKey,
			migratetypes.StoreKey,
		},
	}
)

func GetStoreUpgrades() *store.StoreUpgrades {
	return storeUpgrades
}
