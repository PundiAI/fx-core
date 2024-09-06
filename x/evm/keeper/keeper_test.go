package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/evm/testutil"
)

type KeeperTestSuite struct {
	helpers.BaseSuite
	testutil.EVMSuite
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.BaseSuite.SetupTest()
	s.EVMSuite.Init(s.Require(), s.Ctx, s.App.EvmKeeper, s.BaseSuite.NewSigner())
	s.EVMSuite.WithGasPrice(big.NewInt(500 * 1e9))
}

func (s *KeeperTestSuite) NewERC20Suite() testutil.ERC20Suite {
	return testutil.NewERC20Suite(s.EVMSuite)
}

func (s *KeeperTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, evmtypes.ModuleName, coins)
	s.Require().NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, evmtypes.ModuleName, authtypes.FeeCollectorName, coins)
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

func (s *KeeperTestSuite) AssertContractAddr(singer common.Address, newContractAddr common.Address) {
	nonce, err := s.App.AccountKeeper.GetSequence(s.Ctx, singer.Bytes())
	s.NoError(err)

	contractAddr := crypto.CreateAddress(singer, nonce-1)
	s.Equal(contractAddr, newContractAddr)
}
