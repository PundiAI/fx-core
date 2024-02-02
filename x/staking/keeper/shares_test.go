package keeper_test

import (
	"math/big"

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func (suite *KeeperTestSuite) TestAllowance() {
	vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
	val := vals[0]

	spender := helpers.NewSigner(helpers.NewEthPrivKey())
	allowance := suite.app.StakingKeeper.GetAllowance(suite.ctx, val.GetOperator(), suite.signer.AccAddress(), spender.AccAddress())
	suite.Equal(0, allowance.Cmp(big.NewInt(0)))

	suite.app.StakingKeeper.SetAllowance(suite.ctx, val.GetOperator(), suite.signer.AccAddress(), spender.AccAddress(), big.NewInt(100))

	allowance = suite.app.StakingKeeper.GetAllowance(suite.ctx, val.GetOperator(), suite.signer.AccAddress(), spender.AccAddress())
	suite.Equal(0, allowance.Cmp(big.NewInt(100)))
}
