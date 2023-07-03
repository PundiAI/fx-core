package ante_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v5/ante"
	"github.com/functionx/fx-core/v5/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v5/types"
)

func (suite *AnteTestSuite) TestRejectValidatorGranted() {
	val := helpers.NewEthPrivKey()
	valAddr := sdk.ValAddress(val.PubKey().Address())
	addr := sdk.AccAddress(valAddr)

	grantAcct := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	suite.app.StakingKeeper.UpdateValidatorOperator(suite.ctx, valAddr, grantAcct)

	testCases := []struct {
		name       string
		malleate   func() sdk.Tx
		expectPass bool
	}{
		{
			name: "success",
			malleate: func() sdk.Tx {
				testMsg := banktypes.MsgSend{
					FromAddress: grantAcct.String(),
					ToAddress:   sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: fxtypes.DefaultDenom}},
				}
				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), fxtypes.DefaultDenom, &testMsg)
				return txBuilder.GetTx()
			},
			expectPass: true,
		},
		{
			name: "fail",
			malleate: func() sdk.Tx {
				testMsg := banktypes.MsgSend{
					FromAddress: addr.String(),
					ToAddress:   sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: fxtypes.DefaultDenom}},
				}
				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), fxtypes.DefaultDenom, &testMsg)
				return txBuilder.GetTx()
			},
			expectPass: false,
		},
		{
			name: "fail with multiple msgs",
			malleate: func() sdk.Tx {
				msg1 := banktypes.MsgSend{
					FromAddress: addr.String(),
					ToAddress:   sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: fxtypes.DefaultDenom}},
				}
				msg2 := banktypes.MsgSend{
					FromAddress: addr.String(),
					ToAddress:   sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: fxtypes.DefaultDenom}},
				}
				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), fxtypes.DefaultDenom, &msg1, &msg2)
				return txBuilder.GetTx()
			},
			expectPass: false,
		},
		{
			name: "failed eth tx",
			malleate: func() sdk.Tx {
				from := common.BytesToAddress(addr.Bytes())
				to := helpers.GenerateAddress()
				emptyAccessList := ethtypes.AccessList{}
				msg := suite.BuildTestEthTx(from, to, nil, make([]byte, 0), big.NewInt(0), big.NewInt(100), big.NewInt(50), &emptyAccessList)
				return suite.CreateTestTx(msg, val, 1, false)
			},
			expectPass: false,
		},
	}

	dec := ante.NewRejectValidatorGrantedDecorator(suite.app.StakingKeeper)
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			tx := testCase.malleate()
			_, err := dec.AnteHandle(suite.ctx, tx, false, NextFn)
			if testCase.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
