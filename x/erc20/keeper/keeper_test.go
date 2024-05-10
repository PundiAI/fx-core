package keeper_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v6/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
	localhosttypes "github.com/cosmos/ibc-go/v6/modules/light-clients/09-localhost/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethereumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/evmos/ethermint/x/evm/statedb"
	evm "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/contract"
	fxserverconfig "github.com/functionx/fx-core/v7/server/config"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.App
	queryClient types.QueryClient
	signer      *helpers.Signer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Test helpers
func (suite *KeeperTestSuite) SetupTest() {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(suite.T(), err)
	suite.signer = helpers.NewSigner(priv)

	set, accs, balances := helpers.GenerateGenesisValidator(tmrand.Intn(10)+1, nil)
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), set, accs, balances...)

	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		Height:          suite.app.LastBlockHeight(),
		ChainID:         fxtypes.ChainId(),
		ProposerAddress: set.Proposer.Address,
		Time:            time.Now().UTC(),
	})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1e18))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.Erc20Keeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18)))))
}

func (suite *KeeperTestSuite) Commit() {
	suite.app.EndBlock(abci.RequestEndBlock{
		Height: suite.ctx.BlockHeight(),
	})
	suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height += 1
	header.Time = time.Now().UTC()
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})
	suite.ctx = suite.ctx.WithBlockHeight(header.Height)
}

func (suite *KeeperTestSuite) StateDB() *statedb.StateDB {
	return statedb.New(suite.ctx, suite.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.ctx.HeaderHash().Bytes())))
}

func (suite *KeeperTestSuite) RandSigner() *helpers.Signer {
	privKey := helpers.NewEthPrivKey()
	helpers.AddTestAddr(suite.app, suite.ctx, privKey.PubKey().Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18)))))
	return helpers.NewSigner(privKey)
}

func (suite *KeeperTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) BurnEvmRefundFee(addr sdk.AccAddress, coins sdk.Coins) {
	err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, addr, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)

	bal := suite.app.BankKeeper.GetBalance(suite.ctx, suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName), fxtypes.DefaultDenom)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, authtypes.FeeCollectorName, types.ModuleName, sdk.NewCoins(bal))
	suite.Require().NoError(err)

	err = suite.app.BankKeeper.BurnCoins(suite.ctx, types.ModuleName, sdk.NewCoins(bal))
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) DeployContract(from common.Address) (common.Address, error) {
	contractAddr, err := suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, suite.app.Erc20Keeper.ModuleAddress(), "Test token", "TEST", 18)
	suite.Require().NoError(err)

	_, err = suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.app.Erc20Keeper.ModuleAddress(), contractAddr, nil, contract.GetFIP20().ABI, "transferOwnership", from)
	suite.Require().NoError(err)
	return contractAddr, nil
}

func (suite *KeeperTestSuite) DeployFXRelayToken() (types.TokenPair, banktypes.Metadata) {
	fxToken := fxtypes.GetFXMetaData()

	pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, fxToken)
	suite.Require().NoError(err)
	return *pair, fxToken
}

func (suite *KeeperTestSuite) CrossChainKeepers() map[string]CrossChainKeeper {
	value := reflect.ValueOf(suite.app.CrossChainKeepers)
	keepers := make(map[string]CrossChainKeeper)
	for i := 0; i < value.NumField(); i++ {
		res := value.Field(i).MethodByName("GetGravityID").Call([]reflect.Value{reflect.ValueOf(suite.ctx)})
		gravityID := res[0].String()
		chainName := strings.TrimSuffix(strings.TrimPrefix(gravityID, "fx-"), "-bridge")
		cck := value.Field(i).Interface().(CrossChainKeeper)
		if chainName == "bridge-eth" {
			// keepers["gravity"] = cck
			keepers["eth"] = cck
		} else {
			keepers[chainName] = cck
		}
	}
	return keepers
}

func (suite *KeeperTestSuite) GenerateCrossChainDenoms(addDenoms ...string) Metadata {
	keepers := suite.CrossChainKeepers()
	modules := make([]string, 0, len(keepers))
	for m := range keepers {
		modules = append(modules, m)
	}
	count := tmrand.Intn(len(modules)-1) + 1

	denoms := make([]string, len(modules))
	denomModules := make([]string, len(modules))
	for index, m := range modules {
		address := helpers.GenerateAddressByModule(m)

		denom := crosschaintypes.NewBridgeDenom(m, address)
		denoms[index] = denom
		denomModules[index] = m

		k := keepers[m]
		k.AddBridgeToken(suite.ctx, address, crosschaintypes.NewBridgeDenom(m, address))
	}
	if count >= len(modules) {
		count = len(modules) - 1
	}
	metadata := fxtypes.GetCrossChainMetadataManyToOne("Test Token", helpers.NewRandSymbol(), 18, append(denoms[:count], addDenoms...)...)
	return Metadata{metadata: metadata, modules: denomModules[:count], notModules: denomModules[count:]}
}

func (suite *KeeperTestSuite) MintLockNativeTokenToModule(md banktypes.Metadata, amt sdkmath.Int) *big.Int {
	generateAddress := helpers.GenerateAddress()

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

	return coin.Amount.BigInt()
}

func (suite *KeeperTestSuite) BalanceOf(contractAddr, account common.Address) *big.Int {
	var balanceRes struct{ Value *big.Int }
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, account, contractAddr, contract.GetFIP20().ABI, "balanceOf", &balanceRes, account)
	suite.NoError(err)
	return balanceRes.Value
}

func (suite *KeeperTestSuite) MintERC20Token(signer *helpers.Signer, contractAddr, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := contract.GetFIP20()
	transferData, err := erc20.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(signer, contractAddr, transferData)
}

func (suite *KeeperTestSuite) ModuleMintERC20Token(contractAddr, to common.Address, amount *big.Int) {
	erc20 := contract.GetFIP20()
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.app.Erc20Keeper.ModuleAddress(), contractAddr, nil, erc20.ABI, "mint", to, amount)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
}

func (suite *KeeperTestSuite) TransferERC20Token(signer *helpers.Signer, contractAddr, to common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := contract.GetFIP20()
	transferData, err := erc20.ABI.Pack("transfer", to, amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(signer, contractAddr, transferData)
}

func (suite *KeeperTestSuite) TransferERC20TokenToModule(signer *helpers.Signer, contractAddr common.Address, amount *big.Int) *evm.MsgEthereumTx {
	erc20 := contract.GetFIP20()
	moduleAddress := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	transferData, err := erc20.ABI.Pack("transfer", common.BytesToAddress(moduleAddress.Bytes()), amount)
	suite.Require().NoError(err)
	return suite.sendEvmTx(signer, contractAddr, transferData)
}

func (suite *KeeperTestSuite) TransferERC20TokenToModuleWithoutHook(contractAddr, from common.Address, amount *big.Int) {
	erc20 := contract.GetFIP20()
	moduleAddress := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	_, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, from, contractAddr, nil, erc20.ABI, "transfer", common.BytesToAddress(moduleAddress.Bytes()), amount)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) RandPrefixAndAddress() (string, string) {
	if tmrand.Intn(10)%2 == 0 {
		return "0x", helpers.GenerateAddress().Hex()
	}
	prefix := strings.ToLower(tmrand.Str(5))
	accAddress, err := bech32.ConvertAndEncode(prefix, suite.RandSigner().AccAddress().Bytes())
	suite.NoError(err)
	return prefix, accAddress
}

func (suite *KeeperTestSuite) RandTransferChannel() (portID, channelID string) {
	portID = "transfer"
	channelID = fmt.Sprintf("channel-%d", tmrand.Intn(100))
	connectionID := connectiontypes.FormatConnectionIdentifier(uint64(tmrand.Intn(100)))
	clientID := clienttypes.FormatClientIdentifier(exported.Localhost, uint64(tmrand.Intn(100)))

	revision := clienttypes.ParseChainID(suite.ctx.ChainID())
	localHostClient := localhosttypes.NewClientState(
		suite.ctx.ChainID(), clienttypes.NewHeight(revision, uint64(suite.ctx.BlockHeight())),
	)
	suite.app.IBCKeeper.ClientKeeper.SetClientState(suite.ctx, clientID, localHostClient)

	prevConsState := &ibctmtypes.ConsensusState{
		Timestamp:          suite.ctx.BlockTime(),
		NextValidatorsHash: suite.ctx.BlockHeader().NextValidatorsHash,
	}
	height := clienttypes.NewHeight(0, uint64(suite.ctx.BlockHeight()))
	suite.app.IBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, clientID, height, prevConsState)

	channelCapability, err := suite.app.ScopedIBCKeeper.NewCapability(suite.ctx, host.ChannelCapabilityPath(portID, channelID))
	suite.Require().NoError(err)
	err = suite.app.ScopedTransferKeeper.ClaimCapability(suite.ctx, capabilitytypes.NewCapability(channelCapability.Index), host.ChannelCapabilityPath(portID, channelID))
	suite.Require().NoError(err)

	connectionEnd := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, clientID, connectiontypes.Counterparty{ClientId: "clientId", ConnectionId: "connection-1", Prefix: commitmenttypes.NewMerklePrefix([]byte("prefix"))}, []*connectiontypes.Version{ibctesting.ConnectionVersion}, 500)
	suite.app.IBCKeeper.ConnectionKeeper.SetConnection(suite.ctx, connectionID, connectionEnd)

	channel := channeltypes.NewChannel(channeltypes.OPEN, channeltypes.ORDERED, channeltypes.NewCounterparty(portID, channelID), []string{connectionID}, ibctesting.DefaultChannelVersion)
	suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, portID, channelID, channel)
	suite.app.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.ctx, portID, channelID, uint64(tmrand.Intn(10000)))

	return portID, channelID
}

func (suite *KeeperTestSuite) AddIBCToken(portID, channelID string) string {
	denomTrace := ibctransfertypes.DenomTrace{
		Path:      fmt.Sprintf("%s/%s", portID, channelID),
		BaseDenom: "test",
	}
	suite.app.IBCTransferKeeper.SetDenomTrace(suite.ctx, denomTrace)
	return denomTrace.IBCDenom()
}

func (suite *KeeperTestSuite) sendEvmTx(signer *helpers.Signer, contractAddr common.Address, data []byte) *evm.MsgEthereumTx {
	chainID := suite.app.EvmKeeper.ChainID()
	from := signer.Address()

	args, err := json.Marshal(&evm.TransactionArgs{To: &contractAddr, From: &from, Data: (*hexutil.Bytes)(&data)})
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	res, err := evm.NewQueryClient(queryHelper).EstimateGas(sdk.WrapSDKContext(suite.ctx),
		&evm.EthCallRequest{
			Args:    args,
			GasCap:  fxserverconfig.DefaultGasCap,
			ChainId: suite.app.EvmKeeper.ChainID().Int64(),
		},
	)
	suite.Require().NoError(err)

	totalSupplyBefore := suite.app.BankKeeper.GetSupply(suite.ctx, fxtypes.DefaultDenom)
	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	mintAmount := sdkmath.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64() * int64(res.Gas))
	suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, mintAmount)))

	ercTransferTx := evm.NewTx(
		chainID,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		&contractAddr,
		nil,
		res.Gas,
		nil,
		suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx),
		big.NewInt(1),
		data,
		&ethereumtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = signer.Address().Hex()
	err = ercTransferTx.Sign(ethereumtypes.LatestSignerForChainID(chainID), signer)
	suite.Require().NoError(err)

	rsp, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ercTransferTx)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)

	refundAmount := sdkmath.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64() * int64(res.Gas-rsp.GasUsed))
	suite.BurnEvmRefundFee(signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, refundAmount)))

	totalSupplyAfter := suite.app.BankKeeper.GetSupply(suite.ctx, fxtypes.DefaultDenom)
	suite.Require().Equal(totalSupplyBefore.String(), totalSupplyAfter.String())

	return ercTransferTx
}

func newMetadata() banktypes.Metadata {
	return banktypes.Metadata{
		Description: "description of the token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "usdt",
				Exponent: uint32(0),
				Aliases: []string{
					crosschaintypes.NewBridgeDenom(bsctypes.ModuleName, helpers.GenerateAddress().String()),
					crosschaintypes.NewBridgeDenom(ethtypes.ModuleName, helpers.GenerateAddress().String()),
					// fmt.Sprintf("%s%s", "ibc/", helpers.GenerateAddress().String()),
				},
			}, {
				Denom:    "USDT",
				Exponent: uint32(18),
			},
		},
		Base:    "usdt",
		Display: "display usdt",
		Name:    "Tether USD",
		Symbol:  "USDT",
	}
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

type CrossChainKeeper interface {
	AddBridgeToken(ctx sdk.Context, token, denom string)
	SetIbcDenomTrace(ctx sdk.Context, token, channelIBC string) (string, error)
}
