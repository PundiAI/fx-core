package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

type KeeperTestSuite struct {
	helpers.BaseSuite
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) NewERC20TokenSuite() helpers.ERC20TokenSuite {
	return helpers.NewERC20Suite(s.Require(), s.AddTestSigner(), s.App.EvmKeeper)
}

func (s *KeeperTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, evmtypes.ModuleName, coins)
	s.Require().NoError(err)
	evmModuleAccount := s.App.AccountKeeper.GetModuleAccount(s.Ctx, evmtypes.ModuleName)
	err = s.App.BankKeeper.SendCoinsFromAccountToModuleVirtual(s.Ctx, evmModuleAccount.GetAddress(), authtypes.FeeCollectorName, coins)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) BurnEvmRefundFee(addr sdk.AccAddress, coins sdk.Coins) {
	err := s.App.BankKeeper.SendCoinsFromAccountToModule(s.Ctx, addr, authtypes.FeeCollectorName, coins)
	s.Require().NoError(err)

	bal := s.App.BankKeeper.GetBalance(s.Ctx, s.App.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName), fxtypes.DefaultDenom)
	err = s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, authtypes.FeeCollectorName, evmtypes.ModuleName, sdk.NewCoins(bal))
	s.Require().NoError(err)

	err = s.App.BankKeeper.BurnCoins(s.Ctx, evmtypes.ModuleName, sdk.NewCoins(bal))
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) AssertContractAddr(sender, newContractAddr common.Address) {
	nonce, err := s.App.AccountKeeper.GetSequence(s.Ctx, sender.Bytes())
	s.Require().NoError(err)

	contractAddr := crypto.CreateAddress(sender, nonce-1)
	s.Equal(contractAddr, newContractAddr)
}
