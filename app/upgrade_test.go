package app_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/app"
	nextversion "github.com/functionx/fx-core/v8/app/upgrades/v8"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func Test_UpgradeAndMigrate(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	myApp, chainId := buildApp(t, fxtypes.MainnetChainId)
	myApp.SetStoreLoader(upgradetypes.UpgradeStoreLoader(myApp.LastBlockHeight()+1, nextversion.Upgrade.StoreUpgrades()))
	require.NoError(t, myApp.LoadLatestVersion())

	ctx := newContext(t, myApp, chainId, false)
	require.NoError(t, myApp.UpgradeKeeper.ScheduleUpgrade(ctx, upgradetypes.Plan{
		Name:   nextversion.Upgrade.UpgradeName,
		Height: ctx.BlockHeight() + 1,
	}))

	header := ctx.BlockHeader()
	header.Height = header.Height + 1

	require.NotPanics(t, func() {
		myApp.BeginBlock(abci.RequestBeginBlock{
			Header: header,
		})
	})
	require.NotPanics(t, func() {
		myApp.EndBlock(abci.RequestEndBlock{
			Height: header.Height,
		})
	})

	ctx = newContext(t, myApp, chainId, true)
	ingProposalIds := govDepositAndVote(t, ctx, myApp)

	header.Time = tmtime.Now().Add(21 * 24 * time.Hour) // exec proposal
	require.NotPanics(t, func() {
		myApp.BeginBlock(abci.RequestBeginBlock{
			Header: header,
		})
	})
	require.NotPanics(t, func() {
		myApp.EndBlock(abci.RequestEndBlock{
			Height: header.Height,
		})
	})

	checkProposalPassed(t, ctx, myApp, ingProposalIds)
}

func buildApp(t *testing.T, chainId string) (*app.App, string) {
	fxtypes.SetConfig(true)

	home := filepath.Join(os.Getenv("HOME"), "tmp")
	fxtypes.SetChainId(chainId)

	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(home, "data"))
	require.NoError(t, err)

	makeEncodingConfig := app.MakeEncodingConfig()
	myApp := app.New(log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowAll()),
		db, nil, false, map[int64]bool{}, home, 0,
		makeEncodingConfig, app.EmptyAppOptions{}, baseapp.SetChainID(chainId))
	return myApp, chainId
}

func newContext(t *testing.T, myApp *app.App, chainId string, deliveState bool) sdk.Context {
	header := tmproto.Header{
		ChainID: chainId,
		Height:  myApp.LastBlockHeight(),
		Time:    tmtime.Now(),
	}
	var ctx sdk.Context
	if deliveState {
		ctx = myApp.NewContext(false, header)
	} else {
		ctx = myApp.NewUncachedContext(false, header)
	}
	// set the first validator to proposer
	validators := myApp.StakingKeeper.GetAllValidators(ctx)
	assert.True(t, len(validators) > 0)
	var pubKey cryptotypes.PubKey
	assert.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubKey))
	ctx = ctx.WithProposer(pubKey.Address().Bytes())
	return ctx
}

func govDepositAndVote(t *testing.T, ctx sdk.Context, myApp *app.App) []uint64 {
	var ingProposalIds []uint64
	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntWithDecimal(1_000_000_000, 18)))
	require.NoError(t, myApp.BankKeeper.MintCoins(ctx, minttypes.ModuleName, amount))
	accAddr := helpers.GenAccAddress()
	require.NoError(t, myApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, accAddr, amount))

	endTime := tmtime.Now().Add(21 * 24 * time.Hour)
	myApp.GovKeeper.IterateInactiveProposalsQueue(ctx, endTime, func(proposal v1.Proposal) (stop bool) {
		minDeposit := myApp.GovKeeper.NeedMinDeposit(ctx, proposal)
		if sdk.NewCoins(proposal.TotalDeposit...).IsAllGTE(minDeposit) {
			return false
		}
		insufficientCoins := sdk.NewCoins(minDeposit...).Sub(proposal.TotalDeposit...)
		deposit, err := myApp.GovKeeper.AddDeposit(ctx, proposal.Id, accAddr, insufficientCoins)
		require.NoError(t, err)
		require.True(t, deposit)
		return false
	})

	myApp.GovKeeper.IterateActiveProposalsQueue(ctx, endTime, func(proposal v1.Proposal) (stop bool) {
		if isUpgradeProposal(t, proposal) {
			ctx.Logger().Info("skip upgrade proposal", "id", proposal.Id)
			return false
		}
		ingProposalIds = append(ingProposalIds, proposal.Id)
		powerStoreIterator := myApp.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
		defer powerStoreIterator.Close()
		voteOptions := v1.WeightedVoteOptions{&v1.WeightedVoteOption{Option: v1.VoteOption_VOTE_OPTION_YES, Weight: "1"}}
		for ; powerStoreIterator.Valid(); powerStoreIterator.Next() {
			err := myApp.GovKeeper.AddVote(ctx, proposal.Id, powerStoreIterator.Value(), voteOptions, "")
			require.NoError(t, err)
		}
		return false
	})
	return ingProposalIds
}

func isUpgradeProposal(t *testing.T, proposal v1.Proposal) bool {
	msgs, err := proposal.GetMsgs()
	require.NoError(t, err)

	if len(msgs) != 1 {
		return false
	}
	msg := msgs[0]
	if _, ok := msg.(*upgradetypes.MsgSoftwareUpgrade); ok {
		return true
	}
	legacyMsg, ok := msg.(*v1.MsgExecLegacyContent)
	if !ok {
		return false
	}
	content, err := v1.LegacyContentFromMessage(legacyMsg)
	require.NoError(t, err)
	if _, ok = content.(*upgradetypes.SoftwareUpgradeProposal); ok { // nolint:staticcheck
		return true
	}
	return false
}

func checkProposalPassed(t *testing.T, ctx sdk.Context, myApp *app.App, ids []uint64) {
	for _, id := range ids {
		proposal, found := myApp.GovKeeper.GetProposal(ctx, id)
		require.True(t, found)
		require.Equal(t, proposal.Status, v1.ProposalStatus_PROPOSAL_STATUS_PASSED)
		t.Logf("proposal id:%d, status:%d, title:%s, msgType:%s", id, proposal.Status, proposal.Title, proposal.Messages[0].TypeUrl)
	}
}
