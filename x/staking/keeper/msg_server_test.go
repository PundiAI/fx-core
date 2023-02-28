package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/staking/keeper"
)

func (suite *KeeperTestSuite) TestMsgCreateValidator() {
	valNew := suite.RandSigner()
	valCosKey := tmed25519.GenPrivKey()
	pk, err := cryptocodec.FromTmPubKeyInterface(valCosKey.PubKey())
	suite.Require().NoError(err)
	pkAny, err := codectypes.NewAnyWithValue(pk)
	suite.Require().NoError(err)

	_, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, valNew.AccAddress().Bytes())
	suite.Require().False(found)

	msgServer := keeper.NewMsgServerImpl(suite.app.StakingKeeper)
	_, err = msgServer.CreateValidator(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgCreateValidator{
		Description: stakingtypes.Description{
			Moniker:         "val-test",
			Identity:        "val-test",
			Website:         "val-test",
			SecurityContact: "val-test",
			Details:         "val-test",
		},
		Commission: stakingtypes.CommissionRates{
			Rate:          sdk.MustNewDecFromStr("0.20"),
			MaxRate:       sdk.OneDec(),
			MaxChangeRate: sdk.MustNewDecFromStr("0.10"),
		},
		MinSelfDelegation: sdkmath.NewInt(100).Mul(sdkmath.NewInt(1e18)),
		DelegatorAddress:  valNew.AccAddress().String(),
		ValidatorAddress:  sdk.ValAddress(valNew.AccAddress()).String(),
		Pubkey:            pkAny,
		Value:             sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100).Mul(sdkmath.NewInt(1e18))),
	})
	suite.Require().NoError(err)

	lpToken, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, valNew.AccAddress().Bytes())
	suite.Require().True(found)

	delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, valNew.AccAddress().Bytes(), valNew.AccAddress().Bytes())
	suite.Require().True(found)

	var res struct{ Value *big.Int }
	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken, fxtypes.GetLPToken().ABI, "balanceOf", &res, valNew.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(delegation.Shares.BigInt(), res.Value)
}

func (suite *KeeperTestSuite) TestMsgDelegate() {
	vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)

	del1 := suite.RandSigner()
	del1Amt := suite.app.BankKeeper.GetBalance(suite.ctx, del1.AccAddress(), fxtypes.DefaultDenom)

	val1 := vals[0]
	lpToken, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, val1.GetOperator())
	suite.Require().True(found)

	var res struct{ Value *big.Int }
	err := suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(big.NewInt(0).String(), res.Value.String())

	msgServer := keeper.NewMsgServerImpl(suite.app.StakingKeeper)
	_, err = msgServer.Delegate(sdk.WrapSDKContext(suite.ctx), stakingtypes.NewMsgDelegate(del1.AccAddress(), val1.GetOperator(), del1Amt))
	suite.Require().NoError(err)

	delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, del1.AccAddress(), val1.GetOperator())
	suite.Require().True(found)

	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(delegation.Shares.BigInt(), res.Value)
	suite.Require().True(res.Value.Cmp(big.NewInt(0)) == 1)
}

func (suite *KeeperTestSuite) TestMsgBeginRedelegate() {
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

	msgServer := keeper.NewMsgServerImpl(suite.app.StakingKeeper)
	_, err = msgServer.BeginRedelegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    del1.AccAddress().String(),
		ValidatorSrcAddress: val1.OperatorAddress,
		ValidatorDstAddress: val2.OperatorAddress,
		Amount:              sdk.NewCoin(del1Amt.Denom, del1Amt.Amount.Quo(sdkmath.NewInt(2))),
	})
	suite.Require().NoError(err)

	delegationVal1, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, del1.AccAddress(), val1.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(shares.Quo(sdk.NewDec(2)), delegationVal1.Shares)

	delegationVal2, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, del1.AccAddress(), val2.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(delegationVal1.Shares, delegationVal2.Shares)

	var res struct{ Value *big.Int }
	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken1, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(delegationVal1.Shares.BigInt(), res.Value)

	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken2, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(delegationVal2.Shares.BigInt(), res.Value)
}

func (suite *KeeperTestSuite) TestMsgUndelegate() {
	vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)

	del1 := suite.RandSigner()
	del1Amt := suite.app.BankKeeper.GetBalance(suite.ctx, del1.AccAddress(), fxtypes.DefaultDenom)

	val1 := vals[0]
	lpToken1, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, val1.GetOperator())
	suite.Require().True(found)

	shares, err := suite.app.StakingKeeper.Delegate(suite.ctx, del1.AccAddress(), del1Amt.Amount, stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	msgServer := keeper.NewMsgServerImpl(suite.app.StakingKeeper)

	_, err = msgServer.Undelegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgUndelegate{
		DelegatorAddress: del1.AccAddress().String(),
		ValidatorAddress: val1.OperatorAddress,
		Amount:           sdk.NewCoin(del1Amt.Denom, del1Amt.Amount.Quo(sdkmath.NewInt(2))),
	})
	suite.Require().NoError(err)

	delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, del1.AccAddress(), val1.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(shares.Quo(sdk.NewDec(2)), delegation.Shares)

	var res struct{ Value *big.Int }
	err = suite.app.EvmKeeper.CallContract(suite.ctx, suite.app.StakingKeeper.GetLPTokenModuleAddress(), lpToken1, fxtypes.GetLPToken().ABI, "balanceOf", &res, del1.Address())
	suite.Require().NoError(err)
	suite.Require().Equal(delegation.Shares.BigInt(), res.Value)

	ubd, found := suite.app.StakingKeeper.GetUnbondingDelegation(suite.ctx, del1.AccAddress(), val1.GetOperator())
	suite.Require().True(found)
	suite.Require().True(len(ubd.Entries) == 1)
	suite.Require().Equal(ubd.Entries[0].Balance, del1Amt.Amount.Quo(sdkmath.NewInt(2)))
}
