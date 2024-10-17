package tests

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"

	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
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

func (suite *Erc20TestSuite) Erc20TokenAddress(denom string) common.Address {
	// todo: implement me
	return common.Address{}
}

func (suite *Erc20TestSuite) ToggleTokenConversionProposal(denom string) (*sdk.TxResponse, uint64) {
	msg := &erc20types.MsgToggleTokenConversion{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Token:     denom,
	}
	return suite.BroadcastProposalTx2([]sdk.Msg{msg}, "ToggleTokenConversionProposal", "ToggleTokenConversionProposal")
}

func (suite *Erc20TestSuite) ConvertCoin(recipient common.Address, coin sdk.Coin) *sdk.TxResponse {
	fromAddress := sdk.AccAddress(suite.privKey.PubKey().Address())

	beforeBalance := suite.QueryBalances(fromAddress).AmountOf(coin.Denom)
	erc20TokenAddress := suite.Erc20TokenAddress(coin.Denom)
	beforeBalanceOf := suite.BalanceOf(erc20TokenAddress, recipient)

	msg := erc20types.NewMsgConvertCoin(coin, recipient, sdk.AccAddress(suite.privKey.PubKey().Address()))
	txResponse := suite.BroadcastTx(suite.privKey, msg)

	afterBalance := suite.QueryBalances(fromAddress).AmountOf(coin.Denom)
	afterBalanceOf := suite.BalanceOf(erc20TokenAddress, recipient)
	suite.Require().Equal(beforeBalance.Sub(afterBalance).String(), coin.Amount.String())
	suite.Require().Equal(afterBalanceOf.String(), beforeBalanceOf.String())
	return txResponse
}
