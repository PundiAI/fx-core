package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoenc "github.com/tendermint/tendermint/crypto/encoding"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/staking/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	app    *app.App
	signer *helpers.Signer

	valAccounts     []authtypes.GenesisAccount
	currentVoteInfo []abci.VoteInfo
	nextVoteInfo    []abci.VoteInfo
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
	valNumber := tmrand.Intn(5) + 6
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          suite.app.LastBlockHeight() + 1,
		Time:            time.Now().UTC(),
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})
	suite.ctx = suite.ctx.WithConsensusParams(helpers.ABCIConsensusParams)
	suite.valAccounts = valAccounts

	for _, validator := range valSet.Validators {
		signingInfo := slashingtypes.NewValidatorSigningInfo(
			validator.Address.Bytes(),
			suite.ctx.BlockHeight(),
			100,
			time.Unix(0, 0),
			false,
			0,
		)
		suite.app.SlashingKeeper.SetValidatorSigningInfo(suite.ctx, validator.Address.Bytes(), signingInfo)
	}

	vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
	infos := make([]abci.VoteInfo, 0, len(vals))
	for _, val := range vals {
		addr, err := val.GetConsAddr()
		suite.Require().NoError(err)
		infos = append(infos, abci.VoteInfo{Validator: abci.Validator{Address: addr, Power: 100}})
	}
	suite.currentVoteInfo = infos
	suite.nextVoteInfo = infos

	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.UnbondingTime = time.Second
	stakingParams.MaxValidators = 150
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

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

func (suite *KeeperTestSuite) SetSigningInfo(pk cryptotypes.PubKey, jailed ...bool) {
	consAddr := sdk.ConsAddress(pk.Address())
	jailedUntil := time.Unix(0, 0)
	if len(jailed) > 0 && jailed[0] {
		jailedUntil = suite.ctx.BlockHeader().Time.Add(time.Second)
	}
	signingInfo := slashingtypes.NewValidatorSigningInfo(consAddr, 0, 0, jailedUntil, false, 0)
	suite.app.SlashingKeeper.SetValidatorSigningInfo(suite.ctx, consAddr, signingInfo)

	err := suite.app.SlashingKeeper.AddPubkey(suite.ctx, pk)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) CommitEndBlock() []abci.ValidatorUpdate {
	header := suite.ctx.BlockHeader()
	res := suite.app.EndBlock(abci.RequestEndBlock{Height: header.Height})
	suite.app.Commit()
	return res.ValidatorUpdates
}

func (suite *KeeperTestSuite) CommitBeginBlock(valUpdate []abci.ValidatorUpdate) {
	header := suite.ctx.BlockHeader()

	// begin block
	header.Time = header.Time.Add(5 * time.Second)
	header.Height += 1

	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
		LastCommitInfo: abci.LastCommitInfo{
			Votes: suite.currentVoteInfo,
		},
	})
	suite.ctx = suite.app.NewContext(false, header)
	suite.ctx = suite.ctx.WithConsensusParams(helpers.ABCIConsensusParams)

	pkUpdate := make(map[string]int64, len(valUpdate))
	for _, pk := range valUpdate {
		tmPk, err := cryptoenc.PubKeyFromProto(pk.PubKey)
		suite.Require().NoError(err)
		pkUpdate[sdk.ConsAddress(tmPk.Address()).String()] = pk.Power
	}

	newVoteInfo := make([]abci.VoteInfo, 0, len(suite.currentVoteInfo))
	for _, info := range suite.nextVoteInfo {
		consAddr := sdk.ConsAddress(info.Validator.Address)
		power, ok := pkUpdate[consAddr.String()]
		if ok && power == 0 {
			delete(pkUpdate, consAddr.String())
			continue
		}
		newVoteInfo = append(newVoteInfo, info)
	}
	for addr, power := range pkUpdate {
		consAddr, err := sdk.ConsAddressFromBech32(addr)
		suite.Require().NoError(err)
		newVoteInfo = append(newVoteInfo, abci.VoteInfo{Validator: abci.Validator{Address: consAddr, Power: power}})
	}

	suite.currentVoteInfo, suite.nextVoteInfo = suite.nextVoteInfo, newVoteInfo
}

func (suite *KeeperTestSuite) CurrentVoteFound(pk cryptotypes.PubKey) bool {
	consAddr := sdk.ConsAddress(pk.Address())
	for _, info := range suite.currentVoteInfo {
		if consAddr.Equals(sdk.ConsAddress(info.Validator.Address)) {
			return true
		}
	}
	return false
}

func (suite *KeeperTestSuite) NextVoteFound(pk cryptotypes.PubKey) bool {
	consAddr := sdk.ConsAddress(pk.Address())
	for _, info := range suite.nextVoteInfo {
		if consAddr.Equals(sdk.ConsAddress(info.Validator.Address)) {
			return true
		}
	}
	return false
}

func (suite *KeeperTestSuite) Commit(count ...int) {
	number := 1
	if len(count) > 0 && count[0] > 0 {
		number = count[0]
	}
	for i := 0; i < number; i++ {
		valUpdates := suite.CommitEndBlock()
		suite.CommitBeginBlock(valUpdates)
	}
}

func (suite *KeeperTestSuite) CreateValidatorJailed() ([]abci.ValidatorUpdate, sdk.AccAddress, sdk.ConsAddress) {
	accAddr := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	initBalance := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18))))
	helpers.AddTestAddr(suite.app, suite.ctx, accAddr, initBalance)

	activeMinSelfDelegateCoin := sdkmath.NewInt(99).Mul(sdkmath.NewInt(1e18))
	activeSelfDelegateCoin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100).Mul(sdkmath.NewInt(1e18)))
	des := stakingtypes.Description{Moniker: "test-node"}
	rates := stakingtypes.CommissionRates{
		Rate:          sdk.NewDecWithPrec(1, 2),
		MaxRate:       sdk.NewDecWithPrec(5, 2),
		MaxChangeRate: sdk.NewDecWithPrec(1, 2),
	}

	consPubKey := ed25519.GenPrivKey().PubKey()
	newValMsg, err := stakingtypes.NewMsgCreateValidator(sdk.ValAddress(accAddr), consPubKey, activeSelfDelegateCoin, des, rates, activeMinSelfDelegateCoin)
	suite.Require().NoError(err)
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).CreateValidator(suite.ctx, newValMsg)
	suite.Require().NoError(err)

	suite.Commit(3)
	valUpdates := suite.CommitEndBlock()

	// validator not jailed
	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().True(found)
	suite.Require().False(validator.IsJailed())

	suite.CommitBeginBlock(valUpdates)
	// jailed validator
	undelMsg2 := stakingtypes.NewMsgUndelegate(accAddr, sdk.ValAddress(accAddr), activeSelfDelegateCoin.SubAmount(sdkmath.NewInt(1e18).Mul(sdkmath.NewInt(10))))
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Undelegate(suite.ctx, undelMsg2)
	suite.Require().NoError(err)
	suite.Commit(3)
	valUpdates = suite.CommitEndBlock()

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())
	consAddr, err := validator.GetConsAddr()
	suite.Require().NoError(err)
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, consAddr)
	suite.Require().True(found)
	// if not found, validator maybe not bounded

	return valUpdates, accAddr, consAddr
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
