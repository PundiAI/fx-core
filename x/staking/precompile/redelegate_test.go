package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingRedelegateABI(t *testing.T) {
	redelegateMethod := precompile.NewRedelegationMethod(nil)

	require.Equal(t, 3, len(redelegateMethod.Method.Inputs))
	require.Equal(t, 3, len(redelegateMethod.Method.Outputs))

	require.Equal(t, 6, len(redelegateMethod.Event.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestRedelegate() {
	testABI := contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
	redelegateMethod := precompile.NewRedelegationMethod(nil)
	redelegateV2Method := precompile.NewRedelegateV2Method(nil)

	testCases := []struct {
		name     string
		isV2     bool
		malleate func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.TruncateInt().BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator src",
			malleate: func(_, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.TruncateInt().BigInt(),
				}, fmt.Errorf("validator src not found: %s", valSrc.String())
			},
			result: false,
		},
		{
			name: "failed - invalid validator dst",
			malleate: func(valSrc, _ sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				valDst := sdk.ValAddress(suite.signer.Address().Bytes())
				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.TruncateInt().BigInt(),
				}, fmt.Errorf("validator dst not found: %s", valDst.String())
			},
			result: false,
		},
		{
			name: "failed - no delegation",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				// undelegate all before redelegate
				_, err := suite.App.StakingKeeper.Undelegate(suite.Ctx, suite.signer.AccAddress(), valSrc, shares)
				suite.Require().NoError(err)

				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.TruncateInt().BigInt(),
				}, fmt.Errorf("delegation not found")
			},
			result: false,
		},
		{
			name: "failed - insufficient redelegate shares",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.Add(sdk.NewDec(1e18)).TruncateInt().BigInt(),
				}, fmt.Errorf("insufficient shares to redelegate")
			},
			result: false,
		},
		{
			name: "failed - redelegate limit",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				entries := suite.App.StakingKeeper.MaxEntries(suite.Ctx)
				for i := uint32(0); i < entries; i++ {
					_, err := suite.App.StakingKeeper.BeginRedelegation(suite.Ctx,
						suite.signer.AccAddress(), valSrc, valDst, shares.QuoInt64(10))
					suite.Require().NoError(err)
				}

				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.QuoInt64(10).TruncateInt().BigInt(),
				}, fmt.Errorf("too many redelegation entries for (delegator, src-validator, dst-validator) tuple")
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.TruncateInt().BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator src",
			malleate: func(_, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.TruncateInt().BigInt(),
				}, fmt.Errorf("redelegate failed: validator src not found: %s", valSrc.String())
			},
			result: false,
		},
		{
			name: "contract - failed - invalid validator dst",
			malleate: func(valSrc, _ sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				valDst := sdk.ValAddress(suite.signer.Address().Bytes())

				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.TruncateInt().BigInt(),
				}, fmt.Errorf("redelegate failed: validator dst not found: %s", valDst.String())
			},
			result: false,
		},
		{
			name: "contract - failed - no delegation",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				// undelegate all before redelegate
				_, err := suite.App.StakingKeeper.Undelegate(suite.Ctx, suite.staking.Bytes(), valSrc, shares)
				suite.Require().NoError(err)

				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.TruncateInt().BigInt(),
				}, fmt.Errorf("redelegate failed: delegation not found")
			},
			result: false,
		},
		{
			name: "contract - failed - insufficient redelegate shares",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.Add(sdk.NewDec(1e18)).TruncateInt().BigInt(),
				}, fmt.Errorf("redelegate failed: insufficient shares to redelegate")
			},
			result: false,
		},
		{
			name: "contract - failed - redelegate limit",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				entries := suite.App.StakingKeeper.MaxEntries(suite.Ctx)
				for i := uint32(0); i < entries; i++ {
					_, err := suite.App.StakingKeeper.BeginRedelegation(suite.Ctx,
						suite.staking.Bytes(), valSrc, valDst, shares.QuoInt64(10))
					suite.Require().NoError(err)
				}

				return types.RedelegateArgs{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Shares:       shares.QuoInt64(10).TruncateInt().BigInt(),
				}, fmt.Errorf("redelegate failed: too many redelegation entries for (delegator, src-validator, dst-validator) tuple")
			},
			result: false,
		},

		{
			name: "ok v2",
			isV2: true,
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				return types.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - v2 invalid validator src",
			isV2: true,
			malleate: func(_, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				return types.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},

		{
			name: "contract - ok v2",
			isV2: true,
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				return types.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - v2 invalid validator src",
			isV2: true,
			malleate: func(_, valDst sdk.ValAddress, shares sdk.Dec, delAmount sdkmath.Int) (interface{}, error) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				return types.RedelegateV2Args{
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
			vals := suite.GetValidators()
			val0 := vals[0]
			val1 := vals[1]

			delAmt := helpers.NewRandAmount()

			stakingContract := precompile.GetAddress()
			stakingABI := precompile.GetABI()
			delAddr := suite.signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				stakingABI = testABI
				delAddr = suite.staking
			}

			// delegate to val0
			pack, err := stakingABI.Pack(TestDelegateName, val0.GetOperator().String())
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, delAmt.BigInt(), pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			delegation0, found := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), val0.GetOperator())
			suite.Require().True(found)
			_, found = suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), val1.GetOperator())
			suite.Require().False(found)

			val0, found = suite.App.StakingKeeper.GetValidator(suite.Ctx, val0.GetOperator())
			suite.Require().True(found)

			var packData []byte
			args, errResult := tc.malleate(val0.GetOperator(), val1.GetOperator(), delegation0.Shares, delAmt)
			if !tc.isV2 {
				packData, err = redelegateMethod.PackInput(args.(types.RedelegateArgs))
			} else {
				packData, err = redelegateV2Method.PackInput(args.(types.RedelegateV2Args))
			}
			suite.Require().NoError(err)

			if strings.HasPrefix(tc.name, "contract") {
				if !tc.isV2 {

					argsV1 := args.(types.RedelegateArgs)
					packData, err = testABI.Pack(TestRedelegateName, argsV1.ValidatorSrc, argsV1.ValidatorDst, argsV1.Shares)
				} else {
					argsV2 := args.(types.RedelegateV2Args)
					packData, err = testABI.Pack(TestRedelegateV2Name, argsV2.ValidatorSrc, argsV2.ValidatorDst, argsV2.Amount)
				}
				suite.Require().NoError(err)
			}
			res = suite.EthereumTx(suite.signer, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				if !tc.isV2 {
					unpack, err := stakingABI.Unpack(TestRedelegateName, res.Ret)
					suite.Require().NoError(err)
					// amount,reward,completionTime
					reward := unpack[1].(*big.Int)
					suite.Require().True(reward.Cmp(big.NewInt(0)) == 1, reward.String())
				}

				_, found = suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), val0.GetOperator())
				suite.Require().False(found)
				delegation1New, found := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), val1.GetOperator())
				suite.Require().True(found)
				suite.Require().Equal(delegation0.Shares, delegation1New.Shares)

				redelegates := suite.App.StakingKeeper.GetAllRedelegations(suite.Ctx, delAddr.Bytes(), val0.GetOperator(), val1.GetOperator())
				suite.Require().Equal(1, len(redelegates))

				suite.CheckRedelegateLogs(res.Logs, delAddr, val0.GetOperator().String(), val1.GetOperator().String(),
					delegation0.Shares.TruncateInt().BigInt(), val0.TokensFromShares(delegation0.Shares).TruncateInt().BigInt(),
					redelegates[0].Entries[0].CompletionTime.Unix())

				suite.CheckRedelegateEvents(suite.Ctx, val0.GetOperator().String(), val1.GetOperator().String(),
					val0.TokensFromShares(delegation0.Shares).TruncateInt().BigInt(),
					redelegates[0].Entries[0].CompletionTime)
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
