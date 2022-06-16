package v2

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/ethereum/go-ethereum/crypto"

	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	erc20types "github.com/functionx/fx-core/x/erc20/types"
	migratetypes "github.com/functionx/fx-core/x/migrate/types"
)

const (
	UpgradeName = "v2"

	updateBlockParamsKey   = "BlockParams"
	updateBlockParamsValue = `{"max_bytes":"1048576","max_gas":"3000000000"}`
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
		ibchost.ModuleName:           1,
		bsctypes.ModuleName:          1,
		polygontypes.ModuleName:      1,
		trontypes.ModuleName:         1,
	}

	storeUpgrades = &store.StoreUpgrades{
		Added: []string{
			feemarkettypes.StoreKey,
			evmtypes.StoreKey,
			erc20types.StoreKey,
			migratetypes.StoreKey,
		},
	}

	emptyCodeHash = crypto.Keccak256(nil)
)

func GetStoreUpgrades() *store.StoreUpgrades {
	return storeUpgrades
}
