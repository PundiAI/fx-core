package keeper_test

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/helpers"
	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/crosschain/types"
	"github.com/functionx/fx-core/x/tron/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.App
	ctx sdk.Context

	bridgeAcc    sdk.AccAddress
	bridgeTokens []BridgeToken

	oracleAddressList       []sdk.AccAddress
	orchestratorAddressList []sdk.AccAddress
	externalAccList         []*ExternalAcc

	msgServer   types.MsgServer
	queryClient types.QueryClient
}

type BridgeToken struct {
	token string
	denom string
}

type ExternalAcc struct {
	key     *ecdsa.PrivateKey
	address string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) DoSetupTest() {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := helpers.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initBalances)))

	suite.bridgeAcc = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	suite.bridgeTokens = make([]BridgeToken, 0)
	suite.externalAccList = make([]*ExternalAcc, 0)

	for i := 0; i < 3; i++ {
		suite.bridgeTokens = append(suite.bridgeTokens, BridgeToken{token: GenTronContractAddress()})
	}

	balances = append(balances, banktypes.Balance{
		Address: suite.bridgeAcc.String(),
		Coins: sdk.NewCoins(
			sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(20000)))),
		),
	})

	suite.app = helpers.SetupWithGenesisValSet(suite.T(), validator, genesisAccounts, balances...)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

	suite.oracleAddressList = helpers.AddTestAddrsIncremental(suite.app, suite.ctx, 4, sdk.ZeroInt())
	suite.orchestratorAddressList = helpers.AddTestAddrsIncremental(suite.app, suite.ctx, 4, sdk.ZeroInt())

	suite.app.TronKeeper.SetParams(suite.ctx, &types.Params{
		GravityId:                         "tron",
		AverageBlockTime:                  5000,
		ExternalBatchTimeout:              43200000,
		AverageExternalBlockTime:          3000,
		SignedWindow:                      20000,
		SlashFraction:                     sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
		IbcTransferTimeoutHeight:          10000,
		DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil))),
		DelegateMultiple:                  10,
	})

	for index, oracle := range suite.oracleAddressList {
		dbOracleAddress, found := suite.app.TronKeeper.GetOracleAddressByBridgerKey(suite.ctx, suite.orchestratorAddressList[index])
		require.False(suite.T(), found)
		require.Empty(suite.T(), dbOracleAddress)

		address, key := GenTronAccountAddress()
		suite.externalAccList = append(suite.externalAccList, &ExternalAcc{
			key:     key,
			address: address,
		})

		newOracle := types.Oracle{
			OracleAddress:   oracle.String(),
			BridgerAddress:  suite.orchestratorAddressList[index].String(),
			ExternalAddress: suite.externalAccList[index].address,
			StartHeight:     3,
		}

		if index == 0 {
			suite.app.TronKeeper.SetOracle(suite.ctx, newOracle)
		}
		suite.app.TronKeeper.SetOracleByBridger(suite.ctx, oracle, suite.orchestratorAddressList[index])

		dbOracleAddress, found = suite.app.TronKeeper.GetOracleAddressByBridgerKey(suite.ctx, suite.orchestratorAddressList[index])
		require.True(suite.T(), found)
		require.EqualValues(suite.T(), oracle, dbOracleAddress)

		suite.app.TronKeeper.SetExternalAddressForOracle(suite.ctx, oracle, suite.externalAccList[index].address)
	}

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.TronKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	for index, token := range suite.bridgeTokens {
		channalIbc := ""
		if index == 2 {
			channalIbc = hex.EncodeToString([]byte("transfer/channel-0"))
		}
		err := suite.app.TronKeeper.AttestationHandler(suite.ctx, types.Attestation{}, &types.MsgBridgeTokenClaim{
			TokenContract:  token.token,
			BridgerAddress: suite.orchestratorAddressList[0].String(),
			ChannelIbc:     channalIbc,
			ChainName:      "tron",
		})
		suite.Require().NoError(err)
		denom := suite.app.TronKeeper.GetBridgeTokenDenom(suite.ctx, token.token)
		suite.bridgeTokens[index].denom = denom.Denom
		err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(denom.Denom, initBalances)))
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, suite.bridgeAcc, sdk.NewCoins(sdk.NewCoin(denom.Denom, initBalances)))
		suite.Require().NoError(err)
	}

	suite.msgServer = keeper.NewMsgServerImpl(suite.app.TronKeeper)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest()
}

func GenTronContractAddress() string {
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey)
	rand.Seed(time.Now().Unix())
	contract := crypto.CreateAddress(address, uint64(rand.Intn(100))).Bytes()
	contract = append([]byte{byte(0x41)}, contract...)

	return EncodeCheck(contract)
}

func GenTronAccountAddress() (string, *ecdsa.PrivateKey) {
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey)
	addressBytes := append([]byte{byte(0x41)}, address.Bytes()...)

	return EncodeCheck(addressBytes), key
}

func EncodeCheck(input []byte) string {
	h256h0 := sha256.New()
	h256h0.Write(input)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	inputCheck := input
	inputCheck = append(inputCheck, h1[:4]...)

	return base58.Encode(inputCheck)
}
