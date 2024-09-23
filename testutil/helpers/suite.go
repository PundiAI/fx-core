package helpers

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/app"
)

type BaseSuite struct {
	suite.Suite
	MintValNumber int
	ValSet        *tmtypes.ValidatorSet
	ValAddr       []sdk.ValAddress
	App           *app.App
	Ctx           sdk.Context
}

func (s *BaseSuite) SetupTest() {
	valNumber := s.MintValNumber
	if valNumber <= 0 {
		valNumber = tmrand.Intn(10) + 1
	}
	valSet, valAccounts, valBalances := GenerateGenesisValidator(valNumber, sdk.Coins{})
	s.ValSet = valSet
	for _, account := range valAccounts {
		s.ValAddr = append(s.ValAddr, account.GetAddress().Bytes())
	}
	s.App = SetupWithGenesisValSet(s.T(), valSet, valAccounts, valBalances...)
	s.Ctx = s.App.GetContextForFinalizeBlock(nil)
	s.Ctx = s.Ctx.WithProposer(s.ValSet.Proposer.Address.Bytes())
}

func (s *BaseSuite) AddTestSigner() *Signer {
	signer := NewSigner(NewEthPrivKey())
	s.MintToken(signer.AccAddress(), NewStakingCoin(100, 18))
	return signer
}

func (s *BaseSuite) Commit() {
	s.Ctx = MintBlock(s.App, s.Ctx)
}

func (s *BaseSuite) AddTestSigners(accNum int, coin sdk.Coin) []*Signer {
	signers := make([]*Signer, accNum)
	for i := 0; i < accNum; i++ {
		signer := NewSigner(NewEthPrivKey())
		s.MintToken(signer.AccAddress(), coin)
	}
	return signers
}

func (s *BaseSuite) MintToken(address sdk.AccAddress, amount sdk.Coin) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, sdk.NewCoins(amount))
	s.Require().NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, address.Bytes(), sdk.NewCoins(amount))
	s.Require().NoError(err)
}
