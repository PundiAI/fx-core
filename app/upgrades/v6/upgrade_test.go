package v6_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v6/app"
	v6 "github.com/functionx/fx-core/v6/app/upgrades/v6"
	"github.com/functionx/fx-core/v6/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v6/types"
)

type UpgradeTestSuite struct {
	suite.Suite

	app *app.App
	ctx sdk.Context
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) SetupTest() {
	valNumber := tmrand.Intn(5) + 5
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	s.app = helpers.SetupWithGenesisValSet(s.T(), valSet, valAccounts, valBalances...)
	s.ctx = s.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          s.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})
}

func (s *UpgradeTestSuite) CommitBlock(block int64) {
	_ = s.app.Commit()
	block--

	ctx := s.ctx
	nextBlockHeight := ctx.BlockHeight() + 1

	_ = s.app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			Height: nextBlockHeight,
		},
	})
	s.ctx = s.app.NewContext(false, ctx.BlockHeader())
	s.ctx = s.ctx.WithBlockHeight(nextBlockHeight)
	if block > 0 {
		s.CommitBlock(block)
	}
}

func (s *UpgradeTestSuite) TestUpdateParams() {
	s.NoError(v6.UpdateParams(s.ctx, s.app.AppKeepers))
	s.CommitBlock(10)
}
