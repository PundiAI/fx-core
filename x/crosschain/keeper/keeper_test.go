package keeper_test

import (
	"crypto/ecdsa"
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/crypto"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	trontypes "github.com/pundiai/fx-core/v8/x/tron/types"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	oracleAddrs  []sdk.AccAddress
	bridgerAddrs []sdk.AccAddress
	externalPris []*ecdsa.PrivateKey
	chainName    string

	signer *helpers.Signer

	bridgeFeeSuite  helpers.BridgeFeeSuite
	erc20TokenSuite helpers.ERC20TokenSuite
}

func TestCrosschainKeeperTestSuite(t *testing.T) {
	modules := []string{
		trontypes.ModuleName,
		ethtypes.ModuleName,
	}
	if os.Getenv("TEST_CROSSCHAIN") == "true" {
		modules = fxtypes.GetSupportChains()
	}
	for _, moduleName := range modules {
		suite.Run(t, &KeeperTestSuite{chainName: moduleName})
	}
}

func (suite *KeeperTestSuite) MsgServer() types.MsgServer {
	return keeper.NewMsgServerImpl(suite.Keeper())
}

func (suite *KeeperTestSuite) QueryClient() types.QueryClient {
	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(suite.Keeper()))
	return types.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) Keeper() keeper.Keeper {
	return suite.App.CrosschainKeepers.GetKeeper(suite.chainName)
}

func (suite *KeeperTestSuite) SetupTest() {
	valNumber := tmrand.Intn(types.MaxOracleSize-4) + 4
	suite.MintValNumber = valNumber
	suite.BaseSuite.SetupTest()

	suite.oracleAddrs = suite.AddTestAddress(valNumber, types.NewDelegateAmount(sdkmath.NewInt(300*1e3).MulRaw(1e18)))
	suite.bridgerAddrs = suite.AddTestAddress(valNumber, sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(300*1e3).MulRaw(1e18)))
	suite.externalPris = helpers.CreateMultiECDSA(valNumber)

	suite.signer = suite.AddTestSigner(10_000)

	proposalOracle := &types.ProposalOracle{}
	for _, oracle := range suite.oracleAddrs {
		proposalOracle.Oracles = append(proposalOracle.Oracles, oracle.String())
	}
	suite.Keeper().SetProposalOracle(suite.Ctx, proposalOracle)

	suite.Keeper().SetLastObservedBlockHeight(suite.Ctx, 100, 10)

	suite.bridgeFeeSuite = helpers.NewBridgeFeeSuite(suite.Require(), suite.App.EvmKeeper)
	suite.erc20TokenSuite = helpers.NewERC20Suite(suite.Require(), suite.App.EvmKeeper)
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) ModuleAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(suite.chainName)
}

func (suite *KeeperTestSuite) PubKeyToExternalAddr(publicKey ecdsa.PublicKey) string {
	address := crypto.PubkeyToAddress(publicKey)
	return fxtypes.ExternalAddrToStr(suite.chainName, address.Bytes())
}

func (suite *KeeperTestSuite) SignOracleSetConfirm(external *ecdsa.PrivateKey, oracleSet *types.OracleSet) (string, []byte) {
	externalAddress := crypto.PubkeyToAddress(external.PublicKey).String()
	gravityId := suite.Keeper().GetGravityID(suite.Ctx)
	checkpoint, err := oracleSet.GetCheckpoint(gravityId)
	suite.Require().NoError(err)
	signature, err := ethtypes.NewEthereumSignature(checkpoint, external)
	suite.Require().NoError(err)
	if trontypes.ModuleName == suite.chainName {
		externalAddress = tronaddress.PubkeyToAddress(external.PublicKey).String()

		signature, err = trontypes.NewTronSignature(checkpoint, external)
		suite.Require().NoError(err)
	}
	return externalAddress, signature
}

func (suite *KeeperTestSuite) SendClaim(externalClaim types.ExternalClaim) {
	err := suite.SendClaimReturnErr(externalClaim)
	suite.Require().NoError(err)

	preErr, executeErr := suite.Keeper().ExecuteClaim(suite.Ctx, suite.App.EvmKeeper, externalClaim.GetEventNonce())
	suite.Require().NoError(preErr)
	suite.Require().NoError(executeErr)
}

func (suite *KeeperTestSuite) SendClaimReturnErr(externalClaim types.ExternalClaim) error {
	value, err := codectypes.NewAnyWithValue(externalClaim)
	suite.Require().NoError(err)
	_, err = suite.MsgServer().Claim(suite.Ctx, &types.MsgClaim{Claim: value})
	return err
}

func (suite *KeeperTestSuite) EndBlocker() {
	_, err := suite.App.EndBlocker(suite.Ctx)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) SetIBCDenom(portID, channelID, denom string) ibctransfertypes.DenomTrace {
	sourcePrefix := ibctransfertypes.GetDenomPrefix(portID, channelID)
	prefixedDenom := sourcePrefix + denom
	denomTrace := ibctransfertypes.ParseDenomTrace(prefixedDenom)
	traceHash := denomTrace.Hash()
	if !suite.App.IBCTransferKeeper.HasDenomTrace(suite.Ctx, traceHash) {
		suite.App.IBCTransferKeeper.SetDenomTrace(suite.Ctx, denomTrace)
	}
	return denomTrace
}

func (suite *KeeperTestSuite) GetERC20TokenByBridgeContract(bridgeContract string) *erc20types.ERC20Token {
	bridgeToken, err := suite.Keeper().GetBridgeToken(suite.Ctx, bridgeContract)
	suite.Require().NoError(err)
	erc20Token, err := suite.App.Erc20Keeper.GetERC20Token(suite.Ctx, bridgeToken.Denom)
	suite.Require().NoError(err)
	return &erc20Token
}

func (suite *KeeperTestSuite) BondOracles() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.ValPrivs[i].ValAddress().String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(100).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		suite.Require().NoError(msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(suite.Ctx, msgBondedOracle)
		suite.Require().NoError(err)
	}
}

func (suite *KeeperTestSuite) SetAutoIncrementID(key []byte, id uint64) {
	store := suite.Ctx.KVStore(suite.App.AppKeepers.GetKey(suite.Keeper().ModuleName()))
	store.Set(key, sdk.Uint64ToBigEndian(id))
}
