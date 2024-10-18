package precompile_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	localhost "github.com/cosmos/ibc-go/v8/modules/light-clients/09-localhost"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
)

type PrecompileTestSuite struct {
	helpers.BaseSuite

	signer     *helpers.Signer
	crosschain common.Address
}

func TestPrecompileTestSuite(t *testing.T) {
	suite.Run(t, new(PrecompileTestSuite))
}

func (suite *PrecompileTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()
	suite.Ctx = suite.Ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.Ctx = suite.Ctx.WithBlockGasMeter(storetypes.NewGasMeter(1e18))

	suite.signer = suite.AddTestSigner(10_000)

	crosschainContract, err := suite.App.EvmKeeper.DeployContract(suite.Ctx, suite.signer.Address(), contract.MustABIJson(testscontract.CrosschainTestMetaData.ABI), contract.MustDecodeHex(testscontract.CrosschainTestMetaData.Bin))
	suite.Require().NoError(err)
	suite.crosschain = crosschainContract
}

func (suite *PrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *PrecompileTestSuite) EthereumTx(signer *helpers.Signer, to common.Address, amount *big.Int, data []byte) *evmtypes.MsgEthereumTxResponse {
	ethTx := evmtypes.NewTx(
		fxtypes.EIP155ChainID(suite.Ctx.ChainID()),
		suite.App.EvmKeeper.GetNonce(suite.Ctx, signer.Address()),
		&to,
		amount,
		contract.DefaultGasCap,
		nil,
		nil,
		nil,
		data,
		nil,
	)
	ethTx.From = signer.Address().Bytes()
	err := ethTx.Sign(ethtypes.LatestSignerForChainID(fxtypes.EIP155ChainID(suite.Ctx.ChainID())), signer)
	suite.Require().NoError(err)

	res, err := suite.App.EvmKeeper.EthereumTx(suite.Ctx, ethTx)
	suite.Require().NoError(err)
	return res
}

func (suite *PrecompileTestSuite) Commit() {
	header := suite.Ctx.BlockHeader()
	_, err := suite.App.EndBlocker(suite.Ctx)
	suite.Require().NoError(err)
	_, err = suite.App.Commit()
	suite.Require().NoError(err)
	// after commit Ctx header
	header.Height += 1

	// begin block
	header.Time = time.Now().UTC()
	header.Height += 1
	suite.Ctx = suite.Ctx.WithBlockHeader(header)
	_, err = suite.App.BeginBlocker(suite.Ctx)
	suite.Require().NoError(err)
	suite.Ctx = suite.Ctx.WithBlockHeight(header.Height)
}

func (suite *PrecompileTestSuite) RandSigner() *helpers.Signer {
	privKey := helpers.NewEthPrivKey()
	// suite.MintToken(privKey.PubKey().Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18)))))
	signer := helpers.NewSigner(privKey)
	suite.App.AccountKeeper.SetAccount(suite.Ctx, suite.App.AccountKeeper.NewAccountWithAddress(suite.Ctx, signer.AccAddress()))
	return signer
}

func (suite *PrecompileTestSuite) CrosschainKeepers() map[string]crosschainkeeper.Keeper {
	value := reflect.ValueOf(suite.App.CrosschainKeepers)
	keepers := make(map[string]crosschainkeeper.Keeper)
	for i := 0; i < value.NumField(); i++ {
		res := value.Field(i).MethodByName("GetGravityID").Call([]reflect.Value{reflect.ValueOf(suite.Ctx)})
		gravityID := res[0].String()
		chainName := strings.TrimSuffix(strings.TrimPrefix(gravityID, "fx-"), "-bridge")
		if chainName == "bridge-eth" {
			keepers["eth"] = value.Field(i).Interface().(crosschainkeeper.Keeper)
		} else {
			keepers[chainName] = value.Field(i).Interface().(crosschainkeeper.Keeper)
		}
	}
	return keepers
}

func (suite *PrecompileTestSuite) GenerateModuleName() string {
	keepers := suite.CrosschainKeepers()
	modules := make([]string, 0, len(keepers))
	for m := range keepers {
		modules = append(modules, m)
	}
	if len(modules) == 0 {
		return ""
	}
	return modules[tmrand.Intn(len(modules))]
}

func (suite *PrecompileTestSuite) GenerateOracles(moduleName string, online bool, num int) []Oracle {
	keeper := suite.CrosschainKeepers()[moduleName]
	oracles := make([]Oracle, 0, num)
	for i := 0; i < num; i++ {
		oracle := crosschaintypes.Oracle{
			OracleAddress:     helpers.GenAccAddress().String(),
			BridgerAddress:    helpers.GenAccAddress().String(),
			ExternalAddress:   helpers.GenExternalAddr(moduleName),
			DelegateAmount:    sdkmath.NewInt(1e18).MulRaw(1000),
			StartHeight:       1,
			Online:            online,
			DelegateValidator: sdk.ValAddress(helpers.GenAccAddress()).String(),
			SlashTimes:        0,
		}
		keeper.SetOracle(suite.Ctx, oracle)
		keeper.SetOracleAddrByExternalAddr(suite.Ctx, oracle.ExternalAddress, oracle.GetOracle())
		keeper.SetOracleAddrByBridgerAddr(suite.Ctx, oracle.GetBridger(), oracle.GetOracle())
		oracles = append(oracles, Oracle{
			moduleName: moduleName,
			oracle:     oracle,
		})
	}
	return oracles
}

func (suite *PrecompileTestSuite) GenerateRandOracle(moduleName string, online bool) Oracle {
	oracles := suite.GenerateOracles(moduleName, online, 1)
	return oracles[0]
}

func (suite *PrecompileTestSuite) RandTransferChannel() (portID, channelID string) {
	portID = "transfer"

	channelSequence := suite.App.IBCKeeper.ChannelKeeper.GetNextChannelSequence(suite.Ctx)
	channelID = fmt.Sprintf("channel-%d", channelSequence)
	connectionID := connectiontypes.FormatConnectionIdentifier(uint64(tmrand.Intn(100)))
	clientID := clienttypes.FormatClientIdentifier(exported.Localhost, uint64(tmrand.Intn(100)))

	revision := clienttypes.ParseChainID(suite.Ctx.ChainID())
	localHostClient := localhost.NewClientState(clienttypes.NewHeight(revision, uint64(suite.Ctx.BlockHeight())))
	suite.App.IBCKeeper.ClientKeeper.SetClientState(suite.Ctx, clientID, localHostClient)

	params := suite.App.IBCKeeper.ClientKeeper.GetParams(suite.Ctx)
	params.AllowedClients = append(params.AllowedClients, localHostClient.ClientType())
	suite.App.IBCKeeper.ClientKeeper.SetParams(suite.Ctx, params)

	prevConsState := &ibctm.ConsensusState{
		Timestamp:          suite.Ctx.BlockTime(),
		NextValidatorsHash: suite.Ctx.BlockHeader().NextValidatorsHash,
	}
	height := clienttypes.NewHeight(0, uint64(suite.Ctx.BlockHeight()))
	suite.App.IBCKeeper.ClientKeeper.SetClientConsensusState(suite.Ctx, clientID, height, prevConsState)

	channelCapability, err := suite.App.ScopedIBCKeeper.NewCapability(suite.Ctx, host.ChannelCapabilityPath(portID, channelID))
	suite.Require().NoError(err)
	err = suite.App.ScopedTransferKeeper.ClaimCapability(suite.Ctx, capabilitytypes.NewCapability(channelCapability.Index), host.ChannelCapabilityPath(portID, channelID))
	suite.Require().NoError(err)

	connectionEnd := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, clientID, connectiontypes.Counterparty{ClientId: "clientId", ConnectionId: "connection-1", Prefix: commitmenttypes.NewMerklePrefix([]byte("prefix"))}, connectiontypes.GetCompatibleVersions(), 500)
	suite.App.IBCKeeper.ConnectionKeeper.SetConnection(suite.Ctx, connectionID, connectionEnd)

	channel := channeltypes.NewChannel(channeltypes.OPEN, channeltypes.ORDERED, channeltypes.NewCounterparty(portID, channelID), []string{connectionID}, "mock-version")
	suite.App.IBCKeeper.ChannelKeeper.SetChannel(suite.Ctx, portID, channelID, channel)
	suite.App.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.Ctx, portID, channelID, uint64(tmrand.Intn(10000)+1))
	suite.App.IBCKeeper.ChannelKeeper.SetNextChannelSequence(suite.Ctx, channelSequence+1)
	return portID, channelID
}

func (suite *PrecompileTestSuite) SendEvmTx(signer *helpers.Signer, contractAddr common.Address, data []byte) *evmtypes.MsgEthereumTxResponse {
	from := signer.Address()

	args, err := json.Marshal(&evmtypes.TransactionArgs{To: &contractAddr, From: &from, Data: (*hexutil.Bytes)(&data)})
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelper, suite.App.EvmKeeper)
	res, err := evmtypes.NewQueryClient(queryHelper).EstimateGas(suite.Ctx,
		&evmtypes.EthCallRequest{
			Args:    args,
			GasCap:  contract.DefaultGasCap,
			ChainId: suite.App.EvmKeeper.ChainID().Int64(),
		},
	)
	suite.Require().NoError(err)

	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	// suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx).Int64()*int64(res.Gas)))))

	msg := core.Message{
		To:                &contractAddr,
		From:              signer.Address(),
		Nonce:             suite.App.EvmKeeper.GetNonce(suite.Ctx, signer.Address()),
		Value:             big.NewInt(0),
		GasLimit:          res.Gas,
		GasPrice:          suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx),
		GasFeeCap:         nil,
		GasTipCap:         nil,
		Data:              data,
		AccessList:        nil,
		SkipAccountChecks: false,
	}

	rsp, err := suite.App.EvmKeeper.ApplyMessage(suite.Ctx, &msg, nil, true)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp.VmError)
	return rsp
}

type Oracle struct {
	moduleName string
	oracle     crosschaintypes.Oracle
}

func (o Oracle) GetExternalHexAddr() common.Address {
	return crosschaintypes.ExternalAddrToHexAddr(o.moduleName, o.oracle.ExternalAddress)
}
