package keeper_test

import (
	"crypto/ecdsa"
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	arbitrumtypes "github.com/functionx/fx-core/v8/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v8/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v8/x/bsc/types"
	"github.com/functionx/fx-core/v8/x/crosschain/keeper"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	layer2types "github.com/functionx/fx-core/v8/x/layer2/types"
	optimismtypes "github.com/functionx/fx-core/v8/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v8/x/polygon/types"
	tronkeeper "github.com/functionx/fx-core/v8/x/tron/keeper"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app          *app.App
	ctx          sdk.Context
	oracleAddrs  []sdk.AccAddress
	bridgerAddrs []sdk.AccAddress
	externalPris []*ecdsa.PrivateKey
	valAddrs     []sdk.ValAddress
	chainName    string
}

func TestCrosschainKeeperTestSuite(t *testing.T) {
	mustTestModule := []string{
		trontypes.ModuleName,
		ethtypes.ModuleName,
	}

	subModules := mustTestModule
	if os.Getenv("TEST_CROSSCHAIN") == "true" {
		subModules = append(subModules, []string{
			bsctypes.ModuleName,
			polygontypes.ModuleName,
			avalanchetypes.ModuleName,
			arbitrumtypes.ModuleName,
			optimismtypes.ModuleName,
			layer2types.ModuleName,
		}...)
	}

	for _, moduleName := range subModules {
		suite.Run(t, &KeeperTestSuite{
			chainName: moduleName,
		})
	}
}

func (suite *KeeperTestSuite) MsgServer() types.MsgServer {
	if suite.chainName == trontypes.ModuleName {
		return tronkeeper.NewMsgServerImpl(suite.app.TronKeeper)
	}
	return keeper.NewMsgServerImpl(suite.Keeper())
}

func (suite *KeeperTestSuite) QueryClient() types.QueryClient {
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(suite.Keeper()))
	return types.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) Keeper() keeper.Keeper {
	switch suite.chainName {
	case bsctypes.ModuleName:
		return suite.app.BscKeeper
	case polygontypes.ModuleName:
		return suite.app.PolygonKeeper
	case trontypes.ModuleName:
		return suite.app.TronKeeper.Keeper
	case ethtypes.ModuleName:
		return suite.app.EthKeeper
	case avalanchetypes.ModuleName:
		return suite.app.AvalancheKeeper
	case arbitrumtypes.ModuleName:
		return suite.app.ArbitrumKeeper
	case optimismtypes.ModuleName:
		return suite.app.OptimismKeeper
	case layer2types.ModuleName:
		return suite.app.Layer2Keeper
	default:
		panic("invalid chain name")
	}
}

func (suite *KeeperTestSuite) SetupTest() {
	valNumber := tmrand.Intn(types.MaxOracleSize-4) + 4

	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          suite.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})

	suite.oracleAddrs = helpers.AddTestAddrs(suite.app, suite.ctx, valNumber, sdk.NewCoins(types.NewDelegateAmount(sdkmath.NewInt(300*1e3).MulRaw(1e18))))
	suite.bridgerAddrs = helpers.AddTestAddrs(suite.app, suite.ctx, valNumber, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(300*1e3).MulRaw(1e18))))
	suite.externalPris = helpers.CreateMultiECDSA(valNumber)

	suite.valAddrs = make([]sdk.ValAddress, valNumber)
	for i := 0; i < valNumber; i++ {
		suite.valAddrs[i] = valAccounts[i].GetAddress().Bytes()
	}

	proposalOracle := &types.ProposalOracle{}
	for _, oracle := range suite.oracleAddrs {
		proposalOracle.Oracles = append(proposalOracle.Oracles, oracle.String())
	}
	suite.Keeper().SetProposalOracle(suite.ctx, proposalOracle)
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) PubKeyToExternalAddr(publicKey ecdsa.PublicKey) string {
	address := crypto.PubkeyToAddress(publicKey)
	return types.ExternalAddrToStr(suite.chainName, address.Bytes())
}

func (suite *KeeperTestSuite) Commit(block ...int64) {
	suite.ctx = helpers.MintBlock(suite.app, suite.ctx, block...)
}

func (suite *KeeperTestSuite) SignOracleSetConfirm(external *ecdsa.PrivateKey, oracleSet *types.OracleSet) (string, []byte) {
	externalAddress := crypto.PubkeyToAddress(external.PublicKey).String()
	gravityId := suite.Keeper().GetGravityID(suite.ctx)
	checkpoint, err := oracleSet.GetCheckpoint(gravityId)
	suite.NoError(err)
	signature, err := types.NewEthereumSignature(checkpoint, external)
	suite.NoError(err)
	if trontypes.ModuleName == suite.chainName {
		externalAddress = tronaddress.PubkeyToAddress(external.PublicKey).String()

		checkpoint, err = trontypes.GetCheckpointOracleSet(oracleSet, gravityId)
		suite.Require().NoError(err)

		signature, err = trontypes.NewTronSignature(checkpoint, external)
		suite.Require().NoError(err)
	}
	return externalAddress, signature
}

func (suite *KeeperTestSuite) SendClaim(externalClaim types.ExternalClaim) {
	var err error
	switch claim := externalClaim.(type) {
	case *types.MsgSendToFxClaim:
		_, err = suite.MsgServer().SendToFxClaim(suite.ctx, claim)
	case *types.MsgBridgeCallClaim:
		_, err = suite.MsgServer().BridgeCallClaim(suite.ctx, claim)
	}
	suite.NoError(err)

	err = suite.Keeper().ExecuteClaim(suite.ctx, externalClaim.GetEventNonce())
	suite.Require().NoError(err)
}
