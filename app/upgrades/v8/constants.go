package v8

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/app/upgrades"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v8.0.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{
			Deleted: []string{
				"fxtransfer",
			},
		}
	},
}

const (
	// NOTE: pundix escrow purse amount is the number of purses locked by the ibc channel in pundix
	// Since the IBC cross-chain is always running, this number will keep changing.
	// Therefore, we set a relatively large value here to ensure the normal operation of the BSC module cross-chain,
	// and we will fix this number in the next version.

	testnetPundixEscrowPurseAmount = "64000000000000000000000000000"
	mainnetPundixEscrowPurseAmount = "64000000000000000000000000000"
)

func getPundixEscrowPurseAmount(ctx sdk.Context) (sdkmath.Int, error) {
	var purseAmount sdkmath.Int
	var ok bool
	if ctx.ChainID() == fxtypes.TestnetChainId {
		purseAmount, ok = sdkmath.NewIntFromString(testnetPundixEscrowPurseAmount)
	} else {
		purseAmount, ok = sdkmath.NewIntFromString(mainnetPundixEscrowPurseAmount)
	}
	if !ok {
		return sdkmath.ZeroInt(), fmt.Errorf("pundix escrow purse amount is invalid")
	}
	return purseAmount, nil
}

const (
	testnetOwnerAddress = "0x2DC5f63149F922c8F73A36A7f295Ed2Af269D7d8"
	mainnetOwnerAddress = "0xE77A7EA2F1DC25968b5941a456d99D37b80E98B5"
)

func getContractOwner(ctx sdk.Context) common.Address {
	if ctx.ChainID() == fxtypes.TestnetChainId {
		return common.HexToAddress(testnetOwnerAddress)
	}
	return common.HexToAddress(mainnetOwnerAddress)
}

const (
	pundixSymbol = "PUNDIX"
	purseSymbol  = "PURSE"
)
