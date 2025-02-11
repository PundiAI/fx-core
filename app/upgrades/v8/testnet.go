package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/app/keepers"
)

func upgradeTestnet(ctx sdk.Context, app *keepers.AppKeepers) {
	initBridgeAccount(ctx, app.AccountKeeper)
}
