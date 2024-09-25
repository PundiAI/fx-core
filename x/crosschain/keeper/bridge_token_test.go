package keeper_test

import (
	"fmt"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeTokenCase() {
	// test 1. add bridge token
	// test 2. get bridgeDenom by tokenContract
	// test 3. get tokenContract by bridgeDenom
	testCases := []struct {
		name                     string
		claim                    *types.MsgBridgeTokenClaim
		pass                     bool
		getBridgeDenomByContract map[string]string
		getContractByBridgeDenom map[string]string
		allBridgeTokenLen        int
	}{
		{
			name: "success with FX symbol",
			claim: &types.MsgBridgeTokenClaim{
				Symbol:        fxtypes.DefaultDenom,
				TokenContract: "0x1",
				Decimals:      fxtypes.DenomUnit,
			},
			pass: true,
			getBridgeDenomByContract: map[string]string{
				"0x1": fxtypes.DefaultDenom,
			},
			getContractByBridgeDenom: map[string]string{
				fmt.Sprintf("%s%s", suite.chainName, "0x1"): "0x1",
				fxtypes.DefaultDenom:                        "0x1",
			},
			allBridgeTokenLen: 1,
		},
		{
			name: "success with custom symbol",
			claim: &types.MsgBridgeTokenClaim{
				TokenContract: "0x1",
			},
			pass: true,
			getBridgeDenomByContract: map[string]string{
				"0x1": fmt.Sprintf("%s%s", suite.chainName, "0x1"),
			},
			getContractByBridgeDenom: map[string]string{
				fmt.Sprintf("%s%s", suite.chainName, "0x1"): "0x1",
			},
			allBridgeTokenLen: 1,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := suite.Keeper().AddBridgeTokenExecuted(suite.Ctx, tc.claim)
			if !tc.pass {
				suite.Require().Error(err)
				return
			}

			suite.Require().NoError(err)

			for tokenContract, bridgeDenom := range tc.getBridgeDenomByContract {
				actualBridgeDenom, found := suite.Keeper().GetBridgeDenomByContract(suite.Ctx, tokenContract)
				suite.Require().True(found)
				suite.Require().Equal(bridgeDenom, actualBridgeDenom)
			}

			for bridgeDenom, tokenContract := range tc.getContractByBridgeDenom {
				actualTokenContract, found := suite.Keeper().GetContractByBridgeDenom(suite.Ctx, bridgeDenom)
				suite.Require().True(found)
				suite.Require().Equal(tokenContract, actualTokenContract)
			}

			bridgeTokenLen := 0
			suite.Keeper().IteratorBridgeDenomWithContract(suite.Ctx, func(token *types.BridgeToken) bool {
				bridgeTokenLen++
				return false
			})
			suite.Require().EqualValues(tc.allBridgeTokenLen, bridgeTokenLen)
		})
	}
}
