package tests

import (
	"fmt"
	"math/big"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v3/client"
	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

type Erc20TestSuite struct {
	EvmTestSuite
}

func NewErc20TestSuite(ts *TestSuite) Erc20TestSuite {
	return Erc20TestSuite{
		EvmTestSuite: NewEvmTestSuite(ts),
	}
}

func (suite *Erc20TestSuite) ERC20Query() erc20types.QueryClient {
	return suite.GRPCClient().ERC20Query()
}

func (suite *Erc20TestSuite) CheckRegisterCoin(denom string) {
	aliasesResp, err := suite.ERC20Query().DenomAliases(suite.ctx, &erc20types.QueryDenomAliasesRequest{Denom: denom})
	suite.NoError(err)
	for _, alias := range aliasesResp.Aliases {
		aliasDenom, err := suite.ERC20Query().AliasDenom(suite.ctx, &erc20types.QueryAliasDenomRequest{Alias: alias})
		suite.NoError(err)
		suite.Equal(denom, aliasDenom.Denom)
	}
	pair := suite.TokenPair(denom)
	suite.NotNil(pair)
	metadata := suite.GetMetadata(denom)
	suite.True(len(metadata.DenomUnits) > 0)
	suite.Equal(metadata.DenomUnits[0].Aliases, aliasesResp.Aliases)
}

func (suite *Erc20TestSuite) TokenPair(denom string) *erc20types.TokenPair {
	pairResp, err := suite.ERC20Query().TokenPair(suite.ctx, &erc20types.QueryTokenPairRequest{Token: denom})
	suite.NoError(err)
	return &pairResp.TokenPair
}

func (suite *Erc20TestSuite) Erc20TokenAddress(denom string) common.Address {
	return suite.TokenPair(denom).GetERC20Contract()
}

func (suite *Erc20TestSuite) DenomFromErc20(address common.Address) string {
	pairs, err := suite.ERC20Query().TokenPairs(suite.ctx, &erc20types.QueryTokenPairsRequest{})
	suite.NoError(err)
	for _, pair := range pairs.TokenPairs {
		if pair.Erc20Address == address.String() {
			return pair.Denom
		}
	}
	return ""
}

func (suite *Erc20TestSuite) RegisterCoinProposal(md banktypes.Metadata) (*sdk.TxResponse, uint64) {
	content := &erc20types.RegisterCoinProposal{
		Title:       fmt.Sprintf("register %s denom", md.Base),
		Description: "bar",
		Metadata:    md,
	}
	return suite.BroadcastProposalTx(content)
}

func (suite *Erc20TestSuite) ToggleTokenConversionProposal(denom string) (*sdk.TxResponse, uint64) {
	content := &erc20types.ToggleTokenConversionProposal{
		Title:       fmt.Sprintf("update %s denom", denom),
		Description: "update",
		Token:       denom,
	}
	return suite.BroadcastProposalTx(content)
}

func (suite *Erc20TestSuite) UpdateDenomAliasProposal(denom, alias string) (*sdk.TxResponse, uint64) {
	content := &erc20types.UpdateDenomAliasProposal{
		Title:       fmt.Sprintf("update %s denom %s alias", denom, alias),
		Description: "update",
		Denom:       denom,
		Alias:       alias,
	}
	return suite.BroadcastProposalTx(content)
}

func (suite *Erc20TestSuite) ConvertCoin(private cryptotypes.PrivKey, recipient common.Address, coin sdk.Coin) *sdk.TxResponse {
	suite.True(suite.DenomBalanceGTE(sdk.AccAddress(private.PubKey().Address()), coin))
	msg := erc20types.NewMsgConvertCoin(coin, recipient, sdk.AccAddress(private.PubKey().Address()))
	txResponse := suite.BroadcastTx(suite.privKey, msg)
	suite.TokenBalanceGTE(recipient, suite.Erc20TokenAddress(coin.Denom), coin.Amount.BigInt())
	return txResponse
}

func (suite *Erc20TestSuite) ConvertERC20(private cryptotypes.PrivKey, token common.Address, amount sdk.Int, recipient sdk.AccAddress) *sdk.TxResponse {
	suite.True(suite.TokenBalanceGTE(common.BytesToAddress(private.PubKey().Address().Bytes()), token, amount.BigInt()))
	msg := erc20types.NewMsgConvertERC20(amount, recipient, token, common.BytesToAddress(private.PubKey().Address().Bytes()))
	txResponse := suite.BroadcastTx(private, msg)
	suite.True(suite.DenomBalanceGTE(recipient, sdk.NewCoin(suite.DenomFromErc20(token), amount)))
	return txResponse
}

func (suite *Erc20TestSuite) ConvertDenom(private cryptotypes.PrivKey, receiver sdk.AccAddress, coin sdk.Coin, target string) *sdk.TxResponse {
	suite.True(suite.DenomBalanceGTE(sdk.AccAddress(private.PubKey().Address()), coin))
	txResponse := suite.BroadcastTx(private, &erc20types.MsgConvertDenom{
		Sender:   sdk.AccAddress(private.PubKey().Address()).String(),
		Receiver: receiver.String(),
		Coin:     coin,
		Target:   target,
	})
	return txResponse
}

func (suite *Erc20TestSuite) TransferCrossChain(privateKey cryptotypes.PrivKey, token common.Address, recipient string, amount, fee *big.Int, target string) *ethtypes.Transaction {
	suite.TokenBalanceGTE(common.BytesToAddress(privateKey.PubKey().Address().Bytes()), token, new(big.Int).Add(amount, fee))
	pack, err := fxtypes.GetERC20().ABI.Pack("transferCrossChain", recipient, amount, fee, fxtypes.MustStrToByte32(target))
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *Erc20TestSuite) DenomBalanceGTE(account sdk.AccAddress, coin sdk.Coin) bool {
	balance := suite.QueryBalances(account)
	return balance.AmountOf(coin.Denom).GTE(coin.Amount)
}

func (suite *Erc20TestSuite) TokenBalanceGTE(account common.Address, token common.Address, amount *big.Int) bool {
	balance := suite.BalanceOf(token, account)
	return balance.Cmp(amount) >= 0
}
