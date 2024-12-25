package staking_test

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

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/staking"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func (suite *StakingPrecompileTestSuite) TestUndelegate() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdkmath.LegacyDec, delAmt sdkmath.Int) (contract.UndelegateV2Args, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok v2",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec, delAmt sdkmath.Int) (contract.UndelegateV2Args, error) {
				return contract.UndelegateV2Args{
					Validator: val.String(),
					Amount:    delAmt.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - v2 invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec, delAmt sdkmath.Int) (contract.UndelegateV2Args, error) {
				newVal := val.String() + "1"
				return contract.UndelegateV2Args{
					Validator: newVal,
					Amount:    delAmt.BigInt(),
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			operator := suite.GetFirstValAddr()
			delAmt := helpers.NewRandAmount()
			delAddr := suite.GetDelAddr()

			res := suite.DelegateV2(suite.Ctx, contract.DelegateV2Args{
				Validator: operator.String(),
				Amount:    delAmt.BigInt(),
			})
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			delegation := suite.GetDelegation(delAddr.Bytes(), operator)

			undelegations, err := suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, delAddr.Bytes())
			suite.Require().NoError(err)
			suite.Require().Equal(0, len(undelegations))

			args, expectErr := tc.malleate(operator, delegation.Shares, delAmt)

			res = suite.WithError(expectErr).UndelegateV2(suite.Ctx, args)
			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				undelegations, err := suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, delAddr.Bytes())
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(undelegations))
				suite.Require().Equal(1, len(undelegations[0].Entries))
				suite.Require().Equal(sdk.AccAddress(delAddr.Bytes()).String(), undelegations[0].DelegatorAddress)
				suite.Require().Equal(operator.String(), undelegations[0].ValidatorAddress)
				suite.Require().Equal(delAmt, undelegations[0].Entries[0].Balance)

				suite.CheckUndelegateLogs(res.Logs, delAddr, operator.String(), delegation.Shares.TruncateInt().BigInt(),
					undelegations[0].Entries[0].Balance.BigInt(), undelegations[0].Entries[0].CompletionTime)

				suite.CheckUndelegateEvents(suite.Ctx, operator.String(), undelegations[0].Entries[0].Balance.BigInt(),
					undelegations[0].Entries[0].CompletionTime)
			}
		})
	}
}

func (suite *StakingPrecompileTestSuite) CheckUndelegateLogs(logs []*evmtypes.Log, delAddr common.Address, valAddr string, shares, amount *big.Int, completionTime time.Time) {
	undelegateV2ABI := staking.NewUndelegateV2ABI()
	existLog := false
	for _, log := range logs {
		if log.Topics[0] == undelegateV2ABI.Event.ID.String() {
			suite.Require().Equal(log.Address, contract.StakingAddress)
			event, err := undelegateV2ABI.UnpackEvent(log.ToEthereum())
			suite.Require().NoError(err)
			suite.Require().Equal(event.Sender, delAddr)
			suite.Require().Equal(event.Validator, valAddr)
			suite.Require().Equal(event.Amount.String(), amount.String())
			suite.Require().Equal(event.CompletionTime.Int64(), completionTime.Unix())
			existLog = true
		}
	}
	suite.Require().True(existLog)
}

func (suite *StakingPrecompileTestSuite) CheckUndelegateEvents(ctx sdk.Context, valAddr string, amount *big.Int, completionTime time.Time) {
	existEvent := false
	for _, event := range ctx.EventManager().Events() {
		if event.Type == stakingtypes.EventTypeUnbond {
			for _, attr := range event.Attributes {
				if attr.Key == stakingtypes.AttributeKeyValidator {
					suite.Require().Equal(attr.Value, valAddr)
					existEvent = true
				}
				if attr.Key == sdk.AttributeKeyAmount {
					suite.Require().Equal(strings.TrimSuffix(attr.Value, fxtypes.DefaultDenom), amount.String())
					existEvent = true
				}
				if attr.Key == stakingtypes.AttributeKeyCompletionTime {
					suite.Require().Equal(attr.Value, completionTime.Format(time.RFC3339))
					existEvent = true
				}
			}
		}
	}
	suite.Require().True(existEvent)
}
