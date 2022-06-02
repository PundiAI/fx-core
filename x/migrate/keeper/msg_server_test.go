package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/x/migrate/types"
)

func TestKeeper_MigrateAccount(t *testing.T) {
	myApp, _, delegateAddressArr := initTest(t)
	ctx := myApp.BaseApp.NewContext(false, tmproto.Header{})
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	b1 := myApp.BankKeeper.GetAllBalances(ctx, alice)
	require.False(t, b1.Empty())
	b2 := myApp.BankKeeper.GetAllBalances(ctx, bob)
	require.False(t, b1.Empty())

	_, found := myApp.MigrateKeeper.GetMigrateRecord(ctx, alice)
	require.False(t, found)

	_, found = myApp.MigrateKeeper.GetMigrateRecord(ctx, bob)
	require.False(t, found)

	found = myApp.MigrateKeeper.HasMigratedDirectionFrom(ctx, alice)
	require.False(t, found)

	found = myApp.MigrateKeeper.HasMigratedDirectionTo(ctx, bob)
	require.False(t, found)

	_, err := myApp.MigrateKeeper.MigrateAccount(sdk.WrapSDKContext(ctx), &types.MsgMigrateAccount{
		From:      alice.String(),
		To:        bob.String(),
		Signature: "",
	})
	require.NoError(t, err)

	record, found := myApp.MigrateKeeper.GetMigrateRecord(ctx, alice)
	require.True(t, found)
	require.Equal(t, record.From, alice.String())

	record, found = myApp.MigrateKeeper.GetMigrateRecord(ctx, bob)
	require.True(t, found)
	require.Equal(t, record.To, bob.String())

	found = myApp.MigrateKeeper.HasMigratedDirectionFrom(ctx, alice)
	require.True(t, found)

	found = myApp.MigrateKeeper.HasMigratedDirectionTo(ctx, bob)
	require.True(t, found)

	bb1 := myApp.BankKeeper.GetAllBalances(ctx, alice)
	require.True(t, bb1.Empty())
	bb2 := myApp.BankKeeper.GetAllBalances(ctx, bob)
	require.Equal(t, b1, bb2.Sub(b2))
}
