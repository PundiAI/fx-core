package tests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"

	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

type Erc20TestSuite struct {
	EvmTestSuite
}

func NewErc20TestSuite(ts *TestSuite) Erc20TestSuite {
	return Erc20TestSuite{
		EvmTestSuite: NewEvmTestSuite(ts),
	}
}

func (suite *Erc20TestSuite) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.privKey.PubKey().Address())
}

func (suite *Erc20TestSuite) HexAddress() common.Address {
	return common.BytesToAddress(suite.privKey.PubKey().Address())
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

func (suite *Erc20TestSuite) TokenPairs() []erc20types.TokenPair {
	pairsResp, err := suite.ERC20Query().TokenPairs(suite.ctx, &erc20types.QueryTokenPairsRequest{})
	suite.NoError(err)
	return pairsResp.TokenPairs
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
	msg := &erc20types.MsgRegisterCoin{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Metadata:  md,
	}
	return suite.BroadcastProposalTx2([]sdk.Msg{msg}, "RegisterCoinProposal", "RegisterCoinProposal")
}

func (suite *Erc20TestSuite) RegisterErc20Proposal(erc20Addr string, aliases []string) (*sdk.TxResponse, uint64) {
	msg := &erc20types.MsgRegisterERC20{
		Authority:    authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Erc20Address: erc20Addr,
		Aliases:      aliases,
	}
	return suite.BroadcastProposalTx2([]sdk.Msg{msg}, "RegisterErc20Proposal", "RegisterErc20Proposal")
}

func (suite *Erc20TestSuite) ToggleTokenConversionProposal(denom string) (*sdk.TxResponse, uint64) {
	msg := &erc20types.MsgToggleTokenConversion{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Token:     denom,
	}
	return suite.BroadcastProposalTx2([]sdk.Msg{msg}, "ToggleTokenConversionProposal", "ToggleTokenConversionProposal")
}

func (suite *Erc20TestSuite) UpdateDenomAliasProposal(denom, alias string) (*sdk.TxResponse, uint64) {
	msg := &erc20types.MsgUpdateDenomAlias{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Denom:     denom,
		Alias:     alias,
	}
	return suite.BroadcastProposalTx2([]sdk.Msg{msg}, "UpdateDenomAliasProposal", "UpdateDenomAliasProposal")
}

func (suite *Erc20TestSuite) ConvertCoin(recipient common.Address, coin sdk.Coin) *sdk.TxResponse {
	private := suite.privKey
	fromAddress := sdk.AccAddress(private.PubKey().Address())
	beforeBalance := suite.QueryBalances(fromAddress).AmountOf(coin.Denom)
	beforeBalanceOf := suite.BalanceOf(suite.Erc20TokenAddress(coin.Denom), recipient)
	msg := erc20types.NewMsgConvertCoin(coin, recipient, sdk.AccAddress(private.PubKey().Address()))
	txResponse := suite.BroadcastTx(private, msg)
	afterBalance := suite.QueryBalances(fromAddress).AmountOf(coin.Denom)
	afterBalanceOf := suite.BalanceOf(suite.Erc20TokenAddress(coin.Denom), recipient)
	suite.Require().True(beforeBalance.Sub(afterBalance).Equal(coin.Amount))
	suite.Require().True(new(big.Int).Sub(afterBalanceOf, beforeBalanceOf).Cmp(coin.Amount.BigInt()) == 0)
	return txResponse
}

func (suite *Erc20TestSuite) ConvertERC20(token common.Address, amount sdkmath.Int, recipient sdk.AccAddress) *sdk.TxResponse {
	private := suite.privKey
	beforeBalance := suite.QueryBalances(recipient).AmountOf(suite.DenomFromErc20(token))
	beforeBalanceOf := suite.BalanceOf(token, common.BytesToAddress(private.PubKey().Address().Bytes()))
	msg := erc20types.NewMsgConvertERC20(amount, recipient, token, common.BytesToAddress(private.PubKey().Address().Bytes()))
	txResponse := suite.BroadcastTx(private, msg)
	afterBalance := suite.QueryBalances(recipient).AmountOf(suite.DenomFromErc20(token))
	afterBalanceOf := suite.BalanceOf(token, common.BytesToAddress(private.PubKey().Address().Bytes()))
	suite.Require().True(afterBalance.Sub(beforeBalance).Equal(amount))
	suite.Require().True(new(big.Int).Sub(beforeBalanceOf, afterBalanceOf).Cmp(amount.BigInt()) == 0)
	return txResponse
}

func (suite *Erc20TestSuite) ConvertDenom(receiver sdk.AccAddress, coin sdk.Coin, target string) *sdk.TxResponse {
	private := suite.privKey
	fromAddress := sdk.AccAddress(private.PubKey().Address())
	beforeBalance := suite.QueryBalances(fromAddress).AmountOf(coin.Denom)
	txResponse := suite.BroadcastTx(private, &erc20types.MsgConvertDenom{
		Sender:   fromAddress.String(),
		Receiver: receiver.String(),
		Coin:     coin,
		Target:   target,
	})
	afterBalance := suite.QueryBalances(fromAddress).AmountOf(coin.Denom)
	suite.Require().True(beforeBalance.Sub(afterBalance).Equal(coin.Amount))
	return txResponse
}
