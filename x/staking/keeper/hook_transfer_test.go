package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/x/staking/types"
)

func (suite *KeeperTestSuite) TestHookTransferNativeToken() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, delegate, lpToken common.Address, share sdk.Dec) (types.RelayTransfer, []string)
		result   bool
		error    func(args []string) string
	}{
		{
			name: "ok - transfer lp token",
			malleate: func(val sdk.ValAddress, delegate, lpToken common.Address, share sdk.Dec) (types.RelayTransfer, []string) {
				signer := suite.RandSigner()

				return types.RelayTransfer{
					From:          delegate,
					To:            signer.Address(),
					Amount:        share.TruncateInt().BigInt(),
					TokenContract: lpToken,
					Validator:     val,
				}, []string{}
			},
			result: true,
			error:  nil,
		},
		{
			name: "ok - transfer zero lp token",
			malleate: func(val sdk.ValAddress, delegate, lpToken common.Address, share sdk.Dec) (types.RelayTransfer, []string) {
				signer := suite.RandSigner()

				return types.RelayTransfer{
					From:          delegate,
					To:            signer.Address(),
					Amount:        big.NewInt(0),
					TokenContract: lpToken,
					Validator:     val,
				}, []string{}
			},
			result: true,
			error:  nil,
		},
		{
			name: "failed - transfer lp token",
			malleate: func(val sdk.ValAddress, delegate, lpToken common.Address, share sdk.Dec) (types.RelayTransfer, []string) {
				signer := suite.RandSigner()

				return types.RelayTransfer{
					From:          delegate,
					To:            signer.Address(),
					Amount:        share.Add(sdk.OneDec()).BigInt(),
					TokenContract: lpToken,
					Validator:     val,
				}, []string{share.String()}
			},
			result: false,
			error: func(args []string) string {
				return fmt.Sprintf("%s: not enough delegation shares", args[0])
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			validator, found := suite.app.StakingKeeper.GetValidatorByConsAddr(suite.ctx, suite.ctx.BlockHeader().ProposerAddress)
			suite.Require().True(found)

			signer, val, lpToken, _, share := suite.RandDelegates(validator)

			relay, errArgs := tc.malleate(val, signer.Address(), lpToken, share)

			beforeFromDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, relay.From.Bytes(), val)
			suite.Require().True(found)
			beforeToDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, relay.To.Bytes(), val)
			if !found {
				beforeToDel.Shares = sdk.ZeroDec()
			}

			suite.Commit()

			err := suite.app.StakingKeeper.EVMHooks().HookTransferEvent(suite.ctx, []types.RelayTransfer{relay})
			if tc.result {
				suite.Require().NoError(err)

				afterFromDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, relay.From.Bytes(), val)
				if !found {
					afterFromDel.Shares = sdk.ZeroDec()
				}
				afterToDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, relay.To.Bytes(), val)
				if !found {
					afterToDel.Shares = sdk.ZeroDec()
				}

				suite.Require().Equal(beforeFromDel.Shares.Sub(afterFromDel.Shares).String(), afterToDel.Shares.Sub(beforeToDel.Shares).String())
			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}
