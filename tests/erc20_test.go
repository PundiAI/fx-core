package tests

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v2/app/helpers"
	fxtypes "github.com/functionx/fx-core/v2/types"
	"github.com/functionx/fx-core/v2/types/contract"
	bsctypes "github.com/functionx/fx-core/v2/x/bsc/types"
	erc20types "github.com/functionx/fx-core/v2/x/erc20/types"
	polygontypes "github.com/functionx/fx-core/v2/x/polygon/types"
)

type ERC20TestSuite struct {
	*TestSuite
	privKey cryptotypes.PrivKey
}

func NewERC20TestSuite() ERC20TestSuite {
	return ERC20TestSuite{
		TestSuite: NewTestSuite(),
		privKey:   helpers.NewEthPrivKey(),
	}
}
func NewERC20WithTestSuite(ts *TestSuite) ERC20TestSuite {
	return ERC20TestSuite{
		TestSuite: ts,
		privKey:   helpers.NewEthPrivKey(),
	}
}

func (suite *ERC20TestSuite) SetupSuite() {
	// set test env
	suite.NoError(os.Setenv("GO_ENV", "testing"))

	suite.TestSuite.SetupSuite()
	suite.Send(suite.Address(), suite.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
}

func (suite *ERC20TestSuite) Address() sdk.AccAddress {
	return suite.privKey.PubKey().Address().Bytes()
}

func (suite *ERC20TestSuite) HexAddr() gethcommon.Address {
	return gethcommon.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *ERC20TestSuite) ERC20Query() erc20types.QueryClient {
	return suite.GRPCClient().ERC20Query()
}

func (suite *ERC20TestSuite) RegisterCoinProposal(md banktypes.Metadata) (proposalId uint64) {
	proposal, err := govtypes.NewMsgSubmitProposal(
		&erc20types.RegisterCoinProposal{
			Title:       fmt.Sprintf("register %s denom", md.Base),
			Description: "bar",
			Metadata:    md,
		},
		sdk.NewCoins(suite.NewCoin(sdk.NewInt(10_000).MulRaw(1e18))),
		suite.Address(),
	)
	suite.NoError(err)
	return suite.BroadcastProposalTx(suite.privKey, proposal)
}

func (suite *ERC20TestSuite) ToggleTokenConversionProposal(denom string) (proposalId uint64) {
	proposal, err := govtypes.NewMsgSubmitProposal(
		&erc20types.ToggleTokenConversionProposal{
			Title:       fmt.Sprintf("update %s denom", denom),
			Description: "update",
			Token:       denom,
		},
		sdk.NewCoins(suite.NewCoin(sdk.NewInt(10_000).MulRaw(1e18))),
		suite.Address(),
	)
	suite.NoError(err)
	return suite.BroadcastProposalTx(suite.privKey, proposal)
}

func (suite *ERC20TestSuite) UpdateDenomAliasProposal(denom, alias string) (proposalId uint64) {
	proposal, err := govtypes.NewMsgSubmitProposal(
		&erc20types.UpdateDenomAliasProposal{
			Title:       fmt.Sprintf("update %s denom %s alias", denom, alias),
			Description: "update",
			Denom:       denom,
			Alias:       alias,
		},
		sdk.NewCoins(suite.NewCoin(sdk.NewInt(10_000).MulRaw(1e18))),
		suite.Address(),
	)
	suite.NoError(err)
	return suite.BroadcastProposalTx(suite.privKey, proposal)
}

func (suite *ERC20TestSuite) CheckRegisterCoin(denom string, manyToOne ...bool) {
	_, err := suite.ERC20Query().TokenPair(suite.ctx, &erc20types.QueryTokenPairRequest{Token: denom})
	suite.NoError(err)
	if len(manyToOne) > 0 && manyToOne[0] {
		aliasesResp, err := suite.ERC20Query().DenomAliases(suite.ctx, &erc20types.QueryDenomAliasesRequest{Denom: denom})
		suite.NoError(err)
		suite.T().Log("denom", denom, "alias", aliasesResp.Aliases)
		for _, alias := range aliasesResp.Aliases {
			aliasDenom, err := suite.ERC20Query().AliasDenom(suite.ctx, &erc20types.QueryAliasDenomRequest{Alias: alias})
			suite.NoError(err)
			suite.Equal(denom, aliasDenom.Denom)
		}
	}
}

func (suite *ERC20TestSuite) TokenPair(denom string) erc20types.TokenPair {
	pairResp, err := suite.ERC20Query().TokenPair(suite.ctx, &erc20types.QueryTokenPairRequest{Token: denom})
	suite.NoError(err)
	return pairResp.TokenPair
}

func (suite *ERC20TestSuite) EthClient() *ethclient.Client {
	return suite.GetFirstValidtor().JSONRPCClient
}

func (suite *ERC20TestSuite) EthBalance(address gethcommon.Address) *big.Int {
	amount, err := suite.EthClient().BalanceAt(context.Background(), address, nil)
	suite.NoError(err)
	return amount
}

func (suite *ERC20TestSuite) BalanceOf(contractAddr, address gethcommon.Address) *big.Int {
	caller, err := contract.NewFIP20(contractAddr, suite.EthClient())
	suite.NoError(err)
	balance, err := caller.BalanceOf(nil, address)
	suite.NoError(err)
	return balance
}

func (suite *ERC20TestSuite) ConvertDenom(private cryptotypes.PrivKey, receiver sdk.AccAddress, coin sdk.Coin, target string) {
	suite.BroadcastTx(private, &erc20types.MsgConvertDenom{
		Sender:   sdk.AccAddress(private.PubKey().Address()).String(),
		Receiver: receiver.String(),
		Coin:     coin,
		Target:   target,
	})
}

func (suite *ERC20TestSuite) SendTransaction(tx *ethtypes.Transaction) {
	err := suite.EthClient().SendTransaction(context.Background(), tx)
	require.NoError(suite.T(), err)

	suite.T().Log("pending tx hash", tx.Hash())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.EthClient(), tx)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), receipt.Status, ethtypes.ReceiptStatusSuccessful)
}

func (suite *ERC20TestSuite) Transfer(privateKey cryptotypes.PrivKey, recipient gethcommon.Address, value *big.Int) gethcommon.Hash {
	suite.T().Logf("transfer to %s value %s\n", recipient.String(), value.String())
	ethTx, err := dynamicFeeTx(suite.EthClient(), privateKey, &recipient, value, nil)
	require.NoError(suite.T(), err)

	suite.SendTransaction(ethTx)
	return ethTx.Hash()
}

func (suite *ERC20TestSuite) TransferERC20(privateKey cryptotypes.PrivKey, token, recipient gethcommon.Address, value *big.Int) gethcommon.Hash {
	suite.T().Logf("transfer erc20 %s to %s value %s\n", token, recipient.String(), value.String())

	pack, err := FIP20ABI.Pack("transfer", recipient, value)
	require.NoError(suite.T(), err)

	ethTx, err := dynamicFeeTx(suite.EthClient(), privateKey, &token, nil, pack)
	require.NoError(suite.T(), err)

	suite.SendTransaction(ethTx)
	return ethTx.Hash()
}

func (suite *ERC20TestSuite) TransferCrossChain(privateKey cryptotypes.PrivKey, token gethcommon.Address,
	recipient string, amount, fee *big.Int, target string) gethcommon.Hash {
	suite.T().Log("transfer cross chain", target)
	pack, err := FIP20ABI.Pack("transferCrossChain", recipient, amount, fee, fxtypes.StringToByte32(target))
	require.NoError(suite.T(), err)

	ethTx, err := dynamicFeeTx(suite.EthClient(), privateKey, &token, nil, pack)
	require.NoError(suite.T(), err)

	suite.SendTransaction(ethTx)

	return ethTx.Hash()
}

// crosschain with erc20 test suite
type CrosschainERC20TestSuite struct {
	*TestSuite
	BSCCrossChain     CrosschainTestSuite
	PolygonCrossChain CrosschainTestSuite
	ERC20             ERC20TestSuite
}

func TestCrosschainERC20TestSuite(t *testing.T) {
	testSuite := NewTestSuite()
	crosschainERC20TestSuite := &CrosschainERC20TestSuite{
		TestSuite:         testSuite,
		BSCCrossChain:     NewCrosschainWithTestSuite(bsctypes.ModuleName, testSuite),
		PolygonCrossChain: NewCrosschainWithTestSuite(polygontypes.ModuleName, testSuite),
		ERC20:             NewERC20WithTestSuite(testSuite),
	}
	suite.Run(t, crosschainERC20TestSuite)
}

func NewCrosschainERC20TestSuite(ts *TestSuite) CrosschainERC20TestSuite {
	return CrosschainERC20TestSuite{
		TestSuite:         ts,
		BSCCrossChain:     NewCrosschainWithTestSuite(bsctypes.ModuleName, ts),
		PolygonCrossChain: NewCrosschainWithTestSuite(polygontypes.ModuleName, ts),
		ERC20:             NewERC20WithTestSuite(ts),
	}
}

func (suite *CrosschainERC20TestSuite) SetupSuite() {
	// set test env
	suite.NoError(os.Setenv("GO_ENV", "testing"))

	suite.TestSuite.SetupSuite()

	suite.Send(suite.BSCCrossChain.OracleAddr(), suite.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
	suite.Send(suite.BSCCrossChain.BridgerFxAddr(), suite.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.Send(suite.BSCCrossChain.AccAddr(), suite.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.BSCCrossChain.params = suite.BSCCrossChain.QueryParams()

	suite.Send(suite.PolygonCrossChain.OracleAddr(), suite.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
	suite.Send(suite.PolygonCrossChain.BridgerFxAddr(), suite.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.Send(suite.PolygonCrossChain.AccAddr(), suite.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.PolygonCrossChain.params = suite.PolygonCrossChain.QueryParams()

	suite.Send(suite.ERC20.Address(), suite.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
}

func (suite *CrosschainERC20TestSuite) InitCrossChain() {
	// bsc crosschain
	bscUSDTDenom := fmt.Sprintf("%s%s", suite.BSCCrossChain.chainName, bscUSDToken)

	proposalId := suite.BSCCrossChain.SendUpdateChainOraclesProposal()
	suite.ProposalVote(suite.AdminPrivateKey(), proposalId, govtypes.OptionYes)
	suite.CheckProposal(proposalId, govtypes.StatusPassed)

	suite.BSCCrossChain.BondedOracle()
	suite.BSCCrossChain.SendOracleSetConfirm()

	denom := suite.BSCCrossChain.AddBridgeTokenClaim("Tether USD", "USDT", 18, bscUSDToken, "")
	suite.Equal(denom, bscUSDTDenom)

	suite.BSCCrossChain.SendToFxClaim(bscUSDToken, sdk.NewInt(100).MulRaw(1e18), "")

	// polygon crosschain
	polygonUSDTDenom := fmt.Sprintf("%s%s", suite.PolygonCrossChain.chainName, polygonUSDToken)

	proposalId = suite.PolygonCrossChain.SendUpdateChainOraclesProposal()
	suite.ProposalVote(suite.AdminPrivateKey(), proposalId, govtypes.OptionYes)
	suite.CheckProposal(proposalId, govtypes.StatusPassed)

	suite.PolygonCrossChain.BondedOracle()
	suite.PolygonCrossChain.SendOracleSetConfirm()

	denom = suite.PolygonCrossChain.AddBridgeTokenClaim("Tether USD", "USDT", 18, polygonUSDToken, "")
	suite.Equal(denom, polygonUSDTDenom)

	suite.PolygonCrossChain.SendToFxClaim(polygonUSDToken, sdk.NewInt(100).MulRaw(1e18), "")
}

func (suite *CrosschainERC20TestSuite) InitRegisterCoinUSDT() {
	bscUSDTDenom := fmt.Sprintf("%s%s", suite.BSCCrossChain.chainName, bscUSDToken)
	polygonUSDTDenom := fmt.Sprintf("%s%s", suite.PolygonCrossChain.chainName, polygonUSDToken)
	// erc20
	usdtMetadata := banktypes.Metadata{
		Description: "description of the token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "usdt",
				Exponent: uint32(0),
				Aliases:  []string{bscUSDTDenom, polygonUSDTDenom},
			}, {
				Denom:    "USDT",
				Exponent: uint32(18),
			},
		},
		Base:    "usdt",
		Display: "usdt",
		Name:    "Tether USD",
		Symbol:  "USDT",
	}
	proposalId := suite.ERC20.RegisterCoinProposal(usdtMetadata)
	suite.ProposalVote(suite.AdminPrivateKey(), proposalId, govtypes.OptionYes)
	suite.CheckProposal(proposalId, govtypes.StatusPassed)
	suite.ERC20.CheckRegisterCoin(usdtMetadata.Base, true)

	usdtTokenPair := suite.ERC20.TokenPair("usdt")
	suite.T().Log("token pair", usdtTokenPair.String())
}
