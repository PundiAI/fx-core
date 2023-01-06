package keeper_test

import (
	"crypto/ecdsa"
	"math/rand"
	"reflect"
	"regexp"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	"github.com/functionx/fx-core/v3/x/crosschain/keeper"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	tronkeeper "github.com/functionx/fx-core/v3/x/tron/keeper"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *app.App
	ctx       sdk.Context
	oracles   []sdk.AccAddress
	bridgers  []sdk.AccAddress
	externals []*ecdsa.PrivateKey
	validator []sdk.ValAddress
	chainName string
}

func TestKeeperTestSuite(t *testing.T) {
	compile, err := regexp.Compile("^Test")
	require.NoError(t, err)
	subModules := []string{
		bsctypes.ModuleName, polygontypes.ModuleName, trontypes.ModuleName,
		ethtypes.ModuleName, avalanchetypes.ModuleName,
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
	default:
		panic("invalid chain name")
	}
}

func (suite *KeeperTestSuite) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	valNumber := rand.Intn(types.MaxOracleSize-4) + 4

	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID: fxtypes.MainnetChainId,
		Height:  suite.app.LastBlockHeight() + 1,
	})

	suite.oracles = helpers.AddTestAddrs(suite.app, suite.ctx, valNumber, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(300*1e3).MulRaw(1e18))))
	suite.bridgers = helpers.AddTestAddrs(suite.app, suite.ctx, valNumber, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(300*1e3).MulRaw(1e18))))
	suite.externals = helpers.CreateMultiECDSA(valNumber)

	suite.validator = make([]sdk.ValAddress, valNumber)
	for i := 0; i < valNumber; i++ {
		suite.validator[i] = valAccounts[i].GetAddress().Bytes()
	}

	proposalOracle := &types.ProposalOracle{}
	for _, oracle := range suite.oracles {
		proposalOracle.Oracles = append(proposalOracle.Oracles, oracle.String())
	}
	suite.Keeper().SetProposalOracle(suite.ctx, proposalOracle)
}

func (suite *KeeperTestSuite) PubKeyToExternalAddr(publicKey ecdsa.PublicKey) string {
	if trontypes.ModuleName == suite.chainName {
		return tronAddress.PubkeyToAddress(publicKey).String()
	}
	return crypto.PubkeyToAddress(publicKey).Hex()
}
