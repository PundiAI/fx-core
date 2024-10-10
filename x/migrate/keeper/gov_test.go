package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxgovkeeper "github.com/functionx/fx-core/v8/x/gov/keeper"
	migratekeeper "github.com/functionx/fx-core/v8/x/migrate/keeper"
)

func (suite *KeeperTestSuite) TestMigrateGovInactive() {
	suite.MintToken(suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(2)
	acc1 := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	acc2 := sdk.AccAddress(keys[1].PubKey().Address().Bytes())

	ethKeys := suite.GenerateEthAcc(1)
	ethAcc1 := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	content, _ := govv1beta1.ContentFromProposalType("title", "description", "Text")
	legacyContent, _ := govv1.NewLegacyContent(content, suite.govAddr)
	anys, _ := sdktx.SetMsgs([]sdk.Msg{legacyContent})
	govImpl := fxgovkeeper.NewMsgServerImpl(suite.App.GovKeeper)

	// submit proposal
	proposal, err := govImpl.SubmitProposal(suite.Ctx, &govv1.MsgSubmitProposal{
		Messages:       anys,
		InitialDeposit: helpers.NewStakingCoins(1000, 18),
		Proposer:       acc1.String(),
		Title:          content.GetTitle(),
		Summary:        content.GetDescription(),
	})
	suite.Require().NoError(err)

	// deposit proposal
	_, err = govImpl.Deposit(suite.Ctx, &govv1.MsgDeposit{
		ProposalId: proposal.ProposalId,
		Depositor:  acc2.String(),
		Amount:     helpers.NewStakingCoins(1000, 18),
	})
	suite.Require().NoError(err)

	govParams, _ := suite.App.GovKeeper.Keeper.Params.Get(suite.Ctx)
	suite.Ctx = suite.Ctx.WithBlockTime(suite.Ctx.BlockTime().Add(*govParams.MaxDepositPeriod).Add(time.Second))
	m := migratekeeper.NewGovMigrate(suite.App.GovKeeper, suite.App.AccountKeeper)

	err = m.Validate(suite.Ctx, suite.App.AppCodec(), acc1, ethAcc1)
	suite.Require().ErrorIs(err, sdkerrors.ErrInvalidRequest)

	err = m.Validate(suite.Ctx, suite.App.AppCodec(), acc2, ethAcc1)
	suite.Require().ErrorIs(err, sdkerrors.ErrInvalidRequest)
}

func (suite *KeeperTestSuite) TestMigrateGovActive() {
	suite.MintToken(suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(2)
	acc1 := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	acc2 := sdk.AccAddress(keys[1].PubKey().Address().Bytes())

	ethKeys := suite.GenerateEthAcc(1)
	ethAcc1 := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	content, _ := govv1beta1.ContentFromProposalType("title", "description", "Text")
	legacyContent, _ := govv1.NewLegacyContent(content, suite.govAddr)

	anys, _ := sdktx.SetMsgs([]sdk.Msg{legacyContent})
	govImpl := fxgovkeeper.NewMsgServerImpl(suite.App.GovKeeper)

	// submit proposal
	proposal, err := govImpl.SubmitProposal(suite.Ctx, &govv1.MsgSubmitProposal{
		Messages:       anys,
		InitialDeposit: helpers.NewStakingCoins(30000, 18),
		Proposer:       acc1.String(),
		Metadata:       "",
		Title:          content.GetTitle(),
		Summary:        content.GetDescription(),
		Expedited:      false,
	})
	suite.Require().NoError(err)

	// vote proposal
	_, err = govImpl.Vote(suite.Ctx, &govv1.MsgVote{
		ProposalId: proposal.ProposalId,
		Voter:      acc2.String(),
		Option:     govv1.OptionYes,
	})
	suite.Require().NoError(err)

	// check
	govParams, err := suite.App.GovKeeper.Keeper.Params.Get(suite.Ctx)
	suite.Require().NoError(err)
	suite.Ctx = suite.Ctx.WithBlockTime(suite.Ctx.BlockTime().Add(*govParams.MaxDepositPeriod).Add(time.Second))
	m := migratekeeper.NewGovMigrate(suite.App.GovKeeper, suite.App.AccountKeeper)

	err = m.Validate(suite.Ctx, suite.App.AppCodec(), acc1, ethAcc1)
	suite.Require().ErrorIs(err, sdkerrors.ErrInvalidRequest)

	err = m.Validate(suite.Ctx, suite.App.AppCodec(), acc2, ethAcc1)
	suite.Require().ErrorIs(err, sdkerrors.ErrInvalidRequest)
}
