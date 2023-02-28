package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

func (suite *KeeperTestSuite) TestDelegate() {
	vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)

	del1 := suite.RandSigner()
	del1Amt := suite.app.BankKeeper.GetBalance(suite.ctx, del1.AccAddress(), fxtypes.DefaultDenom)

	val1 := vals[0]
	lpToken, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, val1.GetOperator())
	suite.Require().True(found)

	shares1, err := suite.app.StakingKeeper.Delegate(suite.ctx, del1.AccAddress(), del1Amt.Amount.Quo(sdkmath.NewInt(2)), stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	var res1 struct{ Value *big.Int }
	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken, fxtypes.GetLPToken().ABI, "balanceOf", &res1, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(shares1.BigInt(), res1.Value)

	shares2, err := suite.app.StakingKeeper.Delegate(suite.ctx, del1.AccAddress(), del1Amt.Amount.Quo(sdkmath.NewInt(2)), stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	var res2 struct{ Value *big.Int }
	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken, fxtypes.GetLPToken().ABI, "balanceOf", &res2, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(shares2.Add(shares1).BigInt(), res2.Value)

	suite.Require().Equal(shares1, shares2)
}

func (suite *KeeperTestSuite) TestUndelegate() {
	vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)

	del1 := suite.RandSigner()
	del1Amt := suite.app.BankKeeper.GetBalance(suite.ctx, del1.AccAddress(), fxtypes.DefaultDenom)

	val1 := vals[0]
	lpToken, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, val1.GetOperator())
	suite.Require().True(found)

	shares, err := suite.app.StakingKeeper.Delegate(suite.ctx, del1.AccAddress(), del1Amt.Amount, stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, del1.AccAddress(), val1.GetOperator(), shares.Quo(sdk.NewDec(2)))
	suite.Require().NoError(err)

	var res struct{ Value *big.Int }
	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(shares.Quo(sdk.NewDec(2)).BigInt(), res.Value)

	delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, del1.AccAddress(), val1.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(delegation.Shares.BigInt(), res.Value)

	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, del1.AccAddress(), val1.GetOperator(), shares.Quo(sdk.NewDec(2)))
	suite.Require().NoError(err)

	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(big.NewInt(0).String(), res.Value.String())

	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, del1.AccAddress(), val1.GetOperator())
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestRedelegate() {
	vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)

	del1 := suite.RandSigner()
	del1Amt := suite.app.BankKeeper.GetBalance(suite.ctx, del1.AccAddress(), fxtypes.DefaultDenom)

	val1 := vals[0]
	lpToken1, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, val1.GetOperator())
	suite.Require().True(found)

	val2 := vals[1]
	lpToken2, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, val2.GetOperator())
	suite.Require().True(found)

	shares, err := suite.app.StakingKeeper.Delegate(suite.ctx, del1.AccAddress(), del1Amt.Amount, stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	redelegateShares, err := suite.app.StakingKeeper.ValidateUnbondAmount(suite.ctx, del1.AccAddress(), val1.GetOperator(), del1Amt.Amount.Quo(sdkmath.NewInt(2)))
	suite.Require().NoError(err)

	_, err = suite.app.StakingKeeper.BeginRedelegation(suite.ctx, del1.AccAddress(), val1.GetOperator(), val2.GetOperator(), redelegateShares)
	suite.Require().NoError(err)

	var res struct{ Value *big.Int }
	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken1, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(shares.Quo(sdk.NewDec(2)).BigInt(), res.Value)

	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken2, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(shares.Quo(sdk.NewDec(2)).BigInt(), res.Value)
}
