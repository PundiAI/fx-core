package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	"github.com/ethereum/go-ethereum/common"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	"github.com/functionx/fx-core/v3/x/crosschain"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
	"github.com/functionx/fx-core/v3/x/gravity"
	gravitykeeper "github.com/functionx/fx-core/v3/x/gravity/keeper"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	tronkeeper "github.com/functionx/fx-core/v3/x/tron/keeper"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func (suite *KeeperTestSuite) TestHookChainBSC() {
	suite.mintToken(bsctypes.ModuleName, sdk.NewCoin(purseDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))

	signer1 := helpers.NewSigner(helpers.NewEthPrivKey())
	addr2 := helpers.GenerateEthAddress()

	suite.ctx = testInitBscCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), signer1.Address().Bytes(), addr2)

	purseID := suite.app.Erc20Keeper.GetDenomMap(suite.ctx, purseDenom)
	suite.Require().NotEmpty(purseID)

	tokenPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, purseID)
	suite.Require().True(found)
	suite.Require().NotNil(tokenPair)
	suite.Require().NotEmpty(tokenPair.GetERC20Contract())

	require.Equal(suite.T(), types.TokenPair{
		Erc20Address:  tokenPair.GetErc20Address(),
		Denom:         purseDenom,
		Enabled:       true,
		ContractOwner: types.OWNER_MODULE,
	}, tokenPair)

	fip20, err := suite.app.Erc20Keeper.QueryERC20(suite.ctx, tokenPair.GetERC20Contract())
	suite.Require().NoError(err)
	suite.Require().Equal("PURSE", fip20.Symbol)

	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), signer1.Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt), sdk.NewCoin(purseDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer1.Address().Bytes())
	_ = balances

	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, signer1.Address().Bytes(), signer1.Address(), sdk.NewCoin(purseDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)

	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, tokenPair.GetERC20Contract(), signer1.Address())
	suite.Require().NoError(err)
	_ = balanceOf

	token := tokenPair.GetERC20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, bsctypes.ModuleName)
	transferChainData := packTransferCrossData(suite.T(), addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, signer1.Address(), token, transferChainData)

	transactions := suite.app.BscKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(transactions))
	require.Equal(suite.T(), transactions[0].DestAddress, addr2.String())
	require.Equal(suite.T(), transactions[0].Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Sender, sdk.AccAddress(signer1.Address().Bytes()).String())
}

func (suite *KeeperTestSuite) TestHookChainUSDT() {
	suite.mintToken(bsctypes.ModuleName, sdk.NewCoin(bscDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))
	suite.mintToken(polygontypes.ModuleName, sdk.NewCoin(polygonDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))

	signer1 := helpers.NewSigner(helpers.NewEthPrivKey())
	addr2 := helpers.GenerateEthAddress()
	signer3 := helpers.NewSigner(helpers.NewEthPrivKey())
	addr4 := helpers.GenerateEthAddress()

	suite.ctx = testInitBscCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), signer1.Address().Bytes(), addr2)
	suite.ctx = testInitPolygonCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), signer3.Address().Bytes(), addr4)

	usdtCopy := usdtMatedata
	usdtCopy.DenomUnits[0].Aliases = []string{bscDenom, polygonDenom}

	tokenPair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, usdtCopy)
	suite.Require().NoError(err)
	require.Equal(suite.T(), types.TokenPair{
		Erc20Address:  tokenPair.GetErc20Address(),
		Denom:         "usdt",
		Enabled:       true,
		ContractOwner: types.OWNER_MODULE,
	}, *tokenPair)

	denomBytes := suite.app.Erc20Keeper.GetAliasDenom(suite.ctx, bscDenom)
	require.Equal(suite.T(), usdtCopy.Base, string(denomBytes))

	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), signer1.Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt), sdk.NewCoin(bscDenom, amt)))
	suite.Require().NoError(err)

	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), signer3.Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt), sdk.NewCoin(polygonDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer1.Address().Bytes())
	_ = balances

	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, signer3.Address().Bytes())
	_ = balances

	usdtCoin, err := suite.app.Erc20Keeper.RelayConvertDenomToOne(suite.ctx, signer1.Address().Bytes(), sdk.NewCoin(bscDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)
	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, signer1.Address().Bytes(), signer1.Address(), usdtCoin)
	suite.Require().NoError(err)
	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, tokenPair.GetERC20Contract(), signer1.Address())
	suite.Require().NoError(err)
	_ = balanceOf

	usdtCoin, err = suite.app.Erc20Keeper.RelayConvertDenomToOne(suite.ctx, signer3.Address().Bytes(), sdk.NewCoin(polygonDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)
	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, signer3.Address().Bytes(), signer3.Address(), usdtCoin)
	suite.Require().NoError(err)
	balanceOf, err = suite.app.Erc20Keeper.BalanceOf(suite.ctx, tokenPair.GetERC20Contract(), signer3.Address())
	suite.Require().NoError(err)
	_ = balanceOf

	// addr1 transfer usdt to bsc
	token := tokenPair.GetERC20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, bsctypes.ModuleName)
	transferChainData := packTransferCrossData(suite.T(), addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, signer1.Address(), token, transferChainData)

	transactions := suite.app.BscKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(transactions))
	require.Equal(suite.T(), transactions[0].DestAddress, addr2.String())
	require.Equal(suite.T(), transactions[0].Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Sender, sdk.AccAddress(signer1.Address().Bytes()).String())

	// addr1 transfer usdt to polygon
	crossChainTarget = fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, polygontypes.ModuleName)
	transferChainData = packTransferCrossData(suite.T(), addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, signer1.Address(), token, transferChainData)

	transactions = suite.app.PolygonKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(transactions))
	require.Equal(suite.T(), transactions[0].DestAddress, addr2.String())
	require.Equal(suite.T(), transactions[0].Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Sender, sdk.AccAddress(signer1.Address().Bytes()).String())

	// addr3 transfer usdt to bsc
	crossChainTarget = fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, bsctypes.ModuleName)
	transferChainData = packTransferCrossData(suite.T(), addr4.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer3, signer3.Address(), token, transferChainData)

	transactions = suite.app.BscKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 2, len(transactions))
	transaction := getCrossChainOutgoingTransferTxById(transactions, 2)
	require.NotNil(suite.T(), transaction)
	require.Equal(suite.T(), transaction.DestAddress, addr4.String())
	require.Equal(suite.T(), transaction.Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transaction.Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transaction.Sender, sdk.AccAddress(signer3.Address().Bytes()).String())

	// addr2 transfer usdt to polygon
	crossChainTarget = fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, polygontypes.ModuleName)
	transferChainData = packTransferCrossData(suite.T(), addr4.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer3, signer3.Address(), token, transferChainData)

	transactions = suite.app.PolygonKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 2, len(transactions))
	transaction = getCrossChainOutgoingTransferTxById(transactions, 2)
	require.NotNil(suite.T(), transaction)
	require.Equal(suite.T(), transaction.DestAddress, addr4.String())
	require.Equal(suite.T(), transaction.Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transaction.Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transaction.Sender, sdk.AccAddress(signer3.Address().Bytes()).String())
}

type IBCTransferSimulate struct {
	T *testing.T
}

func (i *IBCTransferSimulate) SendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.Coin, sender sdk.AccAddress,
	receiver string, timeoutHeight ibcclienttypes.Height, timeoutTimestamp uint64, router string, fee sdk.Coin) error {
	require.Equal(i.T, token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(i.T, "transfer", sourcePort)
	require.Equal(i.T, "channel-0", sourceChannel)
	return nil
}

type IBCChannelSimulate struct {
	T *testing.T
}

func (i *IBCChannelSimulate) GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error) {
	require.Equal(i.T, "transfer", portID)
	require.Equal(i.T, "channel-0", channelID)
	return "", &ibctmtypes.ClientState{
		ChainId:         "fxcore",
		TrustLevel:      ibctmtypes.Fraction{},
		TrustingPeriod:  0,
		UnbondingPeriod: 0,
		MaxClockDrift:   0,
		FrozenHeight: ibcclienttypes.Height{
			RevisionHeight: 1000,
			RevisionNumber: 1000,
		},
		LatestHeight: ibcclienttypes.Height{
			RevisionHeight: 10,
			RevisionNumber: 10,
		},
		ProofSpecs:                   nil,
		UpgradePath:                  nil,
		AllowUpdateAfterExpiry:       false,
		AllowUpdateAfterMisbehaviour: false,
	}, nil
}
func (i *IBCChannelSimulate) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	require.Equal(i.T, "transfer", portID)
	require.Equal(i.T, "channel-0", channelID)
	return 1, true
}

func (suite *KeeperTestSuite) TestHookIBC() {
	pairId := suite.app.Erc20Keeper.GetDenomMap(suite.ctx, fxtypes.DefaultDenom)
	suite.Require().Greater(len(pairId), 0)

	pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, pairId)
	suite.Require().True(found)

	signer1 := helpers.NewSigner(helpers.NewEthPrivKey())
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err := suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), signer1.Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer1.Address().Bytes())
	suite.Require().False(balances.IsZero())

	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, signer1.Address().Bytes(), signer1.Address(), sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)

	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, pair.GetERC20Contract(), signer1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(balanceOf.Cmp(big.NewInt(0)), 1)

	//reset ibc
	(&suite.app.Erc20Keeper).IbcTransferKeeper = &IBCTransferSimulate{T: suite.T()}
	(&suite.app.Erc20Keeper).IbcChannelKeeper = &IBCChannelSimulate{T: suite.T()}

	token := pair.GetERC20Contract()
	ibcTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToIBCPrefix, "px/transfer/channel-0")
	transferIBCData := packTransferCrossData(suite.T(), "px16u6kjunrcxkvaln9aetxwjpruply3sgwpr9z8u", big.NewInt(1e18), big.NewInt(0), ibcTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, signer1.Address(), token, transferIBCData)
}

func (suite *KeeperTestSuite) TestHookIBCManyToOne() {
	suite.mintToken(polygontypes.ModuleName, sdk.NewCoin(polygonDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))

	_, usdtTokenPair := suite.setupRegisterCoinUSDT()

	signer1 := helpers.NewSigner(helpers.NewEthPrivKey())
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err := suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), signer1.Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt)))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), signer1.Address().Bytes(), sdk.NewCoins(sdk.NewCoin(polygonDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer1.Address().Bytes())
	suite.Require().False(balances.IsZero())

	usdtCoin, err := suite.app.Erc20Keeper.RelayConvertDenomToOne(suite.ctx, signer1.Address().Bytes(), sdk.NewCoin(polygonDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)

	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, signer1.Address().Bytes(), signer1.Address(), usdtCoin)
	suite.Require().NoError(err)

	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, usdtTokenPair.GetERC20Contract(), signer1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(balanceOf.Cmp(big.NewInt(0)), 1)

	//reset ibc
	(&suite.app.Erc20Keeper).IbcTransferKeeper = &IBCTransferSimulate{T: suite.T()}
	(&suite.app.Erc20Keeper).IbcChannelKeeper = &IBCChannelSimulate{T: suite.T()}

	token := usdtTokenPair.GetERC20Contract()
	ibcTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToIBCPrefix, "px/transfer/channel-0")
	transferIBCData := packTransferCrossData(suite.T(), "px16u6kjunrcxkvaln9aetxwjpruply3sgwpr9z8u", big.NewInt(1e18), big.NewInt(0), ibcTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, signer1.Address(), token, transferIBCData)
}

func (suite *KeeperTestSuite) TestHookIBCOneToMany() {
	suite.mintToken(bsctypes.ModuleName, sdk.NewCoin(bscDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))
	suite.mintToken(gravitytypes.ModuleName, sdk.NewCoin(ethDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))
	suite.mintToken(trontypes.ModuleName, sdk.NewCoin(tronDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))

	addr1 := helpers.GenerateAddress()
	tronExternal, err := tronaddress.Base58ToAddress("THtbMw6byXuiFhsRv1o1BQRtzvube9X1jx")
	suite.Require().NoError(err)

	suite.ctx = testInitGravityChain(suite.T(), suite.ctx, suite.app.GravityKeeper, suite.address.Bytes(), addr1.Bytes(), helpers.GenerateAddress())
	suite.ctx = testInitBscCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), helpers.GenerateAddress().Bytes(), helpers.GenerateAddress())
	suite.ctx = testInitTronCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), helpers.GenerateAddress().Bytes(), tronExternal)

	_, usdtTokenPair := suite.setupRegisterCoinUSDT(ethDenom, bscDenom, tronDenom)

	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt)))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(ethDenom, amt)))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(bscDenom, amt)))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(tronDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1.Bytes())
	suite.Require().False(balances.IsZero())

	_, err = suite.app.Erc20Keeper.RelayConvertDenomToOne(suite.ctx, addr1.Bytes(), sdk.NewCoin(ethDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)
	_, err = suite.app.Erc20Keeper.RelayConvertDenomToOne(suite.ctx, addr1.Bytes(), sdk.NewCoin(bscDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)
	_, err = suite.app.Erc20Keeper.RelayConvertDenomToOne(suite.ctx, addr1.Bytes(), sdk.NewCoin(tronDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)

	sender := sdk.AccAddress(addr1.Bytes()).String()
	hexReceiver := addr1.String()
	tronReceiver := "THtbMw6byXuiFhsRv1o1BQRtzvube9X1jx"
	usdtCoin := sdk.NewCoin(usdtTokenPair.Denom, sdk.NewInt(1))

	err = suite.app.GravityKeeper.TransferAfter(suite.ctx, sender, hexReceiver, usdtCoin, usdtCoin)
	suite.Require().NoError(err)
	gravityTransactions := suite.app.GravityKeeper.GetPoolTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(gravityTransactions))
	require.NotNil(suite.T(), gravityTransactions[0])
	require.Equal(suite.T(), gravityTransactions[0].DestAddress, hexReceiver)
	require.Equal(suite.T(), gravityTransactions[0].Erc20Token.Amount.BigInt(), big.NewInt(1))
	require.Equal(suite.T(), gravityTransactions[0].Erc20Fee.Amount.BigInt(), big.NewInt(1))
	require.Equal(suite.T(), gravityTransactions[0].Sender, sender)

	err = suite.app.BscKeeper.TransferAfter(suite.ctx, sender, hexReceiver, usdtCoin, usdtCoin)
	suite.Require().NoError(err)
	bscTransactions := suite.app.BscKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(bscTransactions))
	bscTransaction := getCrossChainOutgoingTransferTxById(bscTransactions, 1)
	require.NotNil(suite.T(), bscTransaction)
	require.Equal(suite.T(), bscTransaction.DestAddress, hexReceiver)
	require.Equal(suite.T(), bscTransaction.Token.Amount.BigInt(), big.NewInt(1))
	require.Equal(suite.T(), bscTransaction.Fee.Amount.BigInt(), big.NewInt(1))
	require.Equal(suite.T(), bscTransaction.Sender, sender)

	err = suite.app.TronKeeper.TransferAfter(suite.ctx, sender, tronReceiver, usdtCoin, usdtCoin)
	suite.Require().NoError(err)
	tronTransactions := suite.app.TronKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(tronTransactions))
	tronTransaction := getCrossChainOutgoingTransferTxById(tronTransactions, 1)
	require.NotNil(suite.T(), tronTransaction)
	require.Equal(suite.T(), tronTransaction.DestAddress, tronReceiver)
	require.Equal(suite.T(), tronTransaction.Token.Amount.BigInt(), big.NewInt(1))
	require.Equal(suite.T(), tronTransaction.Fee.Amount.BigInt(), big.NewInt(1))
	require.Equal(suite.T(), tronTransaction.Sender, sender)
}

func packTransferCrossData(t *testing.T, to string, amount, fee *big.Int, target string) []byte {
	fip20 := fxtypes.GetERC20()
	targetBytes := fxtypes.StringToByte32(target)
	pack, err := fip20.ABI.Pack("transferCrossChain", to, amount, fee, targetBytes)
	require.NoError(t, err)
	return pack
}

var (
	BSCBridgeTokenContract  = common.HexToAddress("0x29a63F4B209C29B4DC47f06FFA896F32667DAD2C")
	BSCUSDTokenContract     = common.HexToAddress("0x0000000000000000000000000000000000000001")
	PolygonUSDTokenContract = common.HexToAddress("0x0000000000000000000000000000000000000002")
	TronUSDTokenContract, _ = tronaddress.Base58ToAddress("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	EthUSDTokenContract     = common.HexToAddress("0xc2132D05D31c914a87C6611C10748AEb04B58e8F")
	bscDenom                = fmt.Sprintf("bsc%s", BSCUSDTokenContract.String())
	polygonDenom            = fmt.Sprintf("polygon%s", PolygonUSDTokenContract.String())
	tronDenom               = fmt.Sprintf("tron%s", TronUSDTokenContract.String())
	ethDenom                = fmt.Sprintf("eth%s", EthUSDTokenContract.String())
)

func testInitBscCrossChain(t *testing.T, ctx sdk.Context, myApp *app.App, oracleAddress, bridgeAddress sdk.AccAddress, externalAddress common.Address) sdk.Context {
	deposit := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18))))
	err := myApp.BankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, bsctypes.ModuleName, sdk.NewCoins(deposit))
	require.NoError(t, err)

	testCrossChainParamsProposal(t, ctx, myApp.BscKeeper, oracleAddress, bsctypes.ModuleName)

	oracle := crosschaintypes.Oracle{
		OracleAddress:   oracleAddress.String(),
		BridgerAddress:  bridgeAddress.String(),
		ExternalAddress: externalAddress.String(),
		DelegateAmount:  deposit.Amount,
		StartHeight:     ctx.BlockHeight(),
		Online:          true,
		SlashTimes:      0,
	}
	// save oracle
	myApp.BscKeeper.SetOracle(ctx, oracle)

	myApp.BscKeeper.SetOracleByBridger(ctx, bridgeAddress, oracleAddress)
	// set the ethereum address
	myApp.BscKeeper.SetOracleByExternalAddress(ctx, externalAddress.String(), oracleAddress)

	myApp.BscKeeper.CommonSetOracleTotalPower(ctx)

	testCrossChainOracleSetUpdateClaim(t, ctx, myApp.BscKeeper, bridgeAddress, externalAddress, 1, bsctypes.ModuleName)

	testCrossChainBridgeTokenClaim(t, ctx, myApp.BscKeeper, bridgeAddress, 2,
		BSCBridgeTokenContract, "PURSE Token", "PURSE", 18, bsctypes.ModuleName, hex.EncodeToString([]byte("transfer/channel-0")))
	testCrossChainBridgeTokenClaim(t, ctx, myApp.BscKeeper, bridgeAddress, 3,
		BSCUSDTokenContract, "USDT Token", "USDT", 6, bsctypes.ModuleName, "")

	crosschain.EndBlocker(ctx, myApp.BscKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testInitPolygonCrossChain(t *testing.T, ctx sdk.Context, myApp *app.App, oracleAddress, bridgeAddress sdk.AccAddress, externalAddress common.Address) sdk.Context {
	deposit := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18))))
	err := myApp.BankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, polygontypes.ModuleName, sdk.NewCoins(deposit))
	require.NoError(t, err)

	testCrossChainParamsProposal(t, ctx, myApp.PolygonKeeper, oracleAddress, polygontypes.ModuleName)

	oracle := crosschaintypes.Oracle{
		OracleAddress:   oracleAddress.String(),
		BridgerAddress:  bridgeAddress.String(),
		ExternalAddress: externalAddress.String(),
		DelegateAmount:  deposit.Amount,
		StartHeight:     ctx.BlockHeight(),
		Online:          true,
		SlashTimes:      0,
	}
	// save oracle
	myApp.PolygonKeeper.SetOracle(ctx, oracle)

	myApp.PolygonKeeper.SetOracleByBridger(ctx, bridgeAddress, oracleAddress)
	// set the ethereum address
	myApp.PolygonKeeper.SetOracleByExternalAddress(ctx, externalAddress.String(), oracleAddress)

	myApp.PolygonKeeper.CommonSetOracleTotalPower(ctx)

	testCrossChainOracleSetUpdateClaim(t, ctx, myApp.PolygonKeeper, bridgeAddress, externalAddress, 1, polygontypes.ModuleName)

	testCrossChainBridgeTokenClaim(t, ctx, myApp.PolygonKeeper, bridgeAddress, 2,
		PolygonUSDTokenContract, "USDT Token", "USDT", 6, polygontypes.ModuleName, "")

	crosschain.EndBlocker(ctx, myApp.PolygonKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testInitTronCrossChain(t *testing.T, ctx sdk.Context, myApp *app.App, oracleAddress, bridgeAddress sdk.AccAddress, externalAddress tronaddress.Address) sdk.Context {
	deposit := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18))))
	err := myApp.BankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, trontypes.ModuleName, sdk.NewCoins(deposit))
	require.NoError(t, err)

	testTronCrossChainParamsProposal(t, ctx, myApp.TronKeeper, oracleAddress, trontypes.ModuleName)

	oracle := crosschaintypes.Oracle{
		OracleAddress:   oracleAddress.String(),
		BridgerAddress:  bridgeAddress.String(),
		ExternalAddress: externalAddress.String(),
		DelegateAmount:  deposit.Amount,
		StartHeight:     ctx.BlockHeight(),
		Online:          true,
		SlashTimes:      0,
	}
	// save oracle
	myApp.TronKeeper.SetOracle(ctx, oracle)

	myApp.TronKeeper.SetOracleByBridger(ctx, bridgeAddress, oracleAddress)
	// set the ethereum address
	myApp.TronKeeper.SetOracleByExternalAddress(ctx, externalAddress.String(), oracleAddress)

	myApp.TronKeeper.CommonSetOracleTotalPower(ctx)

	testTronCrossChainOracleSetUpdateClaim(t, ctx, myApp.TronKeeper, bridgeAddress, externalAddress, 1, trontypes.ModuleName)

	testTronCrossChainBridgeTokenClaim(t, ctx, myApp.TronKeeper, bridgeAddress, 2,
		TronUSDTokenContract, "USDT Token", "USDT", 6, trontypes.ModuleName, "")

	crosschain.EndBlocker(ctx, myApp.TronKeeper.Keeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testInitGravityChain(t *testing.T, ctx sdk.Context, cck gravitykeeper.Keeper, val sdk.ValAddress, orch sdk.AccAddress, externalAddr common.Address) sdk.Context {
	msg := &gravitytypes.MsgSetOrchestratorAddress{
		Validator:    val.String(),
		Orchestrator: orch.String(),
		EthAddress:   externalAddr.String(),
	}
	impl := gravitykeeper.NewMsgServerImpl(cck)
	_, err := impl.SetOrchestratorAddress(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	gravity.EndBlocker(ctx, cck)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	return ctx
}

func testCrossChainParamsProposal(t *testing.T, ctx sdk.Context, cck crosschainkeeper.Keeper, oracles sdk.AccAddress, chain string) {
	proposal := &crosschaintypes.UpdateChainOraclesProposal{
		Title:       fmt.Sprintf("%s cross chain", chain),
		Description: fmt.Sprintf("%s cross chain oracles init", chain),
		Oracles:     []string{oracles.String()},
		ChainName:   chain,
	}

	k := &crosschainkeeper.EthereumMsgServer{Keeper: cck}
	err := crosschain.HandleUpdateChainOraclesProposal(ctx, k, proposal)
	require.NoError(t, err)
}

func testTronCrossChainParamsProposal(t *testing.T, ctx sdk.Context, cck tronkeeper.Keeper, oracles sdk.AccAddress, chain string) {
	proposal := &crosschaintypes.UpdateChainOraclesProposal{
		Title:       fmt.Sprintf("%s cross chain", chain),
		Description: fmt.Sprintf("%s cross chain oracles init", chain),
		Oracles:     []string{oracles.String()},
		ChainName:   chain,
	}

	k := &tronkeeper.TronMsgServer{EthereumMsgServer: crosschainkeeper.EthereumMsgServer{Keeper: cck.Keeper}}
	err := crosschain.HandleUpdateChainOraclesProposal(ctx, k, proposal)
	require.NoError(t, err)
}

func testCrossChainBridgeTokenClaim(t *testing.T, ctx sdk.Context, cck crosschainkeeper.Keeper,
	orchAddr sdk.AccAddress, eventNonce uint64, contract common.Address,
	name, symbol string, decimals uint64, chain, channelIBC string) {
	msg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     eventNonce,
		BlockHeight:    uint64(ctx.BlockHeight()),
		TokenContract:  contract.String(),
		Name:           name,
		Symbol:         symbol,
		Decimals:       decimals,
		BridgerAddress: orchAddr.String(),
		ChannelIbc:     channelIBC,
		ChainName:      chain,
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = cck.Attest(ctx, msg, anyMsg)
	require.NoError(t, err)
}

func testTronCrossChainBridgeTokenClaim(t *testing.T, ctx sdk.Context, cck tronkeeper.Keeper,
	orchAddr sdk.AccAddress, eventNonce uint64, contract tronaddress.Address,
	name, symbol string, decimals uint64, chain, channelIBC string) {
	msg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     eventNonce,
		BlockHeight:    uint64(ctx.BlockHeight()),
		TokenContract:  contract.String(),
		Name:           name,
		Symbol:         symbol,
		Decimals:       decimals,
		BridgerAddress: orchAddr.String(),
		ChannelIbc:     channelIBC,
		ChainName:      chain,
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = cck.Attest(ctx, msg, anyMsg)
	require.NoError(t, err)
}

func testCrossChainOracleSetUpdateClaim(t *testing.T, ctx sdk.Context, cck crosschainkeeper.Keeper,
	orch sdk.AccAddress, addr common.Address, eventNonce uint64, chain string) {
	msg := &crosschaintypes.MsgOracleSetUpdatedClaim{
		EventNonce:     eventNonce,
		BlockHeight:    uint64(ctx.BlockHeight()),
		OracleSetNonce: 0,
		Members: crosschaintypes.BridgeValidators{
			{
				Power:           uint64(math.MaxUint32),
				ExternalAddress: addr.String(),
			},
		},
		BridgerAddress: orch.String(),
		ChainName:      chain,
	}
	for _, member := range msg.Members {
		_, found := cck.GetOracleByExternalAddress(ctx, member.ExternalAddress)
		require.True(t, found)
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = cck.Attest(ctx, msg, anyMsg)
	require.NoError(t, err)
}

func testTronCrossChainOracleSetUpdateClaim(t *testing.T, ctx sdk.Context, cck tronkeeper.Keeper,
	orch sdk.AccAddress, addr tronaddress.Address, eventNonce uint64, chain string) {
	msg := &crosschaintypes.MsgOracleSetUpdatedClaim{
		EventNonce:     eventNonce,
		BlockHeight:    uint64(ctx.BlockHeight()),
		OracleSetNonce: 0,
		Members: crosschaintypes.BridgeValidators{
			{
				Power:           uint64(math.MaxUint32),
				ExternalAddress: addr.String(),
			},
		},
		BridgerAddress: orch.String(),
		ChainName:      chain,
	}
	for _, member := range msg.Members {
		_, found := cck.GetOracleByExternalAddress(ctx, member.ExternalAddress)
		require.True(t, found)
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = cck.Attest(ctx, msg, anyMsg)
	require.NoError(t, err)
}

func getCrossChainOutgoingTransferTxById(txs []*crosschaintypes.OutgoingTransferTx, id uint64) *crosschaintypes.OutgoingTransferTx {
	for _, tx := range txs {
		if tx.Id == id {
			return tx
		}
	}
	return nil
}
