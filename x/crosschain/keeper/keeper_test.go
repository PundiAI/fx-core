package keeper_test

import (
	"crypto/ecdsa"
	"reflect"
	"regexp"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v6/app"
	"github.com/functionx/fx-core/v6/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v6/types"
	arbitrumtypes "github.com/functionx/fx-core/v6/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v6/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v6/x/bsc/types"
	"github.com/functionx/fx-core/v6/x/crosschain/keeper"
	"github.com/functionx/fx-core/v6/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v6/x/eth/types"
	optimismtypes "github.com/functionx/fx-core/v6/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v6/x/polygon/types"
	tronkeeper "github.com/functionx/fx-core/v6/x/tron/keeper"
	trontypes "github.com/functionx/fx-core/v6/x/tron/types"
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

func TestKeeperTestSuite(t *testing.T) {
	compile := regexp.MustCompile("^Test")
	subModules := []string{
		bsctypes.ModuleName,
		polygontypes.ModuleName,
		trontypes.ModuleName,
		ethtypes.ModuleName,
		avalanchetypes.ModuleName,
		arbitrumtypes.ModuleName,
		optimismtypes.ModuleName,
	}
	for _, moduleName := range subModules {
		methodFinder := reflect.TypeOf(new(KeeperTestSuite))
		for i := 0; i < methodFinder.NumMethod(); i++ {
			method := methodFinder.Method(i)
			if !compile.MatchString(method.Name) {
				continue
			}
			t.Run(method.Name, func(subT *testing.T) {
				mySuite := &KeeperTestSuite{chainName: moduleName}
				mySuite.SetT(subT)
				mySuite.SetupTest()
				method.Func.Call([]reflect.Value{reflect.ValueOf(mySuite)})
			})
		}
	}
}

func (suite *KeeperTestSuite) MsgServer() types.MsgServer {
	if suite.chainName == trontypes.ModuleName {
		return tronkeeper.NewMsgServerImpl(suite.app.TronKeeper)
	}
	return keeper.NewMsgServerImpl(suite.Keeper())
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
	suite.bridgerAddrs = helpers.AddTestAddrs(suite.app, suite.ctx, valNumber, sdk.NewCoins(sdk.NewCoin(types.NativeDenom, sdkmath.NewInt(300*1e3).MulRaw(1e18))))
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

func (suite *KeeperTestSuite) PubKeyToExternalAddr(publicKey ecdsa.PublicKey) string {
	if trontypes.ModuleName == suite.chainName {
		return tronaddress.PubkeyToAddress(publicKey).String()
	}
	return crypto.PubkeyToAddress(publicKey).Hex()
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
		require.NoError(suite.T(), err)

		signature, err = trontypes.NewTronSignature(checkpoint, external)
		require.NoError(suite.T(), err)
	}
	return externalAddress, signature
}
