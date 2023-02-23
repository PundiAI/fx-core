package tests

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/testutil/helpers"
)

func (suite *IntegrationTest) StakingTest() {
	vals := suite.staking.Validators()

	lpTokens := make([]common.Address, 0, len(vals))
	for _, val := range vals {
		lpToken := suite.staking.ValidatorLPToken(val.GetOperator())
		lpTokens = append(lpTokens, lpToken)
	}

	val1, lpToken1 := vals[0], lpTokens[0]

	del1 := helpers.NewEthPrivKey()
	suite.Send(sdk.AccAddress(del1.PubKey().Address()), suite.NewCoin(sdk.NewInt(10_000).Mul(sdk.NewInt(1e18))))

	del1Amt := suite.NewCoin(sdk.NewInt(4_000).Mul(sdk.NewInt(1e18)))
	suite.staking.Delegate(del1, val1.GetOperator(), del1Amt)

	delegation1, del1Coin := suite.staking.GetDelegation(sdk.AccAddress(del1.PubKey().Address()), val1.GetOperator())
	suite.Require().Equal(del1Amt, del1Coin)

	val1Before := suite.staking.Validator(val1.GetOperator())

	shares := suite.staking.BalanceOf(lpToken1, common.BytesToAddress(del1.PubKey().Address().Bytes()))
	suite.Require().Equal(delegation1.Shares.BigInt(), shares)

	del2 := helpers.NewEthPrivKey()
	suite.Send(sdk.AccAddress(del2.PubKey().Address()), suite.NewCoin(sdk.NewInt(10_000).Mul(sdk.NewInt(1e18))))

	suite.staking.LPTokenTransfer(del1, lpToken1, common.BytesToAddress(del2.PubKey().Address().Bytes()), shares)

	delegation2, del2Coin := suite.staking.GetDelegation(sdk.AccAddress(del2.PubKey().Address()), val1.GetOperator())
	suite.Require().Equal(delegation2.Shares, delegation1.Shares)
	suite.Require().Equal(del2Coin, del1Coin)

	val1After := suite.staking.Validator(val1.GetOperator())
	suite.Require().Equal(val1Before.GetTokens(), val1After.GetTokens())
	suite.Require().Equal(val1Before.GetDelegatorShares(), val1After.GetDelegatorShares())
	suite.Require().Equal(val1After.DelegatorShares.Sub(val1.DelegatorShares).BigInt(), shares)

	shares = suite.staking.BalanceOf(lpToken1, common.BytesToAddress(del1.PubKey().Address().Bytes()))
	suite.Require().Equal(shares.String(), big.NewInt(0).String())

	shares = suite.staking.BalanceOf(lpToken1, common.BytesToAddress(del2.PubKey().Address().Bytes()))
	suite.Require().Equal(shares, delegation2.Shares.BigInt())
}
