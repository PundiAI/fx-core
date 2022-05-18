package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/x/migrate/types"
)

func TestKeeper_MigrateAccount(t *testing.T) {
	fxcore, _, delegateAddressArr := initTest(t)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	b1 := fxcore.BankKeeper.GetAllBalances(ctx, alice)
	require.False(t, b1.Empty())
	b2 := fxcore.BankKeeper.GetAllBalances(ctx, bob)
	require.False(t, b1.Empty())

	_, found := fxcore.MigrateKeeper.GetMigrateRecord(ctx, alice)
	require.False(t, found)

	_, found = fxcore.MigrateKeeper.GetMigrateRecord(ctx, bob)
	require.False(t, found)

	found = fxcore.MigrateKeeper.HasMigratedDirectionFrom(ctx, alice)
	require.False(t, found)

	found = fxcore.MigrateKeeper.HasMigratedDirectionTo(ctx, bob)
	require.False(t, found)

	_, err := fxcore.MigrateKeeper.MigrateAccount(sdk.WrapSDKContext(ctx), &types.MsgMigrateAccount{
		From:      alice.String(),
		To:        bob.String(),
		Signature: "",
	})
	require.NoError(t, err)

	record, found := fxcore.MigrateKeeper.GetMigrateRecord(ctx, alice)
	require.True(t, found)
	require.Equal(t, record.From, alice.String())

	record, found = fxcore.MigrateKeeper.GetMigrateRecord(ctx, bob)
	require.True(t, found)
	require.Equal(t, record.To, bob.String())

	found = fxcore.MigrateKeeper.HasMigratedDirectionFrom(ctx, alice)
	require.True(t, found)

	found = fxcore.MigrateKeeper.HasMigratedDirectionTo(ctx, bob)
	require.True(t, found)

	bb1 := fxcore.BankKeeper.GetAllBalances(ctx, alice)
	require.True(t, bb1.Empty())
	bb2 := fxcore.BankKeeper.GetAllBalances(ctx, bob)
	require.Equal(t, b1, bb2.Sub(b2))
}
