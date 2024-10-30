package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
)

func (suite *StakingPrecompileTestSuite) TestRedelegate() {
	testCases := []struct {
		name     string
		malleate func(valSrc, valDst sdk.ValAddress, shares sdkmath.LegacyDec, delAmount sdkmath.Int) (contract.RedelegateV2Args, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok v2",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdkmath.LegacyDec, delAmount sdkmath.Int) (contract.RedelegateV2Args, error) {
				return contract.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - v2 invalid validator src",
			malleate: func(_, valDst sdk.ValAddress, shares sdkmath.LegacyDec, delAmount sdkmath.Int) (contract.RedelegateV2Args, error) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				return contract.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			operator0 := suite.GetFirstValAddr()
			operator1 := suite.GetSecondValAddr()
			delAmt := helpers.NewRandAmount()
			delAddr := suite.GetDelAddr()

			res := suite.DelegateV2(suite.Ctx, contract.DelegateV2Args{
				Validator: operator0.String(),
				Amount:    delAmt.BigInt(),
			})
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			delegation0 := suite.GetDelegation(delAddr.Bytes(), operator0)
			_, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator1)
			suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

			val0 := suite.GetValidator(operator0)

			args, expectErr := tc.malleate(operator0, operator1, delegation0.Shares, delAmt)

			res = suite.WithError(expectErr).RedelegateV2(suite.Ctx, args)
			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator0)
				suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)
				delegation1New, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator1)
				suite.Require().NoError(err)
				suite.Require().Equal(delegation0.Shares, delegation1New.Shares)

				redelegations, err := suite.App.StakingKeeper.GetAllRedelegations(suite.Ctx, delAddr.Bytes(), operator0, operator1)
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(redelegations))

				suite.CheckRedelegateLogs(res.Logs, delAddr, val0.GetOperator(), operator1.String(),
					delegation0.Shares.TruncateInt().BigInt(), val0.TokensFromShares(delegation0.Shares).TruncateInt().BigInt(),
					redelegations[0].Entries[0].CompletionTime.Unix())

				suite.CheckRedelegateEvents(suite.Ctx, val0.GetOperator(), operator1.String(),
					val0.TokensFromShares(delegation0.Shares).TruncateInt().BigInt(),
					redelegations[0].Entries[0].CompletionTime)
			}
		})
	}
}

func (suite *StakingPrecompileTestSuite) CheckRedelegateLogs(logs []*evmtypes.Log, delAddr common.Address, valSrc, valDst string, shares, amount *big.Int, completionTime int64) {
	redelegateV2ABI := precompile.NewRedelegateV2ABI()
	existLog := false
	for _, log := range logs {
		if log.Topics[0] == redelegateV2ABI.Event.ID.String() {
			suite.Require().Equal(log.Address, contract.StakingAddress)
			event, err := redelegateV2ABI.UnpackEvent(log.ToEthereum())
			suite.Require().NoError(err)
			suite.Require().Equal(event.Sender, delAddr)
			suite.Require().Equal(event.ValSrc, valSrc)
			suite.Require().Equal(event.ValDst, valDst)
			suite.Require().Equal(event.Amount.String(), amount.String())
			suite.Require().Equal(event.CompletionTime.Int64(), completionTime)
			existLog = true
		}
	}
	suite.Require().True(existLog)
}

func (suite *StakingPrecompileTestSuite) CheckRedelegateEvents(ctx sdk.Context, valSrc, valDst string, amount *big.Int, completionTime time.Time) {
	existEvent := false
	for _, event := range ctx.EventManager().Events() {
		if event.Type != stakingtypes.EventTypeRedelegate {
			continue
		}
		for _, attr := range event.Attributes {
			if attr.Key == stakingtypes.AttributeKeySrcValidator {
				suite.Require().Equal(attr.Value, valSrc)
			}
			if attr.Key == stakingtypes.AttributeKeyDstValidator {
				suite.Require().Equal(attr.Value, valDst)
			}
			if attr.Key == sdk.AttributeKeyAmount {
				suite.Require().Equal(strings.TrimSuffix(attr.Value, fxtypes.DefaultDenom), amount.String())
			}
			if attr.Key == stakingtypes.AttributeKeyCompletionTime {
				suite.Require().Equal(attr.Value, completionTime.Format(time.RFC3339))
			}
		}
		existEvent = true
	}
	suite.Require().True(existEvent)
}
