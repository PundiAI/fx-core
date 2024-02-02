package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func (suite *KeeperTestSuite) TestKeeper_SetAccount() {
	address := helpers.GenerateAddress()
	suite.Nil(suite.app.AccountKeeper.GetAccount(suite.ctx, address.Bytes()))

	acc := suite.app.EvmKeeper.GetAccountOrEmpty(suite.ctx, address)
	suite.NotNil(acc)
	acc.CodeHash = common.BytesToHash([]byte{1, 2, 3}).Bytes()
	suite.NoError(suite.app.EvmKeeper.SetAccount(suite.ctx, address, acc))

	account := suite.app.AccountKeeper.GetAccount(suite.ctx, address.Bytes())
	ethAcc, ok := account.(ethermint.EthAccountI)
	suite.True(ok)
	suite.Equal(ethAcc.GetCodeHash(), common.BytesToHash(acc.CodeHash))
}
