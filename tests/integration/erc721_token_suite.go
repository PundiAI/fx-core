package integration

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/pundiai/fx-core/v8/tests/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
)

type ERC721TokenSuite struct {
	*EthSuite

	erc721Token *contract.ERC721TokenTest
	signer      *helpers.Signer
}

func NewERC721TokenSuite(suite *EthSuite, token common.Address, signer *helpers.Signer) *ERC721TokenSuite {
	erc20Token, err := contract.NewERC721TokenTest(token, suite.ethCli)
	suite.Require().NoError(err)
	return &ERC721TokenSuite{
		EthSuite:    suite,
		erc721Token: erc20Token,
		signer:      signer,
	}
}

func (suite *ERC721TokenSuite) WithSigner(signer *helpers.Signer) *ERC721TokenSuite {
	return &ERC721TokenSuite{
		EthSuite:    suite.EthSuite,
		erc721Token: suite.erc721Token,
		signer:      signer,
	}
}

func (suite *ERC721TokenSuite) BalanceOf(account common.Address) *big.Int {
	balanceOf, err := suite.erc721Token.BalanceOf(nil, account)
	suite.Require().NoError(err)
	return balanceOf
}

func (suite *ERC721TokenSuite) EqualBalanceOf(account common.Address, value *big.Int) {
	suite.Require().Equal(suite.BalanceOf(account).Cmp(value), 0)
}

func (suite *ERC721TokenSuite) TokenURI(id *big.Int) string {
	uri, err := suite.erc721Token.TokenURI(nil, id)
	suite.Require().NoError(err)
	return uri
}

func (suite *ERC721TokenSuite) IsApprovedForAll(owner, operator common.Address) bool {
	isApproved, err := suite.erc721Token.IsApprovedForAll(nil, owner, operator)
	suite.Require().NoError(err)
	return isApproved
}

func (suite *ERC721TokenSuite) SafeMint(account common.Address) *ethtypes.Transaction {
	ethTx, err := suite.erc721Token.SafeMint(suite.TransactOpts(suite.signer), account, "ipfs://test-url")
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)
	return ethTx
}

func (suite *ERC721TokenSuite) Approve(operator common.Address, id *big.Int) *ethtypes.Transaction {
	ethTx, err := suite.erc721Token.Approve(suite.TransactOpts(suite.signer), operator, id)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)
	return ethTx
}

func (suite *ERC721TokenSuite) SetApprovalForAll(operator common.Address, approved bool) *ethtypes.Transaction {
	ethTx, err := suite.erc721Token.SetApprovalForAll(suite.TransactOpts(suite.signer), operator, approved)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)
	return ethTx
}

func (suite *ERC721TokenSuite) SafeTransferFrom(from, to common.Address, id *big.Int) *ethtypes.Transaction {
	ethTx, err := suite.erc721Token.SafeTransferFrom(suite.TransactOpts(suite.signer), from, to, id)
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)
	return ethTx
}
