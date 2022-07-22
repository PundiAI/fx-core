package tests

import (
	"context"
	"fmt"
	"math/big"
	"os"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/functionx/fx-core/v2/app/helpers"
	fxtypes "github.com/functionx/fx-core/v2/types"
	erc20types "github.com/functionx/fx-core/v2/x/erc20/types"
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
	err := os.Setenv("GO_ENV", "testing")
	suite.NoError(err)
	fxtypes.SetTestingManyToOneBlock(func() int64 { return 5 })

	suite.TestSuite.SetupSuite()
	suite.Send(suite.Address(), helpers.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
}

func (suite *ERC20TestSuite) Address() sdk.AccAddress {
	return suite.privKey.PubKey().Address().Bytes()
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
		sdk.NewCoins(helpers.NewCoin(sdk.NewInt(10_000).MulRaw(1e18))),
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

func (suite *ERC20TestSuite) EthBalance(address common.Address) *big.Int {
	amount, err := suite.EthClient().BalanceAt(context.Background(), address, nil)
	suite.NoError(err)
	return amount
}

func (suite *ERC20TestSuite) BalanceOf(contract, address common.Address) *big.Int {
	caller, err := NewERC20TokenCaller(contract, suite.EthClient())
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
