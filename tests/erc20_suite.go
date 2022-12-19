package tests

import (
	"fmt"
	"math/big"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/client"
	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

type Erc20TestSuite struct {
	metadata banktypes.Metadata
	EvmTestSuite
}

func NewErc20TestSuite(ts *TestSuite) Erc20TestSuite {
	return Erc20TestSuite{
		metadata: banktypes.Metadata{
			Description: "description of the token",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "usdt",
					Exponent: uint32(0),
				}, {
					Denom:    "USDT",
					Exponent: uint32(6),
				},
			},
			Base:    "usdt",
			Display: "usdt",
			Name:    "Tether USD",
			Symbol:  "USDT",
		},
		EvmTestSuite: NewEvmTestSuite(ts),
	}
}

func (suite *Erc20TestSuite) TokenDecimals() uint32 {
	return suite.metadata.DenomUnits[1].Exponent
}

func (suite *Erc20TestSuite) ERC20Query() erc20types.QueryClient {
	return suite.GRPCClient().ERC20Query()
}

func (suite *Erc20TestSuite) CheckRegisterCoin(denom string) {
	aliasesResp, err := suite.ERC20Query().DenomAliases(suite.ctx, &erc20types.QueryDenomAliasesRequest{Denom: denom})
	suite.NoError(err)
	suite.T().Log("denom", denom, "alias", aliasesResp.Aliases)
	for _, alias := range aliasesResp.Aliases {
		aliasDenom, err := suite.ERC20Query().AliasDenom(suite.ctx, &erc20types.QueryAliasDenomRequest{Alias: alias})
		suite.NoError(err)
		suite.Equal(denom, aliasDenom.Denom)
	}
}

func (suite *Erc20TestSuite) TokenPair(denom string) erc20types.TokenPair {
	pairResp, err := suite.ERC20Query().TokenPair(suite.ctx, &erc20types.QueryTokenPairRequest{Token: denom})
	suite.NoError(err)
	return pairResp.TokenPair
}

func (suite *Erc20TestSuite) Erc20TokenAddress(denom string) common.Address {
	return suite.TokenPair(denom).GetERC20Contract()
}

func (suite *Erc20TestSuite) TransferCrossChain(privateKey cryptotypes.PrivKey, token common.Address, recipient string, amount, fee *big.Int, target string) common.Hash {
	suite.T().Log("transfer cross chain", target)
	pack, err := FIP20ABI.Pack("transferCrossChain", recipient, amount, fee, fxtypes.MustStrToByte32(target))
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)

	return ethTx.Hash()
}

func (suite *Erc20TestSuite) ConvertCoin(recipient common.Address, coin sdk.Coin) {
	msg := erc20types.NewMsgConvertCoin(coin, recipient, suite.AccAddress())
	suite.BroadcastTx(suite.privKey, msg)
}

func (suite *Erc20TestSuite) ConvertERC20(token common.Address, amount sdk.Int, recipient sdk.AccAddress) {
	msg := erc20types.NewMsgConvertERC20(amount, recipient, token, suite.HexAddress())
	suite.BroadcastTx(suite.privKey, msg)
}

func (suite *Erc20TestSuite) RegisterCoinProposal(md banktypes.Metadata) (proposalId uint64) {
	content := &erc20types.RegisterCoinProposal{
		Title:       fmt.Sprintf("register %s denom", md.Base),
		Description: "bar",
		Metadata:    md,
	}
	return suite.BroadcastProposalTx(content)
}

func (suite *Erc20TestSuite) ToggleTokenConversionProposal(denom string) (proposalId uint64) {
	content := &erc20types.ToggleTokenConversionProposal{
		Title:       fmt.Sprintf("update %s denom", denom),
		Description: "update",
		Token:       denom,
	}
	return suite.BroadcastProposalTx(content)
}

func (suite *Erc20TestSuite) UpdateDenomAliasProposal(denom, alias string) (proposalId uint64) {
	content := &erc20types.UpdateDenomAliasProposal{
		Title:       fmt.Sprintf("update %s denom %s alias", denom, alias),
		Description: "update",
		Denom:       denom,
		Alias:       alias,
	}
	return suite.BroadcastProposalTx(content)
}

func (suite *Erc20TestSuite) ConvertDenom(private cryptotypes.PrivKey, receiver sdk.AccAddress, coin sdk.Coin, target string) {
	suite.BroadcastTx(private, &erc20types.MsgConvertDenom{
		Sender:   sdk.AccAddress(private.PubKey().Address()).String(),
		Receiver: receiver.String(),
		Coin:     coin,
		Target:   target,
	})
}

func (suite *Erc20TestSuite) TransferERC20(privateKey cryptotypes.PrivKey, token, recipient common.Address, value *big.Int) common.Hash {
	suite.T().Logf("transfer erc20 %s to %s value %s\n", token, recipient.String(), value.String())

	pack, err := FIP20ABI.Pack("transfer", recipient, value)
	suite.Require().NoError(err)

	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)

	suite.SendTransaction(ethTx)
	return ethTx.Hash()
}
