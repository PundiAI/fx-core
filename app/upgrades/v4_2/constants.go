package v4_2

import (
	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v4/app/upgrades"
	v4 "github.com/functionx/fx-core/v4/app/upgrades/v4"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

type Refund struct {
	Address common.Address
	Coins   sdk.Coins
}

var (
	PolygonUSDTDenom = "polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F"
	// SendToExternalRefund send to external refund
	// fxcore tx: https://starscan.io/evm/tx/0x6795d921d532b68b886e8f776143a138688897a8b143f981545475df16f01728
	// polygon fork tx: https://polygonscan.com/tx/0xe7ea70c77960e78703942d895afaac5d6fe0f99cf2b894cbb5b2bcf36e5a82ce
	// amount: 89724233, bridge-fee: 11988
	SendToExternalRefund = Refund{
		Address: common.HexToAddress("0x3Cf9771C961af8727EF5E9b955419cf80b34bcd4"),
		Coins:   sdk.NewCoins(sdk.NewCoin(PolygonUSDTDenom, sdkmath.NewInt(89724233).Add(sdkmath.NewInt(11988)))),
	}
	// SendToFxDelayedPayment send to fx delayed payment
	// polygon tx: https://polygonscan.com/tx/0xebbcc9c9e91a081b7f3717b0d5d0a571e965ce3e0c65c11a9925e1142d829d1a
	// amount: 99892059
	SendToFxDelayedPayment = Refund{
		Address: common.HexToAddress("0x17a3d2EEE4E3558f40c0c7A583A182f44d377759"),
		Coins:   sdk.NewCoins(sdk.NewCoin(PolygonUSDTDenom, sdk.NewInt(99892059))),
	}
	PolygonUSDTRefunds = []Refund{SendToExternalRefund, SendToFxDelayedPayment}
)

func Upgrade() upgrades.Upgrade {
	upgrade := upgrades.Upgrade{
		UpgradeName:          "v4.2.x",
		CreateUpgradeHandler: createUpgradeHandler,
		StoreUpgrades:        v4.Upgrade.StoreUpgrades,
	}

	// if testnet, store has been upgraded in fxv4
	if fxtypes.ChainId() == fxtypes.TestnetChainId {
		upgrade.StoreUpgrades = func() *storetypes.StoreUpgrades {
			return &storetypes.StoreUpgrades{}
		}
	}

	return upgrade
}
