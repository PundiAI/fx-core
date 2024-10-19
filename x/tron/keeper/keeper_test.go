package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	queryServer crosschaintypes.QueryClient
	msgServer   crosschaintypes.MsgServer

	signer *helpers.Signer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.MintValNumber = 1
	suite.BaseSuite.SetupTest()

	err := suite.App.TronKeeper.SetParams(suite.Ctx, &crosschaintypes.Params{
		GravityId:                         "fx-bridge-tron",
		AverageBlockTime:                  5000,
		ExternalBatchTimeout:              43200000,
		AverageExternalBlockTime:          3000,
		SignedWindow:                      20000,
		SlashFraction:                     sdkmath.LegacyNewDec(1).Quo(sdkmath.LegacyNewDec(1000)),
		OracleSetUpdatePowerChangePercent: sdkmath.LegacyNewDec(1).Quo(sdkmath.LegacyNewDec(10)),
		IbcTransferTimeoutHeight:          10000,
		DelegateThreshold: sdk.NewCoin(fxtypes.DefaultDenom,
			sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil))),
		DelegateMultiple:  10,
		BridgeCallTimeout: crosschaintypes.DefBridgeCallTimeout,
	})
	suite.Require().NoError(err)
	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	crosschaintypes.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(suite.App.TronKeeper))
	suite.queryServer = crosschaintypes.NewQueryClient(queryHelper)

	suite.msgServer = keeper.NewMsgServerImpl(suite.App.TronKeeper)
	suite.signer = helpers.NewSigner(helpers.NewEthPrivKey())
	suite.MintToken(suite.signer.AccAddress(), helpers.NewStakingCoin(1000, 18))
}

func (suite *KeeperTestSuite) NewOutgoingTxBatch() *crosschaintypes.OutgoingTxBatch {
	batchNonce := tmrand.Uint64()
	tokenContract := helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex())
	newOutgoingTx := &crosschaintypes.OutgoingTxBatch{
		BatchNonce: batchNonce,
		Transactions: []*crosschaintypes.OutgoingTransferTx{
			{
				Sender:      suite.signer.AccAddress().String(),
				DestAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
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
		FeeReceive:    helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
		Block:         batchNonce,
	}
	err := suite.App.TronKeeper.StoreBatch(suite.Ctx, newOutgoingTx)
	suite.Require().NoError(err)
	return newOutgoingTx
}

func (suite *KeeperTestSuite) NewOracleByBridger() (sdk.AccAddress, sdk.AccAddress, cryptotypes.PrivKey) {
	oracle := helpers.GenAccAddress()
	bridger := helpers.GenAccAddress()
	externalKey := helpers.NewEthPrivKey()
	externalAddress := helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String())
	newOracle := crosschaintypes.Oracle{
		OracleAddress:   oracle.String(),
		BridgerAddress:  bridger.String(),
		ExternalAddress: externalAddress,
	}
	suite.App.TronKeeper.SetOracle(suite.Ctx, newOracle)
	suite.App.TronKeeper.SetOracleAddrByBridgerAddr(suite.Ctx, bridger, oracle)
	oracleAddress, found := suite.App.TronKeeper.GetOracleAddrByBridgerAddr(suite.Ctx, bridger)
	suite.True(found)
	suite.EqualValues(oracle, oracleAddress)
	suite.App.TronKeeper.SetOracleAddrByExternalAddr(suite.Ctx, externalAddress, oracle)
	return oracle, bridger, externalKey
}

func (suite *KeeperTestSuite) NewOracleSet(externalKey cryptotypes.PrivKey) *crosschaintypes.OracleSet {
	newOracleSet := crosschaintypes.NewOracleSet(tmrand.Uint64(), tmrand.Uint64(), crosschaintypes.BridgeValidators{
		{
			Power:           tmrand.Uint64(),
			ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
		},
	})
	suite.App.TronKeeper.StoreOracleSet(suite.Ctx, newOracleSet)
	return newOracleSet
}
