package v6_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v6/app"
	v6 "github.com/functionx/fx-core/v6/app/upgrades/v6"
	"github.com/functionx/fx-core/v6/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v6/types"
	migratetypes "github.com/functionx/fx-core/v6/x/migrate/types"
)

type UpgradeTestSuite struct {
	suite.Suite

	app       *app.App
	ctx       sdk.Context
	valNumber int
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) SetupTest() {
	s.valNumber = tmrand.Intn(5) + 5
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(s.valNumber, sdk.Coins{})
	s.app = helpers.SetupWithGenesisValSet(s.T(), valSet, valAccounts, valBalances...)
	s.ctx = s.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          s.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})
}

func (s *UpgradeTestSuite) CommitBlock(block int64) {
	helpers.MintBlock(s.app, s.ctx, block)
}

func (s *UpgradeTestSuite) TestUpdateParams() {
	s.NoError(v6.UpdateParams(s.ctx, s.app.AppKeepers))
	s.CommitBlock(10)
}

func (s *UpgradeTestSuite) TestAutoUndelegate_And_ExportDelegate() {
	delPrivKey := helpers.NewPriKey()
	delAddr := sdk.AccAddress(delPrivKey.PubKey().Address())
	helpers.AddTestAddr(s.app, s.ctx, delAddr, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10000))))
	account := s.app.AccountKeeper.GetAccount(s.ctx, delAddr)
	s.NoError(account.SetPubKey(delPrivKey.PubKey()))
	s.app.AccountKeeper.SetAccount(s.ctx, account)

	validators := s.app.StakingKeeper.GetAllValidators(s.ctx)
	s.Equal(s.valNumber, len(validators))
	validator := validators[0]

	_, err := s.app.StakingKeeper.Delegate(s.ctx, delAddr, sdkmath.NewInt(100), stakingtypes.Unbonded, validator, true)
	s.NoError(err)

	newPrivKey := helpers.NewEthPrivKey()
	_, err = s.app.MigrateKeeper.MigrateAccount(s.ctx, &migratetypes.MsgMigrateAccount{
		From:      delAddr.String(),
		To:        common.BytesToAddress(newPrivKey.PubKey().Address()).String(),
		Signature: "",
	})
	s.NoError(err)

	delegations := v6.AutoUndelegate(s.ctx, s.app.StakingKeeper)
	s.Equal(s.valNumber+1, len(delegations))

	delegations = v6.ExportDelegate(s.ctx, s.app.MigrateKeeper, delegations)
	s.Equal(1, len(delegations))
}
