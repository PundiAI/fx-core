package v5

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v5/app/upgrades"
)

type SlashPeriod struct {
	Delegator sdk.AccAddress
	Height    uint64
	Period    uint64
}

// ValidatorSlashHeightTestnetFXV4 is a map of testnet validator address to slash height
var ValidatorSlashHeightTestnetFXV4 = map[string][]int64{
	"fxvaloper14lpap6mwytqtnrx6q9cnje2sen5a5wcctuwnsh": {8224664, 8457967, 8552303},
	"fxvaloper1695tvr34jdrmajtc0t6kzumaxxwhulfl7h6mq2": {9469704},
	"fxvaloper16d0jly49xgwm9tyf7lpf0splnfhrnttdejkz9h": {8787488, 8885841, 9416069},
	"fxvaloper1c863577qlc4sq64v5dsfunf5drykz7fnqkgdha": {9128756},
	"fxvaloper1xdqas5ak98us9eljqj5ppj5mhmku4slh2664l8": {8806427},
}

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v5.0.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{}
	},
}
