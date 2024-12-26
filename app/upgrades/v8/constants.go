package v8

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
