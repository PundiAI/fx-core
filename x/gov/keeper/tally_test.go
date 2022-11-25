package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

func TestKeeper_Tally(t *testing.T) {
	val, genAccs, _ := helpers.GenerateGenesisValidator(4, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1e5).MulRaw(1e18))))
	app := helpers.SetupWithGenesisValSet(t, val, genAccs)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{
		ProposerAddress: val.Proposer.Address.Bytes(),
	})
	ctx = ctx.WithBlockHeight(fxtypes.UpgradeTrigonometric2Block())

	proposalContent := erc20types.NewRegisterCoinProposal("register erc20", "foo", banktypes.Metadata{
		Display: "test",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "test",
				Exponent: 0,
			},
			{
				Denom:    "TEST",
				Exponent: 18,
			},
		},
		Base:   "test",
		Name:   "test",
		Symbol: "TEST",
	})
	proposal, err := app.GovKeeper.SubmitProposal(ctx, proposalContent)
	require.NoError(t, err)

	proposalID := proposal.ProposalId
	proposal.Status = types.StatusVotingPeriod
	app.GovKeeper.SetProposal(ctx, proposal)

	poolBalances := app.BankKeeper.GetAllBalances(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress())
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(400).MulRaw(1e18))), poolBalances)

	votoer := genAccs[0].GetAddress().Bytes()
	validator, found := app.StakingKeeper.GetValidator(ctx, votoer)
	require.True(t, found)
	require.Equal(t, validator.Tokens, sdk.NewInt(100).MulRaw(1e18))

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, votoer, types.NewNonSplitVoteOption(types.OptionYes)))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	passes, burnDeposits, tallyResults := app.GovKeeper.Tally(ctx, proposal)
	assert.True(t, passes)
	assert.False(t, burnDeposits)
	t.Log(tallyResults)
}
