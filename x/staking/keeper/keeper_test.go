package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v5/app"
	"github.com/functionx/fx-core/v5/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v5/types"
	"github.com/functionx/fx-core/v5/x/staking/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	app    *app.App
	signer *helpers.Signer

	valAccounts []authtypes.GenesisAccount
}

func TestKeeperTestSuite(t *testing.T) {
	fxtypes.SetConfig(false)
	suite.Run(t, new(KeeperTestSuite))
}

// Test helpers
func (suite *KeeperTestSuite) SetupTest() {
	suite.SetupSubTest()
}

func (suite *KeeperTestSuite) SetupSubTest() {
	valNumber := tmrand.Intn(10) + 1
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	suite.valAccounts = valAccounts

	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          suite.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})

	suite.signer = helpers.NewSigner(helpers.NewEthPrivKey())
	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100).MulRaw(1e18))))
}

func (suite *KeeperTestSuite) GenerateGrantPubkey() (cryptotypes.PrivKey, *codectypes.Any) {
	priKey, _ := ethsecp256k1.GenerateKey()
	pkAny, _ := codectypes.NewAnyWithValue(priKey.PubKey())
	return priKey, pkAny
}

func (suite *KeeperTestSuite) GenerateConsKey() (cryptotypes.PrivKey, *codectypes.Any) {
	priKey := ed25519.GenPrivKey()
	pkAny, _ := codectypes.NewAnyWithValue(priKey.PubKey())
	return priKey, pkAny
}

func (suite *KeeperTestSuite) TestHasValidatorGrant() {
	val := sdk.ValAddress(suite.valAccounts[0].GetAddress())
	addr := sdk.AccAddress(helpers.GenerateAddress().Bytes())

	auth := suite.app.StakingKeeper.HasValidatorGrant(suite.ctx, addr, val)
	suite.Require().False(auth)

	auth = suite.app.StakingKeeper.HasValidatorGrant(suite.ctx, sdk.AccAddress(val), val)
	suite.Require().True(auth)

	suite.app.StakingKeeper.UpdateValidatorOperator(suite.ctx, val, addr)

	auth = suite.app.StakingKeeper.HasValidatorGrant(suite.ctx, addr, val)
	suite.Require().True(auth)

	auth = suite.app.StakingKeeper.HasValidatorGrant(suite.ctx, sdk.AccAddress(val), val)
	suite.Require().False(auth)
}

func (suite *KeeperTestSuite) TestGrantRevokeAuthorization() {
	addr1 := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	addr2 := sdk.AccAddress(helpers.GenerateAddress().Bytes())

	getAuths, err := suite.app.AuthzKeeper.GetAuthorizations(suite.ctx, addr2, addr1)
	suite.Require().NoError(err)
	suite.Require().Len(getAuths, 0)

	auths := make([]authz.Authorization, 0, 1)
	a1 := authz.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgSend{}))
	auths = append(auths, a1)
	err = suite.app.StakingKeeper.GrantAuthorization(suite.ctx, addr2, addr1, auths, types.GrantExpirationTime)
	suite.Require().NoError(err)

	getAuths, err = suite.app.AuthzKeeper.GetAuthorizations(suite.ctx, addr2, addr1)
	suite.Require().NoError(err)
	suite.Require().Len(getAuths, 1)
	suite.Require().Equal(sdk.MsgTypeURL(&banktypes.MsgSend{}), getAuths[0].MsgTypeURL())

	err = suite.app.StakingKeeper.RevokeAuthorization(suite.ctx, addr2, addr1)
	suite.Require().NoError(err)

	getAuths, err = suite.app.AuthzKeeper.GetAuthorizations(suite.ctx, addr2, addr1)
	suite.Require().NoError(err)
	suite.Require().Len(getAuths, 0)
}

func (suite *KeeperTestSuite) TestValidatorOperator() {
	val := sdk.ValAddress(suite.valAccounts[0].GetAddress())
	addr1 := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	addr2 := sdk.AccAddress(helpers.GenerateAddress().Bytes())

	found := suite.app.StakingKeeper.HasValidatorOperator(suite.ctx, val)
	suite.Require().False(found)

	suite.app.StakingKeeper.UpdateValidatorOperator(suite.ctx, val, addr1)

	found = suite.app.StakingKeeper.HasValidatorOperator(suite.ctx, val)
	suite.Require().True(found)

	operAddr, found := suite.app.StakingKeeper.GetValidatorOperator(suite.ctx, val)
	suite.Require().True(found)
	suite.Require().Equal(addr1, operAddr)

	suite.app.StakingKeeper.UpdateValidatorOperator(suite.ctx, val, addr2)

	operAddr, found = suite.app.StakingKeeper.GetValidatorOperator(suite.ctx, val)
	suite.Require().True(found)
	suite.Require().Equal(addr2, operAddr)
}
