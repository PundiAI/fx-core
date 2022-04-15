package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	fxtypes "github.com/functionx/fx-core/types"
	migratekeeper "github.com/functionx/fx-core/x/migrate/keeper"
)

func Test_MigrateGovInactiveFunc(t *testing.T) {
	app, _, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	content := govtypes.ContentFromProposalType("title", "description", "Text")
	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(5000))))

	proposal, err := app.GovKeeper.SubmitProposal(ctx, content)
	require.NoError(t, err)

	_, err = app.GovKeeper.AddDeposit(ctx, proposal.ProposalId, alice, amount)
	require.NoError(t, err)

	deposit1, found := app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, alice)
	require.True(t, found)
	require.Equal(t, amount, deposit1.Amount)

	_, found = app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, bob)
	require.False(t, found)

	migrateKeeper := app.MigrateKeeper
	m := migratekeeper.NewGovMigrate(app.GetKey(govtypes.StoreKey), app.GovKeeper)
	err = m.Validate(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)

	_, found = app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, alice)
	require.False(t, found)

	deposit2, found := app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, bob)
	require.True(t, found)
	require.Equal(t, amount, deposit2.Amount)
}

func Test_MigrateGovActiveFunc(t *testing.T) {
	app, _, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	alice, bob, tom, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	content := govtypes.ContentFromProposalType("title", "description", "Text")
	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(5000))))

	proposal, err := app.GovKeeper.SubmitProposal(ctx, content)
	require.NoError(t, err)

	_, err = app.GovKeeper.AddDeposit(ctx, proposal.ProposalId, alice, amount)
	require.NoError(t, err)

	_, err = app.GovKeeper.AddDeposit(ctx, proposal.ProposalId, bob, amount)
	require.NoError(t, err)

	deposit1, found := app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, alice)
	require.True(t, found)
	require.Equal(t, amount, deposit1.Amount)

	deposit2, found := app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, bob)
	require.True(t, found)
	require.Equal(t, amount, deposit2.Amount)

	_, found = app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, tom)
	require.False(t, found)

	migrateKeeper := app.MigrateKeeper
	m := migratekeeper.NewGovMigrate(app.GetKey(govtypes.StoreKey), app.GovKeeper)
	err = m.Validate(ctx, migrateKeeper, alice, tom)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, alice, tom)
	require.NoError(t, err)

	_, found = app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, alice)
	require.False(t, found)

	deposit3, found := app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, bob)
	require.True(t, found)
	require.Equal(t, amount, deposit3.Amount)

	deposit4, found := app.GovKeeper.GetDeposit(ctx, proposal.ProposalId, tom)
	require.True(t, found)
	require.Equal(t, amount, deposit4.Amount)

}
