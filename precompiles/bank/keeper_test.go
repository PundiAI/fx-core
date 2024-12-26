package bank_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	bsctypes "github.com/pundiai/fx-core/v8/x/bsc/types"
)

func (suite *BankPrecompileTestSuite) TestTransferFromModuleToAccount() {
	moduleName := bsctypes.ModuleName
	coin := sdk.NewCoin(helpers.NewRandDenom(), helpers.NewRandAmount())
	suite.MintTokenToModule(moduleName, coin)

	token := helpers.GenHexAddress()
	suite.SetErc20Token(coin.Denom, token)

	account := suite.NewSigner()
	err := suite.keeper.TransferFromModuleToAccount(suite.Ctx, &contract.TransferFromModuleToAccountArgs{
		Module:  moduleName,
		Account: account.Address(),
		Token:   token,
		Amount:  coin.Amount.BigInt(),
	})
	suite.Require().NoError(err)

	suite.AssertBalance(account.AccAddress(), coin)
	suite.AssertAllBalance(types.NewModuleAddress(moduleName))
}

func (suite *BankPrecompileTestSuite) TestTransferFromAccountToModule() {
	moduleName := bsctypes.ModuleName
	coin := sdk.NewCoin(helpers.NewRandDenom(), helpers.NewRandAmount())

	token := helpers.GenHexAddress()
	suite.SetErc20Token(coin.Denom, token)

	account := suite.NewSigner()
	suite.MintToken(account.AccAddress(), coin)

	err := suite.keeper.TransferFromAccountToModule(suite.Ctx, &contract.TransferFromAccountToModuleArgs{
		Account: account.Address(),
		Module:  moduleName,
		Token:   token,
		Amount:  coin.Amount.BigInt(),
	})
	suite.Require().NoError(err)

	suite.AssertBalance(account.AccAddress())
	suite.AssertAllBalance(types.NewModuleAddress(moduleName), coin)
}
