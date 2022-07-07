package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/tests"
	fxtypes "github.com/functionx/fx-core/types"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	"github.com/functionx/fx-core/x/crosschain"
	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	"github.com/functionx/fx-core/x/erc20/types"
)

func (suite *KeeperTestSuite) TestHookChainBSC() {
	suite.purseBalance = sdk.NewInt(100000).Mul(sdk.NewInt(1e18))
	suite.SetupTest()

	signer1, addr1 := privateSigner()
	_, addr2 := privateSigner()

	suite.ctx = testInitBscCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), addr1.Bytes(), addr2)

	purseID := suite.app.Erc20Keeper.GetDenomMap(suite.ctx, PurseDenom)
	suite.Require().NotEmpty(purseID)

	tokenPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, purseID)
	suite.Require().True(found)
	suite.Require().NotNil(tokenPair)
	suite.Require().NotEmpty(tokenPair.GetERC20Contract())

	require.Equal(suite.T(), types.TokenPair{
		Erc20Address:  tokenPair.GetErc20Address(),
		Denom:         PurseDenom,
		Enabled:       true,
		ContractOwner: types.OWNER_MODULE,
	}, tokenPair)

	fip20, err := suite.app.Erc20Keeper.QueryERC20(suite.ctx, tokenPair.GetERC20Contract())
	suite.Require().NoError(err)
	suite.Require().Equal("PURSE", fip20.Symbol)

	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt), sdk.NewCoin(PurseDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1.Bytes())
	_ = balances

	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, addr1.Bytes(), addr1, sdk.NewCoin(PurseDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)

	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, tokenPair.GetERC20Contract(), addr1)
	suite.Require().NoError(err)
	_ = balanceOf

	token := tokenPair.GetERC20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, bsctypes.ModuleName)
	transferChainData := packTransferCrossData(suite.T(), addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, addr1, token, transferChainData)

	transactions := suite.app.BscKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(transactions))
	require.Equal(suite.T(), transactions[0].DestAddress, addr2.String())
	require.Equal(suite.T(), transactions[0].Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Sender, sdk.AccAddress(addr1.Bytes()).String())
}

type IBCTransferSimulate struct {
	T *testing.T
}

func (it *IBCTransferSimulate) SendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.Coin, sender sdk.AccAddress,
	receiver string, timeoutHeight ibcclienttypes.Height, timeoutTimestamp uint64, router string, fee sdk.Coin) error {
	require.Equal(it.T, token.Amount.BigInt(), ibcTransferAmount)
	return nil
}

type IBCChannelSimulate struct {
}

func (ic *IBCChannelSimulate) GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error) {
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
func (ic *IBCChannelSimulate) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	return 1, true
}

var (
	ibcTransferAmount = big.NewInt(1e18)
)

func (suite *KeeperTestSuite) TestHookIBC() {
	suite.SetupTest()

	pairId := suite.app.Erc20Keeper.GetDenomMap(suite.ctx, "FX")
	suite.Require().Greater(len(pairId), 0)

	pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, pairId)
	suite.Require().True(found)

	signer1, addr1 := privateSigner()
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err := suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1.Bytes())
	suite.Require().False(balances.IsZero())

	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, addr1.Bytes(), addr1, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)

	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, pair.GetERC20Contract(), addr1)
	suite.Require().NoError(err)
	suite.Require().Equal(balanceOf.Cmp(big.NewInt(0)), 1)

	//reset ibc
	suite.app.Erc20Keeper.SetIBCTransferKeeperForTest(&IBCTransferSimulate{T: suite.T()})
	suite.app.Erc20Keeper.SetIBCChannelKeeperForTest(&IBCChannelSimulate{})

	token := pair.GetERC20Contract()
	ibcTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToIBCPrefix, "px/transfer/channel-0")
	transferIBCData := packTransferCrossData(suite.T(), "px16u6kjunrcxkvaln9aetxwjpruply3sgwpr9z8u", ibcTransferAmount, big.NewInt(0), ibcTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, addr1, token, transferIBCData)
}

func packTransferCrossData(t *testing.T, to string, amount, fee *big.Int, target string) []byte {
	fip20 := fxtypes.GetERC20()
	targetBytes := fxtypes.StringToByte32(target)
	pack, err := fip20.ABI.Pack("transferCrossChain", to, amount, fee, targetBytes)
	require.NoError(t, err)
	return pack
}

func privateSigner() (keyring.Signer, common.Address) {
	// account key
	priKey := NewPriKey()
	//ethsecp256k1.GenerateKey()
	ethPriv := &ethsecp256k1.PrivKey{Key: priKey.Bytes()}

	return tests.NewSigner(ethPriv), common.BytesToAddress(ethPriv.PubKey().Address())
}

var (
	BSCBridgeTokenContract = common.HexToAddress("0x29a63F4B209C29B4DC47f06FFA896F32667DAD2C")
)

func testInitBscCrossChain(t *testing.T, ctx sdk.Context, myApp *app.App, oracleAddress, bridgeAddress sdk.AccAddress, externalAddress common.Address) sdk.Context {
	deposit := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18))))
	err := myApp.BankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, bsctypes.ModuleName, sdk.NewCoins(deposit))
	require.NoError(t, err)

	testBSCParamsProposal(t, ctx, myApp, oracleAddress)

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

	testBSCOracleSetUpdateClaim(t, ctx, myApp, bridgeAddress, externalAddress)

	testBSCBridgeTokenClaim(t, ctx, myApp, bridgeAddress)

	crosschain.EndBlocker(ctx, myApp.BscKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testBSCParamsProposal(t *testing.T, ctx sdk.Context, myApp *app.App, oracles sdk.AccAddress) {
	proposal := &crosschaintypes.UpdateChainOraclesProposal{
		Title:       "bsc cross chain",
		Description: "bsc cross chain oracles init",
		Oracles:     []string{oracles.String()},
		ChainName:   bsctypes.ModuleName,
	}

	k := &crosschainkeeper.EthereumMsgServer{Keeper: myApp.BscKeeper}
	err := crosschain.HandleUpdateChainOraclesProposal(ctx, k, proposal)
	require.NoError(t, err)
}

func testBSCBridgeTokenClaim(t *testing.T, ctx sdk.Context, myApp *app.App, orchAddr sdk.AccAddress) {
	msg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     2,
		BlockHeight:    uint64(ctx.BlockHeight()),
		TokenContract:  BSCBridgeTokenContract.String(),
		Name:           "PURSE Token",
		Symbol:         "PURSE",
		Decimals:       18,
		BridgerAddress: orchAddr.String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      bsctypes.ModuleName,
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = myApp.BscKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func testBSCOracleSetUpdateClaim(t *testing.T, ctx sdk.Context, myApp *app.App, orch sdk.AccAddress, addr common.Address) {
	msg := &crosschaintypes.MsgOracleSetUpdatedClaim{
		EventNonce:     1,
		BlockHeight:    uint64(ctx.BlockHeight()),
		OracleSetNonce: 0,
		Members: crosschaintypes.BridgeValidators{
			{
				Power:           uint64(math.MaxUint32),
				ExternalAddress: addr.String(),
			},
		},
		BridgerAddress: orch.String(),
		ChainName:      bsctypes.ModuleName,
	}
	for _, member := range msg.Members {
		_, found := myApp.BscKeeper.GetOracleByExternalAddress(ctx, member.ExternalAddress)
		require.True(t, found)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = myApp.BscKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}
