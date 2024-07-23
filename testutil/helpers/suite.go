package helpers

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v7/app"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

type BaseSuite struct {
	suite.Suite
	App *app.App
	Ctx sdk.Context
}

func (s *BaseSuite) SetupTest() {
	valNumber := tmrand.Intn(10) + 1
	valSet, valAccounts, valBalances := GenerateGenesisValidator(valNumber, sdk.Coins{})

	s.App = SetupWithGenesisValSet(s.T(), valSet, valAccounts, valBalances...)
	s.Ctx = s.App.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.ChainId(),
		Height:          s.App.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})
}

func (s *BaseSuite) NewSigner() *Signer {
	signer := NewSigner(NewEthPrivKey())
	AddTestAddr(s.App, s.Ctx, signer.AccAddress(), NewStakingCoins(100, 18))
	return signer
}
