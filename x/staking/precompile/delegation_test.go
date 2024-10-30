package precompile_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
)

func TestStakingDelegationABI(t *testing.T) {
	delegationABI := precompile.NewDelegationABI()

	require.Len(t, delegationABI.Method.Inputs, 2)
	require.Len(t, delegationABI.Method.Outputs, 2)
}

func (suite *StakingPrecompileTestSuite) TestDelegation() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, del common.Address) (contract.DelegationArgs, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, del common.Address) (contract.DelegationArgs, error) {
				return contract.DelegationArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "ok - zero",
			malleate: func(val sdk.ValAddress, del common.Address) (contract.DelegationArgs, error) {
				return contract.DelegationArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) (contract.DelegationArgs, error) {
				newVal := val.String() + "1"
				return contract.DelegationArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "failed - validator not found",
			malleate: func(val sdk.ValAddress, del common.Address) (contract.DelegationArgs, error) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()

				return contract.DelegationArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			operator0 := suite.GetFirstValAddr()
			delAmount := helpers.NewRandAmount()
			delAddr := suite.GetDelAddr()

			res := suite.DelegateV2(suite.Ctx, contract.DelegateV2Args{
				Validator: operator0.String(),
				Amount:    delAmount.BigInt(),
			})
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			args, expectErr := tc.malleate(operator0, delAddr)

			delegation := suite.GetDelegation(delAddr.Bytes(), operator0)

			delValue, _ := suite.WithError(expectErr).Delegation(suite.Ctx, args)
			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)
				suite.Require().Equal(delegation.GetShares().TruncateInt().String(), delValue.String())
			}
		})
	}
}
