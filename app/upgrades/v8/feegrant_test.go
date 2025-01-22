package v8_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/x/feegrant"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/app/upgrades/v8"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func Test_migrateFeegrant(t *testing.T) {
	app, ctx := helpers.NewAppWithValNumber(t, 1)

	now := time.Now()
	coins := sdk.NewCoins(sdk.NewCoin(fxtypes.LegacyFXDenom, sdkmath.NewInt(1e18)))
	allowanceList := []feegrant.FeeAllowanceI{
		&feegrant.BasicAllowance{
			SpendLimit: coins,
			Expiration: &now,
		},
		&feegrant.PeriodicAllowance{
			Basic: feegrant.BasicAllowance{
				SpendLimit: coins,
				Expiration: &now,
			},
			Period:           1,
			PeriodSpendLimit: coins,
			PeriodCanSpend:   coins,
			PeriodReset:      now,
		},
	}

	allowanceListLen := len(allowanceList)
	for i := 0; i < allowanceListLen; i++ {
		msg, ok := allowanceList[i].(proto.Message)
		require.True(t, ok)
		value, err := codectypes.NewAnyWithValue(msg)
		require.NoError(t, err)
		allowanceList = append(allowanceList,
			&feegrant.AllowedMsgAllowance{
				Allowance: value,
				AllowedMessages: []string{
					codectypes.MsgTypeURL(msg),
				},
			},
		)
	}

	granters := make([]sdk.AccAddress, 0, len(allowanceList))
	grantees := make([]sdk.AccAddress, 0, len(allowanceList))
	for i, allowance := range allowanceList {
		granters = append(granters, helpers.GenAccAddress())
		grantees = append(grantees, helpers.GenAccAddress())
		app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, granters[i]))
		app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, grantees[i]))
		require.NoError(t,
			app.FeeGrantKeeper.GrantAllowance(
				ctx,
				granters[i],
				grantees[i],
				allowance,
			),
		)
	}

	err := v8.MigrateFeegrant(ctx, app.AppCodec(), runtime.NewKVStoreService(app.GetKey(feegrant.StoreKey)), app.AccountKeeper)
	require.NoError(t, err)

	for i, allowance := range allowanceList {
		getAllowance, err := app.FeeGrantKeeper.GetAllowance(ctx, granters[i], grantees[i])
		require.NoError(t, err)
		require.NotEqual(t, allowance, getAllowance)
	}
}
