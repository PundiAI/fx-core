package keeper_test

import (
	"encoding/hex"
	"math/big"
	"math/rand"
	"testing"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	tronkeeper "github.com/functionx/fx-core/v3/x/tron/keeper"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *app.App
	ctx       sdk.Context
	msgServer crosschaintypes.MsgServer

	signer *helpers.Signer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	valNumber := rand.Intn(100-1) + 1
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          suite.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})
	suite.app.TronKeeper.SetParams(suite.ctx, &crosschaintypes.Params{
		GravityId:                         "fx-bridge-tron",
		AverageBlockTime:                  5000,
		ExternalBatchTimeout:              43200000,
		AverageExternalBlockTime:          3000,
		SignedWindow:                      20000,
		SlashFraction:                     sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
		IbcTransferTimeoutHeight:          10000,
		DelegateThreshold: sdk.NewCoin(fxtypes.DefaultDenom,
			sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil))),
		DelegateMultiple: 10,
	})
	suite.msgServer = tronkeeper.NewMsgServerImpl(suite.app.TronKeeper)
	suite.signer = helpers.NewSigner(helpers.NewEthPrivKey())
	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))))
}

func (suite *KeeperTestSuite) NewOutgoingTxBatch() *crosschaintypes.OutgoingTxBatch {
	batchNonce := rand.Uint64()
	tokenContract := trontypes.AddressFromHex(helpers.GenerateAddress().Hex())
	newOutgoingTx := &crosschaintypes.OutgoingTxBatch{
		BatchNonce: batchNonce,
		Transactions: []*crosschaintypes.OutgoingTransferTx{
			{
				Sender:      suite.signer.AccAddress().String(),
				DestAddress: trontypes.AddressFromHex(helpers.GenerateAddress().Hex()),
				Token: crosschaintypes.ERC20Token{
					Contract: tokenContract,
					Amount:   sdk.NewIntFromBigInt(big.NewInt(1e18)),
				},
				Fee: crosschaintypes.ERC20Token{
					Contract: tokenContract,
					Amount:   sdk.NewIntFromBigInt(big.NewInt(1e18)),
				},
			},
		},
		TokenContract: tokenContract,
		FeeReceive:    trontypes.AddressFromHex(helpers.GenerateAddress().Hex()),
		Block:         batchNonce,
	}
	err := suite.app.TronKeeper.StoreBatch(suite.ctx, newOutgoingTx)
	suite.Require().NoError(err)
	return newOutgoingTx
}

func (suite *KeeperTestSuite) NewOracleByBridger() (sdk.AccAddress, sdk.AccAddress, cryptotypes.PrivKey) {
	oracle := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	bridger := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	externalKey := helpers.NewEthPrivKey()
	externalAddress := trontypes.AddressFromHex(externalKey.PubKey().Address().String())
	newOracle := crosschaintypes.Oracle{
		OracleAddress:   oracle.String(),
		BridgerAddress:  bridger.String(),
		ExternalAddress: externalAddress,
	}
	suite.app.TronKeeper.SetOracle(suite.ctx, newOracle)
	suite.app.TronKeeper.SetOracleByBridger(suite.ctx, bridger, oracle)
	oracleAddress, found := suite.app.TronKeeper.GetOracleAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle, oracleAddress)
	suite.app.TronKeeper.SetOracleByExternalAddress(suite.ctx, externalAddress, oracle)
	return oracle, bridger, externalKey
}

func (suite *KeeperTestSuite) NewOracleSet(externalKey cryptotypes.PrivKey) *crosschaintypes.OracleSet {
	newOracleSet := crosschaintypes.NewOracleSet(rand.Uint64(), rand.Uint64(), crosschaintypes.BridgeValidators{
		{
			Power:           rand.Uint64(),
			ExternalAddress: trontypes.AddressFromHex(externalKey.PubKey().Address().String()),
		},
	})
	suite.app.TronKeeper.StoreOracleSet(suite.ctx, newOracleSet)
	return newOracleSet
}

func (suite *KeeperTestSuite) NewBridgeToken(bridger sdk.AccAddress) []crosschaintypes.BridgeToken {
	bridgeTokens := make([]crosschaintypes.BridgeToken, 0)
	for i := 0; i < 3; i++ {
		bridgeTokens = append(bridgeTokens, crosschaintypes.BridgeToken{Token: trontypes.AddressFromHex(helpers.GenerateAddress().Hex())})
		channelIBC := ""
		if i == 2 {
			channelIBC = hex.EncodeToString([]byte("transfer/channel-0"))
		}
		err := suite.app.TronKeeper.AttestationHandler(suite.ctx, &crosschaintypes.MsgBridgeTokenClaim{
			TokenContract:  bridgeTokens[i].Token,
			BridgerAddress: bridger.String(),
			ChannelIbc:     channelIBC,
		})
		suite.Require().NoError(err)
		denom := suite.app.TronKeeper.GetBridgeTokenDenom(suite.ctx, bridgeTokens[i].Token)
		bridgeTokens[i].Denom = denom.Denom
		bridgeDenom := sdk.NewCoins(sdk.NewCoin(bridgeTokens[i].Denom, sdk.NewInt(1e6).MulRaw(1e18)))
		err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, bridgeDenom)
		suite.NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, suite.signer.AccAddress(), bridgeDenom)
		suite.NoError(err)
	}
	return bridgeTokens
}
