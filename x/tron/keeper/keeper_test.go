package keeper_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	tronkeeper "github.com/functionx/fx-core/v7/x/tron/keeper"
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
	valNumber := tmrand.Intn(50) + 1

	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          suite.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})
	err := suite.app.TronKeeper.SetParams(suite.ctx, &crosschaintypes.Params{
		GravityId:                         "fx-bridge-tron",
		AverageBlockTime:                  5000,
		ExternalBatchTimeout:              43200000,
		AverageExternalBlockTime:          3000,
		SignedWindow:                      20000,
		SlashFraction:                     sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
		IbcTransferTimeoutHeight:          10000,
		DelegateThreshold: sdk.NewCoin(fxtypes.DefaultDenom,
			sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil))),
		DelegateMultiple:  10,
		BridgeCallTimeout: crosschaintypes.DefaultBridgeCallTimeout,
	})
	suite.Require().NoError(err)
	suite.msgServer = tronkeeper.NewMsgServerImpl(suite.app.TronKeeper)
	suite.signer = helpers.NewSigner(helpers.NewEthPrivKey())
	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18)))))
}

func (suite *KeeperTestSuite) NewOutgoingTxBatch() *crosschaintypes.OutgoingTxBatch {
	batchNonce := tmrand.Uint64()
	tokenContract := helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex())
	newOutgoingTx := &crosschaintypes.OutgoingTxBatch{
		BatchNonce: batchNonce,
		Transactions: []*crosschaintypes.OutgoingTransferTx{
			{
				Sender:      suite.signer.AccAddress().String(),
				DestAddress: helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
				Token: crosschaintypes.ERC20Token{
					Contract: tokenContract,
					Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
				},
				Fee: crosschaintypes.ERC20Token{
					Contract: tokenContract,
					Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
				},
			},
		},
		TokenContract: tokenContract,
		FeeReceive:    helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex()),
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
	externalAddress := helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String())
	newOracle := crosschaintypes.Oracle{
		OracleAddress:   oracle.String(),
		BridgerAddress:  bridger.String(),
		ExternalAddress: externalAddress,
	}
	suite.app.TronKeeper.SetOracle(suite.ctx, newOracle)
	suite.app.TronKeeper.SetOracleAddrByBridgerAddr(suite.ctx, bridger, oracle)
	oracleAddress, found := suite.app.TronKeeper.GetOracleAddrByBridgerAddr(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle, oracleAddress)
	suite.app.TronKeeper.SetOracleAddrByExternalAddr(suite.ctx, externalAddress, oracle)
	return oracle, bridger, externalKey
}

func (suite *KeeperTestSuite) NewOracleSet(externalKey cryptotypes.PrivKey) *crosschaintypes.OracleSet {
	newOracleSet := crosschaintypes.NewOracleSet(tmrand.Uint64(), tmrand.Uint64(), crosschaintypes.BridgeValidators{
		{
			Power:           tmrand.Uint64(),
			ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
		},
	})
	suite.app.TronKeeper.StoreOracleSet(suite.ctx, newOracleSet)
	return newOracleSet
}

func (suite *KeeperTestSuite) NewBridgeToken(bridger sdk.AccAddress) []crosschaintypes.BridgeToken {
	bridgeTokens := make([]crosschaintypes.BridgeToken, 0)
	for i := 0; i < 3; i++ {
		bridgeTokens = append(bridgeTokens, crosschaintypes.BridgeToken{Token: helpers.HexAddrToTronAddr(helpers.GenerateAddress().Hex())})
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
		bridgeDenom := sdk.NewCoins(sdk.NewCoin(bridgeTokens[i].Denom, sdkmath.NewInt(1e6).MulRaw(1e18)))
		err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, bridgeDenom)
		suite.NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, suite.signer.AccAddress(), bridgeDenom)
		suite.NoError(err)
	}
	return bridgeTokens
}
