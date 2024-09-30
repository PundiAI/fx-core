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
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	localhost "github.com/cosmos/ibc-go/v8/modules/light-clients/09-localhost"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	"github.com/functionx/fx-core/v8/x/erc20/types"
	crossethtypes "github.com/functionx/fx-core/v8/x/eth/types"
)

type PrecompileTestSuite struct {
	suite.Suite
	ctx        sdk.Context
	app        *app.App
	signer     *helpers.Signer
	crosschain common.Address
}

func TestPrecompileTestSuite(t *testing.T) {
	suite.Run(t, new(PrecompileTestSuite))
}

func (suite *PrecompileTestSuite) SetupTest() {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(suite.T(), err)
	suite.signer = helpers.NewSigner(priv)

	set, accs, balances := helpers.GenerateGenesisValidator(tmrand.Intn(10)+1, nil)
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), set, accs, balances...)

	suite.ctx = suite.app.NewContext(false)
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(storetypes.NewGasMeter(1e18))

	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18)))))

	suite.AddFXBridgeToken(helpers.GenExternalAddr(crossethtypes.ModuleName))

	crosschainContract, err := suite.app.EvmKeeper.DeployContract(suite.ctx, suite.signer.Address(), contract.MustABIJson(testscontract.CrossChainTestMetaData.ABI), contract.MustDecodeHex(testscontract.CrossChainTestMetaData.Bin))
	suite.Require().NoError(err)
	suite.crosschain = crosschainContract
}

func (suite *PrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *PrecompileTestSuite) EthereumTx(signer *helpers.Signer, to common.Address, amount *big.Int, data []byte) *evmtypes.MsgEthereumTxResponse {
	ethTx := evmtypes.NewTx(
		fxtypes.EIP155ChainID(suite.ctx.ChainID()),
		suite.app.EvmKeeper.GetNonce(suite.ctx, signer.Address()),
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
	err := ethTx.Sign(ethtypes.LatestSignerForChainID(fxtypes.EIP155ChainID(suite.ctx.ChainID())), signer)
	suite.Require().NoError(err)

	res, err := suite.app.EvmKeeper.EthereumTx(suite.ctx, ethTx)
	suite.Require().NoError(err)
	return res
}

func (suite *PrecompileTestSuite) Commit() {
	header := suite.ctx.BlockHeader()
	_, err := suite.app.EndBlocker(suite.ctx)
	suite.Require().NoError(err)
	_, err = suite.app.Commit()
	suite.Require().NoError(err)
	// after commit ctx header
	header.Height += 1

	// begin block
	header.Time = time.Now().UTC()
	header.Height += 1
	suite.ctx = suite.ctx.WithBlockHeader(header)
	_, err = suite.app.BeginBlocker(suite.ctx)
	suite.Require().NoError(err)
	suite.ctx = suite.ctx.WithBlockHeight(header.Height)
}

func (suite *PrecompileTestSuite) RandSigner() *helpers.Signer {
	privKey := helpers.NewEthPrivKey()
	// helpers.AddTestAddr(suite.app, suite.ctx, privKey.PubKey().Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18)))))
	signer := helpers.NewSigner(privKey)
	suite.app.AccountKeeper.SetAccount(suite.ctx, suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, signer.AccAddress()))
	return signer
}

func (suite *PrecompileTestSuite) Error(res *evmtypes.MsgEthereumTxResponse, errResult error) {
	suite.Require().True(res.Failed())
	if res.VmError != vm.ErrExecutionReverted.Error() {
		suite.Require().Equal(errResult.Error(), res.VmError)
		return
	}

	if len(res.Ret) > 0 {
		reason, err := abi.UnpackRevert(common.CopyBytes(res.Ret))
		suite.Require().NoError(err)

		suite.Require().Equal(errResult.Error(), reason)
		return
	}

	suite.Require().Equal(errResult.Error(), vm.ErrExecutionReverted.Error())
}

func (suite *PrecompileTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) DeployContract(from common.Address) (common.Address, error) {
	contractAddr, err := suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, suite.app.Erc20Keeper.ModuleAddress(), "Test token", "TEST", 18)
	suite.Require().NoError(err)

	_, err = suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.app.Erc20Keeper.ModuleAddress(), contractAddr, nil, contract.GetFIP20().ABI, "transferOwnership", from)
	suite.Require().NoError(err)
	return contractAddr, nil
}

func (suite *PrecompileTestSuite) DeployFXRelayToken() (types.TokenPair, banktypes.Metadata) {
	fxToken := fxtypes.GetFXMetaData()

	pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, fxToken)
	suite.Require().NoError(err)
	return *pair, fxToken
}

func (suite *PrecompileTestSuite) CrossChainKeepers() map[string]crosschainkeeper.Keeper {
	value := reflect.ValueOf(suite.app.CrossChainKeepers)
	keepers := make(map[string]crosschainkeeper.Keeper)
	for i := 0; i < value.NumField(); i++ {
		res := value.Field(i).MethodByName("GetGravityID").Call([]reflect.Value{reflect.ValueOf(suite.ctx)})
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

func (suite *PrecompileTestSuite) SetCorsschainEnablePending(moduleName string, isEnable bool) {
	crosschainKeeper := suite.CrossChainKeepers()[moduleName]
	params := crosschainKeeper.GetParams(suite.ctx)
	params.EnableBridgeCallPending = isEnable
	params.EnableSendToExternalPending = isEnable
	suite.NoError(crosschainKeeper.SetParams(suite.ctx, &params))
}

func (suite *PrecompileTestSuite) GenerateCrossChainDenoms(addDenoms ...string) Metadata {
	keepers := suite.CrossChainKeepers()
	modules := make([]string, 0, len(keepers))
	for m := range keepers {
		modules = append(modules, m)
	}
	count := tmrand.Intn(len(modules)-1) + 1

	denoms := make([]string, len(modules))
	denomModules := make([]string, len(modules))
	for index, m := range modules {
		address := helpers.GenExternalAddr(m)

		denom := crosschaintypes.NewBridgeDenom(m, address)
		denoms[index] = denom
		denomModules[index] = m

		k := keepers[m]
		k.AddBridgeToken(suite.ctx, denom, denom)
	}
	if count >= len(modules) {
		count = len(modules) - 1
	}
	metadata := fxtypes.GetCrossChainMetadataManyToOne("Test Token", helpers.NewRandSymbol(), 18, append(denoms[:count], addDenoms...)...)
	return Metadata{metadata: metadata, modules: denomModules[:count], notModules: denomModules[count:]}
}

func (suite *PrecompileTestSuite) GenerateModuleName() string {
	keepers := suite.CrossChainKeepers()
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
	keeper := suite.CrossChainKeepers()[moduleName]
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
		keeper.SetOracle(suite.ctx, oracle)
		keeper.SetOracleAddrByExternalAddr(suite.ctx, oracle.ExternalAddress, oracle.GetOracle())
		keeper.SetOracleAddrByBridgerAddr(suite.ctx, oracle.GetBridger(), oracle.GetOracle())
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

func (suite *PrecompileTestSuite) InitObservedBlockHeight() {
	keepers := suite.CrossChainKeepers()
	for _, k := range keepers {
		k.SetLastObservedBlockHeight(suite.ctx, 10, uint64(suite.ctx.BlockHeight()))
	}
}

func (suite *PrecompileTestSuite) MintLockNativeTokenToModule(md banktypes.Metadata, amt sdkmath.Int) sdk.Coin {
	generateAddress := helpers.GenHexAddress()

	count := 1
	if len(md.DenomUnits) > 0 && len(md.DenomUnits[0].Aliases) > 0 {
		// add alias to erc20 module
		for _, alias := range md.DenomUnits[0].Aliases {
			// add alias for erc20 module
			coins := sdk.NewCoins(sdk.NewCoin(alias, amt))
			helpers.AddTestAddr(suite.app, suite.ctx, generateAddress.Bytes(), coins)
			err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, generateAddress.Bytes(), types.ModuleName, coins)
			suite.Require().NoError(err)
		}
		count = len(md.DenomUnits[0].Aliases)
	}

	// add denom to erc20 module
	coin := sdk.NewCoin(md.Base, amt.Mul(sdkmath.NewInt(int64(count))))
	helpers.AddTestAddr(suite.app, suite.ctx, generateAddress.Bytes(), sdk.NewCoins(coin))
	err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, generateAddress.Bytes(), types.ModuleName, sdk.NewCoins(coin))
	suite.Require().NoError(err)

	return coin
}

func (suite *PrecompileTestSuite) BalanceOf(contractAddr, account common.Address) *big.Int {
	var balanceRes struct{ Value *big.Int }
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, account, contractAddr, contract.GetFIP20().ABI, "balanceOf", &balanceRes, account)
	suite.Require().NoError(err)
	return balanceRes.Value
}

func (suite *PrecompileTestSuite) MintERC20Token(signer *helpers.Signer, contractAddr, to common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	erc20 := contract.GetFIP20()
	transferData, err := erc20.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(signer, contractAddr, transferData)
}

func (suite *PrecompileTestSuite) ModuleMintERC20Token(contractAddr, to common.Address, amount *big.Int) {
	erc20 := contract.GetFIP20()
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.app.Erc20Keeper.ModuleAddress(), contractAddr, nil, erc20.ABI, "mint", to, amount)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
}

func (suite *PrecompileTestSuite) TransferERC20Token(signer *helpers.Signer, contractAddr, to common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	erc20 := contract.GetFIP20()
	transferData, err := erc20.ABI.Pack("transfer", to, amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(signer, contractAddr, transferData)
}

func (suite *PrecompileTestSuite) ERC20Approve(signer *helpers.Signer, contractAddr, to common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	erc20 := contract.GetFIP20()
	transferData, err := erc20.ABI.Pack("approve", to, amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(signer, contractAddr, transferData)
}

func (suite *PrecompileTestSuite) ERC20Allowance(contractAddr, owner, spender common.Address) *big.Int {
	var allowanceRes struct{ Value *big.Int }
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, owner, contractAddr, contract.GetFIP20().ABI, "allowance", &allowanceRes, owner, spender)
	suite.Require().NoError(err)
	return allowanceRes.Value
}

func (suite *PrecompileTestSuite) TransferERC20TokenToModule(signer *helpers.Signer, contractAddr common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	erc20 := contract.GetFIP20()
	moduleAddress := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	transferData, err := erc20.ABI.Pack("transfer", common.BytesToAddress(moduleAddress.Bytes()), amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(signer, contractAddr, transferData)
}

func (suite *PrecompileTestSuite) TransferERC20TokenToModuleWithoutHook(contractAddr, from common.Address, amount *big.Int) {
	erc20 := contract.GetFIP20()
	moduleAddress := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	_, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, from, contractAddr, nil, erc20.ABI, "transfer", common.BytesToAddress(moduleAddress.Bytes()), amount)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) RandPrefixAndAddress() (string, string) {
	if tmrand.Intn(10)%2 == 0 {
		return "0x", helpers.GenHexAddress().Hex()
	}
	prefix := strings.ToLower(tmrand.Str(5))
	accAddress, err := bech32.ConvertAndEncode(prefix, suite.RandSigner().AccAddress().Bytes())
	suite.Require().NoError(err)
	return prefix, accAddress
}

func (suite *PrecompileTestSuite) RandTransferChannel() (portID, channelID string) {
	portID = "transfer"

	channelSequence := suite.app.IBCKeeper.ChannelKeeper.GetNextChannelSequence(suite.ctx)
	channelID = fmt.Sprintf("channel-%d", channelSequence)
	connectionID := connectiontypes.FormatConnectionIdentifier(uint64(tmrand.Intn(100)))
	clientID := clienttypes.FormatClientIdentifier(exported.Localhost, uint64(tmrand.Intn(100)))

	revision := clienttypes.ParseChainID(suite.ctx.ChainID())
	localHostClient := localhost.NewClientState(clienttypes.NewHeight(revision, uint64(suite.ctx.BlockHeight())))
	suite.app.IBCKeeper.ClientKeeper.SetClientState(suite.ctx, clientID, localHostClient)

	params := suite.app.IBCKeeper.ClientKeeper.GetParams(suite.ctx)
	params.AllowedClients = append(params.AllowedClients, localHostClient.ClientType())
	suite.app.IBCKeeper.ClientKeeper.SetParams(suite.ctx, params)

	prevConsState := &ibctm.ConsensusState{
		Timestamp:          suite.ctx.BlockTime(),
		NextValidatorsHash: suite.ctx.BlockHeader().NextValidatorsHash,
	}
	height := clienttypes.NewHeight(0, uint64(suite.ctx.BlockHeight()))
	suite.app.IBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, clientID, height, prevConsState)

	channelCapability, err := suite.app.ScopedIBCKeeper.NewCapability(suite.ctx, host.ChannelCapabilityPath(portID, channelID))
	suite.Require().NoError(err)
	err = suite.app.ScopedTransferKeeper.ClaimCapability(suite.ctx, capabilitytypes.NewCapability(channelCapability.Index), host.ChannelCapabilityPath(portID, channelID))
	suite.Require().NoError(err)

	connectionEnd := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, clientID, connectiontypes.Counterparty{ClientId: "clientId", ConnectionId: "connection-1", Prefix: commitmenttypes.NewMerklePrefix([]byte("prefix"))}, connectiontypes.GetCompatibleVersions(), 500)
	suite.app.IBCKeeper.ConnectionKeeper.SetConnection(suite.ctx, connectionID, connectionEnd)

	channel := channeltypes.NewChannel(channeltypes.OPEN, channeltypes.ORDERED, channeltypes.NewCounterparty(portID, channelID), []string{connectionID}, "mock-version")
	suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, portID, channelID, channel)
	suite.app.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.ctx, portID, channelID, uint64(tmrand.Intn(10000)+1))
	suite.app.IBCKeeper.ChannelKeeper.SetNextChannelSequence(suite.ctx, channelSequence+1)
	return portID, channelID
}

func (suite *PrecompileTestSuite) AddIBCToken(portID, channelID string) string {
	denomTrace := ibctransfertypes.DenomTrace{
		Path:      fmt.Sprintf("%s/%s", portID, channelID),
		BaseDenom: "test",
	}
	suite.app.IBCTransferKeeper.SetDenomTrace(suite.ctx, denomTrace)
	return denomTrace.IBCDenom()
}

func (suite *PrecompileTestSuite) AddTokenToModule(module string, amt sdk.Coins) {
	tmpAddr := helpers.GenHexAddress()
	helpers.AddTestAddr(suite.app, suite.ctx, tmpAddr.Bytes(), amt)
	err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, tmpAddr.Bytes(), module, amt)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) ConvertOneToManyToken(md banktypes.Metadata) bool {
	emptyAlias := len(md.DenomUnits[0].Aliases) == 0
	if md.Base == fxtypes.DefaultDenom && !emptyAlias {
		return true
	}
	if strings.HasPrefix(md.Base, "ibc/") && !emptyAlias {
		return true
	}
	keepers := suite.CrossChainKeepers()
	for _, k := range keepers {
		if strings.HasPrefix(md.Base, k.ModuleName()) && !emptyAlias {
			return true
		}
	}
	return false
}

func (suite *PrecompileTestSuite) sendEvmTx(signer *helpers.Signer, contractAddr common.Address, data []byte) *evmtypes.MsgEthereumTxResponse {
	from := signer.Address()

	args, err := json.Marshal(&evmtypes.TransactionArgs{To: &contractAddr, From: &from, Data: (*hexutil.Bytes)(&data)})
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	res, err := evmtypes.NewQueryClient(queryHelper).EstimateGas(suite.ctx,
		&evmtypes.EthCallRequest{
			Args:    args,
			GasCap:  contract.DefaultGasCap,
			ChainId: suite.app.EvmKeeper.ChainID().Int64(),
		},
	)
	suite.Require().NoError(err)

	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	// suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64()*int64(res.Gas)))))

	msg := core.Message{
		To:                &contractAddr,
		From:              signer.Address(),
		Nonce:             suite.app.EvmKeeper.GetNonce(suite.ctx, signer.Address()),
		Value:             big.NewInt(0),
		GasLimit:          res.Gas,
		GasPrice:          suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx),
		GasFeeCap:         nil,
		GasTipCap:         nil,
		Data:              data,
		AccessList:        nil,
		SkipAccountChecks: false,
	}

	rsp, err := suite.app.EvmKeeper.ApplyMessage(suite.ctx, &msg, nil, true)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp.VmError)
	return rsp
}

type Metadata struct {
	metadata   banktypes.Metadata
	modules    []string
	notModules []string
}

func (m Metadata) RandModule() string {
	return m.modules[tmrand.Intn(len(m.modules))]
}

func (m Metadata) GetModules() []string {
	return m.modules
}

func (m Metadata) GetDenom(moduleName string) string {
	for _, denom := range m.metadata.DenomUnits[0].Aliases {
		if strings.HasPrefix(denom, moduleName) {
			return denom
		}
	}
	return ""
}

func (m Metadata) GetMetadata() banktypes.Metadata {
	return m.metadata
}

func (suite *PrecompileTestSuite) AddBridgeToken(moduleName, tokenContract string) string {
	bridgeDenom := crosschaintypes.NewBridgeDenom(moduleName, tokenContract)
	suite.CrossChainKeepers()[moduleName].AddBridgeToken(suite.ctx, bridgeDenom, bridgeDenom)
	return bridgeDenom
}

func (suite *PrecompileTestSuite) AddFXBridgeToken(tokenContract string) {
	bridgeDenom := crosschaintypes.NewBridgeDenom(crossethtypes.ModuleName, tokenContract)
	ethKeeper := suite.CrossChainKeepers()[crossethtypes.ModuleName]
	ethKeeper.AddBridgeToken(suite.ctx, bridgeDenom, fxtypes.DefaultDenom)
	ethKeeper.AddBridgeToken(suite.ctx, fxtypes.DefaultDenom, bridgeDenom)
}

type Oracle struct {
	moduleName string
	oracle     crosschaintypes.Oracle
}

func (o Oracle) GetExternalHexAddr() common.Address {
	return crosschaintypes.ExternalAddrToHexAddr(o.moduleName, o.oracle.ExternalAddress)
}
