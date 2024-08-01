package helpers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

type BankSuite struct {
	*require.Assertions
	ctx        sdk.Context
	bankKeeper bankkeeper.Keeper
}

func (s *BankSuite) Init(ass *require.Assertions, ctx sdk.Context, bankKeeper bankkeeper.Keeper) *BankSuite {
	s.Assertions = ass
	s.ctx = ctx
	s.bankKeeper = bankKeeper
	return s
}

func (s *BankSuite) NewSigner(ctx sdk.Context) *Signer {
	signer := NewSigner(NewEthPrivKey())
	err := s.bankKeeper.MintCoins(ctx, banktypes.ModuleName, NewStakingCoins(10_000, 18))
	s.NoError(err)
	err = s.bankKeeper.SendCoinsFromModuleToAccount(ctx, banktypes.ModuleName, signer.AccAddress(), NewStakingCoins(10_000, 18))
	s.NoError(err)
	return signer
}

func (s *BankSuite) GetTotalSupply() sdk.Coins {
	totalSupply, response, err := s.bankKeeper.GetPaginatedTotalSupply(s.ctx, &query.PageRequest{})
	s.NoError(err)
	s.NotNil(response)
	return totalSupply
}

func (s *BankSuite) GetSupply(denom string) sdk.Coin {
	return s.bankKeeper.GetSupply(s.ctx, denom)
}

func (s *BankSuite) GetAllBalances(addr sdk.AccAddress) sdk.Coins {
	return s.bankKeeper.GetAllBalances(s.ctx, addr)
}

func (s *BankSuite) SendCoins(fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) {
	err := s.bankKeeper.SendCoins(s.ctx, fromAddr, toAddr, amt)
	s.NoError(err)
}
