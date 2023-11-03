package v6_test

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v6/app"
	v6 "github.com/functionx/fx-core/v6/app/upgrades/v6"
	"github.com/functionx/fx-core/v6/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v6/types"
	layer2types "github.com/functionx/fx-core/v6/x/layer2/types"
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

	for symbol := range v6.Layer2GenesisTokenAddress {
		v6.Layer2GenesisTokenAddress[symbol] = common.BytesToAddress(helpers.NewEthPrivKey().PubKey().Address().Bytes()).String()
	}
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

func (s *UpgradeTestSuite) TestMigrateMetadata() {
	for symbol := range v6.Layer2GenesisTokenAddress {
		hasDenomMetaData := s.app.BankKeeper.HasDenomMetaData(s.ctx, strings.ToLower(symbol))
		s.False(hasDenomMetaData)

		metadata := fxtypes.GetCrossChainMetadata(symbol, symbol, 18)
		s.app.BankKeeper.SetDenomMetaData(s.ctx, metadata)
	}

	v6.MigrateMetadata(s.ctx, s.app.BankKeeper)

	for symbol, address := range v6.Layer2GenesisTokenAddress {
		metadata, found := s.app.BankKeeper.GetDenomMetaData(s.ctx, strings.ToLower(symbol))
		s.True(found)
		s.True(len(metadata.DenomUnits) > 0)
		s.Subset(metadata.DenomUnits[0].Aliases, []string{fmt.Sprintf("%s%s", layer2types.ModuleName, address)})
	}
}

func (s *UpgradeTestSuite) TestMigrateLayer2Module() {
	v6.MigrateLayer2Module(s.ctx, s.app.Layer2Keeper)
	for _, address := range v6.Layer2GenesisTokenAddress {
		bridgeToken := s.app.Layer2Keeper.GetBridgeTokenDenom(s.ctx, address)
		s.Equal(bridgeToken.Token, address)
		s.Equal(bridgeToken.Denom, fmt.Sprintf("%s%s", layer2types.ModuleName, address))
	}
}
