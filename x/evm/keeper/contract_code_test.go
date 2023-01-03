package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

func (suite *KeeperTestSuite) TestKeeper_CreateContractWithCode() {
	// CreateContractWithCode is in init genesis
	account := suite.app.EvmKeeper.GetAccount(suite.ctx, fxtypes.GetWFX().Address)
	suite.NotNil(account)

	suite.Equal(uint64(0), account.Nonce)
	code := suite.app.EvmKeeper.GetCode(suite.ctx, common.BytesToHash(account.CodeHash))
	suite.Equal(fxtypes.GetWFX().Code, code)
}

func (suite *KeeperTestSuite) TestKeeper_UpdateContractCode() {
	updateCode := []byte{1, 2, 3}
	err := suite.app.EvmKeeper.UpdateContractCode(suite.ctx, fxtypes.GetWFX().Address, updateCode)
	suite.NoError(err)

	account := suite.app.EvmKeeper.GetAccount(suite.ctx, fxtypes.GetWFX().Address)
	suite.NotNil(account)

	suite.Equal(uint64(0), account.Nonce)
	code := suite.app.EvmKeeper.GetCode(suite.ctx, common.BytesToHash(account.CodeHash))
	suite.Equal(updateCode, code)
}
